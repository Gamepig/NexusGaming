package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// PlayerStatus 玩家狀態枚舉
type PlayerStatus string

const (
	PlayerStatusActive    PlayerStatus = "active"
	PlayerStatusInactive  PlayerStatus = "inactive"
	PlayerStatusSuspended PlayerStatus = "suspended"
	PlayerStatusDeleted   PlayerStatus = "deleted"
)

// VerificationLevel 驗證等級枚舉
type VerificationLevel string

const (
	VerificationLevelNone     VerificationLevel = "none"
	VerificationLevelEmail    VerificationLevel = "email"
	VerificationLevelPhone    VerificationLevel = "phone"
	VerificationLevelIdentity VerificationLevel = "identity"
)

// RiskLevel 風險等級枚舉
type RiskLevel string

const (
	RiskLevelLow       RiskLevel = "low"
	RiskLevelMedium    RiskLevel = "medium"
	RiskLevelHigh      RiskLevel = "high"
	RiskLevelBlacklist RiskLevel = "blacklist"
)

// Gender 性別枚舉
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// Player 玩家基本資料模型
type Player struct {
	ID                int64             `json:"id" db:"id"`
	PlayerID          string            `json:"player_id" db:"player_id"`                       // 玩家唯一ID
	Username          string            `json:"username" db:"username"`                         // 玩家帳號
	Email             *string           `json:"email,omitempty" db:"email"`                     // 電子郵件
	Phone             *string           `json:"phone,omitempty" db:"phone"`                     // 電話號碼
	RealName          *string           `json:"real_name,omitempty" db:"real_name"`             // 真實姓名
	Nickname          *string           `json:"nickname,omitempty" db:"nickname"`               // 暱稱
	AvatarURL         *string           `json:"avatar_url,omitempty" db:"avatar_url"`           // 頭像URL
	BirthDate         *time.Time        `json:"birth_date,omitempty" db:"birth_date"`           // 生日
	Gender            *Gender           `json:"gender,omitempty" db:"gender"`                   // 性別
	Country           *string           `json:"country,omitempty" db:"country"`                 // 國家
	Language          string            `json:"language" db:"language"`                         // 偏好語言
	Timezone          string            `json:"timezone" db:"timezone"`                         // 時區
	Status            PlayerStatus      `json:"status" db:"status"`                             // 帳戶狀態
	VerificationLevel VerificationLevel `json:"verification_level" db:"verification_level"`     // 驗證等級
	RiskLevel         RiskLevel         `json:"risk_level" db:"risk_level"`                     // 風險等級
	VIPLevel          int               `json:"vip_level" db:"vip_level"`                       // VIP等級
	ReferrerID        *int64            `json:"referrer_id,omitempty" db:"referrer_id"`         // 推薦人ID
	AgentID           *int              `json:"agent_id,omitempty" db:"agent_id"`               // 所屬代理商ID
	DealerID          *int              `json:"dealer_id,omitempty" db:"dealer_id"`             // 所屬經銷商ID
	RegistrationIP    *string           `json:"registration_ip,omitempty" db:"registration_ip"` // 註冊IP
	LastLoginIP       *string           `json:"last_login_ip,omitempty" db:"last_login_ip"`     // 最後登入IP
	LastLoginAt       *time.Time        `json:"last_login_at,omitempty" db:"last_login_at"`     // 最後登入時間
	LoginCount        int               `json:"login_count" db:"login_count"`                   // 登入次數
	TotalDeposit      float64           `json:"total_deposit" db:"total_deposit"`               // 總儲值金額
	TotalWithdraw     float64           `json:"total_withdraw" db:"total_withdraw"`             // 總提領金額
	TotalBet          float64           `json:"total_bet" db:"total_bet"`                       // 總下注金額
	TotalWin          float64           `json:"total_win" db:"total_win"`                       // 總贏得金額
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`                     // 建立時間
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`                     // 更新時間
	DeletedAt         *time.Time        `json:"deleted_at,omitempty" db:"deleted_at"`           // 刪除時間（軟刪除）
}

// PlayerWallet 玩家錢包模型
type PlayerWallet struct {
	ID            int64     `json:"id" db:"id"`
	PlayerID      int64     `json:"player_id" db:"player_id"`           // 玩家ID
	Currency      string    `json:"currency" db:"currency"`             // 幣別
	Balance       float64   `json:"balance" db:"balance"`               // 可用餘額
	FrozenBalance float64   `json:"frozen_balance" db:"frozen_balance"` // 凍結餘額
	TotalBalance  float64   `json:"total_balance" db:"total_balance"`   // 總餘額（computed）
	CreatedAt     time.Time `json:"created_at" db:"created_at"`         // 建立時間
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`         // 更新時間
}

// PlayerStatusHistory 玩家狀態歷史模型
type PlayerStatusHistory struct {
	ID         int64     `json:"id" db:"id"`
	PlayerID   int64     `json:"player_id" db:"player_id"`               // 玩家ID
	OldStatus  *string   `json:"old_status,omitempty" db:"old_status"`   // 原狀態
	NewStatus  string    `json:"new_status" db:"new_status"`             // 新狀態
	Reason     *string   `json:"reason,omitempty" db:"reason"`           // 變更原因
	OperatorID *int      `json:"operator_id,omitempty" db:"operator_id"` // 操作者ID
	CreatedAt  time.Time `json:"created_at" db:"created_at"`             // 建立時間
}

// PlayerTag 玩家標籤模型
type PlayerTag struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`                         // 標籤名稱
	Description *string   `json:"description,omitempty" db:"description"` // 標籤描述
	Color       string    `json:"color" db:"color"`                       // 標籤顏色
	CreatedAt   time.Time `json:"created_at" db:"created_at"`             // 建立時間
}

// PlayerTagRelation 玩家標籤關聯模型
type PlayerTagRelation struct {
	ID         int64     `json:"id" db:"id"`
	PlayerID   int64     `json:"player_id" db:"player_id"`               // 玩家ID
	TagID      int       `json:"tag_id" db:"tag_id"`                     // 標籤ID
	AssignedBy *int      `json:"assigned_by,omitempty" db:"assigned_by"` // 分配者ID
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`           // 分配時間

	// 關聯資料
	Tag *PlayerTag `json:"tag,omitempty"`
}

// RestrictionType 限制類型枚舉
type RestrictionType string

const (
	RestrictionTypeBetLimit     RestrictionType = "bet_limit"
	RestrictionTypeDepositLimit RestrictionType = "deposit_limit"
	RestrictionTypeGameAccess   RestrictionType = "game_access"
	RestrictionTypeTimeLimit    RestrictionType = "time_limit"
)

// RestrictionValue 限制值結構體（用於 JSON 存儲）
type RestrictionValue map[string]interface{}

// Scan 實現 sql.Scanner 接口
func (rv *RestrictionValue) Scan(value interface{}) error {
	if value == nil {
		*rv = make(RestrictionValue)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, rv)
	case string:
		return json.Unmarshal([]byte(v), rv)
	default:
		return errors.New("cannot scan RestrictionValue")
	}
}

// Value 實現 driver.Valuer 接口
func (rv RestrictionValue) Value() (driver.Value, error) {
	if rv == nil {
		return nil, nil
	}
	return json.Marshal(rv)
}

// PlayerRestriction 玩家限制設定模型
type PlayerRestriction struct {
	ID               int64            `json:"id" db:"id"`
	PlayerID         int64            `json:"player_id" db:"player_id"`                 // 玩家ID
	RestrictionType  RestrictionType  `json:"restriction_type" db:"restriction_type"`   // 限制類型
	RestrictionValue RestrictionValue `json:"restriction_value" db:"restriction_value"` // 限制值
	StartTime        time.Time        `json:"start_time" db:"start_time"`               // 生效時間
	EndTime          *time.Time       `json:"end_time,omitempty" db:"end_time"`         // 結束時間
	IsActive         bool             `json:"is_active" db:"is_active"`                 // 是否啟用
	Reason           *string          `json:"reason,omitempty" db:"reason"`             // 限制原因
	OperatorID       *int             `json:"operator_id,omitempty" db:"operator_id"`   // 設定者ID
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`               // 建立時間
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`               // 更新時間
}

// PlayerWithDetails 玩家詳細資料（包含關聯資料）
type PlayerWithDetails struct {
	Player
	Wallet       *PlayerWallet       `json:"wallet,omitempty"`
	Tags         []PlayerTag         `json:"tags,omitempty"`
	Restrictions []PlayerRestriction `json:"restrictions,omitempty"`
}

// GetTotalBalance 計算總餘額
func (pw *PlayerWallet) GetTotalBalance() float64 {
	return pw.Balance + pw.FrozenBalance
}

// IsVIP 檢查是否為VIP玩家
func (p *Player) IsVIP() bool {
	return p.VIPLevel > 0
}

// IsActive 檢查玩家是否為啟用狀態
func (p *Player) IsActive() bool {
	return p.Status == PlayerStatusActive
}

// CanLogin 檢查玩家是否可以登入
func (p *Player) CanLogin() bool {
	return p.Status == PlayerStatusActive && p.DeletedAt == nil
}

// GetNetProfit 計算淨獲利（贏得金額 - 下注金額）
func (p *Player) GetNetProfit() float64 {
	return p.TotalWin - p.TotalBet
}

// GetWinRate 計算勝率（需要額外的遊戲統計數據）
func (p *Player) GetWinRate() float64 {
	if p.TotalBet == 0 {
		return 0
	}
	return (p.TotalWin / p.TotalBet) * 100
}

// GameSessionStatus 遊戲場次狀態枚舉
type GameSessionStatus string

const (
	GameSessionStatusWaiting   GameSessionStatus = "waiting"
	GameSessionStatusPlaying   GameSessionStatus = "playing"
	GameSessionStatusFinished  GameSessionStatus = "finished"
	GameSessionStatusCancelled GameSessionStatus = "cancelled"
)

// ParticipationStatus 參與狀態枚舉
type ParticipationStatus string

const (
	ParticipationStatusPlaying  ParticipationStatus = "playing"
	ParticipationStatusFinished ParticipationStatus = "finished"
	ParticipationStatusLeft     ParticipationStatus = "left"
)

// GameSession 遊戲場次模型（簡化版，用於玩家歷史）
type GameSession struct {
	ID              int64             `json:"id" db:"id"`
	SessionCode     string            `json:"session_code" db:"session_code"`         // 場次代碼
	RoomID          int64             `json:"room_id" db:"room_id"`                   // 房間ID
	GameID          int               `json:"game_id" db:"game_id"`                   // 遊戲ID
	SessionType     string            `json:"session_type" db:"session_type"`         // 場次類型
	Status          GameSessionStatus `json:"status" db:"status"`                     // 場次狀態
	MaxPlayers      int               `json:"max_players" db:"max_players"`           // 最大玩家數
	CurrentPlayers  int               `json:"current_players" db:"current_players"`   // 當前玩家數
	MinBet          float64           `json:"min_bet" db:"min_bet"`                   // 最低下注
	MaxBet          float64           `json:"max_bet" db:"max_bet"`                   // 最高下注
	TotalPot        float64           `json:"total_pot" db:"total_pot"`               // 總獎池
	HouseCommission float64           `json:"house_commission" db:"house_commission"` // 抽水金額
	StartedAt       *time.Time        `json:"started_at,omitempty" db:"started_at"`   // 開始時間
	FinishedAt      *time.Time        `json:"finished_at,omitempty" db:"finished_at"` // 結束時間
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`             // 建立時間
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`             // 更新時間
}

// GameParticipation 遊戲參與記錄模型
type GameParticipation struct {
	ID           int64               `json:"id" db:"id"`
	SessionID    int64               `json:"session_id" db:"session_id"`             // 場次ID
	PlayerID     int64               `json:"player_id" db:"player_id"`               // 玩家ID
	SeatNumber   *int                `json:"seat_number,omitempty" db:"seat_number"` // 座位號
	JoinTime     time.Time           `json:"join_time" db:"join_time"`               // 加入時間
	LeaveTime    *time.Time          `json:"leave_time,omitempty" db:"leave_time"`   // 離開時間
	InitialChips float64             `json:"initial_chips" db:"initial_chips"`       // 初始籌碼
	FinalChips   float64             `json:"final_chips" db:"final_chips"`           // 最終籌碼
	TotalBet     float64             `json:"total_bet" db:"total_bet"`               // 總下注金額
	TotalWin     float64             `json:"total_win" db:"total_win"`               // 總贏得金額
	NetResult    float64             `json:"net_result" db:"net_result"`             // 淨結果
	Status       ParticipationStatus `json:"status" db:"status"`                     // 參與狀態

	// 關聯資料
	Session *GameSession `json:"session,omitempty"`
}

// PlayerGameHistory 玩家遊戲歷史（包含詳細資訊）
type PlayerGameHistory struct {
	Participation GameParticipation `json:"participation"`
	GameName      string            `json:"game_name"`          // 遊戲名稱
	GameType      string            `json:"game_type"`          // 遊戲類型
	RoomName      string            `json:"room_name"`          // 房間名稱
	Duration      *time.Duration    `json:"duration,omitempty"` // 遊戲時長
}

// PlayerGameStatistics 玩家遊戲統計
type PlayerGameStatistics struct {
	PlayerID         int64         `json:"player_id"`
	GameID           *int          `json:"game_id,omitempty"`            // 特定遊戲ID（nil表示全部遊戲）
	GameType         *string       `json:"game_type,omitempty"`          // 特定遊戲類型
	TotalSessions    int           `json:"total_sessions"`               // 總場次數
	TotalPlayTime    time.Duration `json:"total_play_time"`              // 總遊戲時間
	TotalBet         float64       `json:"total_bet"`                    // 總下注金額
	TotalWin         float64       `json:"total_win"`                    // 總贏得金額
	NetProfit        float64       `json:"net_profit"`                   // 淨獲利
	WinRate          float64       `json:"win_rate"`                     // 勝率（%）
	AvgBetPerSession float64       `json:"avg_bet_per_session"`          // 平均每場下注
	AvgWinPerSession float64       `json:"avg_win_per_session"`          // 平均每場贏得
	BiggestWin       float64       `json:"biggest_win"`                  // 最大單次獲利
	BiggestLoss      float64       `json:"biggest_loss"`                 // 最大單次損失
	WinSessions      int           `json:"win_sessions"`                 // 獲利場次數
	LossSessions     int           `json:"loss_sessions"`                // 虧損場次數
	LastPlayTime     *time.Time    `json:"last_play_time,omitempty"`     // 最後遊戲時間
	FavoriteGameType *string       `json:"favorite_game_type,omitempty"` // 最喜愛的遊戲類型
	PeriodType       string        `json:"period_type"`                  // 統計週期（daily, weekly, monthly, all）
	PeriodStart      *time.Time    `json:"period_start,omitempty"`       // 統計開始時間
	PeriodEnd        *time.Time    `json:"period_end,omitempty"`         // 統計結束時間
}

// GetDuration 計算遊戲時長
func (gp *GameParticipation) GetDuration() *time.Duration {
	if gp.LeaveTime == nil || gp.JoinTime.IsZero() {
		return nil
	}

	duration := gp.LeaveTime.Sub(gp.JoinTime)
	return &duration
}

// IsWinner 判斷是否為獲利
func (gp *GameParticipation) IsWinner() bool {
	return gp.NetResult > 0
}

// GetProfitMargin 計算獲利率
func (gp *GameParticipation) GetProfitMargin() float64 {
	if gp.TotalBet == 0 {
		return 0
	}
	return (gp.NetResult / gp.TotalBet) * 100
}

// CalculateWinRate 計算勝率
func (pgs *PlayerGameStatistics) CalculateWinRate() {
	if pgs.TotalSessions == 0 {
		pgs.WinRate = 0
		return
	}
	pgs.WinRate = (float64(pgs.WinSessions) / float64(pgs.TotalSessions)) * 100
}

// CalculateAverages 計算平均值
func (pgs *PlayerGameStatistics) CalculateAverages() {
	if pgs.TotalSessions == 0 {
		pgs.AvgBetPerSession = 0
		pgs.AvgWinPerSession = 0
		return
	}

	pgs.AvgBetPerSession = pgs.TotalBet / float64(pgs.TotalSessions)
	pgs.AvgWinPerSession = pgs.TotalWin / float64(pgs.TotalSessions)
}

// GetROI 計算投資回報率
func (pgs *PlayerGameStatistics) GetROI() float64 {
	if pgs.TotalBet == 0 {
		return 0
	}
	return (pgs.NetProfit / pgs.TotalBet) * 100
}
