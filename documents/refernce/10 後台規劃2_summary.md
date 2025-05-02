## 文件 10: 後台規劃 2 (儲值、分潤連動、返水)

**核心目標**: 擴充現有後台功能，加入儲值、分潤配置連動及返水機制。

**新增功能需求**:
1.  **開分系統 (儲值)**: 玩家、經銷商、代理商儲值流程管理與審核。
2.  **分潤配置連動**: 代理商/總公司修改分潤比例時，自動更新下級配置。
3.  **返水規則**: 設計並實現基於下注額的返水機制。

**技術棧 (後台)**:
*   **後端**: Golang (Gin 框架, GORM ORM)。
*   **資料庫**: MySQL/PostgreSQL (交易/玩家數據), Redis (快取/即時狀態)。
*   **前端 (後台管理)**: Next.js。
*   **部署**: Docker + Kubernetes。

**功能實現規劃**:

1.  **開分系統 (儲值)**:
    *   **DB**: `deposit_requests` 表記錄申請；擴展 `players`/`agents` 表增加 `balance`/`pending_balance`。
    *   **API (Golang/Gin)**:
        *   `POST /admin/deposit`: 提交儲值申請。
        *   `POST /admin/deposit/:id/review`: 審核申請 (通過/拒絕)，使用事務更新用戶餘額。
    *   **前端 (Next.js)**: `/admin/deposits` 頁面顯示申請列表，提供審核操作按鈕。
    *   **優化**: Redis 快取待審核請求，限制儲值額度，整合第三方支付 API。

2.  **代理商/經銷商分潤配置連動**:
    *   **DB**: `profit_configs` 表儲存配置，`profit_config_logs` 記錄變更歷史。
    *   **API (Golang/Gin)**:
        *   `POST /admin/profit/config`: 更新分潤配置 (驗證權限 - 僅代理商/總公司)。
        *   `GET /admin/profit/configs`: 查詢分潤配置。
    *   **邏輯**: Golang 實現連動更新邏輯 (e.g., `cascadeUpdateProfit` 函數)。
    *   **前端 (Next.js)**: `/admin/profit-config` 頁面管理配置。
    *   **優化**: Redis 快取配置，JWT 驗證權限。

3.  **返水規則**:
    *   **規則設計**: 按遊戲類型、有效下注額計算，可配置 VIP/活躍加成，設定最低返水，防止刷單。
    *   **DB**: `rebates` 表記錄返水發放歷史。
    *   **後端 (Golang)**: 定時任務 (cron `0 0 * * 1`) 調用 `calculateRebates` 函數。
        *   查詢 `game_logs` 計算週期內有效下注。
        *   應用返水規則 (含加成)。
        *   使用事務創建 `rebates` 記錄並更新玩家餘額。
    *   **API (Golang/Gin)**: `GET /admin/rebates`: 查詢返水記錄。
    *   **前端 (Next.js)**: `/admin/rebates` 頁面顯示記錄。
    *   **優化**: Redis 快取計算中間結果，檢測異常下注模式。

**整體技術架構與安全性**:
*   **架構**: 前後端分離，DB/快取分離，延續之前規劃。
*   **安全**: JWT 保護後台 API，AES 加密敏感數據，GORM 防 SQL 注入。

**效能與擴展性**:
*   **高並發**: Goroutines 處理批量操作，Redis 分佈式鎖防衝突。
*   **DB 優化**: 相關表按時間分表，加強索引。
*   **監控**: Prometheus + Grafana 監控，`cron` 執行定時任務。

**與遊戲端整合**:
*   **儲值**: 遊戲端 WebSocket 發送請求 -> 後台處理。
*   **分潤/返水**: 遊戲端記錄下注 (`game_logs`) -> 後台定時計算。
*   **控管**: 遊戲端 API 讀取後台 `game_settings`。

**實現建議 (優先級)**:
1.  開分系統 (核心)。
2.  返水規則 (激勵)。
3.  分潤配置連動 (管理)。

**結論**: 新增的儲值、分潤連動和返水功能可無縫整合到現有 Golang 後台架構中，透過清晰的資料庫設計、API 接口和優化策略，能有效提升後台的管理能力和玩家體驗。 