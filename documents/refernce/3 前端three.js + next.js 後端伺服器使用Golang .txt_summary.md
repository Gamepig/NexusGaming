## 文件 3: 手機 Web 卡牌遊戲 - 前端 Three.js+Next.js, 後端 Golang

**核心需求**: 使用特定技術棧 (前端 Three.js + Next.js, 後端 Golang) 開發手機 Web 卡牌遊戲，包含 AI 對弈，並整合 Python/Cython 遺留 AI 程式碼。

**技術棧分析**: 
*   **前端 (Three.js + Next.js)**:
    *   Three.js: WebGL 3D 渲染庫，適合高級視覺效果 (3D 卡牌、場景)。可能對簡單卡牌略重 (對比 PixiJS)。需優化手機性能。
    *   Next.js: React 框架，支持 SSR/SSG，結構化開發，易部署。
    *   語言: TypeScript (推薦)。
*   **後端 (Golang)**:
    *   Golang: 高效能、高並發，適合 AI 計算、遊戲狀態、WebSocket。
    *   框架: Gin 或 Fiber (推薦)。
*   **AI 邏輯語言**: 
    *   **原生 Go (推薦)**: 最佳性能，利用 goroutines。
    *   Python + Cython (可選): 
        *   **gRPC**: Python 作為獨立服務，Go 呼叫。
        *   **WebAssembly (WASM)**: Cython 編譯為 WASM，Go (wasmtime-go) 或前端執行。
        *   外部進程 (不推薦)。
    *   C/C++ (可選): 極高性能需求，透過 cgo 調用。

**性能優化建議**: (針對手機環境)
*   **前端優化 (Three.js + Next.js)**:
    *   **Three.js 渲染**: 簡化幾何體 (PlaneGeometry)，紋理壓縮 (WebP)，合併渲染 (Sprite 集)，LOD。
    *   **幀率控制**: 限制 30 FPS。
    *   **Next.js 優化**: SSG 預渲染資源，動態載入 Three.js 模組，記憶體管理 (dispose)。
    *   **觸控優化**: `react-three-fiber` 簡化事件。
*   **後端優化 (Golang)**:
    *   **AI 計算**: 
        *   原生 Go (MCTS/Minimax) + goroutines 並行化。
        *   Python/Cython 整合: gRPC (推薦) 或 WASM (wasmtime-go)。
    *   **快取**: 內存快取 (sync.Map) 或 Redis 存儲 AI 結果。
    *   **WebSocket**: `gorilla/websocket` 低延遲通信。
*   **手機環境優化**: 
    *   響應式設計 (Next.js + Three.js 畫布調整)。
    *   觸控支援 (`react-three-fiber`)。
    *   電池優化 (Three.js onDemand 渲染)。
    *   網路優化 (CDN, WebSocket 壓縮)。
*   **進階 AI 優化**: 
    *   Go 實現神經網路 (`gorgonia`)。
    *   模型量化/剪枝 (Python 端處理，導出 ONNX，Go 用 `onnx-go` 推理)。

**性能預期**: 
*   Three.js: 30-60 FPS (中階手機)。
*   Golang 後端: API <50ms, WebSocket <100ms。
*   AI: Go MCTS (數千模擬/秒), Cython+WASM (接近 C), gRPC (+10-20ms)。

**實現建議**: 
1.  MVP: `react-three-fiber` + Gin，簡單 AI (Go 啟發式)。
2.  AI 開發: 優先 Go 原生 MCTS，若需整合 Python 則用 gRPC。
3.  測試: 真機，Chrome DevTools。
4.  部署: 前端 (Vercel), 後端 Go (Cloud Run/ECS), Python 服務 (獨立部署)。

**結論**: 該技術棧可行，但需重點優化 Three.js 手機性能和 Golang AI 計算。原生 Go 是 AI 性能首選，gRPC 是整合 Python 的較優方案。 