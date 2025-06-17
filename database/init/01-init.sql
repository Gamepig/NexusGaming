-- 遊戲管理後台系統資料庫初始化腳本
-- 建立時間: 2024-12-19

USE nexus_gaming;

-- 設定字符集
SET NAMES utf8mb4;
SET character_set_client = utf8mb4;

-- 建立使用者角色表
CREATE TABLE IF NOT EXISTS roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名稱',
    description TEXT COMMENT '角色描述',
    permissions JSON COMMENT '權限列表',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='使用者角色表';

-- 插入預設角色
INSERT INTO roles (name, description, permissions) VALUES 
('admin', '系統管理員', '["*"]'),
('agent', '代理商', '["player.view", "financial.view", "agent.manage"]'),
('dealer', '經銷商', '["player.view", "financial.view"]')
ON DUPLICATE KEY UPDATE name=name;

-- 建立使用者表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE COMMENT '使用者名稱',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT '電子郵件',
    password_hash VARCHAR(255) NOT NULL COMMENT '密碼雜湊',
    role_id INT NOT NULL COMMENT '角色ID',
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active' COMMENT '帳戶狀態',
    last_login_at TIMESTAMP NULL COMMENT '最後登入時間',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='使用者表';

-- 建立預設管理員帳戶（密碼: admin123）
INSERT INTO users (username, email, password_hash, role_id) VALUES 
('admin', 'admin@nexusgaming.com', '$2a$10$rOlF8WF9r7LvbLw5uH2eVeF9r.XrqD3iJ5B6QJzJ7D2E3fG4h5I6K', 1)
ON DUPLICATE KEY UPDATE username=username;

-- 建立操作日誌表
CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL COMMENT '操作者ID',
    action VARCHAR(100) NOT NULL COMMENT '操作動作',
    resource VARCHAR(100) NOT NULL COMMENT '操作資源',
    resource_id VARCHAR(100) COMMENT '資源ID',
    details JSON COMMENT '操作詳情',
    ip_address VARCHAR(45) COMMENT 'IP位址',
    user_agent TEXT COMMENT '使用者代理',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日誌表'; 