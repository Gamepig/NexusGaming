## 文件 16: 初步規劃 (含 AI 對弈與提高莊家勝率)

**基於文件 1-15 的討論，針對後台的玩家行為分析、AI 判斷行為以及新增的 AI 對弈（旨在提高莊家勝率）功能（獨立於遊戲專案，使用 MySQL），制定初步規劃。**

**1. 項目目標**: 
*   **行為分析**: 高效分析玩家行為指標。
*   **AI 判斷**: 使用 AI (Isolation Forest, K-Means) 優化異常檢測和玩家分類。
*   **AI 對弈**: 開發 AI 對弈模組 (MCTS/DQN, 馬爾可夫鏈) 參與遊戲，動態調整策略**提高莊家勝率**。
*   **MySQL 整合**: AI 提取/儲存行為分析和對弈結果。
*   **遊戲分離**: 後台通過 API 與遊戲交互。
*   **報表**: 提供行為分析和對弈結果報表。
*   契合用戶興趣 (AI 優化, 概率計算, 交易分析, Docker)。

**2. 功能範圍**: 
*   **行為分析**: 指標、異常檢測、玩家分類。
*   **AI 判斷**: Isolation Forest (異常), K-Means (分類)。
*   **AI 對弈**: 
    *   德州/梭哈: MCTS/DQN 根據玩家行為調整策略。
    *   百家樂: 馬爾可夫鏈/統計模型調整內部策略或預測。
    *   **莊家勝率提升**: 動態策略 + 風險控制。
*   **報表**: 行為分析 + AI 對弈結果，視覺化，匯出 (CSV)。
*   **MySQL**: `game_logs`, `snapshots`, `ai_anomaly_predictions`, 新增 `ai_game_results`。
*   **遊戲分離**: API (`/game/logs`, `/game/play`) + 可選 WebSocket。

**3. 技術架構**: 
*   **後端**: Golang (Gin, GORM)。
*   **AI**: Python (FastAPI, scikit-learn, PyTorch)。
*   **數據庫**: MySQL, Redis。
*   **前端**: Next.js, Chart.js。
*   **部署**: Docker, Kubernetes。

**4. 資料庫設計 (MySQL)**: 
*   核心表 + 分析表 + **`ai_game_results`** (記錄 AI 對弈行動/結果/參數)。
*   增加 `ai_game_results` 索引。

**5. AI 模型設計**: 
*   **行為分析**: Isolation Forest, K-Means (訓練/部署同文件 15)。
*   **AI 對弈**: 
    *   德州/梭哈: MCTS + DQN (輸入: 牌局狀態, 玩家歷史/特徵; 輸出: 行動)。
    *   百家樂: 馬爾可夫鏈 + 統計模型 (輸入: 歷史結果, 玩家偏好; 輸出: 策略調整/預測)。
    *   **莊家勝率提升**: 動態策略 + 風險控制。
    *   **部署**: Python FastAPI 新增 `/play_game` 端點。

**6. 後端實現 (Golang)**: 
*   `extractPlayerFeatures`。
*   `callAI` 調用 Python (`/predict_anomaly`, `/predict_cluster`, `/play_game`)。
*   API `/admin/ai/detect_anomalies` 儲存行為分析結果。
*   API `/game/play` 獲取玩家特徵，調用 AI 對弈服務獲取行動，記錄 `ai_game_results` (需遊戲專案提供最終結果)。
*   Cron 任務 `runAIDetection` 每日執行行為分析。

**7. 前端實現 (Next.js)**: 
*   `pages/admin/player-behavior.tsx` 增加 AI 對弈結果表格。

**8. 與遊戲專案分離**: 
*   API `/game/logs` 接收日誌。
*   API `/game/play` 提供 AI 行動。

**9. 開發步驟與時間估計**: 
*   8 個階段，增加 AI 對弈模型訓練/部署，總計約 **22 天**。

**10. 風險與緩解**: 
*   新增風險: AI 對弈勝率不足 -> 增加模擬數據，微調模型，定期更新。
*   其他風險同文件 15。

**11. 效能與擴展性**: 
*   AI 對弈性能目標: MCTS < 100ms/行動。
*   其他同文件 15。

**12. 記憶整合**: 契合 AI 優化、概率計算、交易分析、自訂 AI、Docker 經驗。

**13. 結論**: 規劃整合了 AI 對弈以提高莊家勝率的需求，調整了模型、API、資料庫和開發估計，技術可行，風險可控。 