## 文件 7: 純 Golang 後端，專注德州撲克、梭哈、百家樂

**核心需求**: 使用 **Golang** 作為唯一後端語言 (去除 Python/Cython)，專注開發**德州撲克 (Texas Hold'em)**、**傳統梭哈 (Seven-Card Stud)** 和 **百家樂 (Baccarat)** 三款撲克遊戲。前端維持 **Three.js + Next.js**，支援手機 Web 全螢幕。

**三款遊戲 AI 適用性分析 (純 Golang)**:

1.  **德州撲克**: 
    *   AI 算法: **MCTS** (最適合處理不完美資訊和長期收益)。
    *   Golang 實現: 使用 goroutines 並行化 MCTS 模擬，內存快取 (sync.Map) 存儲狀態評估。
    *   挑戰: 高效模擬詐唬、對手建模。
    *   策略深度: 高。

2.  **傳統梭哈**: 
    *   AI 算法: **Minimax + Alpha-Beta 剪枝** (適合分析明牌)，MCTS (更複雜場景)。
    *   Golang 實現: 高效遞迴 Minimax，利用明牌信息。
    *   挑戰: 根據明牌推測底牌。
    *   策略深度: 中高。

3.  **百家樂**: 
    *   AI 算法: **概率計算/期望值分析** (無需複雜搜索)。
    *   Golang 實現: 基於歷史數據或固定概率計算最佳下注 ("閒"/"莊"/"和")。
    *   挑戰: 模擬公平發牌。
    *   策略深度: 低。

**技術實現與優化 (純 Golang 後端)**:

*   **前端 (Three.js + Next.js)**:
    *   全螢幕: CSS (`100vw/vh`) + Fullscreen API。
    *   渲染: 低多邊形卡牌 (PlaneGeometry), 壓縮紋理 (WebP), 限制 30 FPS。
    *   交互: `react-three-fiber` 處理觸控 (下注、選牌)。
    *   手機適配: `viewport` meta, PWA meta (`apple-mobile-web-app-capable`)。

*   **後端 (Golang)**:
    *   **遊戲邏輯**: 為三款遊戲分別實現規則引擎。
    *   **AI 實現**: 
        *   德州撲克: MCTS (goroutines 並行模擬, sync.Map 快取)。
        *   梭哈: Minimax + Alpha-Beta 剪枝。
        *   百家樂: 概率查表或簡單計算。
    *   **API/WebSocket**: Gin 提供 REST API (玩家行動), `gorilla/websocket` 提供實時對戰通道。
    *   **性能**: Goroutines 加速計算密集型 AI (MCTS), 內存快取減少重複計算。

*   **手機優化**: 
    *   電池: 低頻渲染。
    *   網路: WebSocket 重連機制。

**遊戲選擇與實現建議**: 
*   **優先級**: 德州撲克 (最流行，AI 挑戰大) > 傳統梭哈 (經典，策略適中) > 百家樂 (簡單快速)。
*   **實現路徑**: 
    1.  實現三款遊戲的基礎規則引擎 (Golang)。
    2.  開發對應的 Golang AI 模組 (MCTS, Minimax, 概率)。
    3.  整合前端 Three.js 渲染和交互。
    4.  透過 WebSocket 連接前後端。
    5.  優化性能和手機體驗。

**結論**: 純 Golang 後端完全可行，能高效實現德州撲克、梭哈、百家樂的 AI 對弈。MCTS (德州)、Minimax (梭哈)、概率計算 (百家樂) 是推薦的 Golang AI 實現方式。需重點優化 Golang AI 計算效率 (並行、快取) 和前端渲染性能。 