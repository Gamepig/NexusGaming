# 梭哈 (Stud Poker) 遊戲邏輯詳細規劃

本文檔詳細記錄了七張梭哈 (7-Card Stud) 遊戲的核心邏輯設計，涵蓋初始化、發牌、下注、比牌、狀態管理和手機優化等方面，並包含各主要環節的流程圖。

---

## 任務 4.1：梭哈 - 初始化遊戲
**目標**：設置梭哈遊戲開始前的環境，適應 2-8 人桌配置，處理底注 (Ante)。

### 4.1.1：創建遊戲配置
- **4.1.1.1：定義底注 (Ante) 結構**
    - 4.1.1.1.1：設置每位玩家的底注金額 (例如 1 單位)。
    - 4.1.1.1.2：創建配置文件儲存底注值。
- **4.1.1.2：設置玩家上下限**
    - 4.1.1.2.1：設置最小玩家數為 2。
    - 4.1.1.2.2：設置最大玩家數為 8 (可配置)。
    - 4.1.1.2.3：記錄當前玩家數。
    - 4.1.1.2.4：添加檢查確保玩家數在範圍內。
- **4.1.1.3：初始化獎池**
    - 4.1.1.3.1：創建獎池物件 (主池 + 邊池陣列，結構同 Hold'em)。
    - 4.1.1.3.2：初始值為所有玩家底注的總和。
    - 4.1.1.3.3：提供獎池更新接口。
- **4.1.1.4：設置遊戲參數**
    - 4.1.1.4.1：定義輪次/街道名稱 (例如 "Third Street", "Fourth Street", ..., "Seventh Street/River", "Showdown")。
    - 4.1.1.4.2：設置操作時間 (例如 15 秒)。
    - 4.1.1.4.3：設置下注限額結構 (固定限額 Fixed Limit: 小注/大注)。
        - 4.1.1.4.3.1: 定義小注 (Small Bet) 金額 (例如前兩輪下注用)。
        - 4.1.1.4.3.2: 定義大注 (Big Bet) 金額 (例如後幾輪下注用，通常為小注2倍)。
        - 4.1.1.4.3.3: 定義每輪加注次數上限 (例如 1 Bet + 3 Raises)。
    - 4.1.1.4.4：設計參數配置文件。

### 4.1.2：初始化玩家 (結構類似 Hold'em)
- **4.1.2.1：創建玩家數據結構**
    - 4.1.2.1.1：屬性 (ID, 名稱, 籌碼, 手牌[], 亮牌[], 暗牌[], 狀態, 座位索引)。
    - 4.1.2.1.2：亮牌 (Exposed Cards) 陣列。
    - 4.1.2.1.3：暗牌 (Hole Cards) 陣列。
- **4.1.2.2：分配玩家座位**
    - (同 Hold'em 3.1.2.2)
- **4.1.2.3：初始化玩家籌碼**
    - (同 Hold'em 3.1.2.3)
- **4.1.2.4：扣除底注 (Ante)**
    - 4.1.2.4.1：從每位入座玩家籌碼扣除底注金額。
    - 4.1.2.4.2：將底注加入獎池。
    - 4.1.2.4.3：更新玩家籌碼與獎池狀態。
- **4.1.2.5：手機端顯示**
    - (類似 Hold'em 3.1.2.4，需顯示亮牌)

### 4.1.3：初始化牌堆 (同 Hold'em)
- **4.1.3.1：創建撲克牌數據**
- **4.1.3.2：實現洗牌算法**
- **4.1.3.3：設置牌堆狀態**

### 流程圖 (任務 4.1)

#### 流程圖：創建遊戲配置
```
開始
  ↓
讀取遊戲配置文件
  ↓
設置底注 (Ante) 結構
  ↓ Ante 金額 (例如 1)
  ↓ 儲存至配置物件
設置玩家上下限
  ↓ 最小 2 人
  ↓ 最大 8 人 (可配)
  ↓ 記錄當前玩家數
初始化獎池
  ↓ 創建主池 (初始為 Ante 總和)
  ↓ 創建邊池陣列 (空)
設置遊戲參數
  ↓ 定義輪次名稱 (Third St, ..., Seventh St, Showdown)
  ↓ 設置操作時間 (15s)
  ↓ 設置下注限額 (Fixed Limit: 小注/大注, 加注上限)
  ↓ 儲存至配置物件
結束
```

#### 流程圖：初始化玩家
```
開始
  ↓
創建玩家數據結構 (含亮牌[], 暗牌[])
  ↓
分配玩家座位
  ↓
初始化玩家籌碼
  ↓
扣除底注 (Ante)
  ↓ 從每位玩家扣除 Ante
  ↓ 加入獎池
  ↓ 更新籌碼與獎池
手機端顯示 (顯示亮牌)
結束
```

#### 流程圖：初始化牌堆 (同 Hold'em)
```
開始
  ↓
創建撲克牌數據
  ↓
洗牌 (Fisher-Yates)
  ↓
設置牌堆狀態
結束
```
---

## 任務 4.2：梭哈 - 發牌邏輯
**目標**：實現七張梭哈的發牌流程，區分亮牌與暗牌。

### 4.2.1：初始發牌 (Third Street)
- **4.2.1.1：分配初始三張牌**
    - 4.2.1.1.1：從牌堆依次給每位玩家發 **兩張暗牌 (Hole Cards)**。
    - 4.2.1.1.2：再依次給每位玩家發 **一張亮牌 (Door Card / Third Street)**。
    - 4.2.1.1.3：按順時針順序發牌。
    - 4.2.1.1.4：更新牌堆狀態。
- **4.2.1.2：記錄牌面**
    - 4.2.1.2.1：將暗牌 ID 存入玩家 `暗牌[]`。
    - 4.2.1.2.2：將亮牌 ID 存入玩家 `亮牌[]`。
    - 4.2.1.2.3：確保暗牌僅對本人可見。
    - 4.2.1.2.4：提供亮牌查詢接口 (UI 及下注邏輯)。
- **4.2.1.3：確定 "Bring-in"**
    - 4.2.1.3.1：比較所有玩家的第一張亮牌 (Door Card)。
    - 4.2.1.3.2：擁有 **最低** 點數亮牌的玩家必須下 "Bring-in" 注 (點數相同比花色: C < D < H < S)。
    - 4.2.1.3.3：Bring-in 金額通常固定 (例如半個小注)。
    - 4.2.1.3.4：記錄 Bring-in 玩家索引。
    - 4.2.1.3.5：Bring-in 玩家可選擇只下 Bring-in 或直接完成一個小注 (Complete the bet)。

### 4.2.2：後續發牌 (Fourth, Fifth, Sixth Street)
- **4.2.2.1：發放第四張牌 (Fourth Street)**
    - 4.2.2.1.1：給每位 **未棄牌** 的玩家發一張 **亮牌**。
    - 4.2.2.1.2：更新牌堆狀態。
    - 4.2.2.1.3：記錄亮牌 ID。
    - 4.2.2.1.4：觸發 Fourth Street 下注輪。
- **4.2.2.2：發放第五張牌 (Fifth Street)**
    - 4.2.2.2.1：給每位未棄牌玩家發一張 **亮牌**。
    - 4.2.2.2.2：更新牌堆狀態。
    - 4.2.2.2.3：記錄亮牌 ID。
    - 4.2.2.2.4：觸發 Fifth Street 下注輪 (通常從此輪開始使用大注)。
- **4.2.2.3：發放第六張牌 (Sixth Street)**
    - 4.2.2.3.1：給每位未棄牌玩家發一張 **亮牌**。
    - 4.2.2.3.2：更新牌堆狀態。
    - 4.2.2.3.3：記錄亮牌 ID。
    - 4.2.2.3.4：觸發 Sixth Street 下注輪 (使用大注)。

### 4.2.3：發放河牌 (Seventh Street / River)
- **4.2.3.1：發放第七張牌**
    - 4.2.3.1.1：給每位未棄牌玩家發一張 **暗牌**。
    - 4.2.3.1.2：更新牌堆狀態。
    - 4.2.3.1.3：記錄暗牌 ID (僅本人可見)。
    - 4.2.3.1.4：觸發 Seventh Street 最後一輪下注 (使用大注)。

### 4.2.4：發牌動畫接口 (類似 Hold'em)
- **4.2.4.1：提供發牌數據** (需區分亮/暗牌目標位置和狀態)。
- **4.2.4.2：定義動畫參數**。
- **4.2.4.3：手機端顯示** (清晰展示亮牌，隱藏暗牌)。

### 流程圖 (任務 4.2)

#### 流程圖：初始發牌 (Third Street)
```
開始
  ↓
發兩張暗牌/玩家
  ↓
發一張亮牌/玩家 (Door Card)
  ↓
更新牌堆
  ↓
記錄暗牌/亮牌
  ↓
確定 Bring-in 玩家 (最低亮牌)
  ↓ Bring-in 玩家下注 (Bring-in 或 Complete)
結束
```

#### 流程圖：後續發牌 (Fourth, Fifth, Sixth Street)
```
開始
  ↓
循環 (Fourth, Fifth, Sixth Street):
  ↓   發一張亮牌/未棄牌玩家
  ↓   更新牌堆
  ↓   記錄亮牌
  ↓   觸發對應下注輪 (Fifth St 開始可能用大注)
結束
```

#### 流程圖：發放河牌 (Seventh Street)
```
開始
  ↓
發一張暗牌/未棄牌玩家
  ↓
更新牌堆
  ↓
記錄暗牌
  ↓
觸發 Seventh Street 下注輪 (大注)
結束
```

#### 流程圖：發牌動畫接口 (區分亮暗牌)
```
開始
  ↓
提供發牌數據 (含亮/暗狀態)
  ↓
定義動畫參數
  ↓
手機端顯示 (渲染亮/暗牌)
結束
```
---

## 任務 4.3：梭哈 - 下注輪管理
**目標**：實現梭哈多輪下注，處理 Bring-in 和固定限額。

### 4.3.1：設置下注輪結構
- **4.3.1.1：定義輪次** (Third, Fourth, Fifth, Sixth, Seventh Street)。
- **4.3.1.2：設置起始玩家**
    - 4.3.1.2.1：Third Street：Bring-in 玩家左手邊第一位未棄牌玩家開始。
    - 4.3.1.2.2：Fourth Street 及之後：**亮牌牌面最大** 的玩家先開始下注 (若牌面相同，從莊家左手邊最近的開始)。
    - 4.3.1.2.3：記錄起始玩家索引。
- **4.3.1.3：記錄輪次狀態** (同 Hold'em)。

### 4.3.2：實現玩家操作 (固定限額)
- **4.3.2.1：定義操作選項**
    - 4.3.2.1.1：下注 (Bet)：在無人下注時，下一個固定金額 (小注或大注)。
    - 4.3.2.1.2：跟注 (Call)：跟上當前輪的總下注額。
    - 4.3.2.1.3：加注 (Raise)：在 Bet/Call 基礎上增加一個固定額度 (小注或大注)。
    - 4.3.2.1.4：棄牌 (Fold)。
    - 4.3.2.1.5：過牌 (Check)：僅在無人下注時可用 (通常從 Fourth Street 開始)。
- **4.3.2.2：手機端操作**
    - 4.3.2.2.1：按鈕應顯示固定金額 (例如 "Bet 10", "Raise to 20")。
    - 4.3.2.2.2：移除加注滑桿，改為固定加注按鈕。
    - (其他類似 Hold'em)。
- **4.3.2.3：驗證操作**
    - 4.3.2.3.1：檢查籌碼是否足夠 Bet/Call/Raise。
    - 4.3.2.3.2：檢查是否達到每輪加注次數上限。
    - 4.3.2.3.3：檢查操作是否符合當前狀態 (例如不能在有人 Bet 後 Check)。
- **4.3.2.4：記錄操作** (同 Hold'em)。

### 4.3.3：管理下注邏輯
- **4.3.3.1：追蹤下注狀態** (同 Hold'em)。
- **4.3.3.2：結束下注輪**
    - 4.3.3.2.1：檢查所有未棄牌玩家是否已完成操作 (下注額相同或棄牌)。
    - 4.3.3.2.2：若只剩 1 名玩家未棄牌，結束遊戲。
    - 4.3.3.2.3：觸發輪次結束事件。
    - 4.3.3.2.4：重置本輪下注狀態。
- **4.3.3.3：處理全下** (類似 Hold'em，可能需要邊池)。
- **4.3.3.4：手機端提示** (顯示當前下注額、輪次、玩家亮牌)。

### 4.3.4：計時器與提示 (同 Hold'em)

### 流程圖 (任務 4.3)

#### 流程圖：設置下注輪結構
```
開始
  ↓
定義輪次 (Third St, ..., Seventh St)
  ↓
設置起始玩家
  ↓ Third St: Bring-in 左邊
  ↓ Fourth St+: 最大亮牌牌面先
  ↓ 記錄索引
記錄輪次狀態
結束
```

#### 流程圖：實現玩家操作 (固定限額)
```
開始
  ↓
定義操作 (Bet, Call, Raise[固定額], Fold, Check)
  ↓
手機端操作 (固定額按鈕)
  ↓
驗證操作
  ↓ 檢查籌碼
  ↓ 檢查加注上限
  ↓ 檢查狀態合法性
記錄操作
結束
```

#### 流程圖：管理下注邏輯
```
開始
  ↓
追蹤下注狀態
  ↓
結束下注輪判斷
  ↓ 所有玩家操作完成? -> 結束輪次
  ↓ 剩 1 人? -> 結束遊戲
  ↓
處理全下 (邊池邏輯)
  ↓
手機端提示 (亮牌, 限額)
結束
```

#### 流程圖：計時器與提示 (同 Hold'em)
```
開始
  ↓
設置計時器 (15s)
  ↓
處理超時 (棄牌/過牌)
  ↓
手機端提示 (進度條, 高亮)
結束
```
---

## 任務 4.4：梭哈 - 比牌與獎池分配
**目標**：從七張牌中選出最佳五張牌比較大小，分配獎池。

### 4.4.1：實現牌型計算
- **4.4.1.1：定義牌型** (同 Hold'em，10 種)。
- **4.4.1.2：計算玩家最佳牌型**
    - 4.4.1.2.1：從玩家的 7 張牌 (3 暗 + 4 亮) 中枚舉所有 5 張組合 (C(7,5)=21 種)。
    - 4.4.1.2.2：檢查每種組合的牌型等級。
    - 4.4.1.2.3：選出最高牌型 (記錄 5 張牌與等級)。
- **4.4.1.3：優化計算** (同 Hold'em)。

### 4.4.2：判定贏家
- **4.4.2.1：比較牌型** (同 Hold'em)。
- **4.4.2.2：處理平局** (同 Hold'em)。
- **4.4.2.3：手機端顯示** (需顯示最終選出的 5 張最佳牌)。

### 4.4.3：分配獎池
- **4.4.3.1：計算主池** (同 Hold'em)。
- **4.4.3.2：處理邊池** (同 Hold'em)。
- **4.4.3.3：動畫接口** (同 Hold'em)。
- **4.4.3.4：手機端提示** (同 Hold'em)。

### 流程圖 (任務 4.4)

#### 流程圖：實現牌型計算
```
開始
  ↓
定義牌型 (10 種)
  ↓
計算玩家最佳牌型
  ↓ 枚舉 7 選 5
  ↓ 檢查牌型等級
  ↓ 選最高牌型
  ↓ 記錄最佳 5 張
優化計算 (可選)
結束
```

#### 流程圖：判定贏家 (同 Hold'em)
```
開始
  ↓
比較牌型 (等級 -> 點數)
  ↓
處理平局 (記錄平局者)
  ↓
手機端顯示 (顯示最佳 5 張)
結束
```

#### 流程圖：分配獎池 (同 Hold'em)
```
開始
  ↓
分配主池 (給贏家/平分)
  ↓
按序處理邊池
  ↓
動畫接口
  ↓
手機端提示
結束
```

---

## 任務 4.5：梭哈 - 狀態管理
**目標**：追蹤梭哈特定的遊戲與玩家狀態。

### 4.5.1：遊戲狀態
- **4.5.1.1：記錄輪次** (Third Street 到 Seventh Street)。
- **4.5.1.2：追蹤獎池** (同 Hold'em)。
- **4.5.1.3：記錄下注** (記錄當前是小注輪還是大注輪，剩餘加注次數)。

### 4.5.2：玩家狀態
- **4.5.2.1：更新狀態** ("active", "folded", "all-in")。
- **4.5.2.2：記錄牌面** (更新 `亮牌[]` 和 `暗牌[]`)。
- **4.5.2.3：手機端顯示** (即時更新亮牌)。

### 4.5.3：輪次切換 (同 Hold'em，但觸發條件是發牌完成)。
### 4.5.4：遊戲結束與重啟 (同 Hold'em，但需重置亮/暗牌)。

### 流程圖 (任務 4.5)

#### 流程圖：遊戲狀態
```
開始
  ↓
記錄輪次 (Third St - Seventh St)
  ↓
追蹤獎池
  ↓
記錄下注 (限額類型, 剩餘加注數)
結束
```

#### 流程圖：玩家狀態
```
開始
  ↓
更新狀態 (active/folded/all-in)
  ↓
記錄牌面 (亮牌[], 暗牌[])
  ↓
手機端顯示 (更新亮牌)
結束
```

#### 流程圖：輪次切換 (同 Hold'em)
```
開始
  ↓
檢測輪次結束
  ↓
重置狀態 (準備下輪發牌/比牌)
  ↓
手機端提示
結束
```

#### 流程圖：遊戲結束與重啟 (同 Hold'em)
```
開始
  ↓
結束本局
  ↓
開始下一局
  ↓ 重置亮/暗牌
  ↓ 扣除新 Ante
手機端提示
結束
```
---

## 任務 4.6：梭哈 - 手機優先優化
**目標**：確保梭哈界面和操作在手機上流暢友好。

### 4.6.1：簡化操作流程
- **4.6.1.1：提供快速操作** (Bet [固定額], Call, Fold 按鈕)。
- **4.6.1.2：清晰顯示限額** (標明當前是小注還是大注)。

### 4.6.2：觸控友好
- **4.6.2.1：設計按鈕** (同 Hold'em)。
- **4.6.2.2：牌面顯示** (確保多張亮牌在手機上清晰可見，不重疊)。

### 4.6.3：視覺提示
- **4.6.3.1：高亮當前玩家** (同 Hold'em)。
- **4.6.3.2：顯示操作提示** (例如 "Bet 10", "Call 20")。
- **4.6.3.3：亮牌提示** (清晰標示每位玩家的亮牌及其組合潛力)。
- **4.6.3.4：倒計時提示** (同 Hold'em)。

### 流程圖 (任務 4.6)

#### 流程圖：簡化操作流程
```
開始
  ↓
提供快速操作 (固定額按鈕)
  ↓
清晰顯示限額 (小注/大注)
結束
```

#### 流程圖：觸控友好
```
開始
  ↓
設計按鈕 (大尺寸, 間距)
  ↓
牌面顯示優化 (手機端亮牌清晰)
結束
```

#### 流程圖：視覺提示
```
開始
  ↓
高亮當前玩家
  ↓
顯示操作提示 (固定金額)
  ↓
亮牌提示 (潛力分析可選)
  ↓
倒計時提示
結束
```

--- 