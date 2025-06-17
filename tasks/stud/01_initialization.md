# 七張梭哈開發任務：初始化遊戲 (Task 4.1)

**目標**：設置七張梭哈遊戲環境，處理 Ante、玩家配置（最多 8 人）、單副牌牌堆。

---

## 4.1.1：創建遊戲配置 (Backend: Go)

### 4.1.1.1：定義 Ante 和下注限制
-   **4.1.1.1.1-L6**: 在 `config/game_settings.go` 的 `GameConfig` struct 中添加字段 (或使用 map[string]int)。
    -   `AnteAmount int`
    -   `BringInAmount int` (通常是 Ante 的一部分或固定值)
    -   `SmallBetLimit int` (前兩輪下注的基本單位)
    -   `BigBetLimit int` (後三輪下注的基本單位，通常是 SmallBetLimit 的兩倍)
-   **4.1.1.1.2-L6**: 在七張梭哈的配置文件 (例如 `config/stud.json`) 或數據庫 `game_configs` 表中設置這些值 (例如 Ante=1, BringIn=2, SmallBet=5, BigBet=10)。
-   **4.1.1.1.3-L6**: 確保 `LoadGameConfig("stud")` 函數能正確加載這些配置。
-   **4.1.1.1.4-L6**: 考慮固定限注 (Fixed Limit) 和底池限注 (Pot Limit) 的配置差異 (目前按固定限注設計)。

### 4.1.1.2：設置玩家上限
-   **4.1.1.2.1-L6**: 在七張梭哈的配置文件/數據庫條目中設置 `MaxPlayers` 為 8。
-   **4.1.1.2.2-L6**: 確保 `Table.AddPlayer()` 方法使用的 `MaxPlayers` 值來自正確加載的 `GameConfig`。

### 4.1.1.3：初始化獎池
-   **4.1.1.3.1-L6**: 重用德州撲克的 `Pot` struct (`MainPotAmount`, `SidePots`)。邊池邏輯在梭哈中同樣適用於 All-in 情況。
-   **4.1.1.3.2-L6**: `Table.ResetForNewHand()` 方法中同樣初始化 `Pot`。

### 4.1.1.4：設置遊戲參數
-   **4.1.1.4.1-L6**: 在 `GameConfig` struct 中添加梭哈的輪次名稱。
    -   `RoundNames []string` (例如 `["Third Street", "Fourth Street", "Fifth Street", "Sixth Street", "Seventh Street (River)", "Showdown"]`)
-   **4.1.1.4.2-L6**: 梭哈通常不設嚴格的每步操作時間限制，但可以配置一個總時間或輪次時間 (`ActionTimeLimitSeconds`) 以防拖延，加載到 `GameConfig`。

---

## 4.1.2：初始化玩家 (Backend: Go)

### 4.1.2.1：創建玩家數據結構
-   **4.1.2.1.1-L6**: 重用德州撲克的 `Player` struct。
    -   `ID string`
    -   `Name string`
    -   `Chips int`
    -   `Hand []*Card` (梭哈中最多 7 張牌)
    -   `VisibleCards []*Card` **(梭哈特有)**: 存儲玩家亮出的牌 (第 3 到第 6 張)。
    -   `Status string` (同 Hold'em: "waiting", "active", "folded", "all-in", "sitting_out")
    -   `SeatIndex int` (0-7)
-   **4.1.2.1.2-L6**: PlayerStatus 常量復用。

### 4.1.2.2：分配玩家座位
-   **4.1.2.2.1-L6**: 重用德州撲克的座位管理邏輯 (`Table.Seats`, `Table.AssignSeat`)，但 `MaxPlayers` 使用 8。
-   **4.1.2.2.2-L6**: 區分真人/AI 玩家。
-   **4.1.2.2.3-L6**: 加入桌子邏輯。

### 4.1.2.3：初始化玩家籌碼
-   **4.1.2.3.1-L6**: 在 `GameConfig` 中設置梭哈的 `StartingChips`。
-   **4.1.2.3.2-L6**: 玩家加入時設置籌碼 (涉及帳務)。
-   **4.1.2.3.3-L6**: 重用 `Player.AddChips`, `Player.RemoveChips`。
-   **4.1.2.3.4-L6**: 重用籌碼變動通知。

### 4.1.2.4：手機端顯示 (初步接口定義)
-   **4.1.2.4.1-L6**: 獲取玩家基本信息的接口/消息需要額外包含 `VisibleCards` 字段。
    -   返回結構: `PlayerID`, `Name`, `Chips`, `SeatIndex`, `Status`, `VisibleCards: [{"id": "HA"}, {"id": "DK"}, ...]`
-   **4.1.2.4.2-L6**: (Frontend) 實現最多 8 人的座位佈局。
-   **4.1.2.4.3-L6**: (Frontend) 玩家座位旁需要預留空間顯示最多 4 張亮牌。

---

## 4.1.3：初始化牌堆 (Backend: Go, Common Module)

### 4.1.3.1：創建撲克牌數據
-   **4.1.3.1.1-L6**: 重用 `pkg/card` 中的 `Card`, `Suit`, `Rank` 定義。
-   **4.1.3.1.2-L6**: 需要為梭哈牌比較實現特定的 Rank 值（例如 A 通常只算高牌，2 是最低牌）。
    -   可以為 `Card` 添加 `StudRankValue int` 字段，或在比較邏輯中處理。
    -   或者創建一個 `GetStudRankValue(card *Card)` 函數。

### 4.1.3.2：實現洗牌算法
-   **4.1.3.2.1-L6**: 重用 `pkg/deck` 中的 `Deck` struct。
-   **4.1.3.2.2-L6**: 調用 `NewDeck(1)` 創建**一副**標準牌堆。
-   **4.1.3.2.3-L6**: 重用 `Deck.Shuffle()` 方法。
-   **4.1.3.2.4-L6**: 確保隨機數種子初始化。

### 4.1.3.3：設置牌堆狀態與抽牌
-   **4.1.3.3.1-L6**: 重用 `Deck.currentIndex`。
-   **4.1.3.3.2-L6**: 重用 `Deck.DrawCard()` 方法。
-   **4.1.3.3.3-L6**: 重用 `Deck.RemainingCards()` 和 `Deck.Reset()`。
-   **4.1.3.3.4-L6**: 需要注意梭哈在第七張街 (River) 可能因玩家過多而耗盡牌堆的情況 (雖然罕見)。
    -   如果發第七張牌時牌堆耗盡，則發一張公共牌 (Community Card) 代替。需要在發牌邏輯 (4.2) 中處理此特殊情況。

---

## 4.1.4：收取 Ante (Backend: Go)

### 4.1.4.1：實現 Ante 收取邏輯
-   **4.1.4.1.1-L6**: 在 `Table.StartNewHand()` 方法中，發牌**之前**調用 `Table.CollectAnte()`。
-   **4.1.4.1.2-L6**: `CollectAnte()` 方法：
    -   獲取 Ante 金額 (`GameConfig.AnteAmount`)。
    -   遍歷所有即將參與牌局的玩家 (狀態為 Active 或 Waiting)。
    -   為每個玩家調用 `Player.RemoveChips(AnteAmount)`，處理籌碼不足 All-in 的情況。
    -   將實際收取的 Ante 加入底池 `Pot.AddChips(actualAnteAmount, playerID)`。
    -   廣播玩家籌碼更新通知和底池更新通知。
-   **4.1.4.1.3-L6**: 處理玩家籌碼少於 Ante 導致的 All-in 情況，更新玩家狀態。

### 4.1.4.2：Bring-in (稍後確定)
-   **4.1.4.2.1-L6**: Bring-in (強制下注) 是在發完第三張街後才確定的，不由初始化階段處理。
-   **4.1.4.2.2-L6**: 初始化階段只需確保 `BringInAmount` 已從配置中加載。確定 Bring-in 玩家的邏輯放在發牌 (4.2) 或下注 (4.3) 階段。

### 4.1.4.3：手機端提示 (初步接口定義)
-   **4.1.4.3.1-L6**: 廣播 Ante 收取完成事件或在底池更新通知中體現。
-   **4.1.4.3.2-L6**: (Frontend) 實現收取 Ante 的籌碼動畫。
-   **4.1.4.3.3-L6**: (Frontend) 在界面上顯示 Ante 金額。

--- 