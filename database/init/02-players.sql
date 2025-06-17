-- 玩家管理相關表結構
-- 建立時間: 2024-12-19

USE nexus_gaming;

-- 建立玩家表
CREATE TABLE IF NOT EXISTS players (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    player_id VARCHAR(32) NOT NULL UNIQUE COMMENT '玩家唯一ID',
    username VARCHAR(100) NOT NULL UNIQUE COMMENT '玩家帳號',
    email VARCHAR(255) COMMENT '電子郵件',
    phone VARCHAR(20) COMMENT '電話號碼',
    real_name VARCHAR(100) COMMENT '真實姓名',
    nickname VARCHAR(100) COMMENT '暱稱',
    avatar_url VARCHAR(500) COMMENT '頭像URL',
    birth_date DATE COMMENT '生日',
    gender ENUM('male', 'female', 'other') COMMENT '性別',
    country VARCHAR(100) COMMENT '國家',
    language VARCHAR(10) DEFAULT 'zh-TW' COMMENT '偏好語言',
    timezone VARCHAR(50) DEFAULT 'Asia/Taipei' COMMENT '時區',
    status ENUM('active', 'inactive', 'suspended', 'deleted') DEFAULT 'active' COMMENT '帳戶狀態',
    verification_level ENUM('none', 'email', 'phone', 'identity') DEFAULT 'none' COMMENT '驗證等級',
    risk_level ENUM('low', 'medium', 'high', 'blacklist') DEFAULT 'low' COMMENT '風險等級',
    vip_level INT DEFAULT 0 COMMENT 'VIP等級',
    referrer_id BIGINT COMMENT '推薦人ID',
    agent_id INT COMMENT '所屬代理商ID',
    dealer_id INT COMMENT '所屬經銷商ID',
    registration_ip VARCHAR(45) COMMENT '註冊IP',
    last_login_ip VARCHAR(45) COMMENT '最後登入IP',
    last_login_at TIMESTAMP NULL COMMENT '最後登入時間',
    login_count INT DEFAULT 0 COMMENT '登入次數',
    total_deposit DECIMAL(15,2) DEFAULT 0.00 COMMENT '總儲值金額',
    total_withdraw DECIMAL(15,2) DEFAULT 0.00 COMMENT '總提領金額',
    total_bet DECIMAL(15,2) DEFAULT 0.00 COMMENT '總下注金額',
    total_win DECIMAL(15,2) DEFAULT 0.00 COMMENT '總贏得金額',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL COMMENT '刪除時間（軟刪除）',
    INDEX idx_player_id (player_id),
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_agent_id (agent_id),
    INDEX idx_dealer_id (dealer_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (referrer_id) REFERENCES players(id),
    FOREIGN KEY (agent_id) REFERENCES users(id),
    FOREIGN KEY (dealer_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家基本資料表';

-- 建立玩家錢包表
CREATE TABLE IF NOT EXISTS player_wallets (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    currency VARCHAR(10) DEFAULT 'TWD' COMMENT '幣別',
    balance DECIMAL(15,2) DEFAULT 0.00 COMMENT '可用餘額',
    frozen_balance DECIMAL(15,2) DEFAULT 0.00 COMMENT '凍結餘額',
    total_balance DECIMAL(15,2) GENERATED ALWAYS AS (balance + frozen_balance) COMMENT '總餘額',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_player_currency (player_id, currency),
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家錢包表';

-- 建立玩家狀態歷史表
CREATE TABLE IF NOT EXISTS player_status_history (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    old_status VARCHAR(20) COMMENT '原狀態',
    new_status VARCHAR(20) NOT NULL COMMENT '新狀態',
    reason TEXT COMMENT '變更原因',
    operator_id INT COMMENT '操作者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_player_id (player_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE,
    FOREIGN KEY (operator_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家狀態變更歷史表';

-- 建立玩家標籤表
CREATE TABLE IF NOT EXISTS player_tags (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '標籤名稱',
    description TEXT COMMENT '標籤描述',
    color VARCHAR(7) DEFAULT '#007bff' COMMENT '標籤顏色',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家標籤定義表';

-- 建立玩家標籤關聯表
CREATE TABLE IF NOT EXISTS player_tag_relations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    tag_id INT NOT NULL COMMENT '標籤ID',
    assigned_by INT COMMENT '分配者ID',
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_player_tag (player_id, tag_id),
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES player_tags(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家標籤關聯表';

-- 建立玩家限制設定表
CREATE TABLE IF NOT EXISTS player_restrictions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    restriction_type ENUM('bet_limit', 'deposit_limit', 'game_access', 'time_limit') NOT NULL COMMENT '限制類型',
    restriction_value JSON NOT NULL COMMENT '限制值（JSON格式）',
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '生效時間',
    end_time TIMESTAMP NULL COMMENT '結束時間',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否啟用',
    reason TEXT COMMENT '限制原因',
    operator_id INT COMMENT '設定者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_player_id (player_id),
    INDEX idx_restriction_type (restriction_type),
    INDEX idx_is_active (is_active),
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE,
    FOREIGN KEY (operator_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家限制設定表';

-- 插入預設標籤
INSERT INTO player_tags (name, description, color) VALUES 
('新手', '新註冊的玩家', '#28a745'),
('VIP', 'VIP等級玩家', '#ffc107'),
('高風險', '需要特別關注的玩家', '#dc3545'),
('活躍', '經常遊戲的玩家', '#007bff'),
('大戶', '高額投注玩家', '#6f42c1')
ON DUPLICATE KEY UPDATE name=name; 