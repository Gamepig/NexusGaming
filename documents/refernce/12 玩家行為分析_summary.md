## 文件 12: 玩家行為分析

**核心目標**: 深入分析玩家遊戲行為，作為報表功能的一部分，支持運營決策。

**分析維度**:
*   下注行為 (金額、頻率、偏好)。
*   勝率與盈虧。
*   參與頻率 (次數、活躍天數、時長)。
*   活躍時間分佈。
*   異常行為檢測 (刷返水、高勝率)。

**技術棧**:
*   **後端**: Golang (Gin, GORM)。
*   **資料庫**: MySQL/PostgreSQL (`players`, `game_logs`, `point_transactions`), Redis (快取)。
*   **前端 (後台管理)**: Next.js (整合 Chart.js)。
*   **部署**: Docker + Kubernetes。

**功能實現規劃**:

1.  **DB 設計**:
    *   依賴現有表。
    *   新增 `player_behavior_snapshots` 表 (JSON 儲存時間分佈)。
2.  **後端 (Golang)**:
    *   **API**: 
        *   `GET /admin/reports/player_behavior`: 查詢/計算行為快照。
        *   `GET /admin/reports/player_behavior/export`: 匯出 CSV。
        *   `GET /admin/reports/abnormal_behavior`: 查詢異常行為。
    *   **計算邏輯 (`calculatePlayerBehavior`)**: 
        *   聚合 `game_logs` 計算總下注、勝率、盈虧、活躍天數。
        *   計算平均 Session 時長 (e.g., >30 分鐘間隔為新 Session)。
        *   計算每小時下注分佈。
    *   **異常檢測 (`detectAbnormalBehavior`)**: SQL 檢測高頻小額下注、異常高勝率。
    *   **定時任務 (cron `@weekly`)**: 生成 `player_behavior_snapshots` (`generateBehaviorSnapshots`)。
3.  **前端 (Next.js)**:
    *   `pages/admin/player-behavior.tsx`: 展示報表，使用 Chart.js (Bar + Line 混合圖) 顯示下注/勝率，渲染表格。
    *   `components/TimeDistributionHeatmap.tsx`: 視覺化時間分佈。
    *   `pages/admin/abnormal-behavior.tsx`: 展示異常行為列表。
4.  **優化**: 
    *   Redis 快取快照 (`getPlayerBehavior`)。
    *   分頁查詢 (`/admin/reports/player_behavior_paginated`)。
    *   Goroutines 並行計算 (`parallelCalculateBehavior`)。
5.  **安全**: 
    *   JWT 權限控制 (`authBehaviorReport`)。
    *   AES 加密敏感數據。
    *   GORM 防注入。

**整合**: 
*   **遊戲端**: WebSocket 記錄 `game_logs`。
*   **報表功能**: 行為分析作為核心模組嵌入，異常報表單獨展示。

**記憶整合**: 
*   行為分析的模式識別、概率計算、視覺化契合交易分析經驗。
*   未來可引入 ML 模型進行異常檢測，契合 AI 優化經驗。

**實現建議 (優先級)**:
1.  基礎行為分析 (核心報表)。
2.  時間分佈與異常檢測。
3.  視覺化與匯出。

**結論**: 玩家行為分析功能提供了深入的運營洞察，利用 Golang 的並發能力和 Next.js/Chart.js 的視覺化能力，可高效實現並與現有後台整合。 