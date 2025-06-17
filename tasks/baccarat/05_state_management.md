# 百家樂開發任務：狀態管理 (Task 5.5)

**目標**：管理百家樂遊戲的整體流程、玩家狀態、牌局和牌靴的生命週期，並考慮數據持久化。

---

## 5.5.1：遊戲狀態機 (Backend: Go)

### 5.5.1.1：定義遊戲狀態
-   **5.5.1.1.1-L6**: 在 `game/baccarat/table.go` 或 `game/common/statemachine.go` 中定義遊戲狀態枚舉/常量。
    ```go
    type GameState string
    const (
        StateIdle             GameState = "idle"             // 初始狀態或換靴後
        StateWaitingForBets   GameState = "waiting_for_bets" // 剛結束一局，等待新局開始（可選，可能直接進入 BettingOpen）
        StateBettingOpen      GameState = "betting_open"     // 開放接受下注
        StateBettingClosed    GameState = "betting_closed"   // 停止接受下注
        StateDealingInitial   GameState = "dealing_initial"  // 正在發初始四張牌
        StatePlayerTurn       GameState = "player_turn"      // 決定閒家是否補牌
        StateBankerTurn       GameState = "banker_turn"      // 決定莊家是否補牌
        StateCalculatingResult GameState = "calculating_result" // 正在比較點數和對子
        StatePayingOut        GameState = "paying_out"       // 正在計算和分配派彩
        StateRoundOver        GameState = "round_over"       // 本局結束，顯示結果一段時間
        StateShoeChangeNeeded GameState = "shoe_change_needed" // 達到切牌點，需要換靴
        StatePaused           GameState = "paused"           // 遊戲暫停（管理員操作等）
    )
    // Table struct needs a field: CurrentState GameState
    ```

### 5.5.1.2：實現狀態轉換邏輯
-   **5.5.1.2.1-L6**: 創建一個 `Table.SetState(newState GameState)` 方法，內部處理狀態轉換和觸發相關事件/動作。
-   **5.5.1.2.2-L6**: 定義狀態轉換圖：
    -   `Idle` -> `BettingOpen` (開始新局)
    -   `BettingOpen` -> `BettingClosed` (計時器結束)
    -   `BettingClosed` -> `DealingInitial` (開始發牌)
    -   `DealingInitial` -> `PlayerTurn` (發完四張，檢查天牌後) / `CalculatingResult` (如果天牌)
    -   `PlayerTurn` -> `BankerTurn` (閒家動作完成)
    -   `BankerTurn` -> `CalculatingResult` (莊家動作完成)
    -   `CalculatingResult` -> `PayingOut` (結果確定)
    -   `PayingOut` -> `RoundOver` (派彩完成)
    -   `RoundOver` -> `BettingOpen` (短暫延遲後，如果牌靴還有牌) / `ShoeChangeNeeded` (如果達到切牌點)
    -   `ShoeChangeNeeded` -> `Idle` (換靴完成後)
    -   Any -> `Paused` -> Any (管理員控制)
-   **5.5.1.2.3-L6**: 在每個遊戲階段的邏輯完成後，調用 `SetState` 切換到下一個狀態。例如，`CloseBetting()` 應調用 `SetState(StateBettingClosed)`，然後觸發 `DealInitialHands()`。
-   **5.5.1.2.4-L6**: 廣播 `GameStateChanged` 事件，包含新舊狀態和相關數據（如下注結束時間）。
    ```json
    {"event": "game_state_changed", "data": {"new_state": "betting_open", "previous_state": "round_over", "betting_ends_at": "..."}}
    ```
-   **5.5.1.2.5-L6**: (Frontend) 根據遊戲狀態更新 UI（例如，顯示 "請下注", "停止下注", "開牌中", "結算中", 禁用/啟用下注按鈕）。

---

## 5.5.2：玩家狀態管理 (Backend: Go)

### 5.5.2.1：定義玩家狀態
-   **5.5.2.1.1-L6**: 重用 `pkg/player` 中的基礎玩家狀態 (`Connected`, `Authenticated`, `Disconnected`)。
-   **5.5.2.1.2-L6**: 添加特定於遊戲桌的狀態。
    ```go
    type PlayerTableStatus string
    const (
        StatusSpectating PlayerTableStatus = "spectating" // 觀戰
        StatusSeated     PlayerTableStatus = "seated"     // 入座，可以下注
        StatusBetting    PlayerTableStatus = "betting"    // 已下注，等待結果
        StatusWaiting    PlayerTableStatus = "waiting"    // 已入座，但本輪未下注或下注已結算
    )
    // Player struct might have: TableStatus PlayerTableStatus
    ```
-   **5.5.2.1.3-L6**: 玩家加入桌子時設置為 `StatusSeated` 或 `StatusSpectating`。
-   **5.5.2.1.4-L6**: 玩家成功下注後可標記為 `StatusBetting`。
-   **5.5.2.1.5-L6**: 派彩結束後，`StatusBetting` 的玩家變回 `StatusSeated` 或 `StatusWaiting`。

### 5.5.2.2：處理連接與斷線
-   **5.5.2.2.1-L6**: 重用通用的連接管理邏輯 (WebSocket/gRPC 連接建立與斷開)。
-   **5.5.2.2.2-L6**: 玩家斷線時：
    -   標記玩家狀態為 `Disconnected`。
    -   保留其在桌上的座位和籌碼一段時間（可配置）。
    -   如果玩家在斷線前有下注 (`CurrentBaccaratBets` 非空)，該下注**仍然有效**，遊戲繼續。
    -   廣播 `PlayerDisconnected` 事件。
-   **5.5.2.2.3-L6**: 玩家重連時：
    -   驗證身份，恢復連接。
    -   標記玩家狀態為 `Connected`。
    -   發送當前的完整遊戲狀態 (`GameState`, 桌面情況, 玩家自己的下注和籌碼) 給該玩家。
    -   廣播 `PlayerReconnected` 事件。
-   **5.5.2.2.4-L6**: 處理長時間斷線玩家的清理邏輯（例如，超時後自動離桌，退還未結算的下注？ - 需要確認規則，通常下注有效）。

---

## 5.5.3：牌局與牌靴管理 (Backend: Go)

### 5.5.3.1：生成唯一 ID
-   **5.5.3.1.1-L6**: 為每一局 (Round) 生成唯一 ID (例如 UUID 或基於時間戳+序列號)。`Table.CurrentRoundID`。
-   **5.5.3.1.2-L6**: 為每一靴 (Shoe) 生成唯一 ID。`Table.CurrentShoeID`。
-   **5.5.3.1.3-L6**: 在開始新局/新靴時更新這些 ID 並廣播。

### 5.5.3.2：追蹤牌靴進度
-   **5.5.3.2.1-L6**: 在 `Deck` 結構中維護 `initialCardCount` 和 `currentIndex`。
-   **5.5.3.2.2-L6**: 每次 `Deck.DrawCard()` 後 `currentIndex` 增加。
-   **5.5.3.2.3-L6**: 在發牌前檢查剩餘牌數 `Deck.RemainingCards()`。
-   **5.5.3.2.4-L6**: 檢查是否達到切牌點：`if deck.currentIndex >= deck.initialCardCount - gameConfig.BaccaratCuttingCardPosition`。

### 5.5.3.3：換靴邏輯 (Reshuffle/New Shoe)
-   **5.5.3.3.1-L6**: 當達到切牌點時，在 `RoundOver` 狀態後觸發 `SetState(StateShoeChangeNeeded)`。
-   **5.5.3.3.2-L6**: 廣播 `ShoeChangeNeeded` 事件。
-   **5.5.3.3.3-L6**: (Frontend) 顯示 "即將換靴" 的提示。
-   **5.5.3.3.4-L6**: (Backend) 實現 `StartNewShoe()` 方法：
    -   生成新的 `ShoeID`。
    -   創建新的 `Deck` 實例 (`NewDeck(8)`)。
    -   調用 `Deck.Shuffle()`。
    -   執行燒牌邏輯 (5.1.3.4)。
    -   廣播 `NewShoeStarted` 事件，包含新的 `ShoeID` 和燒牌信息。
    -   調用 `SetState(StateBettingOpen)` 開始新靴的第一局。

---

## 5.5.4：數據持久化 (Backend: Go, Database)

### 5.5.4.1：持久化考慮
-   **5.5.4.1.1-L6**: 確定需要持久化的數據：
    -   玩家賬戶信息和餘額 (核心系統，非遊戲邏輯)。
    -   牌局歷史 (Hand History): 每局的 ID, Shoe ID, 閒莊牌, 結果, 玩家下注, 輸贏。
    -   牌靴信息 (可選): Shoe ID, 開始/結束時間。
    -   遊戲配置 (已在 5.1.1 處理)。
-   **5.5.4.1.2-L6**: 選擇數據庫 (例如 PostgreSQL, MongoDB)。
-   **5.5.4.1.3-L6**: 設計數據庫表結構/文檔模型。
    -   `baccarat_rounds` (round_id, shoe_id, start_time, end_time, player_cards, banker_cards, player_total, banker_total, main_result, player_pair_win, banker_pair_win)
    -   `player_round_bets` (round_id, player_id, bet_area, bet_amount, win_amount, net_win_loss)
    -   `player_chip_ledger` (記錄玩家籌碼變動，可能由賬戶系統處理)

### 5.5.4.2：保存數據
-   **5.5.4.2.1-L6**: 在 `PayingOut` 或 `RoundOver` 狀態時，異步地將本局結果和玩家下注/輸贏數據寫入數據庫。
-   **5.5.4.2.2-L6**: 考慮使用消息隊列 (Kafka, RabbitMQ) 將寫入操作解耦，避免阻塞遊戲循環。
-   **5.5.4.2.3-L6**: 玩家籌碼變動應通過事務性操作更新到玩家賬戶餘額。

### 5.5.4.3：服務器重啟恢復 (可選，複雜度高)
-   **5.5.4.3.1-L6**: 策略 1 (簡單): 重啟後所有進行中的遊戲結束，結算當前牌局（如果可能），未下注的返回。不恢復遊戲狀態。
-   **5.5.4.3.2-L6**: 策略 2 (複雜): 定期將 `Table` 的完整狀態 (包括玩家、牌靴、當前狀態、下注) 快照到 Redis 或內存數據庫。重啟後嘗試從快照恢復。需要處理狀態一致性問題。
-   **5.5.4.3.3-L6**: 對於非錦標賽的現金桌，策略 1 通常足夠。重點是確保玩家資金正確結算。

--- 