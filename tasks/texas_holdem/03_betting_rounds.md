# 德州撲克開發任務：下注輪邏輯 (Task 3.3)

**目標**：實現德州撲克 Pre-Flop, Flop, Turn, River 四個下注輪的核心邏輯，包括玩家動作處理、下注驗證和輪次結束條件。

---

## 3.3.1：管理下注輪 (Betting Round) (Backend: Go)

### 3.3.1.1：定義下注輪狀態
-   **3.3.1.1.1-L6**: 在 `Table` struct 中定義管理下注輪的狀態變數。
    -   `CurrentBettingRound string` (例如 "Pre-Flop", "Flop", "Turn", "River")
    -   `CurrentPlayerToActSeatIndex int` (當前輪到哪個座位行動)
    -   `CurrentBetAmount int` (當前輪次需要跟注的最高金額)
    -   `LastRaiserSeatIndex int` (最後一個加注者的座位索引，用於判斷輪次結束)
    -   `PlayerBets map[int]int` (記錄每個座位在本**輪**下注的總金額)
-   **3.3.1.1.2-L6**: 實現 `Table.StartBettingRound(roundName string)` 方法。
    -   設置 `CurrentBettingRound`。
    -   重置 `CurrentBetAmount` 為 0 (除非是 Pre-Flop 輪，此時應為大盲注金額)。
    -   重置 `PlayerBets` map。
    -   確定首個行動玩家 (`CurrentPlayerToActSeatIndex`)。
        -   Pre-Flop: 大盲注左手邊的玩家。
        -   Flop, Turn, River: 小盲注左手邊的玩家 (或 Button 左手邊第一個未蓋牌玩家)。
    -   設置 `LastRaiserSeatIndex` 為初始行動玩家 (或特殊值表示尚未加注)。
    -   觸發輪次開始和輪到玩家行動的通知。

### 3.3.1.2：確定行動順序
-   **3.3.1.2.1-L6**: 實現 `Table.GetNextPlayerToAct(currentSeatIndex int)` 輔助函數。
    -   從 `currentSeatIndex` 開始順時針查找下一個狀態為 "active" 或 "all-in" (但仍需等待結算) 的玩家。
    -   處理座位循環。
    -   跳過狀態為 "folded" 或 "sitting_out" 的玩家。
-   **3.3.1.2.2-L6**: 在 `StartBettingRound` 中使用此函數確定首個行動玩家。
-   **3.3.1.2.3-L6**: 在玩家完成行動後，使用此函數更新 `CurrentPlayerToActSeatIndex`。

### 3.3.1.3：輪次結束條件
-   **3.3.1.3.1-L6**: 在每次玩家行動後，檢查輪次是否結束。
-   **3.3.1.3.2-L6**: 結束條件 1：只剩一個活躍玩家 (未蓋牌且未 All-in)。
    -   實現 `Table.CountActivePlayers()` 輔助函數。
-   **3.3.1.3.3-L6**: 結束條件 2：所有未蓋牌且未 All-in 的玩家都已行動，且他們的下注金額與 `CurrentBetAmount` 相等 (即所有人都跟注或過牌)。
    -   檢查 `CurrentPlayerToActSeatIndex` 是否回到了 `LastRaiserSeatIndex` (或初始行動玩家，如果沒有加注)。
    -   檢查所有活躍玩家的 `PlayerBets[seatIndex]` 是否等於 `CurrentBetAmount` (排除 All-in 金額不足的情況)。
-   **3.3.1.3.4-L6**: 如果輪次結束，調用 `Table.EndBettingRound()` 方法。
    -   收集本輪下注到主池/邊池 (見 3.3.4)。
    -   根據 `CurrentBettingRound` 決定下一步：發牌 (Flop/Turn/River) 或攤牌 (Showdown)。
    -   觸發輪次結束通知。

---

## 3.3.2：處理玩家動作 (Backend: Go)

### 3.3.2.1：接收玩家動作請求
-   **3.3.2.1.1-L6**: 定義 WebSocket 接收玩家動作的消息結構。
    -   `{"action": "player_action", "data": {"action_type": "fold" | "check" | "call" | "bet" | "raise", "amount": 0 | bet_amount}}`
-   **3.3.2.1.2-L6**: 在 WebSocket 處理邏輯中，驗證收到的請求。
    -   檢查是否輪到該玩家 (`requestingPlayerID` 對應的座位是否為 `CurrentPlayerToActSeatIndex`)。
    -   檢查遊戲狀態是否處於下注輪。

### 3.3.2.2：處理蓋牌 (Fold)
-   **3.3.2.2.1-L6**: 實現 `Table.HandleFold(playerSeatIndex int)` 方法。
-   **3.3.2.2.2-L6**: 將玩家狀態 (`Player.Status`) 設置為 "folded"。
-   **3.3.2.2.3-L6**: 廣播玩家蓋牌通知。
    -   `{"event": "player_acted", "data": {"seat_index": ..., "action": "fold"}}`
-   **3.3.2.2.4-L6**: 檢查輪次結束條件。
-   **3.3.2.2.5-L6**: 更新 `CurrentPlayerToActSeatIndex`。
-   **3.3.2.2.6-L6**: 觸發下一個玩家行動的通知。

### 3.3.2.3：處理過牌 (Check)
-   **3.3.2.3.1-L6**: 實現 `Table.HandleCheck(playerSeatIndex int)` 方法。
-   **3.3.2.3.2-L6**: **驗證過牌合法性**: 當前玩家在本輪的下注額 (`PlayerBets[playerSeatIndex]`) 必須等於當前輪最高下注額 (`CurrentBetAmount`)。 (通常意味著 `CurrentBetAmount` 為 0，或是大盲在 Pre-Flop 且無人加注)。
-   **3.3.2.3.3-L6**: 如果不合法，向玩家發送錯誤信息，不改變狀態。
-   **3.3.2.3.4-L6**: 如果合法，廣播玩家過牌通知。
    -   `{"event": "player_acted", "data": {"seat_index": ..., "action": "check"}}`
-   **3.3.2.3.5-L6**: 檢查輪次結束條件。
-   **3.3.2.3.6-L6**: 更新 `CurrentPlayerToActSeatIndex`。
-   **3.3.2.3.7-L6**: 觸發下一個玩家行動的通知。

### 3.3.2.4：處理跟注 (Call)
-   **3.3.2.4.1-L6**: 實現 `Table.HandleCall(playerSeatIndex int)` 方法。
-   **3.3.2.4.2-L6**: **驗證跟注合法性**: `CurrentBetAmount` 必須大於 0。
-   **3.3.2.4.3-L6**: 計算需要跟注的金額 `amountToCall = CurrentBetAmount - PlayerBets[playerSeatIndex]`。
-   **3.3.2.4.4-L6**: 獲取玩家剩餘籌碼 `Player.Chips`。
-   **3.3.2.4.5-L6**: 確定實際跟注金額 `actualCallAmount = min(amountToCall, Player.Chips)`。
-   **3.3.2.4.6-L6**: 從玩家移除籌碼 `Player.RemoveChips(actualCallAmount)`。
-   **3.3.2.4.7-L6**: 更新玩家本輪下注 `PlayerBets[playerSeatIndex] += actualCallAmount`。
-   **3.3.2.4.8-L6**: 檢查是否 All-in (`Player.Chips == 0`)，更新 `Player.Status`。
-   **3.3.2.4.9-L6**: 廣播玩家跟注通知。
    -   `{"event": "player_acted", "data": {"seat_index": ..., "action": "call", "amount": actualCallAmount}}`
-   **3.3.2.4.10-L6**: 廣播玩家籌碼更新通知。
-   **3.3.2.4.11-L6**: 檢查輪次結束條件。
-   **3.3.2.4.12-L6**: 更新 `CurrentPlayerToActSeatIndex`。
-   **3.3.2.4.13-L6**: 觸發下一個玩家行動的通知。

### 3.3.2.5：處理下注 (Bet)
-   **3.3.2.5.1-L6**: 實現 `Table.HandleBet(playerSeatIndex int, betAmount int)` 方法。
-   **3.3.2.5.2-L6**: **驗證下注合法性**:
    -   `CurrentBetAmount` 必須為 0。
    -   `betAmount` 必須大於等於最小下注額 (通常是大盲注金額)。
    -   `betAmount` 必須小於等於玩家剩餘籌碼 `Player.Chips`。
-   **3.3.2.5.3-L6**: 如果不合法，發送錯誤信息。
-   **3.3.2.5.4-L6**: 如果合法，從玩家移除籌碼 `Player.RemoveChips(betAmount)`。
-   **3.3.2.5.5-L6**: 更新玩家本輪下注 `PlayerBets[playerSeatIndex] = betAmount`。
-   **3.3.2.5.6-L6**: 更新當前輪最高下注額 `CurrentBetAmount = betAmount`。
-   **3.3.2.5.7-L6**: 更新最後加注者 `LastRaiserSeatIndex = playerSeatIndex`。
-   **3.3.2.5.8-L6**: 檢查是否 All-in，更新 `Player.Status`。
-   **3.3.2.5.9-L6**: 廣播玩家下注通知。
    -   `{"event": "player_acted", "data": {"seat_index": ..., "action": "bet", "amount": betAmount}}`
-   **3.3.2.5.10-L6**: 廣播玩家籌碼更新通知。
-   **3.3.2.5.11-L6**: 更新 `CurrentPlayerToActSeatIndex`。
-   **3.3.2.5.12-L6**: 觸發下一個玩家行動的通知。

### 3.3.2.6：處理加注 (Raise)
-   **3.3.2.6.1-L6**: 實現 `Table.HandleRaise(playerSeatIndex int, raiseAmount int)` 方法。 (`raiseAmount` 指的是加注**到**的總金額，不是增加的金額)。
-   **3.3.2.6.2-L6**: **驗證加注合法性**:
    -   `CurrentBetAmount` 必須大於 0。
    -   `raiseAmount` 必須大於等於最小加注額 (通常是 `CurrentBetAmount + (CurrentBetAmount - PreviousBetAmount)`，至少是翻倍加注，但有最小加注額限制，通常是 BB)。
    -   `raiseAmount` 必須小於等於玩家剩餘籌碼加上玩家本輪已下注額 (`Player.Chips + PlayerBets[playerSeatIndex]`)。
-   **3.3.2.6.3-L6**: 如果不合法，發送錯誤信息。
-   **3.3.2.6.4-L6**: 計算實際需要投入的籌碼 `chipsToCommit = raiseAmount - PlayerBets[playerSeatIndex]`。
-   **3.3.2.6.5-L6**: 從玩家移除籌碼 `Player.RemoveChips(chipsToCommit)`。
-   **3.3.2.6.6-L6**: 更新玩家本輪下注 `PlayerBets[playerSeatIndex] = raiseAmount`。
-   **3.3.2.6.7-L6**: 更新當前輪最高下注額 `CurrentBetAmount = raiseAmount`。
-   **3.3.2.6.8-L6**: 更新最後加注者 `LastRaiserSeatIndex = playerSeatIndex`。
-   **3.3.2.6.9-L6**: 檢查是否 All-in，更新 `Player.Status`。
-   **3.3.2.6.10-L6**: 廣播玩家加注通知。
    -   `{"event": "player_acted", "data": {"seat_index": ..., "action": "raise", "amount": raiseAmount}}`
-   **3.3.2.6.11-L6**: 廣播玩家籌碼更新通知。
-   **3.3.2.6.12-L6**: 更新 `CurrentPlayerToActSeatIndex`。
-   **3.3.2.6.13-L6**: 觸發下一個玩家行動的通知。

### 3.3.2.7：通知當前行動玩家
-   **3.3.2.7.1-L6**: 在確定下一個行動玩家後，向其發送通知。
    -   `{"event": "your_turn", "data": {"seat_index": ..., "time_limit": 15, "min_raise": ..., "amount_to_call": ...}}`
-   **3.3.2.7.2-L6**: 計算 `min_raise` (最小加注額) 和 `amount_to_call` (需跟注額)。
-   **3.3.2.7.3-L6**: (前端任務) 接收此通知，高亮玩家座位，顯示可用動作按鈕和倒計時。

---

## 3.3.3：處理 All-in (Backend: Go)

### 3.3.3.1：識別 All-in 情況
-   **3.3.3.1.1-L6**: 在處理 Call, Bet, Raise 時，檢查玩家 `Player.Chips` 是否變為 0。
-   **3.3.3.1.2-L6**: 如果 `Player.Chips == 0`，將 `Player.Status` 設置為 "all-in"。

### 3.3.3.2：All-in 對下注輪的影響
-   **3.3.3.2.1-L6**: All-in 玩家不再參與後續的下注輪次，但其籌碼仍在底池中。
-   **3.3.3.2.2-L6**: 如果 All-in 金額小於當前的 `CurrentBetAmount` 或小於最小加注額：
    -   不會增加 `CurrentBetAmount`。
    -   不會成為 `LastRaiserSeatIndex`。
    -   其他玩家只需跟注到之前的 `CurrentBetAmount`。
    -   **需要觸發邊池 (Side Pot) 計算邏輯 (見 3.3.4)**。

### 3.3.3.3：All-in 通知
-   **3.3.3.3.1-L6**: 在廣播玩家動作通知時，如果該動作導致 All-in，可以在消息中添加標記。
    -   `{"event": "player_acted", "data": {"seat_index": ..., "action": "call", "amount": ..., "is_all_in": true}}`
-   **3.3.3.3.2-L6**: (前端任務) 根據 All-in 狀態顯示特殊標記。

---

## 3.3.4：收集籌碼與創建邊池 (Backend: Go)

### 3.3.4.1：輪次結束時收集籌碼
-   **3.3.4.1.1-L6**: 在 `Table.EndBettingRound()` 方法中實現。
-   **3.3.4.1.2-L6**: 遍歷所有參與本輪下注的玩家 (包括蓋牌和 All-in 玩家)。
-   **3.3.4.1.3-L6**: 將每個玩家的 `PlayerBets[seatIndex]` 金額添加到獎池邏輯中 (觸發 `Pot.CollectBets(playerBets map[int]int)` 方法)。

### 3.3.4.2：邊池 (Side Pot) 創建邏輯
-   **3.3.4.2.1-L6**: 在 `Pot.CollectBets()` 方法中實現 (或在玩家 All-in 時觸發計算)。
-   **3.3.4.2.2-L6**: 識別出所有 All-in 玩家及其 All-in 金額。
-   **3.3.4.2.3-L6**: 按 All-in 金額從小到大排序。
-   **3.3.4.2.4-L6**: 創建邊池：
    -   第一個邊池 (或主池，如果沒有更小的 All-in) 的額度 = 最小 All-in 玩家的總下注額 * 參與該池的玩家數。
    -   參與者：所有下注額 >= 該 All-in 玩家總下注額的玩家。
    -   第二個邊池的額度 = (第二小 All-in 玩家總下注額 - 第一小 All-in 玩家總下注額) * 參與該池的玩家數。
    -   參與者：所有下注額 >= 第二小 All-in 玩家總下注額的玩家。
    -   以此類推。
    -   剩餘的籌碼構成主池 (Main Pot)，所有未蓋牌的玩家都有資格贏取。
-   **3.3.4.2.5-L6**: 更新 `Pot.MainPotAmount` 和 `Pot.SidePots` 列表 (`SidePot` 包含 `Amount` 和 `EligiblePlayerIDs`)。
-   **3.3.4.2.6-L6**: 廣播獎池更新通知 (包含主池和所有邊池信息)。
    -   `{"event": "pot_update", "data": {"main_pot": ..., "side_pots": [{"amount": ..., "eligible_seat_indices": [...]}, ...]}}`

### 3.3.4.3：重置輪次狀態
-   **3.3.4.3.1-L6**: 在 `EndBettingRound` 中，獎池收集完成後，重置 `PlayerBets` map (清空或歸零)。
-   **3.3.4.3.2-L6**: 重置 `CurrentBetAmount` 和 `LastRaiserSeatIndex`，準備下一輪或攤牌。

--- 