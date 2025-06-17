# Task: 後台管理介面與核心系統

**目標**: 開發後台所需的核心 API 服務、數據庫結構以及供內部使用的管理介面 (Vanilla JS)。

## 主要步驟

1.  **基礎架構與環境設定**:
    *   **子任務 1.1**: 初始化 Golang 後端專案結構 (Gin 框架, 目錄結構: `cmd`, `internal`, `pkg`, `api`, `configs` 等)。
    *   **子任務 1.2**: 設計 MySQL 資料庫 Schema (會員、帳務、代理、遊戲配置、角色權限等核心表)。
    *   **子任務 1.3**: 設計 MongoDB 集合結構 (操作日誌、遊戲日誌、錯誤日誌)。
    *   **子任務 1.4**: 實現 MySQL 資料庫遷移 (使用 GORM AutoMigrate 或 go-migrate)。
    *   **子任務 1.5**: 初始化 MongoDB 集合 (可選：設定索引)。
    *   **子任務 1.6**: 配置 Golang 資料庫連接 (讀取配置, 實例化 GORM 和 MongoDB Driver)。
    *   **子任務 1.7**: 設定 Golang 基礎 API 路由 (Gin 根路由, 健康檢查端點)。
    *   **子任務 1.8**: 初始化 Vanilla JS 前端專案結構 (Vite/Webpack, 目錄結構: `src`, `public`, `assets`)。
    *   **子任務 1.9**: 設計並實現前端基礎佈局 (HTML/CSS: 頂部導航、側邊欄、主內容區)。

2.  **核心認證與權限系統 (RBAC)**:
    *   **子任務 2.1: 後端: 設計並實現用戶數據模型**
        *   2.1.1: 資料庫: 定義 `users` 表格結構 (欄位: id, username, email, phone, password_hash, status, role_id, created_at, updated_at)，編寫 SQL DDL 或遷移文件。
        *   2.1.2: 後端: 創建 Go struct (`internal/models/user.go`)，使用 GORM 標籤映射 `users` 表。
        *   2.1.3: 後端: 實現密碼加密邏輯 (例如使用 `golang.org/x/crypto/bcrypt`)。
        *   2.1.4: 後端: 創建用戶數據倉庫 (`internal/repository/user_repo.go`)，實現基本 CRUD 操作（`CreateUser`, `GetUserByUsername`, `GetUserByID` 等），包含密碼加密處理。
    *   **子任務 2.2: 後端: 設計並實現角色與權限模型**
        *   2.2.1: 資料庫: 定義 `roles` 表格結構 (欄位: id, name, description, timestamps)。
        *   2.2.2: 資料庫: 定義 `permissions` 表格結構 (欄位: id, action, resource, description, timestamps)。
        *   2.2.3: 資料庫: 定義 `role_permissions` 關聯表格結構 (欄位: role_id, permission_id)。
        *   2.2.4: 後端: 創建對應的 GORM 模型 (`internal/models/role.go`, `internal/models/permission.go`)。
        *   2.2.5: 後端: 創建角色權限數據倉庫 (`internal/repository/rbac_repo.go`)，實現檢查權限的函數 (例如 `HasPermission(roleID, action, resource)`)。
        *   2.2.6: 後端/腳本: 編寫腳本或在初始化代碼中，填充初始的角色和權限數據到數據庫。
    *   **子任務 2.3: 後端: 實現登入認證 API**
        *   2.3.1: 後端: 定義登入請求/響應的 JSON 結構體 (`internal/api/types/auth.go`)。
        *   2.3.2: 後端: 在 Gin 路由中註冊 `/api/auth/login` POST 端點及處理函數 (`internal/api/handlers/auth.go`)。
        *   2.3.3: 後端: 在處理函數中調用 `user_repo` 獲取用戶。
        *   2.3.4: 後端: 調用 `bcrypt.CompareHashAndPassword` 驗證密碼。
        *   2.3.5: 後端: 調用 JWT 生成邏輯，將用戶 ID 和角色 ID 放入 claims。
        *   2.3.6: 後端: 返回成功響應 (含 JWT) 或錯誤信息。
    *   **子任務 2.4: 後端: 實現 JWT 驗證中間件**
        *   2.4.1: 後端: 創建 Gin 中間件函數 (`internal/api/middleware/auth.go`)。
        *   2.4.2: 後端: 提取 `Authorization: Bearer` token。
        *   2.4.3: 後端: 驗證 token 簽名和有效期。
        *   2.4.4: 後端: 解析 claims，獲取用戶 ID 和角色 ID。
        *   2.4.5: 後端: 將用戶信息存儲到 Gin 請求上下文 (`c.Set`)。
        *   2.4.6: 後端: 若 token 無效，返回 401 Unauthorized。
    *   **子任務 2.5: 後端: 實現權限檢查中間件/邏輯**
        *   2.5.1: 後端: 創建 Gin 中間件函數 (`internal/api/middleware/rbac.go`)，接受所需權限參數。
        *   2.5.2: 後端: 從 Gin 上下文中讀取用戶 `roleID`。
        *   2.5.3: 後端: 調用 `rbac_repo.HasPermission` 檢查權限。
        *   2.5.4: 後端: 若無權限，返回 403 Forbidden。
        *   2.5.5: 後端: 將此中間件應用於需要保護的 API 路由。
    *   **子任務 2.6: 前端(Admin): 開發登入頁面**
        *   2.6.1: 前端(Admin): 創建 `login.html` (包含 form, input, button)。
        *   2.6.2: 前端(Admin): 創建 `login.css` 添加樣式。
    *   **子任務 2.7: 前端(Admin): 實現登入 API 調用與 JWT 存儲**
        *   2.7.1: 前端(Admin): 創建 `login.js` 並引入。
        *   2.7.2: 前端(Admin): 添加表單 `submit` 事件監聽器。
        *   2.7.3: 前端(Admin): 阻止默認提交。
        *   2.7.4: 前端(Admin): 獲取輸入值。
        *   2.7.5: 前端(Admin): 使用 `fetch` 調用 `/api/auth/login`。
        *   2.7.6: 前端(Admin): 處理響應 (成功/失敗)。
        *   2.7.7: 前端(Admin): 成功後將 token 存儲到 `localStorage`。
        *   2.7.8: 前端(Admin): 成功後跳轉到主頁面。
    *   **子任務 2.8: 前端(Admin): 實現請求攔截與基礎權限控制**
        *   2.8.1: 前端(Admin): 創建共享 API 客戶端模塊 (`apiClient.js`)。
        *   2.8.2: 前端(Admin): 在請求函數中讀取並添加 `Authorization` header。
        *   2.8.3: 前端(Admin): 在需登入頁面的 JS 入口檢查 token 是否存在。
        *   2.8.4: 前端(Admin): 若 token 不存在，重定向到 `login.html`。
        *   2.8.5: 前端(Admin): 添加登出功能 (清除 token, 重定向)。

3.  **會員管理模組開發**:
    *   **子任務 3.1: 後端: 會員基礎 CRUD API**
        *   3.1.1: 後端: 定義會員列表請求參數(分頁/篩選)和響應結構體。
        *   3.1.2: 後端: 實現獲取會員列表 API (`/api/admin/members` GET)。
        *   3.1.3: 後端: 定義獲取/更新會員詳情請求/響應結構體。
        *   3.1.4: 後端: 實現獲取單個會員詳情 API (`/api/admin/members/:id` GET)。
        *   3.1.5: 後端: 實現更新會員信息 API (`/api/admin/members/:id` PUT)。
        *   3.1.6: 後端: (可選)實現手動創建會員 API (`/api/admin/members` POST)。
    *   **子任務 3.2: 後端: 會員狀態管理 API**
        *   3.2.1: 後端: 定義更新會員狀態請求體。
        *   3.2.2: 後端: 實現更新會員狀態 API (`/api/admin/members/:id/status` PUT)。
        *   3.2.3: 後端: 在 `user_repo` 添加更新狀態方法。
        *   3.2.4: 後端: 添加操作日誌記錄。
    *   **子任務 3.3: 後端: KYC 管理 API**
        *   3.3.1: 資料庫: 設計 `kyc_submissions` 表 (user_id, doc_type, file_path, status, timestamps, notes)。
        *   3.3.2: 後端: 創建 KYC GORM 模型。
        *   3.3.3: 後端: (可選)實現管理員上傳 KYC 文件邏輯。
        *   3.3.4: 後端: 實現獲取用戶 KYC 列表 API (`/api/admin/members/:id/kyc` GET)。
        *   3.3.5: 後端: 實現更新 KYC 狀態 API (`/api/admin/kyc/:submission_id/status` PUT)。
        *   3.3.6: 後端: 創建 `kyc_repo`。
    *   **子任務 3.4: 後端: 等級/VIP 系統 API**
        *   3.4.1: 資料庫: 設計 `levels` 表 (level, name, threshold, benefits_json)。
        *   3.4.2: 後端: 創建 Level GORM 模型。
        *   3.4.3: 後端: 實現等級規則 CRUD API (`/api/admin/levels` GET/POST/PUT/DELETE)。
        *   3.4.4: 資料庫: 在 `users` 表添加 `level_id` 並更新遷移。
        *   3.4.5: 後端: 實現手動關聯用戶等級 API (`/api/admin/members/:id/level` PUT)。
        *   3.4.6: 後端: 創建 `level_repo`。
    *   **子任務 3.5: 後端: 會員行為記錄 API (基礎)**
        *   3.5.1: MongoDB: 確認 `login_logs` 集合結構 (user_id, ip, device, timestamp)。
        *   3.5.2: 後端: 在登入成功邏輯中寫入登入日誌到 MongoDB。
        *   3.5.3: 後端: 實現查詢會員登入記錄 API (`/api/admin/members/:id/login-logs` GET)。
        *   3.5.4: 後端: 創建 `log_repo` (操作 MongoDB)。
    *   **子任務 3.6: 後端: 會員限制設定 API**
        *   3.6.1: 資料庫: 在 `users` 表添加限制相關欄位 (deposit/withdrawal limits etc.) 並更新遷移。
        *   3.6.2: 後端: 定義更新會員限制請求體。
        *   3.6.3: 後端: 實現更新會員限制 API (`/api/admin/members/:id/limits` PUT)。
        *   3.6.4: 後端: 在 `user_repo` 添加更新限制方法。
    *   **子任務 3.7: 前端(Admin): 會員列表與查詢介面**
        *   3.7.1: 前端(Admin): 創建 `members.html` 頁面結構。
        *   3.7.2: 前端(Admin): 添加搜索/篩選表單。
        *   3.7.3: 前端(Admin): 創建會員數據表格和分頁控件。
        *   3.7.4: 前端(Admin): JS: 調用 API 獲取會員列表數據。
        *   3.7.5: 前端(Admin): JS: 將數據填充到表格。
        *   3.7.6: 前端(Admin): JS: 實現分頁邏輯。
        *   3.7.7: 前端(Admin): JS: 實現搜索/篩選邏輯。
    *   **子任務 3.8: 前端(Admin): 會員詳情與編輯介面**
        *   3.8.1: 前端(Admin): 添加列表頁跳轉到詳情頁 (`member_detail.html?id=xxx`) 的鏈接。
        *   3.8.2: 前端(Admin): 創建 `member_detail.html` 頁面結構。
        *   3.8.3: 前端(Admin): JS: 頁面加載時獲取會員 ID 並調用 API 獲取詳情。
        *   3.8.4: 前端(Admin): JS: 將會員信息顯示在頁面（表單元素）。
        *   3.8.5: 前端(Admin): JS: 實現編輯信息表單提交邏輯 (調用更新 API)。
        *   3.8.6: 前端(Admin): JS: 實現狀態變更按鈕邏輯 (調用狀態更新 API)。
        *   3.8.7: 前端(Admin): JS: 實現限制設定表單提交邏輯 (調用限制更新 API)。
    *   **子任務 3.9: 前端(Admin): KYC 管理介面**
        *   3.9.1: 前端(Admin): 在會員詳情頁添加 KYC 顯示區域。
        *   3.9.2: 前端(Admin): JS: 調用 API 獲取並顯示 KYC 提交列表。
        *   3.9.3: 前端(Admin): 添加 KYC 狀態更新按鈕 (批准/拒絕)。
        *   3.9.4: 前端(Admin): JS: 實現按鈕點擊調用更新 KYC 狀態 API。
    *   **子任務 3.10: 前端(Admin): 等級/VIP 配置介面**
        *   3.10.1: 前端(Admin): 創建 `levels.html` 頁面結構。
        *   3.10.2: 前端(Admin): JS: 調用 API 獲取並顯示等級列表。
        *   3.10.3: 前端(Admin): 添加創建/編輯等級的表單。
        *   3.10.4: 前端(Admin): JS: 實現表單提交調用創建/更新 API。
        *   3.10.5: 前端(Admin): JS: 實現刪除按鈕調用刪除 API。

4.  **帳務管理模組開發**:
    *   **子任務 4.1: 後端: 支付方式配置 API**
        *   4.1.1: 資料庫: 設計 `payment_methods` 表 (id, name, type[bank, third_party, crypto], config_json, status, timestamps)。
        *   4.1.2: 後端: 創建 PaymentMethod GORM 模型。
        *   4.1.3: 後端: 實現支付方式 CRUD API (`/api/admin/payment-methods`)。
        *   4.1.4: 後端: 創建 `payment_method_repo`。
    *   **子任務 4.2: 後端: 充值管理 API**
        *   4.2.1: 資料庫: 設計 `deposits` 表 (id, user_id, amount, payment_method_id, transaction_id, status[pending, completed, failed], timestamps, notes)。
        *   4.2.2: 後端: 創建 Deposit GORM 模型。
        *   4.2.3: 後端: 實現獲取充值記錄列表 API (`/api/admin/deposits` GET)。
        *   4.2.4: 後端: 實現手動審核/更新充值狀態 API (`/api/admin/deposits/:id/status` PUT)。
        *   4.2.5: 後端: 實現處理支付回調的邏輯 (更新訂單狀態, 觸發資金增加, 記錄流水)。
        *   4.2.6: 後端: 創建 `deposit_repo`。
    *   **子任務 4.3: 後端: 提現管理 API**
        *   4.3.1: 資料庫: 設計 `withdrawals` 表 (id, user_id, amount, bank_info/wallet_address, status[pending, processing, completed, failed], timestamps, notes)。
        *   4.3.2: 後端: 創建 Withdrawal GORM 模型。
        *   4.3.3: 後端: 實現獲取提現記錄列表 API (`/api/admin/withdrawals` GET)。
        *   4.3.4: 後端: 實現審核/更新提現狀態 API (`/api/admin/withdrawals/:id/status` PUT)。
        *   4.3.5: 後端: 創建 `withdrawal_repo`。
    *   **子任務 4.4: 後端: 資金流水記錄與查詢 API**
        *   4.4.1: 資料庫: 設計 `transactions` 表 (id, user_id, type[deposit, withdrawal, game_bet, game_win, bonus, commission, manual_add, manual_sub], amount, before_balance, after_balance, related_id, timestamp, notes)。
        *   4.4.2: 後端: 創建 Transaction GORM 模型。
        *   4.4.3: 後端: 創建記錄資金流水的服務函數 (`RecordTransaction`)，確保餘額更新和流水記錄的原子性。
        *   4.4.4: 後端: 在所有涉及資金變動的地方（充值、提現、開分、發紅利等）調用 `RecordTransaction`。
        *   4.4.5: 後端: 實現查詢資金流水列表 API (`/api/admin/transactions` GET，支持按用戶/類型/時間篩選)。
        *   4.4.6: 後端: 創建 `transaction_repo`。
        *   4.4.7: 後端: 更新 `users` 表，增加 `balance` 欄位。
    *   **子任務 4.5: 後端: 玩家開分 API (安全強化)**
        *   4.5.1: 後端: 定義開分請求體 (user_id, amount, type[add/sub], **reason: string [Mandatory]**)。
        *   4.5.2: 後端: 實現開分 API (`/api/admin/members/adjust-balance` POST)。
        *   4.5.3: 後端(API Handler): **輸入驗證**: 驗證 amount, user_id, reason。
        *   4.5.4: 後端(API Handler): **權限檢查**: 調用 RBAC 中間件/邏輯，檢查是否有 `ADJUST_BALANCE_ADD` 或 `ADJUST_BALANCE_SUB` 權限。
        *   4.5.5: 後端(API Handler): **額度檢查**: 讀取操作員角色對應的單次/單日額度限制，檢查是否超限。
        *   4.5.6: 後端(API Handler): **二次驗證 (2FA)**: 檢查此操作是否需要 2FA，如果需要，驗證請求中提供的 2FA code。
        *   4.5.7: 後端(API Handler): 調用 `RecordTransaction` 服務函數處理餘額和流水。
        *   4.5.8: 後端(Service): 在 `RecordTransaction` 或獨立日誌服務中，確保記錄操作者 ID, IP, 理由 到 `transactions` 表和詳細的操作日誌 (MongoDB `admin_operations`)。
        *   4.5.9: 後端(API Endpoint): 應用**速率限制**。
    *   **子任務 4.6: 後端: 開分安全配置管理**
        *   4.6.1: 資料庫/配置: 設計存儲角色開分額度限制的結構 (e.g., `role_limits` 表或配置文件)。
        *   4.6.2: 資料庫/配置: 設計存儲哪些操作需要 2FA 的配置 (e.g., `operation_security_config` 表或配置文件)。
        *   4.6.3: 後端: 實現管理這些安全配置的 CRUD API (`/api/admin/security-configs/limits`, `/api/admin/security-configs/2fa`)。
        *   4.6.4: 後端: (可選)實現觸發異常開分警報的邏輯 (e.g., 在 `RecordTransaction` 後檢查規則並發送通知)。
    *   **子任務 4.7: 前端(Admin): 強化開分介面**
        *   4.7.1: 前端(Admin): 更新開分表單 (`adjust_balance.html` 或會員詳情)，`reason` 欄位設為必填。
        *   4.7.2: 前端(Admin): 實現處理 2FA 驗證的輸入和流程（如果需要 2FA）。
        *   4.7.3: 前端(Admin): (可選) 顯示當前操作員的剩餘可用開分額度。
        *   4.7.4: 前端(Admin): JS: 實現表單提交調用新的安全強化版開分 API。
    *   **子任務 4.8 (原 4.6): 後端: 返水/紅利規則配置 API**
        *   4.8.1 (原 4.6.1): 資料庫: 設計 `rebate_rules` 表 (id, level_id/user_id, game_type, rate, status)。
        *   4.8.2 (原 4.6.2): 資料庫: 設計 `bonus_rules` 表 (id, name, type, config_json, start_time, end_time, status)。
        *   4.8.3 (原 4.6.3): 後端: 創建對應 GORM 模型。
        *   4.8.4 (原 4.6.4): 後端: 實現返水規則 CRUD API (`/api/admin/rebate-rules`)。
        *   4.8.5 (原 4.6.5): 後端: 實現紅利規則 CRUD API (`/api/admin/bonus-rules`)。
        *   4.8.6 (原 4.6.6): 後端: 創建 `rebate_rule_repo`, `bonus_rule_repo`。
    *   **子任務 4.9 (原 4.7): 後端: 返水/紅利計算與發放 API/定時任務**
        *   4.9.1 (原 4.7.1): 後端: 設計計算返水的邏輯 (基於遊戲流水和規則)。
        *   4.9.2 (原 4.7.2): 後端: 創建定時任務 (e.g., daily) 觸發返水計算與發放 (調用 `RecordTransaction`)。
        *   4.9.3 (原 4.7.3): 後端: 設計計算/發放特定紅利的邏輯 (可能由 API 觸發或定時任務觸發)。
        *   4.9.4 (原 4.7.4): 後端: 實現觸發紅利發放的 API (可選，用於手動觸發)。
    *   **子任務 4.10 (原 4.8): 後端: 財務報表生成 API**
        *   4.10.1 (原 4.8.1): 後端: 設計獲取日/月收支報表數據的查詢邏輯 (聚合 `transactions` 表)。
        *   4.10.2 (原 4.8.2): 後端: 實現獲取收支報表 API (`/api/admin/reports/financial-summary`)。
        *   4.10.3 (原 4.8.3): 後端: 設計獲取遊戲盈虧分析數據的查詢邏輯。
        *   4.10.4 (原 4.8.4): 後端: 實現獲取遊戲盈虧報表 API (`/api/admin/reports/game-profit-loss`)。
    *   **子任務 4.11 (原 4.9): 前端(Admin): 充值/提現管理介面**
        *   4.11.1 (原 4.9.1): 前端(Admin): 創建 `deposits.html`, `withdrawals.html` 頁面。
        *   4.11.2 (原 4.9.2): 前端(Admin): JS: 調用 API 獲取列表數據並填充表格。
        *   4.11.3 (原 4.9.3): 前端(Admin): 添加審核按鈕 (批准/拒絕)。
        *   4.11.4 (原 4.9.4): 前端(Admin): JS: 實現審核按鈕調用更新狀態 API。
    *   **子任務 4.12 (原 4.10): 前端(Admin): 資金流水查詢介面**
        *   4.12.1 (原 4.10.1): 前端(Admin): 創建 `transactions.html` 頁面。
        *   4.12.2 (原 4.10.2): 前端(Admin): 添加篩選表單 (用戶 ID, 類型, 日期範圍)。
        *   4.12.3 (原 4.10.3): 前端(Admin): JS: 調用 API 獲取數據並顯示。
    *   **子任務 4.13 (原 4.11): 前端(Admin): 玩家開分介面 (安全強化對應)**
        *   4.13.1 (原 4.11.1): 前端(Admin): 更新 `adjust_balance.html` 或會員詳情頁的開分表單，添加必填的 `reason` 輸入框。
        *   4.13.2 (原 4.11.2): 前端(Admin): JS: 更新表單提交邏輯，包含 `reason`，並處理可能的 2FA 流程。
        *   4.13.3 (原 4.11.3): 前端(Admin): 根據登入用戶權限決定是否顯示此功能，並可選顯示可用額度。
    *   **子任務 4.14 (原 4.12): 前端(Admin): 返水/紅利配置介面**
        *   4.14.1 (原 4.12.1): 前端(Admin): 創建 `rebate_rules.html`, `bonus_rules.html` 頁面。
        *   4.14.2 (原 4.12.2): 前端(Admin): JS: 調用 API 獲取規則列表。
        *   4.14.3 (原 4.12.3): 前端(Admin): 添加創建/編輯規則的表單。
        *   4.14.4 (原 4.12.4): 前端(Admin): JS: 實現表單提交調用 CRUD API。
    *   **子任務 4.15 (原 4.13): 前端(Admin): 財務報表展示介面**
        *   4.15.1 (原 4.13.1): 前端(Admin): 創建 `reports.html` 頁面。
        *   4.15.2 (原 4.13.2): 前端(Admin): 添加日期範圍或其他篩選條件。
        *   4.15.3 (原 4.13.3): 前端(Admin): JS: 調用報表 API 獲取數據。
        *   4.15.4 (原 4.13.4): 前端(Admin): JS: 使用 Chart.js 將數據渲染成圖表。

5.  **代理與經銷商管理模組開發**:
    *   **子任務 5.1: 後端: 代理層級結構管理 API**
        *   5.1.1: 資料庫: 在 `users` 表中添加 `parent_agent_id` 欄位 (nullable, FK to users.id)。更新遷移。
        *   5.1.2: 資料庫: 確保 `roles` 表包含「代理」角色。
        *   5.1.3: 後端: 實現創建代理用戶 API (`/api/admin/agents` POST)，設置角色為代理並關聯上級。
        *   5.1.4: 後端: 實現查詢代理列表 API (`/api/admin/agents` GET)，支持按層級/上級篩選。
        *   5.1.5: 後端: 實現查詢代理詳情 API (`/api/admin/agents/:id` GET)。
        *   5.1.6: 後端: 實現修改代理信息 API (`/api/admin/agents/:id` PUT)。
        *   5.1.7: 後端: 實現遞歸查詢下線代理列表的邏輯 (服務層)。
        *   5.1.8: 後端: 創建 `agent_repo` (或擴展 `user_repo`)。
    *   **子任務 5.2: 後端: 佣金規則配置 API**
        *   5.2.1: 資料庫: 設計 `commission_rules` 表 (id, agent_level/agent_id, game_type, commission_type[revenue_share, turnover], rate, status, timestamps)。
        *   5.2.2: 後端: 創建 CommissionRule GORM 模型。
        *   5.2.3: 後端: 實現佣金規則 CRUD API (`/api/admin/commission-rules`)。
        *   5.2.4: 後端: 創建 `commission_rule_repo`。
    *   **子任務 5.3: 後端: 佣金計算與結算 API/定時任務**
        *   5.3.1: 資料庫: 設計 `commission_settlements` 表 (id, agent_id, period_start, period_end, turnover_amount, profit_amount, rate, amount, status[pending, paid], settled_at, timestamps)。
        *   5.3.2: 後端: 創建 CommissionSettlement GORM 模型。
        *   5.3.3: 後端: 設計佣金計算邏輯 (聚合下線流水/盈虧, 應用規則)。
        *   5.3.4: 後端: 創建定時任務 (e.g., weekly/monthly) 觸發佣金計算並寫入結算表。
        *   5.3.5: 後端: 實現手動觸發計算 API (可選)。
        *   5.3.6: 後端: 實現查詢佣金結算記錄 API (`/api/admin/commission-settlements` GET)。
        *   5.3.7: 後端: 創建 `commission_settlement_repo`。
    *   **子任務 5.4: 後端: 推廣工具管理 API**
        *   5.4.1: 資料庫: 設計 `promotion_codes` 表 (id, agent_id, code, link, status, timestamps)。
        *   5.4.2: 後端: 創建 PromotionCode GORM 模型。
        *   5.4.3: 後端: 實現生成推廣碼/連結邏輯。
        *   5.4.4: 後端: 實現推廣碼/連結 CRUD API (`/api/admin/promotion-codes`)。
        *   5.4.5: 後端: 創建 `promotion_code_repo`。
        *   5.4.6: 後端: 修改用戶註冊流程，允許關聯推廣碼 (將 agent_id 寫入新用戶記錄)。
    *   **子任務 5.5: 後端: 代理績效報表 API**
        *   5.5.1: 後端: 設計查詢代理新增下線數的邏輯。
        *   5.5.2: 後端: 設計查詢代理下線總投注額/總盈虧的邏輯 (聚合流水)。
        *   5.5.3: 後端: 設計查詢代理總佣金的邏輯 (聚合結算表)。
        *   5.5.4: 後端: 實現獲取代理績效報表 API (`/api/admin/reports/agent-performance` GET)。
    *   **子任務 5.6: 前端(Admin): 代理列表與層級管理介面**
        *   5.6.1: 前端(Admin): 創建 `agents.html` 頁面。
        *   5.6.2: 前端(Admin): JS: 調用 API 獲取代理列表數據。
        *   5.6.3: 前端(Admin): 實現代理層級的樹狀或嵌套列表顯示。
        *   5.6.4: 前端(Admin): 添加創建/編輯代理的表單/彈窗。
        *   5.6.5: 前端(Admin): JS: 實現創建/編輯代理的 API 調用。
    *   **子任務 5.7: 前端(Admin): 佣金規則配置介面**
        *   5.7.1: 前端(Admin): 創建 `commission_rules.html` 頁面。
        *   5.7.2: 前端(Admin): JS: 調用 API 獲取規則列表。
        *   5.7.3: 前端(Admin): 添加創建/編輯規則的表單。
        *   5.7.4: 前端(Admin): JS: 實現表單提交調用 CRUD API。
    *   **子任務 5.8: 前端(Admin): 佣金結算記錄查詢介面**
        *   5.8.1: 前端(Admin): 創建 `commission_settlements.html` 頁面。
        *   5.8.2: 前端(Admin): 添加篩選條件 (代理 ID, 時間範圍)。
        *   5.8.3: 前端(Admin): JS: 調用 API 獲取結算記錄並顯示。
    *   **子任務 5.9: 前端(Admin): 推廣工具管理介面**
        *   5.9.1: 前端(Admin): 創建 `promotion_codes.html` 頁面。
        *   5.9.2: 前端(Admin): JS: 調用 API 獲取指定代理的推廣碼/連結。
        *   5.9.3: 前端(Admin): 添加生成新推廣碼的功能。
        *   5.9.4: 前端(Admin): JS: 實現生成/刪除推廣碼的 API 調用。
    *   **子任務 5.10: 前端(Admin): 代理績效報表展示介面**
        *   5.10.1: 前端(Admin): 在 `reports.html` 或創建新頁面 `agent_reports.html`。
        *   5.10.2: 前端(Admin): 添加篩選條件 (代理 ID, 時間範圍)。
        *   5.10.3: 前端(Admin): JS: 調用績效報表 API 獲取數據。
        *   5.10.4: 前端(Admin): 將績效數據顯示在表格或圖表中。

6.  **遊戲管理模組 (後台配置)**:
    *   **子任務 6.1: 後端: 遊戲與供應商數據模型設計**
        *   6.1.1: 資料庫: 設計 `game_providers` 表 (id, name, code, api_config_json, status, timestamps)。
        *   6.1.2: 資料庫: 設計 `game_categories` 表 (id, name, code, description, status, timestamps)。
        *   6.1.3: 資料庫: 設計 `games` 表 (id, provider_id, category_id, name, code, description, image_url, status[enabled, disabled, maintenance], config_json[e.g., bet limits, features], timestamps)。
        *   6.1.4: 資料庫: (可選) 設計更細化的配置表，例如 `game_odds`, `game_rooms` (如果遊戲類型複雜)。
        *   6.1.5: 後端: 創建對應的 GORM 模型 (`game_provider.go`, `game_category.go`, `game.go` 等)。
    *   **子任務 6.2: 後端: 遊戲供應商管理 API**
        *   6.2.1: 後端: 實現遊戲供應商 CRUD API (`/api/admin/game-providers`)。
        *   6.2.2: 後端: 創建 `game_provider_repo`。
    *   **子任務 6.3: 後端: 遊戲分類管理 API**
        *   6.3.1: 後端: 實現遊戲分類 CRUD API (`/api/admin/game-categories`)。
        *   6.3.2: 後端: 創建 `game_category_repo`。
    *   **子任務 6.4: 後端: 遊戲管理 API**
        *   6.4.1: 後端: 實現遊戲列表查詢 API (`/api/admin/games` GET, 支持按供應商/分類/狀態篩選)。
        *   6.4.2: 後端: 實現創建遊戲 API (`/api/admin/games` POST)。
        *   6.4.3: 後端: 實現獲取單個遊戲詳情 API (`/api/admin/games/:id` GET)。
        *   6.4.4: 後端: 實現更新遊戲信息 (包括基本信息、配置 `config_json`) API (`/api/admin/games/:id` PUT)。
        *   6.4.5: 後端: 實現更新遊戲狀態 API (`/api/admin/games/:id/status` PUT)。
        *   6.4.6: 後端: 創建 `game_repo`。
    *   **子任務 6.5: 後端: 遊戲配置獲取 API (供前端或遊戲服務使用)**
        *   6.5.1: 後端: 實現獲取啟用狀態的遊戲列表 API (`/api/games` GET - 無需管理員權限，可能需要緩存)。
        *   6.5.2: 後端: 實現獲取單個遊戲詳細配置 API (`/api/games/:code/config` GET - 無需管理員權限)。
    *   **子任務 6.6: 前端(Admin): 遊戲供應商管理介面**
        *   6.6.1: 前端(Admin): 創建 `game_providers.html` 頁面。
        *   6.6.2: 前端(Admin): JS: 調用 API 獲取供應商列表並顯示。
        *   6.6.3: 前端(Admin): 添加創建/編輯供應商的表單。
        *   6.6.4: 前端(Admin): JS: 實現表單提交調用 CRUD API。
    *   **子任務 6.7: 前端(Admin): 遊戲分類管理介面**
        *   6.7.1: 前端(Admin): 創建 `game_categories.html` 頁面。
        *   6.7.2: 前端(Admin): JS: 調用 API 獲取分類列表並顯示。
        *   6.7.3: 前端(Admin): 添加創建/編輯分類的表單。
        *   6.7.4: 前端(Admin): JS: 實現表單提交調用 CRUD API。
    *   **子任務 6.8: 前端(Admin): 遊戲列表與管理介面**
        *   6.8.1: 前端(Admin): 創建 `games.html` 頁面。
        *   6.8.2: 前端(Admin): 添加篩選條件 (供應商, 分類, 狀態)。
        *   6.8.3: 前端(Admin): JS: 調用 API 獲取遊戲列表並顯示 (包含狀態)。
        *   6.8.4: 前端(Admin): 添加啟用/禁用/維護的按鈕。
        *   6.8.5: 前端(Admin): 添加跳轉到遊戲編輯頁面的鏈接。
        *   6.8.6: 前端(Admin): JS: 實現狀態變更按鈕調用更新狀態 API。
    *   **子任務 6.9: 前端(Admin): 遊戲編輯介面**
        *   6.9.1: 前端(Admin): 創建 `game_edit.html` 頁面。
        *   6.9.2: 前端(Admin): JS: 頁面加載時獲取遊戲 ID 並調用 API 獲取詳情。
        *   6.9.3: 前端(Admin): 顯示遊戲基本信息表單 (名稱, 描述, 分類, 供應商等)。
        *   6.9.4: 前端(Admin): 提供編輯 `config_json` 的方式 (可能是文本域或更結構化的表單)。
        *   6.9.5: 前端(Admin): JS: 實現表單提交調用更新遊戲信息 API。

7.  **系統管理模組開發**:
    *   **子任務 7.1: 後端: 角色與權限管理 API (Admin CRUD)**
        *   7.1.1: 後端: 實現角色 CRUD API (`/api/admin/roles`) - 擴展 `2.2.5` 的 `rbac_repo`。
        *   7.1.2: 後端: 實現權限 CRUD API (`/api/admin/permissions`) - 擴展 `2.2.5` 的 `rbac_repo`。
        *   7.1.3: 後端: 實現管理角色權限關聯的 API (`/api/admin/roles/:id/permissions` GET/PUT)。
        *   7.1.4: 後端: 實現分配用戶角色的 API (`/api/admin/users/:id/role` PUT) - 擴展 `3.1.5` 或 `user_repo`。
    *   **子任務 7.2: 後端: 日誌查詢 API**
        *   7.2.1: 後端: 實現查詢操作日誌 API (`/api/admin/logs/operations` GET, from MongoDB `admin_operations`, 支持篩選)。
        *   7.2.2: 後端: 實現查詢登入日誌 API (`/api/admin/logs/login` GET, from MongoDB `login_logs`, 支持篩選)。
        *   7.2.3: 後端: 實現查詢錯誤日誌 API (`/api/admin/logs/errors` GET, from MongoDB `error_logs`, 支持篩選)。
        *   7.2.4: 後端: 確保 `log_repo` (子任務 `3.5.4`) 包含這些查詢方法。
    *   **子任務 7.3: 後端: 安全設定 API**
        *   7.3.1: 資料庫: (可選) 設計 `security_configs` 表 (id, key, value, description, type[bool, int, string], timestamps) 或使用配置文件。
        *   7.3.2: 後端: 創建 SecurityConfig GORM 模型 (如果使用數據庫)。
        *   7.3.3: 後端: 實現安全設定讀取 API (`/api/admin/security-configs` GET)。
        *   7.3.4: 後端: 實現安全設定更新 API (`/api/admin/security-configs` PUT)。
        *   7.3.5: 後端: 創建 `security_config_repo` 或配置讀取服務。
    *   **子任務 7.4: 後端: 通知管理 API**
        *   7.4.1: 資料庫: 設計 `notifications` 表 (id, type[system, user], title, content, target_user_id[nullable], status[unread, read], created_at)。
        *   7.4.2: 後端: 創建 Notification GORM 模型。
        *   7.4.3: 後端: 實現發送系統通知/指定用戶通知的 API (`/api/admin/notifications` POST)。
        *   7.4.4: 後端: 實現獲取通知列表 API (`/api/admin/notifications` GET)。
        *   7.4.5: 後端: 創建 `notification_repo`。
    *   **子任務 7.5: 後端: 多語言/貨幣配置 API**
        *   7.5.1: 資料庫: 設計 `languages` 表 (id, code, name, status[enabled, disabled])。
        *   7.5.2: 資料庫: 設計 `currencies` 表 (id, code, name, symbol, rate_to_base, status[enabled, disabled])。
        *   7.5.3: 後端: 創建 Language, Currency GORM 模型。
        *   7.5.4: 後端: 實現語言 CRUD API (`/api/admin/languages`)。
        *   7.5.5: 後端: 實現貨幣 CRUD API (`/api/admin/currencies`)。
        *   7.5.6: 後端: 創建 `language_repo`, `currency_repo`。
    *   **子任務 7.6: 前端(Admin): 角色與權限管理介面**
        *   7.6.1: 前端(Admin): 創建 `roles.html`, `permissions.html` 頁面。
        *   7.6.2: 前端(Admin): JS: 調用 API 獲取角色/權限列表並顯示。
        *   7.6.3: 前端(Admin): 添加創建/編輯角色/權限的表單。
        *   7.6.4: 前端(Admin): JS: 實現表單提交調用 CRUD API。
        *   7.6.5: 前端(Admin): 實現管理角色權限關聯的介面 (e.g., 在角色編輯頁勾選權限)。
        *   7.6.6: 前端(Admin): 實現修改用戶角色的介面 (e.g., 在用戶詳情頁下拉選擇)。
    *   **子任務 7.7: 前端(Admin): 日誌查詢介面**
        *   7.7.1: 前端(Admin): 創建 `logs.html` 頁面 (可分頁籤顯示不同日誌)。
        *   7.7.2: 前端(Admin): 添加篩選條件 (日誌類型, 用戶 ID, 時間範圍等)。
        *   7.7.3: 前端(Admin): JS: 調用對應 API 獲取日誌並顯示。
    *   **子任務 7.8: 前端(Admin): 安全設定介面**
        *   7.8.1: 前端(Admin): 創建 `security_configs.html` 頁面。
        *   7.8.2: 前端(Admin): JS: 調用 API 獲取設定列表並顯示。
        *   7.8.3: 前端(Admin): 提供編輯設定值的表單。
        *   7.8.4: 前端(Admin): JS: 實現表單提交調用更新 API。
    *   **子任務 7.9: 前端(Admin): 通知管理介面**
        *   7.9.1: 前端(Admin): 創建 `notifications.html` 頁面。
        *   7.9.2: 前端(Admin): JS: 調用 API 獲取通知列表並顯示。
        *   7.9.3: 前端(Admin): 添加發送新通知的表單 (可選目標用戶)。
        *   7.9.4: 前端(Admin): JS: 實現表單提交調用創建 API。
    *   **子任務 7.10: 前端(Admin): 多語言/貨幣配置介面**
        *   7.10.1: 前端(Admin): 創建 `languages.html`, `currencies.html` 頁面。
        *   7.10.2: 前端(Admin): JS: 調用 API 獲取語言/貨幣列表並顯示。
        *   7.10.3: 前端(Admin): 添加創建/編輯語言/貨幣的表單 (包含狀態切換)。
        *   7.10.4: 前端(Admin): JS: 實現表單提交調用 CRUD API。

8.  **活動與行銷管理模組開發**:
    *   **子任務 8.1: 後端: 活動數據模型與管理 API**
        *   8.1.1: 資料庫: 設計 `activities` 表 (id, name, type[bonus, tournament, special_offer], description, rules_config_json, banner_image_url, start_time, end_time, status[draft, active, expired], timestamps)。
        *   8.1.2: 後端: 創建 Activity GORM 模型。
        *   8.1.3: 後端: 實現活動 CRUD API (`/api/admin/activities`)。
        *   8.1.4: 後端: 創建 `activity_repo`。
    *   **子任務 8.2: 後端: (可選) 活動參與記錄與資格檢查**
        *   8.2.1: 資料庫: 設計 `activity_participation` 表 (id, user_id, activity_id, participated_at, status[qualified, disqualified], notes)。
        *   8.2.2: 後端: 創建 ActivityParticipation GORM 模型。
        *   8.2.3: 後端: 設計檢查用戶是否符合活動資格的邏輯 (根據 `rules_config_json`)。
        *   8.2.4: 後端: (可選) 實現記錄用戶參與活動的 API 或內部服務。
    *   **子任務 8.3: 前端(Admin): 活動列表與創建/編輯介面**
        *   8.3.1: 前端(Admin): 創建 `activities.html` 頁面。
        *   8.3.2: 前端(Admin): JS: 調用 API 獲取活動列表並顯示。
        *   8.3.3: 前端(Admin): 添加創建/編輯活動的表單 (包含名稱、類型、描述、規則配置輸入、時間、狀態等)。
        *   8.3.4: 前端(Admin): JS: 實現表單提交調用 CRUD API。
        *   8.3.5: 前端(Admin): 實現活動狀態切換 (草稿/啟用/過期)。
    *   **(待定) 子任務 8.4 及後續**: 根據具體的活動類型 (例如：首儲紅利、簽到活動、排行榜競賽) 可能需要更特定的後端邏輯和前端介面。這部分可以在基礎活動管理框架完成後再細化。

**(後續將每個大步驟細分為更小的任務)** 