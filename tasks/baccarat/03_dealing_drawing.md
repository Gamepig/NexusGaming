# 百家樂開發任務：發牌與補牌邏輯 (Task 5.3)

**目標**：實現百家樂的發牌和補牌規則，包括初始兩張牌的派發以及根據閒家和莊家點數決定是否補第三張牌。

---

## 5.3.1：初始發牌 (Backend: Go)

### 5.3.1.1：定義手牌結構
-   **5.3.1.1.1-L6**: 在 `game/baccarat/table.go` 或 `game/baccarat/hand.go` 中定義手牌結構。
    ```go
    type BaccaratHand struct {
        Cards []*card.Card
        Total int // Calculated according to Baccarat rules
    }
    // Table state might include:
    // PlayerHand *BaccaratHand
    // BankerHand *BaccaratHand
    ```
-   **5.3.1.1.2-L6**: 實現 `BaccaratHand.CalculateTotal()` 方法。
    -   遍歷 `Cards`。
    -   使用 `card.GetValue("baccarat")` (A=1, 2-9, T/J/Q/K=0)。
    -   累加點數。
    -   結果對 10 取模 (`total % 10`)。
    -   更新 `Hand.Total`。

### 5.3.1.2：實現發牌順序
-   **5.3.1.2.1-L6**: 在 `Table` 或 `GameManager` 中實現 `DealInitialHands()` 方法。
-   **5.3.1.2.2-L6**: 確保在 `BETTING_CLOSED` 狀態後調用。
-   **5.3.1.2.3-L6**: 從 `Deck` 中按順序發牌：
    1.  `PlayerHand.Cards = append(PlayerHand.Cards, Deck.DrawCard())`
    2.  `BankerHand.Cards = append(BankerHand.Cards, Deck.DrawCard())`
    3.  `PlayerHand.Cards = append(PlayerHand.Cards, Deck.DrawCard())`
    4.  `BankerHand.Cards = append(BankerHand.Cards, Deck.DrawCard())`
-   **5.3.1.2.4-L6**: **牌面朝上** 發出。
-   **5.3.1.2.5-L6**: 每次發牌後，廣播 `CardDealt` 事件，包含牌的信息和發往的位置 (Player/Banker)。
    ```json
    {"event": "card_dealt", "data": {"card": {"id": "S5", "suit":"S", "rank":"5"}, "position": "player", "card_index": 0}}
    {"event": "card_dealt", "data": {"card": {"id": "HA", "suit":"H", "rank":"A"}, "position": "banker", "card_index": 0}}
    // ... and so on for the first 4 cards
    ```
-   **5.3.1.2.6-L6**: (Frontend) 接收事件並顯示發牌動畫，將牌放置到對應的閒/莊區域。
-   **5.3.1.2.7-L6**: 計算初始點數 `PlayerHand.CalculateTotal()`, `BankerHand.CalculateTotal()`。
-   **5.3.1.2.8-L6**: 廣播初始點數 `InitialTotalsCalculated` 事件。
    ```json
    {"event": "initial_totals", "data": {"player_total": 7, "banker_total": 4}}
    ```
-   **5.3.1.2.9-L6**: (Frontend) 顯示閒家和莊家的初始點數。

---

## 5.3.2：閒家補牌規則 (Player's Third Card Rule) (Backend: Go)

### 5.3.2.1：檢查是否需要補牌
-   **5.3.2.1.1-L6**: 在 `DealInitialHands()` 完成後，實現 `DeterminePlayerAction()` 方法。
-   **5.3.2.1.2-L6**: 檢查天牌 (Natural)：`if PlayerHand.Total >= 8 || BankerHand.Total >= 8`，則都不補牌，直接進入結算 (Task 5.4)。廣播 `NaturalWin` 事件。
-   **5.3.2.1.3-L6**: 如果沒有天牌，根據閒家規則判斷：
    -   `if PlayerHand.Total <= 5`: 閒家需要補牌。
    -   `if PlayerHand.Total == 6 || PlayerHand.Total == 7`: 閒家停牌 (Stand)。

### 5.3.2.2：執行閒家補牌
-   **5.3.2.2.1-L6**: 如果需要補牌，調用 `DrawThirdCard("player")`。
-   **5.3.2.2.2-L6**: `DrawThirdCard` 方法：
    -   `thirdCard := Deck.DrawCard()`
    -   `PlayerHand.Cards = append(PlayerHand.Cards, thirdCard)`
    -   `PlayerHand.CalculateTotal()`
    -   廣播 `CardDealt` 事件 (position: "player", card_index: 2)。
    -   廣播 `HandTotalUpdated` 事件 (position: "player", new_total: ...)。
    -   記錄 `playerThirdCardValue = thirdCard.GetValue("baccarat")` (注意 JQK=0)。如果沒補牌，此值為特殊標記，例如 -1。
-   **5.3.2.2.3-L6**: (Frontend) 顯示閒家補第三張牌的動畫和更新後的點數。
-   **5.3.2.2.4-L6**: 如果閒家停牌，廣播 `PlayerStands` 事件。

---

## 5.3.3：莊家補牌規則 (Banker's Third Card Rule) (Backend: Go)

### 5.3.3.1：檢查是否需要補牌
-   **5.3.3.1.1-L6**: 在閒家動作完成後 (補牌或停牌)，實現 `DetermineBankerAction()` 方法。
-   **5.3.3.1.2-L6**: **重要**: 如果初始發牌是天牌 (Natural 8 或 9)，莊家不補牌 (已在 5.3.2.1 處理)。
-   **5.3.3.1.3-L6**: 根據莊家規則判斷：
    -   **Case 1: 閒家停牌 (Player Stood on 6 or 7)**
        -   `if BankerHand.Total <= 5`: 莊家補牌。
        -   `if BankerHand.Total == 6 || BankerHand.Total == 7`: 莊家停牌。
    -   **Case 2: 閒家補了第三張牌 (Player Drew a Third Card)**
        -   使用 `playerThirdCardValue` (閒家第三張牌的點數, A=1, ..., 9=9, T/J/Q/K=0) 和 `BankerHand.Total` (莊家前兩張牌的點數) 查表決定莊家是否補牌。
        -   `bankerShouldDraw := CheckBankerDrawingRule(BankerHand.Total, playerThirdCardValue)`
        -   實現 `CheckBankerDrawingRule` 函數，包含詳細的查表邏輯：
            -   莊家 2 點或以下：**補**
            -   莊家 3 點：閒家第三張牌**不是** 8 時，**補**
            -   莊家 4 點：閒家第三張牌是 2-7 時，**補**
            -   莊家 5 點：閒家第三張牌是 4-7 時，**補**
            -   莊家 6 點：閒家第三張牌是 6 或 7 時，**補**
            -   莊家 7 點：**停**
            -   莊家 8, 9 點：天牌，不補 (已處理)

### 5.3.3.2：執行莊家補牌
-   **5.3.3.2.1-L6**: 如果 `bankerShouldDraw` 為 true，調用 `DrawThirdCard("banker")`。
-   **5.3.3.2.2-L6**: `DrawThirdCard("banker")` 邏輯類似閒家補牌：
    -   抽牌、添加、計算總點數。
    -   廣播 `CardDealt` (position: "banker", card_index: 2)。
    -   廣播 `HandTotalUpdated` (position: "banker", new_total: ...)。
-   **5.3.3.2.3-L6**: (Frontend) 顯示莊家補第三張牌的動畫和更新後的點數。
-   **5.3.3.2.4-L6**: 如果莊家停牌，廣播 `BankerStands` 事件。

### 5.3.3.3：完成發牌
-   **5.3.3.3.1-L6**: 莊家動作完成後，發牌階段結束。
-   **5.3.3.3.2-L6**: 廣播 `DealingComplete` 事件，包含最終的閒莊手牌和點數。
    ```json
    {
      "event": "dealing_complete",
      "data": {
        "player_hand": {"cards": [{"id":"S5"},{"id":"C2"},{"id":"D1"}], "total": 8},
        "banker_hand": {"cards": [{"id":"HA"},{"id":"H3"},{"id":"S9"}], "total": 3}
      }
    }
    ```
-   **5.3.3.3.3-L6**: 進入結果判定階段 (Task 5.4)。

---

## 5.3.4：卡牌點數計算 (Backend: Go, Common Module)

### 5.3.4.1：實現點數獲取
-   **5.3.4.1.1-L6**: 確保 `pkg/card/card.go` 中的 `Card` 結構或相關函數能返回特定遊戲的點數。
    ```go
    // Example approach 1: Method on Card
    func (c *Card) GetValue(gameType string) int {
        switch gameType {
        case "baccarat":
            if c.Rank >= TEN && c.Rank <= KING {
                return 0
            }
            if c.Rank == ACE {
                return 1
            }
            return int(c.Rank) // Assumes Rank maps directly to 2-9
        case "blackjack":
            // ... Blackjack logic ...
        default:
            return 0 // Or handle error
        }
    }

    // Example approach 2: Separate utility function
    func GetCardValue(card *card.Card, gameType string) int {
         // ... logic ...
    }
    ```
-   **5.3.4.1.2-L6**: 在 `BaccaratHand.CalculateTotal()` 中調用 `card.GetValue("baccarat")`。
-   **5.3.4.1.3-L6**: 添加單元測試驗證百家樂點數計算的正確性。

--- 