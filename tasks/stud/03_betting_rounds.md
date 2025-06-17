# 七張梭哈開發任務：下注輪邏जिक (Task 4.3)

**目標**：實現七張梭哈各個下注輪 (Third Street 到 Seventh Street) 的核心邏輯，包括 Bring-in、大小注切換、玩家動作處理和輪次結束。

---

## 4.3.1：管理下注輪 (Betting Round) (Backend: Go)

### 4.3.1.1：定義下注輪狀態
-   **4.3.1.1.1-L6**: 重用 Hold'em 的 `Table` 狀態變數，但需調整其含義和初始值。
    -   `CurrentBettingRound string` (例如 "Third Street", "Fourth Street", ...)
    -   `CurrentPlayerToActSeatIndex int`
    -   `CurrentBetAmount int` (梭哈中是**本輪**相對下注額，不是累計)
    -   `LastBettorOrRaiserSeatIndex int` (記錄最後一個下注或加注者)
    -   `PlayerBetsInCurrentRound map[int]int` (記錄每個座位在**本輪**投入的總額)
    -   `CurrentBetLimit int` (記錄本輪是使用 `SmallBetLimit` 還是 `BigBetLimit`)
-   **4.3.1.1.2-L6**: 實現 `Table.StartStudBettingRound(roundName string)` 方法。
    -   設置 `CurrentBettingRound`。
    -   重置 `CurrentBetAmount` 為 0。
    -   重置 `PlayerBetsInCurrentRound` map。
    -   設置 `CurrentBetLimit` (Third/Fourth Street 用 `SmallBetLimit`, Fifth/Sixth/Seventh 用 `BigBetLimit`)。
    -   確定首個行動玩家 (`CurrentPlayerToActSeatIndex`)。
        -   **Third Street**: Bring-in 玩家 (`BringInSeatIndex`)，他必須先行動。
        -   **Fourth Street 及之後**: 亮牌牌面**最大**的玩家先行動。如果最大牌面有多人，則按座位順序 (從 Bring-in 玩家之後算起) 決定先手。
    -   設置 `LastBettorOrRaiserSeatIndex` 為初始行動玩家 (或特殊值)。
    -   觸發輪次開始和輪到玩家行動的通知。

### 4.3.1.2：確定行動順序 (Fourth Street 及之後)
-   **4.3.1.2.1-L6**: 實現 `Table.DetermineFirstPlayerToAct(round string)` 輔助函數 (用於 Fourth Street 及之後)。
-   **4.3.1.2.2-L7**: 邏輯：
    -   獲取所有未蓋牌的玩家列表。
    -   對每個玩家，找出其當前所有亮牌 (`Player.VisibleCards`) 中能組成的最佳牌面 (例如一對、兩對、高牌等，但**只看亮牌**)。
    -   比較這些玩家的亮牌牌面等級 (需要一個簡化的 `EvaluateVisibleHand` 函數)。
    -   如果牌面等級不同，最高者先行動。
    -   如果牌面等級相同，則比較關鍵牌 Rank (例如 A 對子 > K 對子)。
    -   如果關鍵牌 Rank 仍相同，則按座位順序決定 (從 Bring-in 玩家之後算起，順時針第一個達到該牌面的玩家先行動)。
-   **4.3.1.2.3-L6**: 在 `StartStudBettingRound` 中調用此函數設置 `CurrentPlayerToActSeatIndex`。
-   **4.3.1.2.4-L6**: 玩家完成行動後，下一個行動者是其左手邊第一個未蓋牌的玩家 (使用類似 Hold'em 的 `GetNextPlayerToAct`)。

### 4.3.1.3：輪次結束條件
-   **4.3.1.3.1-L6**: 重用 Hold'em 的輪次結束判斷邏輯，但注意梭哈的下注結構。
-   **4.3.1.3.2-L6**: 條件 1：只剩一個未蓋牌玩家。
-   **4.3.1.3.3-L6**: 條件 2：所有未蓋牌且未 All-in 的玩家都已行動，且他們在本輪投入的總額 (`PlayerBetsInCurrentRound`) 都相等 (且等於最後一個 Bet/Raise 後的金額)。
    -   檢查 `CurrentPlayerToActSeatIndex` 是否回到了 `LastBettorOrRaiserSeatIndex` (或首個行動玩家，如果沒有 Bet/Raise)。
-   **4.3.1.3.4-L6**: 如果輪次結束，調用 `Table.EndStudBettingRound()`。
    -   收集本輪下注到獎池 (邏輯類似，但基於 `PlayerBetsInCurrentRound`)。
    -   根據 `CurrentBettingRound` 決定下一步：發牌或攤牌。
    -   觸發輪次結束通知。

---

## 4.3.2：處理 Bring-in (Third Street 特有) (Backend: Go)

### 4.3.2.1：強制 Bring-in 下注
-   **4.3.2.1.1-L6**: 在 `StartStudBettingRound("Third Street")` 確定 Bring-in 玩家後，觸發其行動。
-   **4.3.2.1.2-L6**: Bring-in 玩家**必須**下注，有兩種選擇：
    -   選擇 1：下注 `BringInAmount`。
    -   選擇 2：**Complete the bet**，直接下注 `SmallBetLimit`。
-   **4.3.2.1.3-L6**: 在發送給 Bring-in 玩家的 `your_turn` 通知中，需要包含這兩個選項和對應金額。
-   **4.3.2.1.4-L6**: 處理 Bring-in 玩家的動作請求 (強制 Bet)。
    -   驗證其下注金額是否為 `BringInAmount` 或 `SmallBetLimit`。
    -   從其籌碼扣除，更新 `PlayerBetsInCurrentRound`。
    -   更新 `CurrentBetAmount` (等於玩家下注額)。
    -   設置 `LastBettorOrRaiserSeatIndex` 為 Bring-in 玩家。
    -   檢查 All-in。
    -   廣播其動作 (`action: "bring_in"` 或 `"complete_bet"`?)。
    -   更新 `CurrentPlayerToActSeatIndex` 到其左手邊玩家。
    -   觸發下一玩家行動通知。

---

## 4.3.3：處理玩家動作 (通用輪次) (Backend: Go)

### 4.3.3.1：接收玩家動作
-   **4.3.3.1.1-L6**: 重用 Hold'em 的 WebSocket 消息結構: `{"action": "player_action", "data": {"action_type": ..., "amount": ...}}`。
-   **4.3.3.1.2-L6**: 驗證請求玩家和遊戲狀態。

### 4.3.3.2：處理蓋牌 (Fold)
-   **4.3.3.2.1-L6**: 重用 Hold'em 的 `HandleFold` 邏輯。
    -   設置狀態 "folded"。
    -   廣播通知。
    -   檢查結束條件。
    -   更新並通知下一玩家。

### 4.3.3.3：處理過牌 (Check)
-   **4.3.3.3.1-L6**: 實現 `Table.HandleStudCheck(playerSeatIndex int)`。
-   **4.3.3.3.2-L6**: **驗證過牌合法性**: 當前輪 `CurrentBetAmount` 必須為 0 (即前面沒有人下注)。
-   **4.3.3.3.3-L6**: 其他邏輯同 Hold'em 的 `HandleCheck`。

### 4.3.3.4：處理跟注 (Call)
-   **4.3.3.4.1-L6**: 實現 `Table.HandleStudCall(playerSeatIndex int)`。
-   **4.3.3.4.2-L6**: **驗證跟注合法性**: `CurrentBetAmount` 必須大於 0。
-   **4.3.3.4.3-L6**: 計算需要跟注的金額 `amountToCall = CurrentBetAmount - PlayerBetsInCurrentRound[playerSeatIndex]`。
-   **4.3.3.4.4-L6**: 其他邏輯（扣籌碼、更新 PlayerBets、處理 All-in、廣播、更新下一玩家）同 Hold'em 的 `HandleCall`。

### 4.3.3.5：處理下注 (Bet)
-   **4.3.3.5.1-L6**: 實現 `Table.HandleStudBet(playerSeatIndex int)`。
-   **4.3.3.5.2-L6**: **驗證下注合法性**:
    -   `CurrentBetAmount` 必須為 0。
    -   下注金額必須等於當前輪的下注限額 (`CurrentBetLimit`，即 `SmallBetLimit` 或 `BigBetLimit`)。**注意：固定限注！**
    -   玩家有足夠籌碼。
-   **4.3.3.5.3-L6**: 如果合法：
    -   扣除 `CurrentBetLimit` 籌碼。
    -   更新 `PlayerBetsInCurrentRound[playerSeatIndex] = CurrentBetLimit`。
    -   更新 `CurrentBetAmount = CurrentBetLimit`。
    -   更新 `LastBettorOrRaiserSeatIndex`。
    -   處理 All-in。
    -   廣播通知。
    -   更新並通知下一玩家。

### 4.3.3.6：處理加注 (Raise)
-   **4.3.3.6.1-L6**: 實現 `Table.HandleStudRaise(playerSeatIndex int)`。
-   **4.3.3.6.2-L6**: **驗證加注合法性**:
    -   `CurrentBetAmount` 必須大於 0。
    -   加注後的總金額必須是 `CurrentBetAmount + CurrentBetLimit`。**注意：固定限注！**
    -   玩家有足夠籌碼完成此次加注 (需要投入 `CurrentBetLimit` 的籌碼)。
    -   通常有限制每輪加注次數 (例如 1 bet + 3 raises)，需要在 `Table` 狀態中添加 `RaiseCountInCurrentRound` 並檢查。**確認：是否需要限制加注次數？** 假設限制為 3 次加注。
-   **4.3.3.6.3-L6**: 如果合法：
    -   計算需投入籌碼 `chipsToCommit = CurrentBetAmount + CurrentBetLimit - PlayerBetsInCurrentRound[playerSeatIndex]`。
    -   扣除籌碼。
    -   更新 `PlayerBetsInCurrentRound[playerSeatIndex] = CurrentBetAmount + CurrentBetLimit`。
    -   更新 `CurrentBetAmount = CurrentBetAmount + CurrentBetLimit`。
    -   更新 `LastBettorOrRaiserSeatIndex`。
    -   增加 `RaiseCountInCurrentRound`。
    -   處理 All-in。
    -   廣播通知 (action: "raise", amount: `CurrentBetAmount`)。
    -   更新並通知下一玩家。

### 4.3.3.7：特殊規則：第四街對子開牌 (Open Pair on Fourth Street)
-   **4.3.3.7.1-L6**: 在第四街發牌後 (`DealFourthStreet` 完成後)，檢查是否有玩家的兩張亮牌組成一對 (Open Pair)。
-   **4.3.3.7.2-L6**: 如果有 Open Pair，則該輪第一個行動的玩家可以選擇下小注 (`SmallBetLimit`) 或 大注 (`BigBetLimit`)。
-   **4.3.3.7.3-L6**: 需要修改第四街 `StartStudBettingRound` 的邏जिक：
    -   檢測 Open Pair。
    -   如果檢測到，確定第一個行動玩家 (亮牌最高者)。
    -   在發送給該玩家的 `your_turn` 通知中，提供 Bet Small 和 Bet Big 兩個選項。
    -   修改 `HandleStudBet` (僅第四街) 以接受兩種下注額。
    -   一旦有人下注 (無論大小)，後續加注都以該下注額為基礎單位 (`CurrentBetLimit` 設置為玩家選擇的下注額)。

### 4.3.3.8：通知當前行動玩家
-   **4.3.3.8.1-L6**: 在確定下一個行動玩家後發送 `your_turn` 通知。
-   **4.3.3.8.2-L6**: 通知數據中需要包含：
    -   `seat_index`
    -   `time_limit` (如果有的話)
    -   `amount_to_call = CurrentBetAmount - PlayerBetsInCurrentRound[seatIndex]`
    -   `can_check` (布爾值, `CurrentBetAmount == PlayerBetsInCurrentRound[seatIndex]`)
    -   `can_bet` (布爾值, `CurrentBetAmount == 0`)
    -   `bet_amount = CurrentBetLimit`
    -   `can_raise` (布爾值, `CurrentBetAmount > 0` 且 `RaiseCountInCurrentRound < MaxRaises`)
    -   `raise_amount = CurrentBetAmount + CurrentBetLimit` (加注到的總額)
    -   (第四街特有) `can_bet_big` (如果適用)
    -   (第四街特有) `big_bet_amount = BigBetLimit`
-   **4.3.3.8.3-L6**: (Frontend) 根據這些數據動態生成可用的操作按鈕。

---

## 4.3.4：收集籌碼與邊池 (Backend: Go)

### 4.3.4.1：輪次結束時收集籌碼
-   **4.3.4.1.1-L6**: 在 `Table.EndStudBettingRound()` 中實現。
-   **4.3.4.1.2-L6**: 遍歷 `PlayerBetsInCurrentRound` map。
-   **4.3.4.1.3-L6**: 將每個玩家本輪投入的籌碼添加到獎池邏輯中 (觸發 `Pot.CollectBets(playerBets map[int]int)` 方法，此方法需要能處理多輪累積下注以正確計算邊池)。

### 4.3.4.2：邊池創建邏輯
-   **4.3.4.2.1-L6**: 重用 Hold'em 的邊池計算邏輯 (`Pot.CollectBets` 内部實現)。
    -   關鍵是 `CollectBets` 需要知道每個玩家**到目前為止**的總投入，而不僅僅是當前輪。因此 `Table` 需要維護一個 `PlayerTotalBetsInHand map[int]int`，在每輪結束收集籌碼時更新。`CollectBets` 應基於這個總投入來計算邊池。
-   **4.3.4.2.2-L6**: 廣播獎池更新通知 (同 Hold'em)。

### 4.3.4.3：重置輪次狀態
-   **4.3.4.3.1-L6**: 在 `EndStudBettingRound` 中，獎池收集完成後：
    -   重置 `PlayerBetsInCurrentRound` map。
    -   重置 `CurrentBetAmount`。
    -   重置 `RaiseCountInCurrentRound`。
    -   重置 `LastBettorOrRaiserSeatIndex`。

--- 