# 百家樂開發任務：初始化遊戲 (Task 5.1)

**目標**：設置百家樂遊戲環境，包括牌靴（8 副牌）、賠率、佣金、下注區域和限額。

---

## 5.1.1：創建遊戲配置 (Backend: Go)

### 5.1.1.1：定義賠率與佣金
-   **5.1.1.1.1-L6**: 在 `config/game_settings.go` 的 `GameConfig` struct 中添加百家樂相關配置字段。
    -   `BaccaratPayouts map[string]float64` (例如 `{"player": 1.0, "banker": 0.95, "tie": 8.0, "player_pair": 11.0, "banker_pair": 11.0}`)
    -   `BankerCommissionRate float64` (例如 0.05)
-   **5.1.1.1.2-L6**: 在百家樂的配置文件 (例如 `config/baccarat.json`) 或數據庫 `game_configs` 表中設置這些賠率和佣金值。
-   **5.1.1.1.3-L6**: 確保 `LoadGameConfig("baccarat")` 能正確加載。
-   **5.1.1.1.4-L6**: 考慮不同賭場規則變體（例如免佣百家樂）的配置靈活性。

### 5.1.1.2：設置下注限額
-   **5.1.1.2.1-L6**: 在 `GameConfig` struct 中添加下注限額字段。
    -   `BaccaratBetLimits map[string]struct{ Min int; Max int }` (例如 `{"player": {Min: 10, Max: 1000}, "banker": {Min: 10, Max: 1000}, "tie": {Min: 5, Max: 100}, "player_pair": {Min: 1, Max: 50}, "banker_pair": {Min: 1, Max: 50}}`)
-   **5.1.1.2.2-L6**: 在百家樂配置文件/數據庫中設置各區域的限額。
-   **5.1.1.2.3-L6**: 下注驗證邏輯 (5.2) 將使用這些限額。

### 5.1.1.3：設置牌靴 (Shoe)
-   **5.1.1.3.1-L6**: 在 `GameConfig` struct 中添加牌靴相關配置。
    -   `BaccaratNumDecks int` (設置為 8)
    -   `BaccaratCuttingCardPosition int` (例如 14，表示剩餘 14 張牌時提示換靴)
-   **5.1.1.3.2-L6**: 牌靴的實際狀態管理 (剩餘牌數、已發牌) 在 `Deck` 或 `Table` 狀態中處理 (見 5.1.3)。
-   **5.1.1.3.3-L6**: 換靴邏輯 (Reshuffle) 在狀態管理 (5.5) 中處理。

---

## 5.1.2：初始化玩家 (Backend: Go)

### 5.1.2.1：創建玩家數據結構
-   **5.1.2.1.1-L6**: 重用基礎 `Player` struct (`ID`, `Name`, `Chips`, `Status`, `SeatIndex`)。座位在百家樂中通常不影響遊戲邏輯，但可用於界面顯示。
-   **5.1.2.1.2-L6**: 添加特定於百家樂的玩家數據字段，可能放在 `Player.GameSpecificData` 中或直接添加。
    -   `CurrentBaccaratBets map[string]int` (例如 `{"player": 100, "banker_pair": 10}`)
-   **5.1.2.1.3-L6**: 玩家狀態 (`Waiting`, `Active`, `SittingOut`)。百家樂通常沒有 `Folded` 或 `AllIn` 狀態（除非是特殊錦標賽）。

### 5.1.2.2：初始化玩家籌碼
-   **5.1.2.2.1-L6**: 百家樂通常沒有固定起始籌碼，玩家根據自己的餘額下注。加入桌子時顯示玩家帳戶餘額。
-   **5.1.2.2.2-L6**: 重用 `Player.AddChips`, `Player.RemoveChips`。
-   **5.1.2.2.3-L6**: 重用籌碼變動通知。

### 5.1.2.3：手機端顯示
-   **5.1.2.3.1-L6**: (Backend) 提供獲取玩家列表及其籌碼的接口/消息。
-   **5.1.2.3.2-L6**: (Frontend) 顯示玩家列表和籌碼。
-   **5.1.2.3.3-L6**: (Frontend) 顯示主要的下注區域 (閒、莊、和、閒對、莊對)。

---

## 5.1.3：初始化牌堆 (多副牌牌靴) (Backend: Go, Common Module)

### 5.1.3.1：創建多副牌數據
-   **5.1.3.1.1-L6**: 重用 `pkg/card` 中的 `Card`, `Suit`, `Rank`。
-   **5.1.3.1.2-L6**: 實現 `card.GetValue(gameType string)` 方法或類似機制，獲取牌在不同遊戲中的點數。
    -   百家樂點數: A=1, 2-9=面值, 10/J/Q/K=0。
-   **5.1.3.1.3-L6**: 重用 `pkg/deck`。

### 5.1.3.2：實現洗牌算法
-   **5.1.3.2.1-L6**: `NewDeck(numDecks int)` 函數接收 `GameConfig.BaccaratNumDecks` (值為 8)。
-   **5.1.3.2.2-L6**: 重用 `Deck.Shuffle()`，對 8 副牌共 416 張牌進行洗牌。

### 5.1.3.3：設置牌靴狀態
-   **5.1.3.3.1-L6**: `Deck` struct 需要能處理多副牌。
-   **5.1.3.3.2-L6**: `Deck.currentIndex` 正常工作。
-   **5.1.3.3.3-L6**: `Deck.DrawCard()` 正常工作。
-   **5.1.3.3.4-L6**: `Deck.RemainingCards()` 返回剩餘牌數。
-   **5.1.3.3.5-L6**: 在 `Table` 或遊戲管理器中持有 `Deck` 實例。
-   **5.1.3.3.6-L6**: 需要在發牌過程中檢查是否達到切牌標記位置 (`Deck.currentIndex >= totalCards - GameConfig.BaccaratCuttingCardPosition`)，觸發換靴提示 (見 5.5)。

### 5.1.3.4：首輪燒牌 (Burn Card)
-   **5.1.3.4.1-L6**: 在 `Table.StartNewShoe()` 方法中實現 (僅在新牌靴啟用時執行一次)。
-   **5.1.3.4.2-L6**: 調用 `Deck.Shuffle()` 後。
-   **5.1.3.4.3-L6**: 調用 `Deck.DrawCard()` 抽出第一張牌，**將其亮出**。
-   **5.1.3.4.4-L6**: 獲取該牌的百家樂點數 `burnValue = card.GetValue("baccarat")` (10/J/Q/K 點數為 10)。
-   **5.1.3.4.5-L6**: 再連續調用 `Deck.DrawCard()` `burnValue` 次，將這些牌丟棄 (不加入任何手牌，不顯示)。
-   **5.1.3.4.6-L6**: 廣播燒牌完成和第一張亮牌的信息。
    -   `{"event": "shoe_prepared", "data": {"first_card": {"id": "H7"}, "cards_burnt": 7}}`
-   **5.1.3.4.7-L6**: (Frontend) 顯示燒牌動畫和第一張亮牌。

--- 