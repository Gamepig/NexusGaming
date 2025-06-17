# 德州撲克開發任務：狀態管理 (Task 3.5)

**目標**：實現對整個遊戲牌局狀態的精確追蹤、轉換和重置，確保遊戲流程的正確性和數據一致性。

---

## 3.5.1：遊戲狀態機 (Game State Machine) (Backend: Go)

### 3.5.1.1：定義遊戲狀態
-   **3.5.1.1.1-L6**: 在 `game/table.go` 或 `game/gamestate.go` 中定義遊戲主狀態枚舉或常量。
    -   `GameStateWaitingForPlayers`
    -   `GameStateDealingHoleCards`
    -   `GameStatePreFlopBetting`
    -   `GameStateDealingFlop`
    -   `GameStateFlopBetting`
    -   `GameStateDealingTurn`
    -   `GameStateTurnBetting`
    -   `GameStateDealingRiver`
    -   `GameStateRiverBetting`
    -   `GameStateShowdown`
    -   `GameStateHandOver`
-   **3.5.1.1.2-L6**: 在 `Table` struct 中添加 `CurrentGameState string` (或枚舉類型) 字段。

### 3.5.1.2：狀態轉換邏輯
-   **3.5.1.2.1-L6**: 在每個關鍵節點實現狀態轉換。
    -   `StartNewHand()` -> `GameStateDealingHoleCards`
    -   `DealHoleCards()` 完成 -> `GameStatePreFlopBetting`
    -   `EndBettingRound()` (Pre-Flop) -> `GameStateDealingFlop`
    -   `DealFlop()` 完成 -> `GameStateFlopBetting`
    -   `EndBettingRound()` (Flop) -> `GameStateDealingTurn`
    -   ... (依此類推) ...
    -   `EndBettingRound()` (River) -> `GameStateShowdown` (如果需要) 或 `GameStateHandOver` (如果只有一人)
    -   `DistributePots()` 完成 -> `GameStateHandOver`
    -   `ResetForNewHand()` -> `GameStateDealingHoleCards` (如果自動開始下一局) 或 `GameStateWaitingForPlayers`
-   **3.5.1.2.2-L6**: 確保狀態轉換的原子性和正確性，避免競爭條件。
-   **3.5.1.2.3-L6**: 每次狀態轉換時，廣播遊戲狀態更新通知。
    -   `{"event": "game_state_update", "data": {"state": "FlopBetting", "round": "Flop", ...}}` (可附加輪次等信息)

### 3.5.1.3：管理遊戲主循環 (Game Loop)
-   **3.5.1.3.1-L6**: 設計遊戲主邏輯的驅動方式。可能是：
    -   **事件驅動**: 玩家動作、定時器觸發狀態轉換和後續操作。
    -   **狀態機循環**: 一個 goroutine 根據當前狀態執行相應操作，並等待事件來觸發轉換。
-   **3.5.1.3.2-L6**: 選擇一種方式並在 `Table` 或專門的 `GameManager` 中實現。
-   **3.5.1.3.3-L6**: 處理超時自動動作 (例如超時自動蓋牌或過牌)。
    -   在 `your_turn` 通知時啟動定時器 (`Player.ActionTimer`)。
    -   如果玩家在超時前行動，取消定時器。
    -   如果定時器觸發，執行默認動作 (通常是蓋牌，或在可過牌時過牌)。

---

## 3.5.2：玩家狀態管理 (Backend: Go)

### 3.5.2.1：定義玩家狀態
-   **3.5.2.1.1-L6**: 回顧 `Player` struct 中的 `Status` 字段及其可能的值：
    -   `PlayerStatusWaiting` (等待加入遊戲或下一局)
    -   `PlayerStatusSittingOut` (暫時離開，保留座位)
    -   `PlayerStatusActive` (參與當前牌局，可以行動)
    -   `PlayerStatusFolded` (在本局已蓋牌)
    -   `PlayerStatusAllIn` (在本局已全下)
    -   `PlayerStatusEliminated` (籌碼為 0，已淘汰 - 主要用於錦標賽)
-   **3.5.2.1.2-L6**: 確保 PlayerStatus 常量定義清晰。

### 3.5.2.2：玩家狀態轉換
-   **3.5.2.2.1-L6**: 在相應的遊戲邏輯中更新玩家狀態。
    -   加入桌子 -> `Waiting` 或 `Active` (取決於是否立即開始)
    -   `HandleFold()` -> `Folded`
    -   籌碼變為 0 (Call/Bet/Raise/PostBlinds) -> `AllIn` (如果本局還在進行) 或 `Eliminated` (如果本局結束時籌碼為 0)
    -   新一局開始 (`ResetForNewHand()`) -> 將 `Folded`, `AllIn` 狀態重置為 `Active` (如果籌碼 > 0)。`Eliminated` 狀態不變。
    -   處理玩家主動 "Sit Out" 或 "Back" 的請求。
-   **3.5.2.2.2-L6**: 每次玩家狀態變化時，廣播玩家狀態更新通知。
    -   `{"event": "player_status_update", "data": {"seat_index": ..., "status": "folded"}}`

### 3.5.2.3：處理玩家斷線與重連
-   **3.5.2.3.1-L6**: WebSocket 連接斷開時，檢測玩家是否在牌局中。
-   **3.5.2.3.2-L6**: 如果在牌局中，可以將玩家狀態暫時標記為 `Disconnected` 或類似狀態，並啟動一個較長的重連計時器。
-   **3.5.2.3.3-L6**: 如果輪到斷線玩家行動，可以自動執行蓋牌或過牌。
-   **3.5.2.3.4-L6**: 玩家重連時，恢復其連接狀態，並發送當前完整的遊戲狀態。
-   **3.5.2.3.5-L6**: 如果重連超時，可以將玩家標記為 `SittingOut` 或強制蓋牌。

---

## 3.5.3：牌局序列管理 (Hand Sequencing) (Backend: Go)

### 3.5.3.1：牌局 ID
-   **3.5.3.1.1-L6**: 為每局牌生成一個唯一的 ID (例如 `HandID string`，可以是 UUID 或自增序號)。
-   **3.5.3.1.2-L6**: 在 `Table.StartNewHand()` 時生成並存儲。
-   **3.5.3.1.3-L6**: 在所有與該局相關的日誌和通知中包含 HandID，方便追溯。

### 3.5.3.2：牌局開始與結束
-   **3.5.3.2.1-L6**: `Table.StartNewHand()` 方法應包含：
    -   增加 HandID。
    -   重置牌局相關狀態 (玩家手牌、公共牌、底池、玩家本輪下注、玩家狀態從 Folded/AllIn 到 Active)。
    -   移動 Button。
    -   確定並扣除盲注。
    -   洗牌 (如果需要，Deck 可能需要重置)。
    -   發底牌。
    -   啟動 Pre-Flop 下注輪。
    -   廣播新牌局開始通知。
        - `{"event": "new_hand", "data": {"hand_id": ..., "button_seat": ..., "sb_seat": ..., "bb_seat": ...}}`
-   **3.5.3.2.2-L6**: `Table.EndHand()` 方法應包含：
    -   標記牌局結束。
    -   記錄牌局結果 (獲勝者、牌型、底池大小等，用於歷史記錄或統計)。
    -   清理臨時數據。
    -   檢查是否有玩家籌碼歸零並處理淘汰。
    -   決定是否自動開始下一局或等待。

### 3.5.3.3：準備下一局
-   **3.5.3.3.1-L6**: 在 `EndHand` 後，可以有一個短暫的延遲 (`time.Sleep`)。
-   **3.5.3.3.2-L6**: 檢查是否有足夠的玩家 (至少 2 人) 且狀態不是 `SittingOut` 或 `Eliminated`。
-   **3.5.3.3.3-L6**: 如果滿足條件，調用 `StartNewHand()` 開始下一局。
-   **3.5.3.3.4-L6**: 如果不滿足條件，轉換到 `GameStateWaitingForPlayers` 狀態。

---

## 3.5.4：數據持久化與恢復 (考慮點)

### 3.5.4.1：遊戲狀態持久化 (可選)
-   **3.5.4.1.1-L6**: 考慮是否需要在關鍵節點 (例如每輪下注結束、每局結束) 將 `Table` 的核心狀態序列化並存儲到數據庫 (例如 Redis 或 MongoDB)。
-   **3.5.4.1.2-L6**: 這對於服務器重啟後的遊戲恢復可能有用，但會增加複雜性。對於非關鍵應用，可能不需要。
-   **3.5.4.1.3-L6**: 設計序列化格式 (JSON 或其他)。

### 3.5.4.2：牌局歷史記錄
-   **3.5.4.2.1-L6**: 在 `EndHand` 時，將牌局的關鍵信息 (HandID, 玩家, 底牌(可選), 公共牌, 下注歷史摘要, 結果, 獲勝者, 底池大小) 存儲到數據庫 (例如 MySQL 或 MongoDB 的 `hand_history` 表)。
-   **3.5.4.2.2-L6**: 設計 `hand_history` 表結構。
-   **3.5.4.2.3-L6**: 提供 API 或後台功能查詢牌局歷史。

### 3.5.4.3：服務器重啟恢復 (如果實現持久化)
-   **3.5.4.3.1-L6**: 服務器啟動時，檢查是否有未完成的牌局狀態存儲在數據庫中。
-   **3.5.4.3.2-L6**: 加載狀態，重建 `Table` 對象。
-   **3.5.4.3.3-L6**: 重新建立與玩家的 WebSocket 連接 (需要玩家重連機制)。
-   **3.5.4.3.4-L6**: 從中斷的點繼續遊戲。

--- 