# 百家樂開發任務：結果判定與派彩 (Task 5.4)

**目標**：比較閒家和莊家最終手牌點數，確定勝負結果（閒贏、莊贏、和），檢查對子賭注，計算並分配玩家的輸贏。

---

## 5.4.1：比較手牌與確定勝負 (Backend: Go)

### 5.4.1.1：實現比較邏輯
-   **5.4.1.1.1-L6**: 在 `Table` 或 `GameManager` 中實現 `DetermineOutcome()` 方法，在 `DealingComplete` (5.3.3.3) 後調用。
-   **5.4.1.1.2-L6**: 獲取最終點數 `playerFinalTotal = PlayerHand.Total`, `bankerFinalTotal = BankerHand.Total`。
-   **5.4.1.1.3-L6**: 比較點數確定主要結果：
    -   `if playerFinalTotal > bankerFinalTotal`: 閒家贏 (`result = BetAreaPlayer`)。
    -   `else if bankerFinalTotal > playerFinalTotal`: 莊家贏 (`result = BetAreaBanker`)。
    -   `else`: 和局 (`result = BetAreaTie`)。
-   **5.4.1.1.4-L6**: 將結果存儲在 `Table` 狀態中，例如 `Table.CurrentOutcome.MainResult = result`。

### 5.4.1.2：檢查對子結果
-   **5.4.1.2.1-L6**: 在 `DetermineOutcome()` 中，檢查閒家和莊家**前兩張牌**是否構成對子。
    -   `playerIsPair := PlayerHand.Cards[0].Rank == PlayerHand.Cards[1].Rank`
    -   `bankerIsPair := BankerHand.Cards[0].Rank == BankerHand.Cards[1].Rank`
-   **5.4.1.2.2-L6**: 將對子結果存儲在 `Table` 狀態中，例如 `Table.CurrentOutcome.PlayerPairWins = playerIsPair`, `Table.CurrentOutcome.BankerPairWins = bankerIsPair`。

### 5.4.1.3：廣播結果
-   **5.4.1.3.1-L6**: 廣播 `GameResult` 事件，包含所有結果信息。
    ```json
    {
      "event": "game_result",
      "data": {
        "main_result": "player", // "banker", "tie"
        "player_total": 8,
        "banker_total": 3,
        "player_pair_win": false,
        "banker_pair_win": false,
        "player_hand": ["S5", "C2", "D1"], // Optional: Final cards again
        "banker_hand": ["HA", "H3", "S9"]  // Optional: Final cards again
      }
    }
    ```
-   **5.4.1.3.2-L6**: (Frontend) 接收事件，高亮顯示獲勝區域 (閒/莊/和)，顯示最終點數，並標示對子是否中獎。

---

## 5.4.2：計算賠付 (Backend: Go)

### 5.4.2.1：實現賠付計算函數
-   **5.4.2.1.1-L6**: 在 `Table` 或 `GameManager` 中實現 `CalculatePayouts()` 方法，在 `DetermineOutcome()` 後調用。
-   **5.4.2.1.2-L6**: 遍歷所有下注的玩家 (`Table.Players`)。
-   **5.4.2.1.3-L6**: 對於每個玩家，遍歷其下注記錄 `player.CurrentBaccaratBets`。
    ```go
    playerTotalWinLoss := 0
    payoutDetails := []map[string]interface{}{} // To store details for broadcasting

    for betArea, betAmount := range player.CurrentBaccaratBets {
        winAmount := 0
        payoutRate := gameConfig.BaccaratPayouts[betArea]
        commissionRate := gameConfig.BankerCommissionRate

        switch betArea {
        case BetAreaPlayer:
            if table.CurrentOutcome.MainResult == BetAreaPlayer {
                winAmount = int(float64(betAmount) * payoutRate)
            }
        case BetAreaBanker:
            if table.CurrentOutcome.MainResult == BetAreaBanker {
                payout = float64(betAmount) * payoutRate // Payout rate already accounts for standard commission (e.g., 0.95)
                // OR: payout = float64(betAmount) * 1.0
                //     commission = payout * commissionRate
                //     winAmount = int(payout - commission)
                // Choose ONE consistent approach based on how payoutRate is defined. Assume payoutRate is 0.95 here.
                 winAmount = int(float64(betAmount) * payoutRate) // Assuming 0.95 payout implies commission included
            }
        case BetAreaTie:
            if table.CurrentOutcome.MainResult == BetAreaTie {
                 winAmount = int(float64(betAmount) * payoutRate) // e.g., 8.0
            }
        case BetAreaPlayerPair:
            if table.CurrentOutcome.PlayerPairWins {
                 winAmount = int(float64(betAmount) * payoutRate) // e.g., 11.0
            }
        case BetAreaBankerPair:
            if table.CurrentOutcome.BankerPairWins {
                 winAmount = int(float64(betAmount) * payoutRate) // e.g., 11.0
            }
        }

        netWinLoss := winAmount - betAmount // If lost, winAmount is 0, net is negative betAmount
        playerTotalWinLoss += netWinLoss

        // Add details for this bet area
        payoutDetails = append(payoutDetails, map[string]interface{}{
            "area": betArea,
            "bet": betAmount,
            "win": winAmount,
            "net": netWinLoss,
        })
    }
    // Store results for the player
    player.LastRoundWinLoss = playerTotalWinLoss
    player.LastRoundPayoutDetails = payoutDetails
    ```
-   **5.4.2.1.4-L6**: **重要**: 處理和局 (Tie) 時的閒/莊下注：通常是**退回本金** (Push)。
    -   在 `case BetAreaPlayer` 和 `case BetAreaBanker` 中添加檢查：
        ```go
        if table.CurrentOutcome.MainResult == BetAreaTie {
            winAmount = betAmount // Push - return original bet
            netWinLoss = 0        // Net change is zero
        } else if table.CurrentOutcome.MainResult == betArea { // Win condition
            winAmount = int(float64(betAmount) * payoutRate)
            netWinLoss = winAmount - betAmount
        } else { // Lose condition
            winAmount = 0
            netWinLoss = -betAmount
        }
        ```
    -   重新計算 `playerTotalWinLoss` 和 `payoutDetails`。

### 5.4.2.2：處理佣金 (如果適用)
-   **5.4.2.2.1-L6**: 確認賠率定義 (`BaccaratPayouts["banker"]`) 是否已包含佣金。
    -   如果賠率是 1:1 (或 1.0)，則需要在計算莊贏時手動扣除佣金。
    -   如果賠率是 0.95:1 (或 0.95)，則計算出的 `winAmount` 已包含佣金。
-   **5.4.2.2.2-L6**: 推薦使用已包含佣金的賠率 (0.95) 以簡化計算。
-   **5.4.2.2.3-L6**: 如果需要支持免佣模式，應在 `GameConfig` 中有標識，並在計算莊贏時使用不同的賠率或跳過佣金計算。

---

## 5.4.3：分配彩金與更新籌碼 (Backend: Go)

### 5.4.3.1：更新玩家籌碼
-   **5.4.3.1.1-L6**: 在 `CalculatePayouts()` 計算完每個玩家的 `playerTotalWinLoss` 後，更新其籌碼。
-   **5.4.3.1.2-L6**: `if playerTotalWinLoss > 0 { player.AddChips(playerTotalWinLoss) }`。注意：`AddChips` 應只增加贏得的金額，本金已在下注時扣除，並在計算 `netWinLoss` 時考慮了返還。更準確地說，應該是 `player.AddChips(player.OriginalBetTotal + playerTotalWinLoss)`，或者更簡單：直接 `player.AddChips(totalReturn)` 其中 `totalReturn = originalBet + netWin`。
    -   **修正**：最清晰的方式是，在計算循環中累加總返還金額 `totalReturnAmount`。
        ```go
        totalReturnAmount := 0
        for _, detail := range payoutDetails {
            if detail["net"].(int) >= 0 { // If it's a win or push
               totalReturnAmount += detail["bet"].(int) + detail["net"].(int) // Return bet + net win
            }
             // If net < 0 (loss), return nothing for this bet
        }
        player.AddChips(totalReturnAmount) // Add the total amount returned
        ```
-   **5.4.3.1.3-L6**: 記錄籌碼變動日誌。

### 5.4.3.2：廣播派彩結果
-   **5.4.3.2.1-L6**: 遍歷所有玩家，發送個性化的 `PayoutResult` 事件。
    ```json
    // Sent to each player individually or broadcast with player-specific details
    {
      "event": "payout_result",
      "data": {
         "player_id": "player123",
         "total_win_loss": 55, // Net amount won or lost (-ve for loss)
         "new_chip_balance": 1055,
         "details": [
           {"area": "player", "bet": 100, "win": 100, "net": 0}, // Push example
           {"area": "banker", "bet": 50, "win": 0, "net": -50}, // Loss example
           {"area": "player_pair", "bet": 5, "win": 55, "net": 50}  // Win example (11:1)
         ]
      }
    }
    ```
-   **5.4.3.2.2-L6**: (Frontend) 接收 `PayoutResult`，顯示玩家的總輸贏和各項目的明細。顯示籌碼從莊家飛向獲勝玩家的動畫。更新界面上的玩家籌碼餘額。

---

## 5.4.4：清理桌面 (Backend: Go)

### 5.4.4.1：清除本輪數據
-   **5.4.4.1.1-L6**: 在派彩完成後，實現 `ClearRoundData()` 方法。
-   **5.4.4.1.2-L6**: 清除閒家和莊家手牌 `PlayerHand = nil`, `BankerHand = nil` 或 `PlayerHand.Cards = []*card.Card{}`, etc.。
-   **5.4.4.1.3-L6**: 清除桌面和玩家的下注記錄 `Table.TotalBets = make(map[string]int)`, `player.CurrentBaccaratBets = make(map[string]int)`。
-   **5.4.4.1.4-L6**: 清除本輪結果 `Table.CurrentOutcome = nil`。
-   **5.4.4.1.5-L6**: 廣播 `RoundEnd` 或 `ReadyForNextRound` 事件。
-   **5.4.4.1.6-L6**: (Frontend) 清理桌面上的牌和籌碼（除了玩家自己的餘額籌碼）。準備下一輪下注。

--- 