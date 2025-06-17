# 七張梭哈開發任務：手機優先優化 (Task 4.6)

**目標**：確保七張梭哈遊戲在手機端的界面佈局適應亮牌顯示、操作流暢、信息清晰。

---

## 4.6.1：簡化操作流程 (Frontend Focus, Backend Support)

### 4.6.1.1：佈局設計 (適應亮牌)
-   **4.6.1.1.1-L6**: (Frontend) 設計能清晰展示每個玩家最多 4 張亮牌的座位佈局。
    -   可能需要比 Hold'em 更緊湊的佈局或將亮牌堆疊顯示。
-   **4.6.1.1.2-L6**: (Frontend) 優先顯示自己的暗牌、所有玩家的亮牌、底池、輪到誰行動。
-   **4.6.1.1.3-L6**: (Frontend) 操作按鈕區域設計 (固定底部，同 Hold'em)。

### 4.6.1.2：動作按鈕優化 (固定限注)
-   **4.6.1.2.1-L6**: (Backend) `your_turn` 通知需要明確可用動作和對應的**固定**下注/加注額。
    -   包含 `amount_to_call`。
    -   包含 `can_check`, `can_bet`, `can_raise`。
    -   包含 `bet_amount = CurrentBetLimit`。
    -   包含 `raise_amount = CurrentBetAmount + CurrentBetLimit` (加注**到**的總額)。
    -   (Third St Bring-in) 特殊選項 `bring_in_amount`, `complete_bet_amount`。
    -   (Fourth St Open Pair) 特殊選項 `can_bet_big`, `big_bet_amount`。
-   **4.6.1.2.2-L6**: (Frontend) 根據通知動態生成按鈕。
    -   Bet 按鈕直接顯示固定額度 (例如 "Bet $10")。
    -   Raise 按鈕直接顯示固定額度 (例如 "Raise to $20")。
    -   **不需要** 下注滑塊。
-   **4.6.1.2.3-L6**: (Frontend) Check/Call 按鈕合併。
-   **4.6.1.2.4-L6**: (Frontend) Fold 按鈕。

### 4.6.1.3：快捷操作
-   **4.6.1.3.1-L6**: (Frontend) 可以提供 "預選動作" (Check/Fold, Call Any)。

---

## 4.6.2：觸控友好 (Frontend Focus)

### 4.6.2.1：交互元素設計
-   **4.6.2.1.1-L6**: (Frontend) 確保按鈕、座位、亮牌區域有足夠觸控空間。
-   **4.6.2.1.2-L6**: (Frontend) 清晰的觸控反饋。

### 4.6.2.2：手勢支持 (可選)
-   **4.6.2.2.1-L6**: (Frontend) 同 Hold'em，考慮滑動 Fold 等。

### 4.6.2.3：響應性
-   **4.6.2.3.1-L6**: (Frontend) 快速響應，優化性能。

---

## 4.6.3：視覺提示與信息傳達 (Frontend Focus, Backend Support)

### 4.6.3.1：狀態高亮
-   **4.6.3.1.1-L6**: (Backend) 通知需包含行動玩家、Bring-in 玩家、亮牌最高玩家等信息。
-   **4.6.3.1.2-L6**: (Frontend) 清晰高亮當前行動玩家。
-   **4.6.3.1.3-L6**: (Frontend) 明顯標示 Bring-in 玩家 (Third St)。
-   **4.6.3.1.4-L6**: (Frontend) 視覺上區分 Folded/All-in 玩家。
-   **4.6.3.1.5-L6**: (Frontend) 清晰展示每個玩家的所有亮牌 (`VisibleCards`)。

### 4.6.3.2：信息顯示
-   **4.6.3.2.1-L6**: (Frontend) 清晰顯示籌碼、下注額、底池大小。
-   **4.6.3.2.2-L6**: (Frontend) 牌面顯示清晰。
-   **4.6.3.2.3-L6**: (Frontend) 攤牌時顯示最佳牌型。
-   **4.6.3.2.4-L6**: (Frontend) 簡潔的文字提示 (例如 "Player X Brings In", "Player Y Bets $10")。
-   **4.6.3.2.5-L6**: (Frontend) 顯示當前下注限額 (Small Bet / Big Bet)。

### 4.6.3.3：動畫效果
-   **4.6.3.3.1-L6**: (Backend) 發送觸發動畫的事件通知。
-   **4.6.3.3.2-L6**: (Frontend) 實現梭哈的發牌動畫 (兩暗一明，後續逐張明牌，第七張暗牌/公共牌)。
-   **4.6.3.3.3-L6**: (Frontend) 實現籌碼移動動畫 (Ante, Bring-in, Bet, Raise, Call, 收池, 分池)。
-   **4.6.3.3.4-L6**: (Frontend) 實現攤牌動畫 (翻開暗牌)。
-   **4.6.3.3.5-L6**: (Frontend) 允許關閉/加速動畫。

### 4.6.3.4：網絡狀態提示
-   **4.6.3.4.1-L6**: (Frontend) 同 Hold'em。

--- 