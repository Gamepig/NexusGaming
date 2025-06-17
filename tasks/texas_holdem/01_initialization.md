# 德州撲克開發任務：初始化遊戲 (Task 3.1)

**目標**：實現遊戲開始前的所有必要設置和數據結構初始化，支援 9 人桌和手機端。

---

## 3.1.1：創建遊戲配置 (Backend: Go)

### 3.1.1.1：定義盲注結構
-   **3.1.1.1.1-L6**: 在 `config/game_settings.go` (或類似文件) 中定義 `BlindStructure` struct。
    -   `SmallBlindAmount int`
    -   `BigBlindAmount int`
-   **3.1.1.1.2-L6**: 在配置文件 (例如 `config/texas_holdem.json` 或數據庫表 `game_configs`) 中設置默認小盲注值 (10)。
-   **3.1.1.1.3-L6**: 在配置文件/數據庫中設置默認大盲注值 (20)。
-   **3.1.1.1.4-L6**: 創建 `LoadGameConfig(gameType string)` 函數，從配置文件/數據庫加載指定遊戲的配置，返回包含 `BlindStructure` 的 `GameConfig` struct。
-   **3.1.1.1.5-L6**: (可選) 在 `BlindStructure` struct 中添加 `IncreaseRule` 字段 (例如 `{"interval_minutes": 10, "multiplier": 2}`)。
-   **3.1.1.1.6-L6**: (可選) 實現盲注遞增的計時器邏輯 (可能放在遊戲主循環或獨立 goroutine 中)。
-   **3.1.1.1.7-L6**: 設計 `GameConfig` struct 時考慮通用性，允許梭哈的 Ante 或百家樂無盲注。

### 3.1.1.2：設置玩家上限
-   **3.1.1.2.1-L6**: 在 `GameConfig` struct 中添加 `MaxPlayers int` 字段。
-   **3.1.1.2.2-L6**: 在德州撲克的配置文件/數據庫條目中設置 `MaxPlayers` 為 9。
-   **3.1.1.2.3-L6**: 在 `game/table.go` (或類似文件) 中定義 `Table` struct，包含 `Players []*Player` 和 `CurrentPlayerCount int`。
-   **3.1.1.2.4-L6**: 實現 `Table.AddPlayer(player *Player)` 方法，內部檢查 `CurrentPlayerCount < MaxPlayers`，超限則返回錯誤。
-   **3.1.1.2.5-L6**: 確保梭哈可以配置不同的 `MaxPlayers` 值 (例如 8)。

### 3.1.1.3：初始化獎池
-   **3.1.1.3.1-L6**: 在 `game/pot.go` (或類似文件) 中定義 `Pot` struct。
    -   `MainPotAmount int`
    -   `SidePots []SidePot`
-   **3.1.1.3.2-L6**: 定義 `SidePot` struct。
    -   `Amount int`
    -   `EligiblePlayerIDs []string` (或其他唯一標識符)
-   **3.1.1.3.3-L6**: 在 `Table` struct 中添加 `Pot Pot` 成員。
-   **3.1.1.3.4-L6**: 在 `Table.ResetForNewHand()` 方法中初始化 `Pot.MainPotAmount = 0` 和 `Pot.SidePots = []SidePot{}`。
-   **3.1.1.3.5-L6**: 實現 `Pot.AddChips(amount int, playerID string)` 方法 (處理籌碼加入主池或邊池的基礎邏輯，邊池創建邏輯在下注環節處理)。
-   **3.1.1.3.6-L6**: 實現 `Pot.GetTotalPot()` 方法，計算總獎池。
-   **3.1.1.3.7-L6**: 定義獎池狀態更新時通知前端的 WebSocket 消息結構 (例如 `{"event": "pot_update", "data": {"main_pot": 100, "side_pots": [...]}}`)。

### 3.1.1.4：設置遊戲參數
-   **3.1.1.4.1-L6**: 在 `GameConfig` struct 中添加 `RoundNames []string` (例如 `["Pre-Flop", "Flop", "Turn", "River", "Showdown"]`)。
-   **3.1.1.4.2-L6**: 在 `GameConfig` struct 中添加 `ActionTimeLimitSeconds int`。
-   **3.1.1.4.3-L6**: 在德州撲克配置文件/數據庫中設置 `ActionTimeLimitSeconds` 為 15。
-   **3.1.1.4.4-L6**: 在 `game/player.go` (或類似文件) 中定義 `Player` struct，包含 `ActionTimer *time.Timer` 和 `RemainingActionTime time.Duration`。 (計時器管理可能放在 `Table` 或遊戲主循環中)。
-   **3.1.1.4.5-L6**: 確保 `GameConfig` 結構支持不同遊戲調整輪次和時間。

---

## 3.1.2：初始化玩家 (Backend: Go)

### 3.1.2.1：創建玩家數據結構
-   **3.1.2.1.1-L6**: 在 `game/player.go` 中定義 `Player` struct。
    -   `ID string` (唯一標識，例如 UUID 或來自用戶系統的 ID)
    -   `Name string`
    -   `Chips int`
    -   `Hand []*Card` (手牌，Card 結構見 3.1.3)
    -   `Status string` (例如 "waiting", "active", "folded", "all-in", "sitting_out")
    -   `SeatIndex int` (0-8)
-   **3.1.2.1.2-L6**: 定義 PlayerStatus 常量 (例如 `PlayerStatusActive`, `PlayerStatusFolded` 等)。
-   **3.1.2.1.3-L6**: (可選擴展) 添加 `GameSpecificData interface{}` 以支持百家樂押注類型等。

### 3.1.2.2：分配玩家座位
-   **3.1.2.2.1-L6**: 在 `Table` struct 中添加 `Seats [MaxPlayers]*Player` 數組或 Map `map[int]*Player` 來管理座位。
-   **3.1.2.2.2-L6**: 實現 `Table.AssignSeat(player *Player)` 方法，為玩家分配一個空閒座位 (隨機或按順序)，並設置其 `SeatIndex`。
-   **3.1.2.2.3-L6**: 需要區分真實玩家和 AI 玩家 (可能通過 `Player` struct 中的 `IsAI bool` 字段)。
-   **3.1.2.2.4-L6**: 實現加入桌子/遊戲房間的邏輯，調用 `AddPlayer` 和 `AssignSeat`。
-   **3.1.2.2.5-L6**: 遊戲開始前，需要有機制確認哪些座位有玩家加入。

### 3.1.2.3：初始化玩家籌碼
-   **3.1.2.3.1-L6**: 在 `GameConfig` struct 中添加 `StartingChips int`。
-   **3.1.2.3.2-L6**: 在德州撲克配置文件/數據庫中設置 `StartingChips` 為 1000。
-   **3.1.2.3.3-L6**: 在玩家加入桌子時，從其帳戶扣除買入額 (Buy-in) 或直接設置為 `StartingChips` (取決於是錦標賽還是現金局)。**注意：此處涉及帳務系統，需要與帳務模塊接口交互**。
-   **3.1.2.3.4-L6**: 添加 `Player.AddChips(amount int)` 和 `Player.RemoveChips(amount int)` 方法，包含籌碼 > 0 的檢查。
-   **3.1.2.3.5-L6**: 定義玩家籌碼變動時通知前端的 WebSocket 消息結構 (例如 `{"event": "player_chips_update", "data": {"player_id": "...", "chips": 980}}`)。

### 3.1.2.4：手機端顯示 (初步接口定義)
-   **3.1.2.4.1-L6**: 定義獲取所有玩家基本資訊 (供 UI 面板顯示) 的 API 端點或 WebSocket 消息。
    -   返回數據結構應包含 `PlayerID`, `Name`, `Chips`, `SeatIndex`, `Status`。
-   **3.1.2.4.2-L6**: (前端任務) 確保 UI 字體大小符合要求。
-   **3.1.2.4.3-L6**: (前端任務) 實現 9 人桌的滾動/縮放視圖。
-   **3.1.2.4.4-L6**: (前端任務) 設計點擊座位顯示玩家詳細資訊的交互。

---

## 3.1.3：初始化牌堆 (Backend: Go, Common Module)

(建議將 Card 和 Deck 相關邏輯放到通用的 `pkg/card` 或 `pkg/deck` 包中)

### 3.1.3.1：創建撲克牌數據
-   **3.1.3.1.1-L6**: 在 `pkg/card/card.go` 中定義 `Suit` (Spade, Heart, Club, Diamond) 和 `Rank` (Two, Three, ..., King, Ace) 的枚舉或常量。
-   **3.1.3.1.2-L6**: 定義 `Card` struct。
    -   `Suit Suit`
    -   `Rank Rank`
    -   `ID string` (例如 "S2", "HA"，方便序列化和前端使用)
    -   (可選) `Value int` (百家樂計點用)
-   **3.1.3.1.3-L6**: 實現 `NewCard(suit Suit, rank Rank)` 構造函數，自動生成 ID 和 Value。
-   **3.1.3.1.4-L6**: 提供 Card 的 String() 方法方便調試 (例如 `"Ace of Spades"`）。

### 3.1.3.2：實現洗牌算法
-   **3.1.3.2.1-L6**: 在 `pkg/deck/deck.go` 中定義 `Deck` struct。
    -   `Cards []*card.Card`
    -   `DealtCards []*card.Card` (或使用 index 追蹤)
-   **3.1.3.2.2-L6**: 實現 `NewDeck(numDecks int)` 函數，生成包含 N * 52 張標準牌的牌堆。
-   **3.1.3.2.3-L6**: 實現 `Deck.Shuffle()` 方法，使用 `math/rand` 和 Fisher-Yates 算法。
-   **3.1.3.2.4-L6**: 確保使用 `rand.Seed(time.Now().UnixNano())` 或更可靠的隨機源進行初始化。
-   **3.1.3.2.5-L6**: (單元測試) 編寫測試驗證洗牌後的隨機性 (難以完美驗證，但可做基本檢查)。
-   **3.1.3.2.6-L6**: 定義洗牌完成的內部事件或回調。

### 3.1.3.3：設置牌堆狀態與抽牌
-   **3.1.3.3.1-L6**: 在 `Deck` struct 中維護一個 `currentIndex int` 來表示下一張要發的牌。
-   **3.1.3.3.2-L6**: 實現 `Deck.DrawCard()` 方法。
    -   檢查 `currentIndex` 是否越界 (牌已發完)。
    -   返回 `Cards[currentIndex]`。
    -   `currentIndex++`。
    -   返回錯誤如果牌已發完。
-   **3.1.3.3.3-L6**: 實現 `Deck.RemainingCards()` 方法。
-   **3.1.3.3.4-L6**: 實現 `Deck.Reset()` 方法，將 `currentIndex` 設為 0 並重新洗牌。
-   **3.1.3.3.5-L6**: 確保 Deck 狀態可以被遊戲主邏輯 (例如 `Table` struct) 持有和管理。

---

## 3.1.4：分配座位與盲注 (Backend: Go)

### 3.1.4.1：確定莊家位置 (Button)
-   **3.1.4.1.1-L6**: 在 `Table` struct 中添加 `ButtonSeatIndex int`。
-   **3.1.4.1.2-L6**: 在新一局開始時 (例如 `Table.StartNewHand()` 方法)，隨機選擇一個活躍玩家的座位索引作為初始 Button 位置。
-   **3.1.4.1.3-L6**: 實現 `Table.MoveButton()` 方法。
    -   找到當前 Button 的下一個活躍玩家座位索引 (順時針)。
    -   處理循環 (座位 8 -> 座位 0)。
    -   更新 `ButtonSeatIndex`。
-   **3.1.4.1.4-L6**: 定義莊家位置變動時通知前端的 WebSocket 消息結構。

### 3.1.4.2：設置盲注玩家
-   **3.1.4.2.1-L6**: 在 `Table` struct 中添加 `SmallBlindSeatIndex int` 和 `BigBlindSeatIndex int`。
-   **3.1.4.2.2-L6**: 實現 `Table.DetermineBlinds()` 方法 (在 `MoveButton` 後調用)。
    -   找到 Button 左手邊第一個活躍玩家作為小盲。
    -   找到小盲左手邊第一個活躍玩家作為大盲。
    -   處理玩家數不足 (例如兩人單挑 Heads-up) 的特殊盲注規則。
    -   更新 `SmallBlindSeatIndex` 和 `BigBlindSeatIndex`。
-   **3.1.4.2.3-L6**: 定義大小盲位置確定後通知前端的 WebSocket 消息結構。

### 3.1.4.3：扣除盲注
-   **3.1.4.3.1-L6**: 在 `Table.PostBlinds()` 方法中實現。
    -   獲取小盲注金額 (`GameConfig.BlindStructure.SmallBlindAmount`)。
    -   從小盲玩家處 `RemoveChips`，處理籌碼不足 All-in 的情況。
    -   將實際扣除的籌碼 `Pot.AddChips`。
    -   記錄小盲玩家本輪已下注金額。
    -   對大盲玩家執行類似操作 (使用大盲金額)。
-   **3.1.4.3.2-L6**: 調用玩家籌碼和獎池更新的通知。
-   **3.1.4.3.3-L6**: 處理玩家籌碼少於盲注金額導致的 All-in 情況，標記玩家狀態。

### 3.1.4.4：手機端提示 (初步接口定義)
-   **3.1.4.4.1-L6**: WebSocket 消息應包含 Button, SB, BB 的座位索引，前端據此高亮。
-   **3.1.4.4.2-L6**: 盲注金額應包含在遊戲配置或牌桌狀態信息中發送給前端。
-   **3.1.4.4.3-L6**: (前端任務) 實現高亮和金額顯示。
-   **3.1.4.4.4-L6**: (前端任務) 實現盲注籌碼移動動畫 (基於後端通知)。

--- 