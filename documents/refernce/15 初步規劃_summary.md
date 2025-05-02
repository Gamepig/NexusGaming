## 文件 15: 初步規劃 (玩家行為分析與 AI 整合)

**基於文件 1-14 的討論，針對後台的玩家行為分析與 AI 整合功能（獨立於遊戲專案，使用 MySQL），制定初步規劃。**

**1. 項目目標**: 
*   開發高效的玩家行為分析系統。
*   使用 AI (Isolation Forest, K-Means) 減少計算負擔，提高異常識別準確性。
*   確保後台與獨立遊戲專案分離 (MySQL)。
*   提供視覺化報表 (Chart.js) 和匯出。
*   契合用戶興趣 (AI 優化, 概率計算, 交易分析, Docker)。

**2. 功能範圍**: 
*   **行為分析**: 下注/勝率/頻率/時間分佈指標，異常檢測 (高頻/高勝率/連贏)，玩家分類 (休閒/VIP/高風險)。
*   **AI 功能**: 異常檢測 (Isolation Forest), 玩家分類 (K-Means), 計算優化。
*   **報表**: 篩選，視覺化 (柱狀/散點/熱力圖)，匯出 (CSV)。
*   **MySQL 整合**: AI 提取 `game_logs`，結果存入 `ai_anomaly_predictions`。
*   **遊戲分離**: 後台通過 API 與遊戲交互。

**3. 技術架構**: 
*   **後端**: Golang (Gin, GORM)。
*   **AI**: Python (FastAPI, scikit-learn)。
*   **數據庫**: MySQL, Redis (快取)。
*   **前端**: Next.js, Chart.js。
*   **部署**: Docker, Kubernetes。

**4. 資料庫設計 (MySQL)**: 
*   **核心表**: `players`, `game_logs`, `point_transactions`。
*   **分析表**: `player_behavior_snapshots`, `ai_anomaly_predictions` (含 `anomaly_score`)。
*   **索引**: 優化 `player_id`, `game_type`, `created_at`。

**5. AI 模型設計**: 
*   **異常檢測**: Isolation Forest (特徵: 下注額/次數/勝率/連勝/頻率/時間)。
*   **玩家分類**: K-Means (特徵: 下注額/活躍天數/時長/盈虧/偏好)。
*   **部署**: Python FastAPI 提供 `/predict_anomaly` 和 `/predict_cluster` 端點。

**6. 後端實現 (Golang)**: 
*   `extractPlayerFeatures` 從 MySQL 提取特徵。
*   `callAIPrediction` 調用 Python AI 服務。
*   API `/admin/ai/detect_anomalies` 儲存預測結果。
*   Cron 任務 `runAIDetection` 每日執行。

**7. 前端實現 (Next.js)**: 
*   `pages/admin/player-behavior.tsx` 展示行為報表 (Bar/Line 圖) 和 AI 異常檢測 (Scatter 圖)。
*   表格結合基礎指標和 AI 預測結果。

**8. 與遊戲專案分離**: 
*   API `/game/logs` 接收日誌。
*   API `/game/player_behavior` 提供分析結果。

**9. 開發步驟與時間估計**: 
*   8 個階段，總計約 **17 天** (DB設計 -> 特徵提取 -> AI訓練 -> AI部署 -> 後端整合 -> 前端 -> 測試 -> 部署)。

**10. 風險與緩解**: 
*   MySQL 性能 -> 分表/索引/快取。
*   AI 準確性 -> 足夠標記數據，定期重訓。
*   API 不相容 -> 提前定義規範，模擬測試。
*   AI 服務延遲 -> 輕量模型，K8s 擴展。

**11. 效能與擴展性**: 
*   計算: AI (O(n)) + 快照 + 批量。
*   MySQL: 分表 + 索引 + 批量插入。
*   部署: Docker/K8s。
*   監控: Prometheus。

**12. 記憶整合**: 契合 AI 優化、概率計算、交易分析、Docker 經驗。

**13. 結論**: 初步規劃可行，整合了 AI 優化、MySQL 和獨立專案需求，技術棧成熟，時間估計合理，風險可控。 