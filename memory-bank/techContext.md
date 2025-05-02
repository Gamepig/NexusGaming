# 技術背景

## 後端服務

-   **語言**: Golang (v1.18+)
-   **Web 框架**: Gin
-   **ORM**: GORM
-   **主要職責**: API 接口、業務邏輯、數據處理協調、與 AI 服務/DB 交互、後台管理功能。

## AI 服務

-   **語言**: Python (v3.9+)
-   **API 框架**: FastAPI
-   **機器學習庫**:
    -   Scikit-learn: 用於行為分析模型 (Isolation Forest, K-Means)。
    -   PyTorch: 用於對弈模型 (DQN)。
-   **主要職責**: 提供 AI 模型訓練和推理的 API 端點。

## 前端管理介面

-   **框架**: Next.js (v13+)
-   **UI 庫**: React (v18+)
-   **圖表庫**: Chart.js (v4+)
-   **主要職責**: 提供給管理員使用的視覺化報表和管理工具。

## 數據存儲

-   **關係型數據庫**: MySQL (v8.0+) - 持久化存儲核心業務數據、分析結果、對弈記錄、管理信息。
-   **內存數據庫/快取**: Redis (v6.0+) - 緩存常用數據、特徵、會話狀態，加速讀取。

## 部署與運維

-   **容器化**: Docker
-   **容器編排**: Kubernetes
-   **監控**: Prometheus, Grafana (規劃中)

## 身份驗證

-   **機制**: JSON Web Tokens (JWT)

## 項目路徑

-   `/Users/vichuang/projects/NexusGaming` 