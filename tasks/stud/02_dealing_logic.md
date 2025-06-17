# 七張梭哈開發任務：發牌邏輯 (Task 4.2)

**目標**：實現七張梭哈獨特的發牌流程，包括 Third Street (兩暗一明)，Fourth 到 Sixth Street (一明)，以及 Seventh Street (一暗或公共牌)。

---

## 4.2.1：發第三街 (Third Street) (Backend: Go)

### 4.2.1.1：發牌順序
-   **4.2.1.1.1-L6**: 在 `Table.DealThirdStreet()` 方法中實現 (收取 Ante 後)。
-   **4.2.1.1.2-L6**: 從莊家按鈕 (如果使用) 左手邊第一個玩家開始，按順時針順序。梭哈通常沒有固定莊家按鈕移動，發牌起始點固定或隨機選擇一個玩家開始。**確認：是否需要隨機起始玩家？** 假設從座位 0 開始。
-   **4.2.1.1.3-L6**: 獲取所有已支付 Ante 的活躍玩家列表。

### 4.2.1.2：執行發牌 (兩暗一明)
-   **4.2.1.2.1-L6**: 循環發三輪牌。
-   **4.2.1.2.2-L6**: 第一輪 (暗牌 1)：順時針給每個活躍玩家發一張暗牌 (`Deck.DrawCard()`)，加入 `Player.Hand`。
-   **4.2.1.2.3-L6**: 第二輪 (暗牌 2)：順時針給每個活躍玩家再發一張暗牌 (`Deck.DrawCard()`)，加入 `Player.Hand`。
-   **4.2.1.2.4-L6**: 第三輪 (明牌 - Door Card)：順時針給每個活躍玩家發一張**明牌** (`Deck.DrawCard()`)。
    -   將此牌同時加入 `Player.Hand` 和 `Player.VisibleCards`。
-   **4.2.1.2.5-L6**: 檢查 `DrawCard` 錯誤 (牌堆耗盡，不應發生)。

### 4.2.1.3：更新狀態與通知
-   **4.2.1.3.1-L6**: 更新遊戲階段狀態 (`GameStateThirdStreetBetting`)。
-   **4.2.1.3.2-L6**: **私密地** 將兩張暗牌信息發送給每個對應玩家。
    -   `{"event": "hole_cards_stud", "data": {"cards": [{"id": "SA"}, {"id": "S2"}]}}`
-   **4.2.1.3.3-L6**: **公開地** 廣播所有玩家的第三張明牌 (Door Card)。
    -   `{"event": "dealing_complete_stud", "data": {"round": "Third Street", "visible_cards": [{"seat_index": 0, "card": {"id": "HK"}}, {"seat_index": 1, "card": {"id": "D5"}}, ...]}}`
-   **4.2.1.3.4-L6**: (前端任務) 接收暗牌並顯示給自己。
-   **4.2.1.3.5-L6**: (Frontend) 接收明牌並顯示在對應玩家的亮牌區域。
-   **4.2.1.3.6-L6**: (Frontend) 實現發牌動畫。

### 4.2.1.4：確定 Bring-in 玩家
-   **4.2.1.4.1-L6**: 在發完第三街後，立即調用 `Table.DetermineBringInPlayer()`。
-   **4.2.1.4.2-L6**: 遍歷所有活躍玩家，找到擁有**最低**點數明牌 (Door Card) 的玩家。
    -   需要 `GetStudRankValue(card *Card)` 函數 (A 為高牌，2 為最低)。
-   **4.2.1.4.3-L6**: 如果最低點數有多個玩家，則根據花色比較 (例如：黑桃 > 紅心 > 方塊 > 梅花，**確認：使用哪種花色順序？** 假設 Spades > Hearts > Diamonds > Clubs，需要明確定義)。
-   **4.2.1.4.4-L6**: 記錄 Bring-in 玩家的座位索引 (`BringInSeatIndex`)。
-   **4.2.1.4.5-L6**: 觸發 Bring-in 玩家確定的通知 (用於前端高亮和後續強制下注)。
    -   `{"event": "bring_in_determined", "data": {"seat_index": ...}}`

---

## 4.2.2：發第四街 (Fourth Street) (Backend: Go)

### 4.2.2.1：執行發牌 (一明牌)
-   **4.2.2.1.1-L6**: 在 `Table.DealFourthStreet()` 方法中實現 (Third Street 下注結束後)。
-   **4.2.2.1.2-L6**: **不燒牌** (傳統梭哈規則通常不燒牌，**確認：是否遵循此規則？** 假設不燒牌)。
-   **4.2.2.1.3-L6**: 從 Bring-in 玩家**之後**第一個未蓋牌的玩家開始，順時針給每個未蓋牌的玩家發一張**明牌** (`Deck.DrawCard()`)。
-   **4.2.2.1.4-L6**: 將此牌加入 `Player.Hand` 和 `Player.VisibleCards`。
-   **4.2.2.1.5-L6**: 檢查 `DrawCard` 錯誤。

### 4.2.2.2：更新狀態與通知
-   **4.2.2.2.1-L6**: 更新遊戲階段狀態 (`GameStateFourthStreetBetting`)。
-   **4.2.2.2.2-L6**: **公開地** 廣播所有玩家的第四張明牌。
    -   `{"event": "dealing_complete_stud", "data": {"round": "Fourth Street", "visible_cards": [{"seat_index": 0, "card": {"id": "CA"}}, {"seat_index": 1, "card": {"id": "S8"}}, ...]}}` (只發送新增的牌)
-   **4.2.2.2.3-L6**: (Frontend) 接收明牌並追加顯示在亮牌區域。
-   **4.2.2.2.4-L6**: (Frontend) 實現發牌動畫。

---

## 4.2.3：發第五街 (Fifth Street) (Backend: Go)

### 4.2.3.1：執行發牌 (一明牌)
-   **4.2.3.1.1-L6**: 在 `Table.DealFifthStreet()` 方法中實現 (Fourth Street 下注結束後)。
-   **4.2.3.1.2-L6**: 不燒牌。
-   **4.2.3.1.3-L6**: 從 Bring-in 玩家**之後**第一個未蓋牌的玩家開始，順時針給每個未蓋牌的玩家發一張**明牌** (`Deck.DrawCard()`)。
-   **4.2.3.1.4-L6**: 將此牌加入 `Player.Hand` 和 `Player.VisibleCards`。
-   **4.2.3.1.5-L6**: 檢查 `DrawCard` 錯誤。

### 4.2.3.2：更新狀態與通知
-   **4.2.3.2.1-L6**: 更新遊戲階段狀態 (`GameStateFifthStreetBetting`)。
-   **4.2.3.2.2-L6**: **公開地** 廣播所有玩家的第五張明牌。
    -   `{"event": "dealing_complete_stud", "data": {"round": "Fifth Street", ...}}` (結構同上)
-   **4.2.3.2.3-L6**: (Frontend) 接收明牌並追加顯示。
-   **4.2.3.2.4-L6**: (Frontend) 實現發牌動畫。
-   **4.2.3.2.5-L6**: 注意：從第五街開始，下注額通常變為大注 (`BigBetLimit`)。發牌後應觸發下注輪，下注輪邏輯需處理大小注切換。

---

## 4.2.4：發第六街 (Sixth Street) (Backend: Go)

### 4.2.4.1：執行發牌 (一明牌)
-   **4.2.4.1.1-L6**: 在 `Table.DealSixthStreet()` 方法中實現 (Fifth Street 下注結束後)。
-   **4.2.4.1.2-L6**: 不燒牌。
-   **4.2.4.1.3-L6**: 從 Bring-in 玩家**之後**第一個未蓋牌的玩家開始，順時針給每個未蓋牌的玩家發一張**明牌** (`Deck.DrawCard()`)。
-   **4.2.4.1.4-L6**: 將此牌加入 `Player.Hand` 和 `Player.VisibleCards`。
-   **4.2.4.1.5-L6**: 檢查 `DrawCard` 錯誤。

### 4.2.4.2：更新狀態與通知
-   **4.2.4.2.1-L6**: 更新遊戲階段狀態 (`GameStateSixthStreetBetting`)。
-   **4.2.4.2.2-L6**: **公開地** 廣播所有玩家的第六張明牌。
    -   `{"event": "dealing_complete_stud", "data": {"round": "Sixth Street", ...}}` (結構同上)
-   **4.2.4.2.3-L6**: (Frontend) 接收明牌並追加顯示。
-   **4.2.4.2.4-L6**: (Frontend) 實現發牌動畫。
-   **4.2.4.2.5-L6**: 下注額繼續使用大注 (`BigBetLimit`)。

---

## 4.2.5：發第七街 (Seventh Street / River) (Backend: Go)

### 4.2.5.1：執行發牌 (一暗牌 或 公共牌)
-   **4.2.5.1.1-L6**: 在 `Table.DealSeventhStreet()` 方法中實現 (Sixth Street 下注結束後)。
-   **4.2.5.1.2-L6**: 不燒牌。
-   **4.2.5.1.3-L6**: **檢查牌堆剩餘數量**: `Deck.RemainingCards()` 是否小於未蓋牌玩家數量。
-   **4.2.5.1.4-L7**: **情況 A：牌夠用**
    -   從 Bring-in 玩家**之後**第一個未蓋牌的玩家開始，順時針給每個未蓋牌的玩家發一張**暗牌** (`Deck.DrawCard()`)。
    -   將此牌只加入 `Player.Hand` (不加入 `VisibleCards`)。
    -   **私密地** 將第七張暗牌信息發送給每個對應玩家。
        - `{"event": "seventh_street_card", "data": {"card": {"id": "C3"}}}`
    -   廣播發牌完成通知 (不含牌面)。
        - `{"event": "dealing_complete_stud", "data": {"round": "Seventh Street"}}`
-   **4.2.5.1.5-L7**: **情況 B：牌不夠用** (罕見)
    -   從牌堆頂抽一張牌 (`Deck.DrawCard()`) 作為**公共牌** (Community Card)。
    -   將此公共牌存儲在 `Table` 的一個新字段中，例如 `SeventhStreetCommunityCard *Card`。
    -   所有未蓋牌的玩家共享這張公共牌。
    -   **公開地** 廣播這張公共牌。
        - `{"event": "seventh_street_community_card", "data": {"card": {"id": "S9"}}}`
    -   廣播發牌完成通知。
        - `{"event": "dealing_complete_stud", "data": {"round": "Seventh Street", "used_community_card": true}}`
-   **4.2.5.1.6-L6**: 檢查 `DrawCard` 錯誤。

### 4.2.5.2：更新狀態與通知
-   **4.2.5.2.1-L6**: 更新遊戲階段狀態 (`GameStateSeventhStreetBetting`)。
-   **4.2.5.2.2-L6**: (Frontend) 接收第七張牌信息（暗牌或公共牌）並顯示。
-   **4.2.5.2.3-L6**: (Frontend) 實現發牌動畫。
-   **4.2.5.2.4-L6**: 下注額繼續使用大注 (`BigBetLimit`)。

--- 