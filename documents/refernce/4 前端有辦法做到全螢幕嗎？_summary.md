## 文件 4: 前端實現全螢幕顯示

**核心問題**: 前端 (特別是 Three.js + Next.js 技術棧) 如何實現全螢幕顯示，尤其是在手機上？

**可行性**: 可以。透過 CSS 和 Fullscreen API。

**實現方式**: 
1.  **CSS 視窗全螢幕**: 
    *   方法: 設定畫布容器寬高為 `100vw`, `100vh`，移除 body 邊距和滾動條。
    *   效果: 畫布填滿視窗，但瀏覽器 UI (地址欄、狀態列) 仍在。
    *   手機優化: 添加 `viewport` meta 標籤控制縮放。
2.  **Fullscreen API 真全螢幕**: 
    *   方法: 使用 `element.requestFullscreen()` (需用戶觸發，如按鈕點擊)，監聽 `fullscreenchange` 事件。
    *   效果: 隱藏瀏覽器 UI，提供沉浸體驗。
    *   手機優化: 處理瀏覽器前綴 (`webkitRequestFullscreen`)，可選用 Screen Orientation API 鎖定方向。
3.  **響應式設計**: 
    *   方法: `resize` 事件監聽器動態調整 Three.js 渲染器和相機。
    *   手機優化: CSS `env(safe-area-inset-*)` 避開「劉海」等區域，`aspect-ratio` 固定比例。

**與 Three.js + Next.js 整合**: 
*   在 Next.js component (如 `GameCanvas.tsx`) 中管理 Three.js 實例和全螢幕邏輯。
*   使用 `useEffect` 處理事件監聽器和清理。
*   全局 CSS (`globals.css`) 和 `_document.tsx` (meta 標籤) 配合。

**對性能的影響與優化**: 
*   **Three.js 渲染**: 全螢幕可能增加 GPU 負載。
    *   優化: 簡化場景/模型 (PlaneGeometry), 壓縮紋理 (WebP), 限制幀率 (30 FPS), 全螢幕時使用 `screen.width/height` 更新渲染器。
*   **AI 計算 (後端 Golang)**: 全螢幕不直接影響後端，但需確保 AI 回應速度。
    *   優化: 高效 Go AI (MCTS), gRPC 整合 Python/Cython (快取結果), WebSocket 低延遲通信。
*   **手機效能**: 
    *   電池: 低功耗渲染 (onDemand `useFrame`)。
    *   網路: WebSocket 重連。
    *   觸控: `react-three-fiber` 保證精準。

**手機全螢幕挑戰與解決**: 
*   **iOS Safari**: 可能保留底部欄 -> 添加 PWA meta 標籤 (`apple-mobile-web-app-capable`)。
*   **Android Chrome**: 可能不隱藏狀態列 -> 依賴 Fullscreen API。
*   **方向變化**: `resize` 事件處理。
*   **測試**: 在中低階手機測試性能 (30 FPS)。

**結論**: 全螢幕在 Three.js + Next.js 中是可行的，需結合 CSS、Fullscreen API 和響應式設計，並針對手機渲染、AI 交互和硬體限制進行優化。 