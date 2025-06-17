-- 代理商/經銷商管理相關表結構
-- 建立時間: 2024-12-19
-- 層級結構: 總公司 -> 代理商 -> 經銷商 -> 玩家

USE nexus_gaming;

-- 建立代理商表
CREATE TABLE IF NOT EXISTS agents (
    id INT AUTO_INCREMENT PRIMARY KEY,
    agent_code VARCHAR(32) NOT NULL UNIQUE COMMENT '代理商編號',
    agent_name VARCHAR(100) NOT NULL COMMENT '代理商名稱',
    contact_person VARCHAR(100) NOT NULL COMMENT '聯絡人',
    email VARCHAR(255) NOT NULL COMMENT '電子郵件',
    phone VARCHAR(20) COMMENT '電話號碼',
    address TEXT COMMENT '地址',
    business_license VARCHAR(100) COMMENT '營業執照號碼',
    tax_id VARCHAR(50) COMMENT '統一編號',
    bank_account JSON COMMENT '銀行帳戶資訊',
    contract_start_date DATE COMMENT '合約開始日期',
    contract_end_date DATE COMMENT '合約結束日期',
    status ENUM('active', 'inactive', 'suspended', 'terminated') DEFAULT 'active' COMMENT '狀態',
    user_id INT NOT NULL COMMENT '關聯的使用者ID',
    parent_agent_id INT NULL COMMENT '上級代理商ID（總公司為NULL）',
    level INT NOT NULL DEFAULT 1 COMMENT '層級（1=總代理商，2=子代理商）',
    commission_rate DECIMAL(5,4) DEFAULT 0.0000 COMMENT '基礎佣金比率',
    max_dealers INT DEFAULT 0 COMMENT '最大經銷商數量限制（0=無限制）',
    current_dealers INT DEFAULT 0 COMMENT '當前經銷商數量',
    total_players BIGINT DEFAULT 0 COMMENT '總玩家數量',
    total_revenue DECIMAL(20,2) DEFAULT 0.00 COMMENT '總營收',
    total_commission DECIMAL(20,2) DEFAULT 0.00 COMMENT '總佣金',
    last_settlement_date DATE COMMENT '最後結算日期',
    notes TEXT COMMENT '備註',
    created_by INT COMMENT '建立者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_agent_code (agent_code),
    INDEX idx_status (status),
    INDEX idx_parent_agent_id (parent_agent_id),
    INDEX idx_level (level),
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (parent_agent_id) REFERENCES agents(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代理商表';

-- 建立經銷商表
CREATE TABLE IF NOT EXISTS dealers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    dealer_code VARCHAR(32) NOT NULL UNIQUE COMMENT '經銷商編號',
    dealer_name VARCHAR(100) NOT NULL COMMENT '經銷商名稱',
    contact_person VARCHAR(100) NOT NULL COMMENT '聯絡人',
    email VARCHAR(255) NOT NULL COMMENT '電子郵件',
    phone VARCHAR(20) COMMENT '電話號碼',
    address TEXT COMMENT '地址',
    id_number VARCHAR(50) COMMENT '身分證號碼',
    bank_account JSON COMMENT '銀行帳戶資訊',
    agent_id INT NOT NULL COMMENT '所屬代理商ID',
    user_id INT NOT NULL COMMENT '關聯的使用者ID',
    status ENUM('active', 'inactive', 'suspended', 'terminated') DEFAULT 'active' COMMENT '狀態',
    commission_rate DECIMAL(5,4) DEFAULT 0.0000 COMMENT '佣金比率',
    max_players INT DEFAULT 0 COMMENT '最大玩家數量限制（0=無限制）',
    current_players INT DEFAULT 0 COMMENT '當前玩家數量',
    total_revenue DECIMAL(20,2) DEFAULT 0.00 COMMENT '總營收',
    total_commission DECIMAL(20,2) DEFAULT 0.00 COMMENT '總佣金',
    last_settlement_date DATE COMMENT '最後結算日期',
    territory VARCHAR(255) COMMENT '負責區域',
    notes TEXT COMMENT '備註',
    created_by INT COMMENT '建立者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_dealer_code (dealer_code),
    INDEX idx_agent_id (agent_id),
    INDEX idx_status (status),
    INDEX idx_user_id (user_id),
    FOREIGN KEY (agent_id) REFERENCES agents(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='經銷商表';

-- 建立層級關係表
CREATE TABLE IF NOT EXISTS agent_hierarchy (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ancestor_id INT NOT NULL COMMENT '祖先節點ID（代理商）',
    descendant_id INT NOT NULL COMMENT '後代節點ID（代理商或經銷商）',
    descendant_type ENUM('agent', 'dealer') NOT NULL COMMENT '後代類型',
    level_difference INT NOT NULL COMMENT '層級差距',
    path_length INT NOT NULL COMMENT '路徑長度',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_hierarchy (ancestor_id, descendant_id, descendant_type),
    INDEX idx_ancestor_id (ancestor_id),
    INDEX idx_descendant (descendant_id, descendant_type),
    INDEX idx_level_difference (level_difference),
    FOREIGN KEY (ancestor_id) REFERENCES agents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代理商層級關係表';

-- 建立分潤配置表
CREATE TABLE IF NOT EXISTS commission_configs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    config_name VARCHAR(100) NOT NULL COMMENT '配置名稱',
    target_type ENUM('agent', 'dealer') NOT NULL COMMENT '目標類型',
    target_id INT NOT NULL COMMENT '目標ID',
    game_type VARCHAR(50) DEFAULT 'all' COMMENT '遊戲類型（all=全部）',
    commission_type ENUM('revenue_share', 'cpa', 'hybrid') DEFAULT 'revenue_share' COMMENT '佣金類型',
    revenue_share_rate DECIMAL(5,4) DEFAULT 0.0000 COMMENT '營收分成比率',
    cpa_amount DECIMAL(10,2) DEFAULT 0.00 COMMENT 'CPA固定金額',
    min_players INT DEFAULT 0 COMMENT '最低玩家數要求',
    min_revenue DECIMAL(15,2) DEFAULT 0.00 COMMENT '最低營收要求',
    tier_config JSON COMMENT '階梯式佣金配置',
    bonus_config JSON COMMENT '獎金配置',
    settlement_period ENUM('daily', 'weekly', 'monthly') DEFAULT 'monthly' COMMENT '結算週期',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否啟用',
    effective_from DATE NOT NULL COMMENT '生效日期',
    effective_to DATE NULL COMMENT '失效日期',
    created_by INT COMMENT '建立者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_target (target_type, target_id),
    INDEX idx_game_type (game_type),
    INDEX idx_is_active (is_active),
    INDEX idx_effective_period (effective_from, effective_to),
    FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分潤配置表';

-- 建立代理商結算表
CREATE TABLE IF NOT EXISTS agent_settlements (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    settlement_id VARCHAR(64) NOT NULL UNIQUE COMMENT '結算ID',
    agent_id INT NOT NULL COMMENT '代理商ID',
    settlement_period_start DATE NOT NULL COMMENT '結算期間開始',
    settlement_period_end DATE NOT NULL COMMENT '結算期間結束',
    total_revenue DECIMAL(20,2) NOT NULL COMMENT '總營收',
    total_bets DECIMAL(20,2) NOT NULL COMMENT '總下注',
    total_wins DECIMAL(20,2) NOT NULL COMMENT '總贏得',
    net_revenue DECIMAL(20,2) NOT NULL COMMENT '淨營收',
    commission_rate DECIMAL(5,4) NOT NULL COMMENT '佣金比率',
    commission_amount DECIMAL(15,2) NOT NULL COMMENT '佣金金額',
    bonus_amount DECIMAL(10,2) DEFAULT 0.00 COMMENT '獎金金額',
    adjustment_amount DECIMAL(10,2) DEFAULT 0.00 COMMENT '調整金額',
    total_payout DECIMAL(15,2) NOT NULL COMMENT '總支付金額',
    player_count INT NOT NULL COMMENT '玩家數量',
    active_player_count INT NOT NULL COMMENT '活躍玩家數量',
    new_player_count INT NOT NULL COMMENT '新玩家數量',
    status ENUM('pending', 'calculated', 'approved', 'paid', 'disputed') DEFAULT 'pending' COMMENT '結算狀態',
    calculation_details JSON COMMENT '計算明細',
    payment_method VARCHAR(50) COMMENT '支付方式',
    payment_reference VARCHAR(100) COMMENT '支付參考號',
    paid_at TIMESTAMP NULL COMMENT '支付時間',
    approved_by INT COMMENT '審核者ID',
    approved_at TIMESTAMP NULL COMMENT '審核時間',
    notes TEXT COMMENT '備註',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_settlement_id (settlement_id),
    INDEX idx_agent_id (agent_id),
    INDEX idx_settlement_period (settlement_period_start, settlement_period_end),
    INDEX idx_status (status),
    FOREIGN KEY (agent_id) REFERENCES agents(id),
    FOREIGN KEY (approved_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代理商結算表';

-- 建立經銷商結算表
CREATE TABLE IF NOT EXISTS dealer_settlements (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    settlement_id VARCHAR(64) NOT NULL UNIQUE COMMENT '結算ID',
    dealer_id INT NOT NULL COMMENT '經銷商ID',
    agent_id INT NOT NULL COMMENT '所屬代理商ID',
    settlement_period_start DATE NOT NULL COMMENT '結算期間開始',
    settlement_period_end DATE NOT NULL COMMENT '結算期間結束',
    total_revenue DECIMAL(20,2) NOT NULL COMMENT '總營收',
    total_bets DECIMAL(20,2) NOT NULL COMMENT '總下注',
    total_wins DECIMAL(20,2) NOT NULL COMMENT '總贏得',
    net_revenue DECIMAL(20,2) NOT NULL COMMENT '淨營收',
    commission_rate DECIMAL(5,4) NOT NULL COMMENT '佣金比率',
    commission_amount DECIMAL(15,2) NOT NULL COMMENT '佣金金額',
    bonus_amount DECIMAL(10,2) DEFAULT 0.00 COMMENT '獎金金額',
    adjustment_amount DECIMAL(10,2) DEFAULT 0.00 COMMENT '調整金額',
    total_payout DECIMAL(15,2) NOT NULL COMMENT '總支付金額',
    player_count INT NOT NULL COMMENT '玩家數量',
    active_player_count INT NOT NULL COMMENT '活躍玩家數量',
    new_player_count INT NOT NULL COMMENT '新玩家數量',
    status ENUM('pending', 'calculated', 'approved', 'paid', 'disputed') DEFAULT 'pending' COMMENT '結算狀態',
    calculation_details JSON COMMENT '計算明細',
    payment_method VARCHAR(50) COMMENT '支付方式',
    payment_reference VARCHAR(100) COMMENT '支付參考號',
    paid_at TIMESTAMP NULL COMMENT '支付時間',
    approved_by INT COMMENT '審核者ID',
    approved_at TIMESTAMP NULL COMMENT '審核時間',
    notes TEXT COMMENT '備註',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_settlement_id (settlement_id),
    INDEX idx_dealer_id (dealer_id),
    INDEX idx_agent_id (agent_id),
    INDEX idx_settlement_period (settlement_period_start, settlement_period_end),
    INDEX idx_status (status),
    FOREIGN KEY (dealer_id) REFERENCES dealers(id),
    FOREIGN KEY (agent_id) REFERENCES agents(id),
    FOREIGN KEY (approved_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='經銷商結算表';

-- 建立代理商業績統計表
CREATE TABLE IF NOT EXISTS agent_performance_stats (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    agent_id INT NOT NULL COMMENT '代理商ID',
    stat_date DATE NOT NULL COMMENT '統計日期',
    total_players INT DEFAULT 0 COMMENT '總玩家數',
    active_players INT DEFAULT 0 COMMENT '活躍玩家數',
    new_players INT DEFAULT 0 COMMENT '新增玩家數',
    total_deposits DECIMAL(15,2) DEFAULT 0.00 COMMENT '總儲值',
    total_withdrawals DECIMAL(15,2) DEFAULT 0.00 COMMENT '總提領',
    total_bets DECIMAL(15,2) DEFAULT 0.00 COMMENT '總下注',
    total_wins DECIMAL(15,2) DEFAULT 0.00 COMMENT '總贏得',
    gross_revenue DECIMAL(15,2) DEFAULT 0.00 COMMENT '總營收',
    net_revenue DECIMAL(15,2) DEFAULT 0.00 COMMENT '淨營收',
    commission_earned DECIMAL(15,2) DEFAULT 0.00 COMMENT '已賺取佣金',
    dealer_count INT DEFAULT 0 COMMENT '經銷商數量',
    avg_bet_per_player DECIMAL(10,2) DEFAULT 0.00 COMMENT '每玩家平均下注',
    player_retention_rate DECIMAL(5,2) DEFAULT 0.00 COMMENT '玩家留存率',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_agent_date (agent_id, stat_date),
    INDEX idx_stat_date (stat_date),
    FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代理商業績統計表';

-- 建立經銷商業績統計表
CREATE TABLE IF NOT EXISTS dealer_performance_stats (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    dealer_id INT NOT NULL COMMENT '經銷商ID',
    agent_id INT NOT NULL COMMENT '所屬代理商ID',
    stat_date DATE NOT NULL COMMENT '統計日期',
    total_players INT DEFAULT 0 COMMENT '總玩家數',
    active_players INT DEFAULT 0 COMMENT '活躍玩家數',
    new_players INT DEFAULT 0 COMMENT '新增玩家數',
    total_deposits DECIMAL(15,2) DEFAULT 0.00 COMMENT '總儲值',
    total_withdrawals DECIMAL(15,2) DEFAULT 0.00 COMMENT '總提領',
    total_bets DECIMAL(15,2) DEFAULT 0.00 COMMENT '總下注',
    total_wins DECIMAL(15,2) DEFAULT 0.00 COMMENT '總贏得',
    gross_revenue DECIMAL(15,2) DEFAULT 0.00 COMMENT '總營收',
    net_revenue DECIMAL(15,2) DEFAULT 0.00 COMMENT '淨營收',
    commission_earned DECIMAL(15,2) DEFAULT 0.00 COMMENT '已賺取佣金',
    avg_bet_per_player DECIMAL(10,2) DEFAULT 0.00 COMMENT '每玩家平均下注',
    player_retention_rate DECIMAL(5,2) DEFAULT 0.00 COMMENT '玩家留存率',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_dealer_date (dealer_id, stat_date),
    INDEX idx_stat_date (stat_date),
    INDEX idx_agent_id (agent_id),
    FOREIGN KEY (dealer_id) REFERENCES dealers(id) ON DELETE CASCADE,
    FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='經銷商業績統計表';

-- 插入預設總公司代理商（管理者角色）
INSERT INTO agents (agent_code, agent_name, contact_person, email, user_id, parent_agent_id, level, status) VALUES 
('COMPANY001', '總公司', '系統管理員', 'admin@nexusgaming.com', 1, NULL, 0, 'active')
ON DUPLICATE KEY UPDATE agent_code=agent_code;

-- 建立觸發器：自動更新代理商的經銷商數量
DELIMITER //
CREATE TRIGGER IF NOT EXISTS update_agent_dealer_count_insert
AFTER INSERT ON dealers
FOR EACH ROW
BEGIN
    UPDATE agents 
    SET current_dealers = current_dealers + 1 
    WHERE id = NEW.agent_id;
END//

CREATE TRIGGER IF NOT EXISTS update_agent_dealer_count_delete
AFTER DELETE ON dealers
FOR EACH ROW
BEGIN
    UPDATE agents 
    SET current_dealers = current_dealers - 1 
    WHERE id = OLD.agent_id;
END//

CREATE TRIGGER IF NOT EXISTS update_agent_dealer_count_update
AFTER UPDATE ON dealers
FOR EACH ROW
BEGIN
    IF OLD.agent_id != NEW.agent_id THEN
        UPDATE agents SET current_dealers = current_dealers - 1 WHERE id = OLD.agent_id;
        UPDATE agents SET current_dealers = current_dealers + 1 WHERE id = NEW.agent_id;
    END IF;
END//

-- 建立觸發器：自動更新經銷商的玩家數量
CREATE TRIGGER IF NOT EXISTS update_dealer_player_count_insert
AFTER INSERT ON players
FOR EACH ROW
BEGIN
    IF NEW.dealer_id IS NOT NULL THEN
        UPDATE dealers 
        SET current_players = current_players + 1 
        WHERE id = NEW.dealer_id;
    END IF;
END//

CREATE TRIGGER IF NOT EXISTS update_dealer_player_count_delete
AFTER DELETE ON players
FOR EACH ROW
BEGIN
    IF OLD.dealer_id IS NOT NULL THEN
        UPDATE dealers 
        SET current_players = current_players - 1 
        WHERE id = OLD.dealer_id;
    END IF;
END//

CREATE TRIGGER IF NOT EXISTS update_dealer_player_count_update
AFTER UPDATE ON players
FOR EACH ROW
BEGIN
    IF OLD.dealer_id != NEW.dealer_id THEN
        IF OLD.dealer_id IS NOT NULL THEN
            UPDATE dealers SET current_players = current_players - 1 WHERE id = OLD.dealer_id;
        END IF;
        IF NEW.dealer_id IS NOT NULL THEN
            UPDATE dealers SET current_players = current_players + 1 WHERE id = NEW.dealer_id;
        END IF;
    END IF;
END//

DELIMITER ; 