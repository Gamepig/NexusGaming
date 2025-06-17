# 百家樂開發任務：下注階段 (Task 5.2)

**目標**：實現百家樂的下注流程，允許玩家在指定時間內將籌碼放置於不同的下注區域。

---

## 5.2.1：定義下注區域 (Backend: Go, Frontend: React/Three.js)

### 5.2.1.1：後端定義
-   **5.2.1.1.1-L6**: (Backend) 在 `game/baccarat/table.go` (或類似文件) 中定義常量或枚舉表示下注區域。
    ```go
    const (
        BetAreaPlayer      = "player"
        BetAreaBanker      = "banker"
        BetAreaTie         = "tie"
        BetAreaPlayerPair  = "player_pair"
        BetAreaBankerPair  = "banker_pair"
    )
    var ValidBetAreas = []string{BetAreaPlayer, BetAreaBanker, BetAreaTie, BetAreaPlayerPair, BetAreaBankerPair}
    ```
-   **5.2.1.1.2-L6**: (Backend) 確保 `GameConfig.BaccaratPayouts` 和 `GameConfig.BaccaratBetLimits` 的鍵 (keys) 與這些定義一致。

### 5.2.1.2：前端顯示
-   **5.2.1.2.1-L6**: (Frontend) 在 Three.js 場景中創建代表各下注區域的可交互 3D 模型或平面。
-   **5.2.1.2.2-L6**: (Frontend) 為每個區域添加清晰的標識（例如，紋理上的文字 "閒 Player", "莊 Banker", "和 Tie 8:1", "閒對 Player Pair 11:1", "莊對 Banker Pair 11:1"）。
-   **5.2.1.2.3-L6**: (Frontend) 確保手機端觸控交互良好 (足夠大的觸控區域)。

---

## 5.2.2：實現下注邏輯 (Backend: Go)

### 5.2.2.1：處理玩家下注請求
-   **5.2.2.1.1-L6**: (Backend) 定義 WebSocket/gRPC 消息結構用於接收玩家下注。
    ```protobuf
    // Incoming message from client
    message PlaceBaccaratBetRequest {
      string bet_area = 1; // e.g., "player", "banker_pair"
      int64 amount = 2;
    }
    // Outgoing message to clients
    message BaccaratBetPlaced {
      string player_id = 1;
      string bet_area = 2;
      int64 amount = 3;
      int64 total_player_bet_on_area = 4; // Player's total bet on this area
      int64 total_table_bet_on_area = 5; // All players' total bet on this area
    }
    message BetRejected {
      string player_id = 1;
      string reason = 2; // e.g., "INVALID_AREA", "BELOW_MIN", "ABOVE_MAX", "INSUFFICIENT_CHIPS", "BETTING_CLOSED"
      string bet_area = 3;
      int64 amount = 4;
    }
    ```
-   **5.2.2.1.2-L6**: (Backend) 在 `Table` 或 `GameManager` 中實現 `HandlePlaceBet(playerID string, betArea string, amount int64)` 方法。

### 5.2.2.2：驗證下注
-   **5.2.2.2.1-L6**: (Backend) 檢查遊戲狀態是否為 `BETTING_OPEN`。
-   **5.2.2.2.2-L6**: (Backend) 檢查 `betArea` 是否為有效區域 (`ValidBetAreas`)。
-   **5.2.2.2.3-L6**: (Backend) 檢查玩家是否有足夠籌碼 (`Player.Chips >= amount`)。
-   **5.2.2.2.4-L6**: (Backend) 檢查下注金額是否符合該區域的最小/最大限額 (`GameConfig.BaccaratBetLimits`)。
    -   注意：限額通常是針對**單次下注**還是**玩家在該區域的總下注**？需要明確。假設是**總下注**。
    -   `currentTotalBet = Player.CurrentBaccaratBets[betArea]`
    -   `newTotalBet = currentTotalBet + amount`
    -   `if newTotalBet < limits.Min || newTotalBet > limits.Max { reject }`
-   **5.2.2.2.5-L6**: (Backend) 如果驗證失敗，發送 `BetRejected` 消息。

### 5.2.2.3：接受下注
-   **5.2.2.3.1-L6**: (Backend) 從玩家籌碼中扣除金額 `Player.RemoveChips(amount)`。
-   **5.2.2.3.2-L6**: (Backend) 更新玩家的當前下注記錄 `Player.CurrentBaccaratBets[betArea] += amount`。
-   **5.2.2.3.3-L6**: (Backend) 更新桌面的總下注記錄 `Table.TotalBets[betArea] += amount`。
-   **5.2.2.3.4-L6**: (Backend) 廣播 `BaccaratBetPlaced` 消息給所有玩家，包含玩家 ID、區域、本次下注額、玩家在該區域總額、桌面在該區域總額。

---

## 5.2.3：下注計時器 (Backend: Go, Frontend: React)

### 5.2.3.1：後端計時器管理
-   **5.2.3.1.1-L6**: (Backend) 在 `Table` 狀態中添加 `BettingTimer *time.Timer` 和 `BettingEndTime time.Time`。
-   **5.2.3.1.2-L6**: (Backend) 在進入 `BETTING_OPEN` 狀態時，啟動計時器。
    ```go
    bettingDuration := 15 * time.Second // Configurable?
    table.BettingEndTime = time.Now().Add(bettingDuration)
    table.BettingTimer = time.AfterFunc(bettingDuration, table.CloseBetting)
    ```
-   **5.2.3.1.3-L6**: (Backend) 廣播 `BettingOpened` 事件，包含結束時間戳。
    ```json
    {"event": "betting_opened", "data": {"ends_at": "2023-10-27T10:30:15Z"}}
    ```
-   **5.2.3.1.4-L6**: (Backend) 實現 `CloseBetting()` 方法，該方法將遊戲狀態更改為 `BETTING_CLOSED`，停止接受新下注，並廣播 `BettingClosed` 事件。
-   **5.2.3.1.5-L6**: (Backend) 如果需要，在遊戲中途取消計時器（例如所有玩家都準備好了？百家樂通常不需要）。

### 5.2.3.2：前端計時器顯示
-   **5.2.3.2.1-L6**: (Frontend) 接收 `BettingOpened` 事件和 `ends_at` 時間戳。
-   **5.2.3.2.2-L6**: (Frontend) 使用 React 狀態和 `setInterval` 或 `requestAnimationFrame` 實現一個倒計時顯示器。
-   **5.2.3.2.3-L6**: (Frontend) 在計時器結束或收到 `BettingClosed` 事件時，停止倒計時並禁用下注交互。
-   **5.2.3.2.4-L6**: (Frontend) 在計時器快結束時（例如最後 5 秒）提供視覺或聽覺提示。

---

## 5.2.4：前端交互 (Frontend: React/Three.js)

### 5.2.4.1：選擇籌碼
-   **5.2.4.1.1-L6**: (Frontend) 提供不同面額的籌碼供玩家選擇 (例如 1, 5, 25, 100, 500)。
-   **5.2.4.1.2-L6**: (Frontend) 顯示玩家當前選擇的籌碼面額。

### 5.2.4.2：放置籌碼
-   **5.2.4.2.1-L6**: (Frontend) 玩家點擊（或拖放）選定的籌碼到 Three.js 場景中的下注區域。
-   **5.2.4.2.2-L6**: (Frontend) 點擊下注區域後，向後端發送 `PlaceBaccaratBetRequest` 消息，帶上選擇的籌碼面額和目標區域。
-   **5.2.4.2.3-L6**: (Frontend) 提供 "撤銷" (Undo) 和 "清除" (Clear) 下注的按鈕（僅在下注階段有效）。
    -   撤銷：向後端發送 `UndoLastBetRequest`。
    -   清除：向後端發送 `ClearAllBetsRequest`。
    -   後端需要實現對應的處理邏輯，包括驗證、更新玩家籌碼和下注記錄、廣播更新。

### 5.2.4.3：顯示下注
-   **5.2.4.3.1-L6**: (Frontend) 收到 `BaccaratBetPlaced` 消息後，在對應的下注區域顯示籌碼動畫（例如，從玩家位置飛向區域）。
-   **5.2.4.3.2-L6**: (Frontend) 在每個下注區域顯示該玩家下注的總額。
-   **5.2.4.3.3-L6**: (Frontend) 可選：顯示該區域所有玩家下注的總額。
-   **5.2.4.3.4-L6**: (Frontend) 收到 `BetRejected` 消息時，向玩家顯示拒絕原因（例如，短暫的彈出提示）。
-   **5.2.4.3.5-L6**: (Frontend) 優化大量籌碼的顯示性能（例如，合併籌碼圖像，使用 Instanced Meshes）。

--- 