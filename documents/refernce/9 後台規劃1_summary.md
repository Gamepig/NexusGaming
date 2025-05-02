## 文件 9: 後台規劃 1 (功能需求與初步實現)

**核心目標**: 為手機 Web 撲克遊戲 (Three.js+Next.js 前端, Golang 後端) 規劃後台系統。

**後台功能需求**: 
1.  **玩家點數統計**: 追蹤餘額、歷史交易 (充值/消費/輸贏)，查詢，排行榜，異常防止。
2.  **代理商/經銷商分潤與損失**: 管理層級，計算分潤/損失，帳目報表，結算。
3.  **遊戲控管**: 監控運行狀態 (玩家數/頻率/異常)，設置遊戲參數 (下注/桌數/AI 難度)，風險控制 (作弊檢測)，日誌/報表。

**技術棧 (後台)**:
*   **後端**: Golang (Gin 框架, GORM ORM)。
*   **資料庫**: MySQL/PostgreSQL (交易/玩家數據), Redis (快取/即時狀態)。
*   **前端 (後台管理)**: Next.js。
*   **部署**: Docker + Kubernetes。

**功能實現規劃**: 

1.  **玩家點數統計**: 
    *   DB: `players` (餘額), `point_transactions` (歷史)。
    *   API (Golang/Gin): `/admin/players/:id/points` (查餘額), `/admin/players/:id/transactions` (查歷史), `/admin/players/:id/update_points` (更新點數 - 使用事務)。
    *   前端 (Next.js): 顯示點數/歷史，排行榜。
    *   優化: DB 索引 (`player_id`, `created_at`), Redis 快取排行榜。

2.  **代理商/經銷商分潤與損失**: 
    *   DB: `agents` (層級, `parent_id`), `agent_transactions` (分潤/損失記錄)。
    *   邏輯: 根據玩家消費計算分潤比例 (e.g., 代理 5%, 經銷 3%)。
    *   API (Golang/Gin): `calculateProfit` (計算單次分潤), `/admin/agents/:id/report` (查帳目)。
    *   前端 (Next.js): 顯示帳目報表。
    *   優化: Redis 快取報表, goroutines 批量結算分潤。

3.  **遊戲控管**: 
    *   DB: `game_logs` (玩家行為), `game_settings` (參數)。
    *   API (Golang/Gin): `/admin/games/status` (監控狀態), `/admin/games/settings` (更新設置)。
    *   風險控制: Golang 定時任務 (cron) 檢測異常玩家 (e.g., 高頻/高額下注)。
    *   前端 (Next.js): 顯示狀態，修改設置。
    *   優化: Redis 存儲即時狀態, 定時任務檢測風險, Elasticsearch/Logrotate 處理日誌。

**技術架構與安全性**: 
*   **架構**: 前後端分離，DB/快取分離。
*   **安全**: JWT 保護後台 API, AES 加密敏感數據, GORM 防 SQL 注入。

**效能與擴展性**: 
*   **高並發**: Goroutines 處理請求, Redis 分佈式鎖 (防點數更新衝突)。
*   **DB 優化**: 交易表按時間分表, 加索引。
*   **監控**: Prometheus + Grafana 監控後台性能。

**與遊戲端整合**: 
*   **點數**: WebSocket 同步遊戲結果 -> 後台更新 `point_transactions`。
*   **控管**: 遊戲端 API 獲取 `game_settings`。

**實現建議**: 
*   **優先級**: 點數統計 > 遊戲控管 > 代理商分潤。
*   **步驟**: DB -> Golang API -> Next.js 後台 -> 分潤邏輯 -> 部署與監控。
*   **測試**: 高並發點數更新, 分潤準確性, 手機 UI。

**結論**: Golang 後端配合 MySQL/Redis 和 Next.js 前端，可構建功能完善、高效能、可擴展的遊戲後台，滿足玩家點數、代理分潤和遊戲控管需求。 