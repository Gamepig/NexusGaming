# 七張梭哈開發任務：狀態管理 (Task 4.5)

**目標**：實現七張梭哈遊戲狀態的精確追蹤、轉換和重置，特別注意梭哈特有的狀態和流程。

---

## 4.5.1：遊戲狀態機 (Game State Machine) (Backend: Go)

### 4.5.1.1：定義遊戲狀態
-   **4.5.1.1.1-L6**: 在 `game/gamestate.go` 中定義梭哈的遊戲主狀態。
    -   `GameStateStudWaitingForPlayers`
    -   `GameStateStudCollectingAnte`
    -   `GameStateStudDealingThirdStreet`
    -   `GameStateStudThirdStreetBetting`
    -   `GameStateStudDealingFourthStreet`
    -   `GameStateStudFourthStreetBetting`
    -   `GameStateStudDealingFifthStreet`
    -   `GameStateStudFifthStreetBetting`
    -   `GameStateStudDealingSixthStreet`
    -   `GameStateStudSixthStreetBetting`
    -   `GameStateStudDealingSeventhStreet`
    -   `GameStateStudSeventhStreetBetting`
    -   `GameStateStudShowdown`
    -   `GameStateStudHandOver`
-   **4.5.1.1.2-L6**: `Table.CurrentGameState` 使用這些狀態。

### 4.5.1.2：狀態轉換邏輯
-   **4.5.1.2.1-L6**: 實現梭哈的狀態轉換流程。
    -   `StartNewHand()` -> `GameStateStudCollectingAnte`
    -   `CollectAnte()` 完成 -> `GameStateStudDealingThirdStreet`
    -   `DealThirdStreet()` 完成 -> `GameStateStudThirdStreetBetting` (確定 Bring-in 後)
    -   `EndStudBettingRound()` (Third St) -> `GameStateStudDealingFourthStreet`
    -   `DealFourthStreet()` 完成 -> `GameStateStudFourthStreetBetting`
    -   ... (依此類推) ...
    -   `EndStudBettingRound()` (Seventh St) -> `GameStateStudShowdown` 或 `GameStateStudHandOver`
    -   `DistributePots()` 完成 -> `GameStateStudHandOver`
    -   `ResetForNewHand()` -> `GameStateStudCollectingAnte` 或 `GameStateStudWaitingForPlayers`
-   **4.5.1.2.2-L6**: 確保原子性。
-   **4.5.1.2.3-L6**: 廣播遊戲狀態更新通知 (使用梭哈的狀態名)。

### 4.5.1.3：管理遊戲主循環
-   **4.5.1.3.1-L6**: 重用 Hold'em 的事件驅動或狀態機循環模式。
-   **4.5.1.3.2-L6**: 處理超時自動動作 (梭哈中通常是蓋牌)。

---

## 4.5.2：玩家狀態管理 (Backend: Go)

### 4.5.2.1：定義玩家狀態
-   **4.5.2.1.1-L6**: 重用 Hold'em 的 `PlayerStatus` 枚舉/常量 (`Waiting`, `SittingOut`, `Active`, `Folded`, `AllIn`, `Eliminated`)。

### 4.5.2.2：玩家狀態轉換
-   **4.5.2.2.1-L6**: 在梭哈的遊戲邏輯中更新玩家狀態。
    -   Ante 不足 / 下注/跟注/加注導致籌碼為 0 -> `AllIn` 或 `Eliminated`。
    -   `HandleFold()` -> `Folded`。
    -   `ResetForNewHand()` -> `Folded`, `AllIn` 重置為 `Active`。
    -   處理 Sit Out / Back。
-   **4.5.2.2.2-L6**: 廣播玩家狀態更新通知。

### 4.5.2.3：處理玩家斷線與重連
-   **4.5.2.3.1-L6**: 重用 Hold'em 的斷線重連邏輯。
    -   標記 `Disconnected`。
    -   超時自動蓋牌。
    -   重連後發送完整狀態。

---

## 4.5.3：牌局序列管理 (Hand Sequencing) (Backend: Go)

### 4.5.3.1：牌局 ID
-   **4.5.3.1.1-L6**: 重用 Hold'em 的 HandID 生成和記錄機制。

### 4.5.3.2：牌局開始與結束
-   **4.5.3.2.1-L6**: `Table.StartNewHand()` (梭哈版本) 應包含：
    -   增加 HandID。
    -   重置牌局狀態 (玩家手牌 `Hand`, **亮牌 `VisibleCards`**, 底池, 玩家本輪/總下注, 玩家狀態)。
    -   **收取 Ante** (`CollectAnte()`)。
    -   洗牌 (`Deck.Reset()` - 梭哈用單副牌)。
    -   **發第三街** (`DealThirdStreet()`)。
    -   **確定 Bring-in 玩家** (`DetermineBringInPlayer()`)。
    -   啟動 Third Street 下注輪 (`StartStudBettingRound("Third Street")`)。
    -   廣播新牌局開始通知 (可包含 Ante 金額)。
        - `{"event": "new_hand_stud", "data": {"hand_id": ..., "ante": ...}}`
-   **4.5.3.2.2-L6**: `Table.EndHand()` (梭哈版本) 應包含：
    -   標記結束。
    -   記錄歷史 (梭哈需要記錄亮牌過程)。
    -   清理。
    -   處理淘汰。
    -   決定下一局。

### 4.5.3.3：準備下一局
-   **4.5.3.3.1-L6**: 重用 Hold'em 的邏輯。
    -   延遲。
    -   檢查玩家數量和狀態。
    -   調用 `StartNewHand()` 或等待。

---

## 4.5.4：數據持久化與恢復 (考慮點)

### 4.5.4.1：遊戲狀態持久化 (可選)
-   **4.5.4.1.1-L6**: 如果需要，持久化 `Table` 狀態，需包含梭哈特有字段 (如 `VisibleCards`, `BringInSeatIndex`, `CurrentBetLimit`)。

### 4.5.4.2：牌局歷史記錄
-   **4.5.4.2.1-L6**: 在 `EndHand` 時存儲牌局歷史。
    -   梭哈歷史記錄應包含每條街發的牌（暗牌對應玩家自己可見，亮牌公開），Bring-in，下注輪次摘要，結果等。
-   **4.5.4.2.2-L6**: 設計 `stud_hand_history` 表結構。

### 4.5.4.3：服務器重啟恢復
-   **4.5.4.3.1-L6**: 如果實現持久化，需要能從存儲的狀態恢復梭哈牌局。

--- 