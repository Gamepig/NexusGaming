## 文件 2: 手機 Web 卡牌遊戲語言選擇與性能優化

**核心需求**: 開發一款可在手機上以 Web 形式運行的卡牌遊戲，包含 AI 對弈，需考慮語言選擇和性能優化 (結合 Cython 和 AI 加速經驗)。

**推薦技術棧**: 
*   **前端**: TypeScript + PixiJS (2D 渲染) 或 Phaser (遊戲框架)。
    *   理由: 高效 WebGL 渲染，手機兼容，觸控支持。
    *   替代: React + Three.js (若需 3D)。
*   **後端**: Python + FastAPI。
    *   理由: 熟悉 Python，利於整合 Cython 加速 AI，FastAPI 支持異步和 WebSocket。
    *   替代: Node.js (統一語言，但不利於 Cython)。
*   **AI 計算**: Python + Cython (後端)，或編譯為 WebAssembly (WASM) 在前端/後端運行。
*   **通信**: WebSocket (實時對戰) + REST API。
*   **部署**: 後端 (雲端)，前端 (CDN)。

**性能優化策略**: (針對手機環境)
*   **前端渲染優化**: 
    *   使用 WebGL (PixiJS)。
    *   合併貼圖 (Sprite Sheet)。
    *   分層渲染。
    *   觸控事件優化 (`pointerdown`)。
    *   限制幀率 (30 FPS)。
*   **AI 計算優化**: 
    *   **Cython 加速**: 編譯後端 AI 核心邏輯 (MCTS)。
    *   **WebAssembly**: 若 AI 需前端運行，將 Cython 轉 WASM (Emscripten/Pyodide)。
    *   **快取/記憶化**: 儲存重複狀態評估。
    *   **異步計算**: FastAPI 非同步路由避免阻塞。
*   **網路通信優化**: 
    *   WebSocket 減少延遲。
    *   壓縮數據 (JSON/Protobuf)。
    *   前端執行簡單 AI 邏輯 (離線計算)。
*   **手機硬體適配**: 
    *   響應式設計 (螢幕)。
    *   電池優化 (低頻更新, `requestAnimationFrame`)。
    *   記憶體管理 (清除物件)。
*   **進階 AI 優化**: 
    *   **輕量模型**: 神經網路 AI 模型量化 (8-bit ONNX)。
    *   **剪枝**: 移除不重要層。
    *   **GPU 加速**: 若後端有 GPU。

**性能預期**: 
*   Cython 加速 5-50x。
*   WASM 接近原生 C (JS 的 2-10x)。
*   WebGL 渲染可達 60 FPS。
*   網路延遲 100-200ms。

**實現建議**: 
1.  MVP: Phaser/PixiJS + FastAPI，簡單 AI (Python)。
2.  優化: Cython/WASM 加速 AI，PixiJS 優化渲染。
3.  多人: WebSocket。
4.  測試部署: 真機測試，CDN + 雲端。 