-- 玩家分析系統相關表結構（修正版）
-- 建立時間: 2025-01-18

USE nexus_gaming;

-- 檢查並刪除可能已存在的分析表（小心使用）
DROP TABLE IF EXISTS player_game_sessions;
DROP TABLE IF EXISTS player_behavior_analysis;
DROP TABLE IF EXISTS player_game_preference_analysis;
DROP TABLE IF EXISTS player_value_score_analysis;
DROP TABLE IF EXISTS player_spending_habits_analysis;

-- 建立玩家遊戲會話分析表（重命名以避免衝突）
CREATE TABLE IF NOT EXISTS player_game_sessions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    game_id INT NOT NULL COMMENT '遊戲ID',
    game_type VARCHAR(50) NOT NULL COMMENT '遊戲類型',
    session_id VARCHAR(64) COMMENT '會話ID',
    start_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '開始時間',
    end_time TIMESTAMP NULL COMMENT '結束時間',
    session_duration INT NULL COMMENT '遊戲時長（秒）',
    total_bets INT DEFAULT 0 COMMENT '總下注次數',
    total_bet_amount DECIMAL(15,2) DEFAULT 0.00 COMMENT '總下注金額',
    total_win_amount DECIMAL(15,2) DEFAULT 0.00 COMMENT '總贏得金額',
    net_result DECIMAL(15,2) GENERATED ALWAYS AS (total_win_amount - total_bet_amount) COMMENT '淨結果',
    max_bet_amount DECIMAL(15,2) DEFAULT 0.00 COMMENT '最大單次下注',
    min_bet_amount DECIMAL(15,2) DEFAULT 0.00 COMMENT '最小單次下注',
    device_type ENUM('desktop', 'mobile', 'tablet') COMMENT '設備類型',
    is_completed BOOLEAN DEFAULT FALSE COMMENT '是否完成',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_player_id (player_id),
    INDEX idx_game_id (game_id),
    INDEX idx_game_type (game_type),
    INDEX idx_start_time (start_time),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家遊戲會話分析表';

-- 建立玩家行為分析結果表
CREATE TABLE IF NOT EXISTS player_behavior_analysis (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    analysis_date DATE NOT NULL COMMENT '分析日期',
    behavior_score DECIMAL(5,2) DEFAULT 0.00 COMMENT '行為評分',
    gaming_frequency JSON COMMENT '遊戲頻率分析',
    betting_pattern JSON COMMENT '下注模式分析',
    time_preference JSON COMMENT '時間偏好分析',
    session_behavior JSON COMMENT '會話行為分析',
    risk_profile VARCHAR(20) COMMENT '風險檔案',
    recommendations JSON COMMENT '建議',
    analysis_version VARCHAR(10) DEFAULT '1.0' COMMENT '分析版本',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_player_date (player_id, analysis_date),
    INDEX idx_player_id (player_id),
    INDEX idx_analysis_date (analysis_date),
    INDEX idx_behavior_score (behavior_score),
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家行為分析結果表';

-- 建立玩家遊戲偏好分析結果表
CREATE TABLE IF NOT EXISTS player_game_preference_analysis (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    analysis_date DATE NOT NULL COMMENT '分析日期',
    time_range VARCHAR(10) NOT NULL COMMENT '分析時間範圍',
    total_games_played INT DEFAULT 0 COMMENT '總遊戲次數',
    unique_game_types INT DEFAULT 0 COMMENT '遊戲類型數量',
    favorite_game_type VARCHAR(50) COMMENT '最喜愛遊戲類型',
    game_type_stats JSON COMMENT '遊戲類型統計',
    time_distribution JSON COMMENT '時間分佈',
    trend_analysis JSON COMMENT '趨勢分析',
    betting_habits JSON COMMENT '下注習慣',
    preference_metrics JSON COMMENT '偏好指標',
    recommendations JSON COMMENT '建議',
    graph_data JSON COMMENT '圖表數據',
    analysis_version VARCHAR(10) DEFAULT '1.0' COMMENT '分析版本',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_player_date_range (player_id, analysis_date, time_range),
    INDEX idx_player_id (player_id),
    INDEX idx_analysis_date (analysis_date),
    INDEX idx_favorite_game_type (favorite_game_type),
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家遊戲偏好分析結果表';

-- 建立玩家價值評分分析結果表
CREATE TABLE IF NOT EXISTS player_value_score_analysis (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    analysis_date DATE NOT NULL COMMENT '分析日期',
    time_range VARCHAR(10) NOT NULL COMMENT '分析時間範圍',
    overall_score DECIMAL(5,2) DEFAULT 0.00 COMMENT '總體評分',
    value_category VARCHAR(20) COMMENT '價值類別',
    activity_score JSON COMMENT '活躍度評分',
    loyalty_score JSON COMMENT '忠誠度評分',
    spending_score JSON COMMENT '消費力評分',
    risk_score JSON COMMENT '風險評分',
    profitability_score JSON COMMENT '盈利性評分',
    trend_analysis JSON COMMENT '趨勢分析',
    competitor_analysis JSON COMMENT '同類玩家比較',
    retention_risk JSON COMMENT '留存風險分析',
    value_potential JSON COMMENT '價值潛力分析',
    recommendations JSON COMMENT '建議',
    weight_config JSON COMMENT '權重配置',
    analysis_version VARCHAR(10) DEFAULT '1.0' COMMENT '分析版本',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_player_date_range (player_id, analysis_date, time_range),
    INDEX idx_player_id (player_id),
    INDEX idx_analysis_date (analysis_date),
    INDEX idx_overall_score (overall_score),
    INDEX idx_value_category (value_category),
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家價值評分分析結果表';

-- 建立玩家消費習慣分析結果表
CREATE TABLE IF NOT EXISTS player_spending_habits_analysis (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    analysis_date DATE NOT NULL COMMENT '分析日期',
    spending_frequency JSON COMMENT '消費頻率分析',
    spending_amount JSON COMMENT '消費金額分析',
    spending_time_pattern JSON COMMENT '消費時間模式',
    spending_channel JSON COMMENT '消費管道分析',
    spending_risk JSON COMMENT '消費風險評估',
    spending_capacity JSON COMMENT '消費能力評估',
    recommended_actions JSON COMMENT '建議行動',
    summary TEXT COMMENT '分析總結',
    analysis_version VARCHAR(10) DEFAULT '1.0' COMMENT '分析版本',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_player_date (player_id, analysis_date),
    INDEX idx_player_id (player_id),
    INDEX idx_analysis_date (analysis_date),
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家消費習慣分析結果表';

-- 插入測試數據
INSERT INTO player_game_sessions (player_id, game_id, game_type, start_time, end_time, session_duration, total_bets, total_bet_amount, total_win_amount, max_bet_amount, min_bet_amount, device_type) VALUES
(1, 1, 'texas_holdem', DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY) + INTERVAL 45 MINUTE, 2700, 25, 2500.00, 1800.00, 200.00, 50.00, 'desktop'),
(1, 2, 'baccarat', DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY) + INTERVAL 30 MINUTE, 1800, 15, 1500.00, 2100.00, 150.00, 50.00, 'mobile'),
(2, 2, 'baccarat', DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY) + INTERVAL 1 HOUR, 3600, 20, 2000.00, 1700.00, 150.00, 50.00, 'mobile'),
(2, 3, 'stud_poker', DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY) + INTERVAL 2 HOUR, 7200, 40, 4000.00, 3800.00, 200.00, 50.00, 'desktop'); 