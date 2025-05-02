## 文件 11: 後台規劃 3 (報表、賠率與機率控管)

**核心目標**: 新增報表功能和遊戲賠率/機率控管機制。

**新增功能需求**:
1.  **報表功能**: 生成玩家行為、財務、遊戲運營報表，支援篩選、匯出、視覺化。
2.  **遊戲賠率/機率控管**: 設定賠率/RTP，監控機率分佈，動態調整，確保公平與盈利。

**技術棧 (後台)**:
*   **後端**: Golang (Gin, GORM)。
*   **資料庫**: MySQL/PostgreSQL, Redis (快取), Elasticsearch (日誌)。
*   **前端 (後台管理)**: Next.js (整合 Chart.js)。
*   **部署**: Docker + Kubernetes。

**功能實現規劃**:

1.  **報表功能**:
    *   **DB**: 利用現有表，新增 `report_snapshots` (JSON 儲存預計算數據)。
    *   **API (Golang/Gin)**:
        *   `GET /admin/reports/...`: 生成各類報表 (玩家/財務/遊戲)，支援篩選。
        *   `GET /admin/reports/export`: 匯出報表 (CSV)。
    *   **後端邏輯**: 定時任務 (cron `@daily`) 生成報表快照 (`generateDailyReports`)。
    *   **前端 (Next.js)**: `/admin/reports` 頁面使用 Chart.js 展示圖表，渲染表格，提供匯出功能。
    *   **優化**: Redis 快取報表數據 (`getPlayerReport`)，分頁查詢大數據量。
    *   **安全**: JWT 限制報表存取權限 (`authReportAccess`)。

2.  **遊戲賠率、機率控管**:
    *   **設計**: 設定抽佣 (德州/梭哈) 或固定賠率 (百家樂)，目標 RTP (e.g., 90-95%)。
    *   **DB**: 新增 `game_odds` 表儲存賠率配置，擴展 `game_logs` 記錄當局賠率 (JSON)。
    *   **API (Golang/Gin)**:
        *   `POST /admin/odds`: 設置/更新賠率配置。
        *   (內部) 定時任務 (cron `@hourly`) 監控機率分佈 (`monitorProbability`)。
        *   (內部) 根據實際 RTP 動態調整賠率 (`adjustOdds`)。
        *   `GET /api/game/deck`: 提供加密隨機牌庫給遊戲端。
    *   **後端邏輯**: 使用 `crypto/rand` 實現公平洗牌 (`shuffleDeck`)。
    *   **前端 (Next.js)**: `/admin/odds` 頁面管理賠率設置。
    *   **優化**: Redis 快取賠率設置 (`getGameOdds`)，異常日誌記錄至 Elasticsearch (`logAnomaly`)。
    *   **安全**: JWT 限制賠率修改權限 (僅總公司 `authOddsControl`)，審計賠率變更日誌。

**整體技術架構與安全性**:
*   **架構**: 延續前後端分離，加入 Elasticsearch。
*   **安全**: JWT, AES 加密, GORM 防注入。

**效能與擴展性**:
*   **高並發**: Goroutines 批量處理，Redis 分佈式鎖。
*   **DB 優化**: 按月分表 (`game_logs`, `report_snapshots`)，加強索引。
*   **監控**: Prometheus + Grafana (API 延遲, 異常)。

**與遊戲端整合**:
*   **報表**: 遊戲端 WebSocket 傳送 `game_logs`。
*   **賠率**: 遊戲端 API 讀取 `game_odds`，記錄於 `game_logs`。
*   **機率**: 遊戲端 API 獲取隨機牌庫 (`/api/game/deck`)。

**實現建議 (優先級)**:
1.  報表功能 (基礎)。
2.  賠率設置與隨機性 (公平)。
3.  機率監控與動態調整 (控管)。

**結論**: 報表和賠率控管功能強化了後台的分析與風險管理能力，與 Golang/Next.js 技術棧及 Redis/Elasticsearch 整合良好，並契合了用戶在交易分析和概率計算方面的經驗。 