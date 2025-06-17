-- 財務管理相關表結構
-- 建立時間: 2024-12-19

USE nexus_gaming;

-- 建立交易記錄表
CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL UNIQUE COMMENT '交易流水號',
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    transaction_type ENUM('deposit', 'withdrawal', 'bet', 'win', 'bonus', 'commission', 'refund', 'adjustment') NOT NULL COMMENT '交易類型',
    amount DECIMAL(15,2) NOT NULL COMMENT '交易金額',
    currency VARCHAR(10) DEFAULT 'TWD' COMMENT '幣別',
    balance_before DECIMAL(15,2) NOT NULL COMMENT '交易前餘額',
    balance_after DECIMAL(15,2) NOT NULL COMMENT '交易後餘額',
    status ENUM('pending', 'completed', 'failed', 'cancelled', 'processing') DEFAULT 'pending' COMMENT '交易狀態',
    payment_method VARCHAR(50) COMMENT '付款方式',
    reference_id VARCHAR(100) COMMENT '參考編號（遊戲場次、訂單等）',
    reference_type VARCHAR(50) COMMENT '參考類型',
    description TEXT COMMENT '交易描述',
    metadata JSON COMMENT '額外資料',
    processor_fee DECIMAL(10,2) DEFAULT 0.00 COMMENT '手續費',
    processed_at TIMESTAMP NULL COMMENT '處理時間',
    operator_id INT COMMENT '操作員ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_transaction_id (transaction_id),
    INDEX idx_player_id (player_id),
    INDEX idx_transaction_type (transaction_type),
    INDEX idx_status (status),
    INDEX idx_reference_id (reference_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (player_id) REFERENCES players(id),
    FOREIGN KEY (operator_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易記錄表';

-- 建立儲值記錄表
CREATE TABLE IF NOT EXISTS deposits (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    deposit_id VARCHAR(64) NOT NULL UNIQUE COMMENT '儲值訂單號',
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    amount DECIMAL(15,2) NOT NULL COMMENT '儲值金額',
    currency VARCHAR(10) DEFAULT 'TWD' COMMENT '幣別',
    payment_method ENUM('credit_card', 'debit_card', 'bank_transfer', 'e_wallet', 'cryptocurrency', 'voucher') NOT NULL COMMENT '付款方式',
    payment_provider VARCHAR(50) COMMENT '支付提供商',
    payment_account VARCHAR(100) COMMENT '付款帳戶',
    status ENUM('pending', 'processing', 'completed', 'failed', 'cancelled', 'expired') DEFAULT 'pending' COMMENT '儲值狀態',
    gateway_transaction_id VARCHAR(100) COMMENT '支付閘道交易ID',
    gateway_response JSON COMMENT '支付閘道回應',
    bonus_amount DECIMAL(10,2) DEFAULT 0.00 COMMENT '贈送金額',
    bonus_type VARCHAR(50) COMMENT '贈送類型',
    processor_fee DECIMAL(10,2) DEFAULT 0.00 COMMENT '手續費',
    net_amount DECIMAL(15,2) GENERATED ALWAYS AS (amount - processor_fee) COMMENT '淨額',
    processed_at TIMESTAMP NULL COMMENT '處理完成時間',
    expired_at TIMESTAMP NULL COMMENT '過期時間',
    notes TEXT COMMENT '備註',
    ip_address VARCHAR(45) COMMENT '請求IP',
    user_agent TEXT COMMENT '用戶代理',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_deposit_id (deposit_id),
    INDEX idx_player_id (player_id),
    INDEX idx_status (status),
    INDEX idx_payment_method (payment_method),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (player_id) REFERENCES players(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='儲值記錄表';

-- 建立提領記錄表
CREATE TABLE IF NOT EXISTS withdrawals (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    withdrawal_id VARCHAR(64) NOT NULL UNIQUE COMMENT '提領訂單號',
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    amount DECIMAL(15,2) NOT NULL COMMENT '提領金額',
    currency VARCHAR(10) DEFAULT 'TWD' COMMENT '幣別',
    withdrawal_method ENUM('bank_transfer', 'e_wallet', 'cryptocurrency', 'check') NOT NULL COMMENT '提領方式',
    account_info JSON NOT NULL COMMENT '提領帳戶資訊',
    status ENUM('pending', 'reviewing', 'approved', 'processing', 'completed', 'rejected', 'cancelled') DEFAULT 'pending' COMMENT '提領狀態',
    review_status ENUM('auto_approved', 'manual_review', 'risk_review', 'rejected') COMMENT '審核狀態',
    reviewer_id INT COMMENT '審核員ID',
    reviewed_at TIMESTAMP NULL COMMENT '審核時間',
    review_notes TEXT COMMENT '審核備註',
    processor_fee DECIMAL(10,2) DEFAULT 0.00 COMMENT '手續費',
    net_amount DECIMAL(15,2) GENERATED ALWAYS AS (amount - processor_fee) COMMENT '實際到帳金額',
    processed_at TIMESTAMP NULL COMMENT '處理完成時間',
    gateway_transaction_id VARCHAR(100) COMMENT '支付閘道交易ID',
    gateway_response JSON COMMENT '支付閘道回應',
    risk_score DECIMAL(5,2) DEFAULT 0.00 COMMENT '風險評分',
    ip_address VARCHAR(45) COMMENT '請求IP',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_withdrawal_id (withdrawal_id),
    INDEX idx_player_id (player_id),
    INDEX idx_status (status),
    INDEX idx_review_status (review_status),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (player_id) REFERENCES players(id),
    FOREIGN KEY (reviewer_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='提領記錄表';

-- 建立支付方式表
CREATE TABLE IF NOT EXISTS payment_methods (
    id INT AUTO_INCREMENT PRIMARY KEY,
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    method_type ENUM('credit_card', 'debit_card', 'bank_account', 'e_wallet', 'cryptocurrency') NOT NULL COMMENT '方式類型',
    provider_name VARCHAR(50) NOT NULL COMMENT '提供商名稱',
    account_info JSON NOT NULL COMMENT '帳戶資訊（加密存儲）',
    display_name VARCHAR(100) COMMENT '顯示名稱',
    is_verified BOOLEAN DEFAULT FALSE COMMENT '是否已驗證',
    is_default BOOLEAN DEFAULT FALSE COMMENT '是否預設',
    status ENUM('active', 'inactive', 'blocked') DEFAULT 'active' COMMENT '狀態',
    last_used_at TIMESTAMP NULL COMMENT '最後使用時間',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_player_id (player_id),
    INDEX idx_method_type (method_type),
    INDEX idx_status (status),
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付方式表';

-- 建立獎金記錄表
CREATE TABLE IF NOT EXISTS bonuses (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    bonus_id VARCHAR(64) NOT NULL UNIQUE COMMENT '獎金ID',
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    bonus_type ENUM('welcome', 'deposit', 'cashback', 'referral', 'loyalty', 'promotion', 'manual') NOT NULL COMMENT '獎金類型',
    amount DECIMAL(15,2) NOT NULL COMMENT '獎金金額',
    currency VARCHAR(10) DEFAULT 'TWD' COMMENT '幣別',
    wagering_requirement DECIMAL(5,2) DEFAULT 0.00 COMMENT '流水要求倍數',
    wagered_amount DECIMAL(15,2) DEFAULT 0.00 COMMENT '已完成流水',
    remaining_wagering DECIMAL(15,2) GENERATED ALWAYS AS (amount * wagering_requirement - wagered_amount) COMMENT '剩餘流水要求',
    status ENUM('pending', 'active', 'completed', 'forfeited', 'expired') DEFAULT 'pending' COMMENT '獎金狀態',
    source_reference VARCHAR(100) COMMENT '來源參考（儲值ID、推薦等）',
    valid_from TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '有效開始時間',
    valid_to TIMESTAMP NULL COMMENT '有效結束時間',
    completed_at TIMESTAMP NULL COMMENT '完成時間',
    description TEXT COMMENT '獎金描述',
    terms_conditions JSON COMMENT '條款條件',
    granted_by INT COMMENT '發放者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_bonus_id (bonus_id),
    INDEX idx_player_id (player_id),
    INDEX idx_bonus_type (bonus_type),
    INDEX idx_status (status),
    INDEX idx_valid_to (valid_to),
    FOREIGN KEY (player_id) REFERENCES players(id),
    FOREIGN KEY (granted_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='獎金記錄表';

-- 建立財務報表快照表
CREATE TABLE IF NOT EXISTS financial_reports (
    id INT AUTO_INCREMENT PRIMARY KEY,
    report_date DATE NOT NULL COMMENT '報表日期',
    report_type ENUM('daily', 'weekly', 'monthly') NOT NULL COMMENT '報表類型',
    total_deposits DECIMAL(20,2) DEFAULT 0.00 COMMENT '總儲值',
    total_withdrawals DECIMAL(20,2) DEFAULT 0.00 COMMENT '總提領',
    total_bets DECIMAL(20,2) DEFAULT 0.00 COMMENT '總下注',
    total_wins DECIMAL(20,2) DEFAULT 0.00 COMMENT '總贏得',
    house_profit DECIMAL(20,2) DEFAULT 0.00 COMMENT '莊家利潤',
    bonus_paid DECIMAL(20,2) DEFAULT 0.00 COMMENT '已支付獎金',
    commission_earned DECIMAL(20,2) DEFAULT 0.00 COMMENT '佣金收入',
    active_players INT DEFAULT 0 COMMENT '活躍玩家數',
    new_players INT DEFAULT 0 COMMENT '新玩家數',
    currency VARCHAR(10) DEFAULT 'TWD' COMMENT '幣別',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_report (report_date, report_type, currency)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='財務報表快照表';

-- 建立佣金記錄表
CREATE TABLE IF NOT EXISTS commissions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    commission_id VARCHAR(64) NOT NULL UNIQUE COMMENT '佣金ID',
    agent_id INT NOT NULL COMMENT '代理商ID',
    player_id BIGINT NOT NULL COMMENT '玩家ID',
    source_transaction_id BIGINT NOT NULL COMMENT '源交易ID',
    commission_type ENUM('bet_commission', 'loss_commission', 'deposit_commission') NOT NULL COMMENT '佣金類型',
    rate DECIMAL(5,4) NOT NULL COMMENT '佣金比率',
    base_amount DECIMAL(15,2) NOT NULL COMMENT '基礎金額',
    commission_amount DECIMAL(15,2) NOT NULL COMMENT '佣金金額',
    currency VARCHAR(10) DEFAULT 'TWD' COMMENT '幣別',
    status ENUM('pending', 'confirmed', 'paid', 'cancelled') DEFAULT 'pending' COMMENT '佣金狀態',
    period_start DATE NOT NULL COMMENT '結算期間開始',
    period_end DATE NOT NULL COMMENT '結算期間結束',
    paid_at TIMESTAMP NULL COMMENT '支付時間',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_commission_id (commission_id),
    INDEX idx_agent_id (agent_id),
    INDEX idx_player_id (player_id),
    INDEX idx_status (status),
    INDEX idx_period (period_start, period_end),
    FOREIGN KEY (agent_id) REFERENCES users(id),
    FOREIGN KEY (player_id) REFERENCES players(id),
    FOREIGN KEY (source_transaction_id) REFERENCES transactions(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='佣金記錄表'; 