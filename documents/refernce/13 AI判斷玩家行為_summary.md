## 文件 13: AI 判斷玩家行為

**核心目標**: 探討使用 AI 分析玩家行為以減少計算負擔，並設計 AI 與資料庫的整合方案。

**AI 優勢**:
*   **減少計算**: 模型推理 (O(n)) 比複雜 SQL 聚合 (O(n*log(n)+)) 更高效。
*   **提高準確性**: 識別隱性模式，檢測複雜異常。
*   **動態適應**: 模型可持續學習。
*   **應用**: 異常檢測、玩家分類、行為預測。

**AI 與資料庫整合方案**:

1.  **數據來源與特徵提取**:
    *   **來源**: `game_logs`, `point_transactions`, `player_behavior_snapshots`。
    *   **特徵**: 基本 (下注/勝率/活躍度)、進階 (頻率/時間分佈/序列)、異常 (波動/連勝敗)。
    *   **實現**: Golang `extractPlayerFeatures` 從 DB 提取數據。
2.  **AI 模型設計**:
    *   **選擇**: 異常檢測 (Isolation Forest), 分類 (K-Means/Random Forest), 預測 (LSTM)。
    *   **訓練**: 使用歷史數據，標記異常樣本。
    *   **部署**: Python (scikit-learn/TF/PyTorch) 開發，保存模型 (.pkl/ONNX)。
3.  **整合架構**:
    *   **數據流**: Golang (GORM) -> DB (MySQL/PG) -> Golang (提取特徵 API) -> Python (FastAPI 推理服務) -> Golang (調用推理) -> DB (儲存結果)。
    *   **通信**: Golang 與 Python 服務間使用 REST API 或 gRPC。
4.  **結果儲存**: 新增 `ai_anomaly_predictions` 表儲存預測結果 (is_anomaly, features JSON)。
5.  **減少計算**: 
    *   AI 推理取代複雜 SQL。
    *   批量推理 + 增量更新 (`incrementalAIPrediction`)。
    *   定時任務 (cron `@daily`) 執行。

**前端與報表整合**:
*   **展示**: Next.js 頁面 (`pages/admin/ai-anomalies.tsx`) 顯示異常列表。
*   **視覺化**: Chart.js 散點圖 (`AnomalyScatter`) 展示異常玩家特徵分佈。

**效能與擴展性**:
*   AI 推理速度優勢。
*   利用快照和 Redis 快取 (`cachePlayerFeatures`) 減輕 DB 負載。
*   獨立部署 Python AI 服務 (Docker/K8s)。
*   Prometheus 監控推理時間。

**記憶整合**: 
*   契合交易分析 (模式識別)、概率計算 (風險評估)、數據處理 (API)、部署 (Docker/K8s) 的經驗。

**實現建議 (優先級)**:
1.  數據提取 API。
2.  異常檢測模型 (Isolation Forest)。
3.  玩家分類模型 (K-Means)。
4.  前端整合。

**結論**: AI 能有效減少玩家行為分析的計算負擔，提高準確性。透過 Golang 和 Python 服務的整合，結合現有資料庫和技術棧，可以構建高效、可擴展的 AI 分析系統。 