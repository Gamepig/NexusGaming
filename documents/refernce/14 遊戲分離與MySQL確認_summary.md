## 文件 14: 遊戲分離與 MySQL 確認

**核心確認**:
1.  遊戲將是**獨立專案**。
2.  後台資料庫確認使用 **MySQL**。

**後台規劃調整 (聚焦玩家行為分析與 AI 整合)**:

1.  **MySQL 資料庫設計**:
    *   **核心表**: `players`, `game_logs`, `point_transactions` (與遊戲專案共用，MySQL 語法)。
    *   **分析表**: `player_behavior_snapshots`, `ai_anomaly_predictions` (含 `anomaly_score`)。
    *   **索引**: 針對 `player_id`, `game_type`, `created_at` 優化查詢。
2.  **AI 模型設計**:
    *   **異常檢測**: Isolation Forest。
    *   **玩家分類**: K-Means。
    *   **部署**: Python (FastAPI + scikit-learn) + `.pkl` 模型。
3.  **後端 (Golang)**:
    *   **特徵提取**: `extractPlayerFeatures` (從 MySQL 提取，含 `ConsecutiveWins`)。
    *   **AI 調用**: `callAIPrediction` (調用 Python FastAPI 服務，處理含 score 的返回)。
    *   **結果儲存**: 將 `is_anomaly` 和 `anomaly_score` 存入 `ai_anomaly_predictions`。
    *   **定時任務 (cron `@daily`)**: `runAIDetection` (提取、預測、儲存)。
4.  **Python AI 服務**: 
    *   FastAPI `/predict_anomaly` 接口，加載模型，返回 `is_anomaly` 和 `score`。
    *   提供離線訓練腳本。
5.  **前端 (Next.js)**:
    *   `pages/admin/ai-anomalies.tsx` 展示異常列表和分數，使用 Chart.js 散點圖視覺化。
6.  **與遊戲專案分離**: 
    *   **交互**: 後台提供 API (`/game/logs` 接收日誌, `/game/player_behavior` 提供分析結果) 或 WebSocket 與遊戲專案通信。
7.  **效能與擴展性**:
    *   **計算**: AI 推理 + 批量處理 + 快照機制。
    *   **MySQL**: 分表 (`game_logs`) + 索引 + 批量插入 (`CreateInBatches`)。
    *   **快取**: Redis (`cachePlayerFeatures`)。
    *   **部署**: Docker (提供 Python 服務 Dockerfile) + Kubernetes。
    *   **監控**: Prometheus。

**記憶整合**: 契合 AI 優化、概率計算、交易分析、部署經驗。

**實現建議**: 優先實現特徵提取 -> 異常檢測 -> 玩家分類 -> 前端整合。

**結論**: 確認遊戲分離和 MySQL 後，AI 驅動的玩家行為分析方案依然可行，透過清晰的 API 設計、MySQL 優化和獨立部署，可構建高效、解耦的後台系統。 