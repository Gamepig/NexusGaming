-- 遊戲管理相關表結構
-- 建立時間: 2024-12-19

USE nexus_gaming;

-- 建立遊戲基本資訊表
CREATE TABLE IF NOT EXISTS games (
    id INT AUTO_INCREMENT PRIMARY KEY,
    game_code VARCHAR(50) NOT NULL UNIQUE COMMENT '遊戲代碼',
    name VARCHAR(100) NOT NULL COMMENT '遊戲名稱',
    name_en VARCHAR(100) COMMENT '英文名稱',
    description TEXT COMMENT '遊戲描述',
    game_type ENUM('texas_holdem', 'stud_poker', 'baccarat', 'blackjack', 'roulette', 'slots') NOT NULL COMMENT '遊戲類型',
    category VARCHAR(50) COMMENT '遊戲分類',
    thumbnail_url VARCHAR(500) COMMENT '縮圖URL',
    banner_url VARCHAR(500) COMMENT '橫幅圖URL',
    min_bet DECIMAL(10,2) DEFAULT 1.00 COMMENT '最低下注金額',
    max_bet DECIMAL(10,2) DEFAULT 10000.00 COMMENT '最高下注金額',
    house_edge DECIMAL(5,4) DEFAULT 0.0250 COMMENT '莊家優勢（預設2.5%）',
    rtp_rate DECIMAL(5,4) DEFAULT 0.9750 COMMENT '玩家回報率（預設97.5%）',
    status ENUM('active', 'inactive', 'maintenance', 'testing') DEFAULT 'inactive' COMMENT '遊戲狀態',
    is_featured BOOLEAN DEFAULT FALSE COMMENT '是否精選遊戲',
    sort_order INT DEFAULT 0 COMMENT '排序順序',
    version VARCHAR(20) DEFAULT '1.0.0' COMMENT '遊戲版本',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_game_code (game_code),
    INDEX idx_game_type (game_type),
    INDEX idx_status (status),
    INDEX idx_sort_order (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='遊戲基本資訊表';

-- 建立遊戲配置表
CREATE TABLE IF NOT EXISTS game_configs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    game_id INT NOT NULL COMMENT '遊戲ID',
    config_key VARCHAR(100) NOT NULL COMMENT '配置鍵',
    config_value JSON NOT NULL COMMENT '配置值',
    description TEXT COMMENT '配置描述',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否啟用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_game_config (game_id, config_key),
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='遊戲配置表';

-- 建立遊戲房間表
CREATE TABLE IF NOT EXISTS game_rooms (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    room_code VARCHAR(32) NOT NULL UNIQUE COMMENT '房間代碼',
    game_id INT NOT NULL COMMENT '遊戲ID',
    name VARCHAR(100) NOT NULL COMMENT '房間名稱',
    description TEXT COMMENT '房間描述',
    room_type ENUM('public', 'private', 'vip') DEFAULT 'public' COMMENT '房間類型',
    max_players INT DEFAULT 6 COMMENT '最大玩家數',
    current_players INT DEFAULT 0 COMMENT '當前玩家數',
    min_bet DECIMAL(10,2) COMMENT '最低下注（覆蓋遊戲設定）',
    max_bet DECIMAL(10,2) COMMENT '最高下注（覆蓋遊戲設定）',
    status ENUM('active', 'inactive', 'full', 'maintenance') DEFAULT 'active' COMMENT '房間狀態',
    ai_enabled BOOLEAN DEFAULT TRUE COMMENT '是否啟用AI',
    ai_difficulty ENUM('easy', 'medium', 'hard', 'expert') DEFAULT 'medium' COMMENT 'AI難度',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_room_code (room_code),
    INDEX idx_game_id (game_id),
    INDEX idx_status (status),
    FOREIGN KEY (game_id) REFERENCES games(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='遊戲房間表';

-- 建立遊戲場次表
CREATE TABLE IF NOT EXISTS game_sessions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    session_code VARCHAR(64) NOT NULL UNIQUE COMMENT '場次代碼',
    room_id BIGINT NOT NULL COMMENT '房間ID',
    game_id INT NOT NULL COMMENT '遊戲ID',
    session_type ENUM('practice', 'normal', 'tournament') DEFAULT 'normal' COMMENT '場次類型',
    status ENUM('waiting', 'playing', 'finished', 'cancelled') DEFAULT 'waiting' COMMENT '場次狀態',
    max_players INT DEFAULT 6 COMMENT '最大玩家數',
    current_players INT DEFAULT 0 COMMENT '當前玩家數',
    min_bet DECIMAL(10,2) NOT NULL COMMENT '最低下注',
    max_bet DECIMAL(10,2) NOT NULL COMMENT '最高下注',
    total_pot DECIMAL(15,2) DEFAULT 0.00 COMMENT '總獎池',
    house_commission DECIMAL(15,2) DEFAULT 0.00 COMMENT '抽水金額',
    game_data JSON COMMENT '遊戲數據（牌局、下注等）',
    ai_players JSON COMMENT 'AI玩家資訊',
    started_at TIMESTAMP NULL COMMENT '開始時間',
    finished_at TIMESTAMP NULL COMMENT '結束時間',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_session_code (session_code),
    INDEX idx_room_id (room_id),
    INDEX idx_game_id (game_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (room_id) REFERENCES game_rooms(id),
    FOREIGN KEY (game_id) REFERENCES games(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='遊戲場次表';

-- 建立遊戲參與記錄表
CREATE TABLE IF NOT EXISTS game_participations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    session_id BIGINT NOT NULL COMMENT '場次ID',
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    seat_number INT COMMENT '座位號',
    join_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '加入時間',
    leave_time TIMESTAMP NULL COMMENT '離開時間',
    initial_chips DECIMAL(15,2) DEFAULT 0.00 COMMENT '初始籌碼',
    final_chips DECIMAL(15,2) DEFAULT 0.00 COMMENT '最終籌碼',
    total_bet DECIMAL(15,2) DEFAULT 0.00 COMMENT '總下注金額',
    total_win DECIMAL(15,2) DEFAULT 0.00 COMMENT '總贏得金額',
    net_result DECIMAL(15,2) GENERATED ALWAYS AS (total_win - total_bet) COMMENT '淨結果',
    status ENUM('playing', 'finished', 'left') DEFAULT 'playing' COMMENT '參與狀態',
    UNIQUE KEY unique_session_player (session_id, player_id),
    INDEX idx_session_id (session_id),
    INDEX idx_player_id (player_id),
    INDEX idx_join_time (join_time),
    FOREIGN KEY (session_id) REFERENCES game_sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES players(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='遊戲參與記錄表';

-- 建立遊戲賠率表
CREATE TABLE IF NOT EXISTS game_odds (
    id INT AUTO_INCREMENT PRIMARY KEY,
    game_id INT NOT NULL COMMENT '遊戲ID',
    bet_type VARCHAR(50) NOT NULL COMMENT '下注類型',
    odds_value DECIMAL(10,4) NOT NULL COMMENT '賠率值',
    min_bet DECIMAL(10,2) DEFAULT 1.00 COMMENT '最低下注',
    max_bet DECIMAL(10,2) DEFAULT 1000.00 COMMENT '最高下注',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否啟用',
    effective_from TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '生效時間',
    effective_to TIMESTAMP NULL COMMENT '失效時間',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_game_id (game_id),
    INDEX idx_bet_type (bet_type),
    INDEX idx_is_active (is_active),
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='遊戲賠率表';

-- 插入預設遊戲資料
INSERT INTO games (game_code, name, name_en, description, game_type, min_bet, max_bet, house_edge, rtp_rate, status) VALUES 
('texas_holdem', '德州撲克', 'Texas Holdem', '經典的德州撲克遊戲', 'texas_holdem', 5.00, 1000.00, 0.0250, 0.9750, 'active'),
('stud_poker', '梭哈撲克', 'Five Card Stud', '傳統五張牌梭哈遊戲', 'stud_poker', 2.00, 500.00, 0.0300, 0.9700, 'active'),
('baccarat', '百家樂', 'Baccarat', '經典百家樂遊戲', 'baccarat', 10.00, 5000.00, 0.0106, 0.9894, 'active')
ON DUPLICATE KEY UPDATE name=name;

-- 插入預設遊戲配置
INSERT INTO game_configs (game_id, config_key, config_value, description) VALUES 
(1, 'blind_structure', '{"small_blind": 5, "big_blind": 10}', '德州撲克盲注結構'),
(1, 'max_rounds', '{"value": 50}', '最大回合數'),
(1, 'ai_strategy', '{"aggression": 0.6, "bluff_rate": 0.15}', 'AI策略參數'),
(2, 'ante_amount', '{"value": 2}', '梭哈底注金額'),
(2, 'max_rounds', '{"value": 5}', '最大回合數'),
(3, 'commission_rate', '{"banker": 0.05, "player": 0.0}', '百家樂抽水率'),
(3, 'min_cards', '{"value": 6}', '最少發牌數')
ON DUPLICATE KEY UPDATE config_key=config_key;

-- 插入預設賠率
INSERT INTO game_odds (game_id, bet_type, odds_value, min_bet, max_bet) VALUES 
-- 百家樂賠率
(3, 'banker', 1.95, 10.00, 5000.00),
(3, 'player', 2.00, 10.00, 5000.00),
(3, 'tie', 9.00, 10.00, 500.00),
(3, 'banker_pair', 12.00, 5.00, 200.00),
(3, 'player_pair', 12.00, 5.00, 200.00)
ON DUPLICATE KEY UPDATE bet_type=bet_type; 