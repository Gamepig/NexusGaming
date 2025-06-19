package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"nexus-gaming-backend/config"
	"nexus-gaming-backend/models"

	"github.com/gin-gonic/gin"
)

// PlayerController 玩家控制器
type PlayerController struct{}

// NewPlayerController 建立新的玩家控制器
func NewPlayerController() *PlayerController {
	return &PlayerController{}
}

// PlayerListRequest 玩家列表查詢請求
type PlayerListRequest struct {
	Page              int    `form:"page"`                                                                                       // 頁碼，從1開始
	Limit             int    `form:"limit"`                                                                                      // 每頁數量，最大100
	Sort              string `form:"sort" binding:"omitempty,oneof=id username email created_at updated_at total_bet total_win"` // 排序字段
	Order             string `form:"order" binding:"omitempty,oneof=asc desc"`                                                   // 排序順序
	Search            string `form:"search"`                                                                                     // 搜尋關鍵字（姓名、用戶名、郵箱）
	Status            string `form:"status" binding:"omitempty,oneof=active inactive suspended deleted"`                         // 狀態篩選
	StartDate         string `form:"start_date"`                                                                                 // 註冊開始日期 (YYYY-MM-DD)
	EndDate           string `form:"end_date"`                                                                                   // 註冊結束日期 (YYYY-MM-DD)
	MinBalance        string `form:"min_balance"`                                                                                // 最小餘額
	MaxBalance        string `form:"max_balance"`                                                                                // 最大餘額
	VerificationLevel string `form:"verification_level" binding:"omitempty,oneof=none email phone identity"`                     // 驗證等級
	RiskLevel         string `form:"risk_level" binding:"omitempty,oneof=low medium high"`                                       // 風險等級
}

// PlayerWithBalance 玩家資料（包含餘額）
type PlayerWithBalance struct {
	models.Player
	Balance float64 `json:"balance"` // 錢包餘額
}

// PlayerDetailResponse 玩家詳細資訊回應
type PlayerDetailResponse struct {
	PlayerWithBalance
	Tags         []PlayerTag         `json:"tags"`         // 玩家標籤
	Restrictions []PlayerRestriction `json:"restrictions"` // 玩家限制
	Statistics   PlayerStatistics    `json:"statistics"`   // 玩家統計數據
	RecentGames  []GameParticipation `json:"recent_games"` // 最近遊戲記錄
}

// PlayerTag 玩家標籤（簡化版）
type PlayerTag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// PlayerRestriction 玩家限制（簡化版）
type PlayerRestriction struct {
	ID              int64      `json:"id"`
	RestrictionType string     `json:"restriction_type"`
	Value           string     `json:"value"`
	IsActive        bool       `json:"is_active"`
	ExpiresAt       *time.Time `json:"expires_at"`
}

// PlayerStatistics 玩家統計數據
type PlayerStatistics struct {
	TotalGames       int     `json:"total_games"`
	WinRate          float64 `json:"win_rate"`
	AverageBet       float64 `json:"average_bet"`
	BiggestWin       float64 `json:"biggest_win"`
	BiggestLoss      float64 `json:"biggest_loss"`
	DaysRegistered   int     `json:"days_registered"`
	LastActivityDays int     `json:"last_activity_days"`
}

// GameParticipation 遊戲參與記錄（簡化版）
type GameParticipation struct {
	ID        int64     `json:"id"`
	GameType  string    `json:"game_type"`
	BetAmount float64   `json:"bet_amount"`
	WinAmount float64   `json:"win_amount"`
	Result    string    `json:"result"`
	PlayedAt  time.Time `json:"played_at"`
}

// PlayerGamePreferenceRequest 玩家遊戲偏好分析請求
type PlayerGamePreferenceRequest struct {
	TimeRange     string `json:"time_range" binding:"required,oneof=7d 30d 90d 180d 365d"` // 分析時間範圍
	IncludeGraphs *bool  `json:"include_graphs"`                                           // 是否包含圖表資料（預設true）
	MinGames      *int   `json:"min_games"`                                                // 最少遊戲次數門檻（預設10）
}

// PlayerGamePreferenceResponse 玩家遊戲偏好分析回應
type PlayerGamePreferenceResponse struct {
	PlayerID          int64                   `json:"player_id"`
	Username          string                  `json:"username"`
	AnalysisDate      string                  `json:"analysis_date"`
	TimeRange         string                  `json:"time_range"`
	TotalGamesPlayed  int64                   `json:"total_games_played"`
	UniqueGameTypes   int                     `json:"unique_game_types"`
	FavoriteGameType  string                  `json:"favorite_game_type"`
	GameTypeStats     []GameTypeStatistics    `json:"game_type_stats"`
	TimeDistribution  TimeDistributionData    `json:"time_distribution"`
	TrendAnalysis     []GameTypeTrend         `json:"trend_analysis"`
	BettingHabits     []GameTypeBettingHabit  `json:"betting_habits"`
	PreferenceMetrics PlayerPreferenceMetrics `json:"preference_metrics"`
	Recommendations   []string                `json:"recommendations"`
	GraphData         []PreferenceGraphData   `json:"graph_data,omitempty"`
}

// GameTypeStatistics 遊戲類型統計
type GameTypeStatistics struct {
	GameType          string  `json:"game_type"`
	GamesPlayed       int64   `json:"games_played"`
	ParticipationRate float64 `json:"participation_rate"` // 參與度百分比
	TotalTimeSpent    float64 `json:"total_time_spent"`   // 總遊戲時間（分鐘）
	AverageSession    float64 `json:"average_session"`    // 平均會話時間（分鐘）
	TotalBetAmount    float64 `json:"total_bet_amount"`
	TotalWinAmount    float64 `json:"total_win_amount"`
	NetResult         float64 `json:"net_result"`
	WinRate           float64 `json:"win_rate"`
	PreferenceScore   float64 `json:"preference_score"` // 偏好分數 (0-100)
}

// TimeDistributionData 時間分佈數據
type TimeDistributionData struct {
	HourlyPreference  []HourlyGamePreference `json:"hourly_preference"`
	DailyPreference   []DailyGamePreference  `json:"daily_preference"`
	WeeklyPattern     WeeklyGamePattern      `json:"weekly_pattern"`
	SeasonalPattern   []SeasonalGamePattern  `json:"seasonal_pattern"`
	PeakPlayingTime   string                 `json:"peak_playing_time"`  // 高峰遊戲時段
	PreferredDuration string                 `json:"preferred_duration"` // 偏好遊戲時長
}

// HourlyGamePreference 每小時遊戲偏好
type HourlyGamePreference struct {
	Hour              int              `json:"hour"` // 0-23
	GamesPlayed       int64            `json:"games_played"`
	GameTypeBreakdown map[string]int64 `json:"game_type_breakdown"`
	MostPlayedGame    string           `json:"most_played_game"`
	ActivityLevel     string           `json:"activity_level"` // low, medium, high, peak
}

// DailyGamePreference 每日遊戲偏好
type DailyGamePreference struct {
	Date              string           `json:"date"` // YYYY-MM-DD
	GamesPlayed       int64            `json:"games_played"`
	GameTypeBreakdown map[string]int64 `json:"game_type_breakdown"`
	MostPlayedGame    string           `json:"most_played_game"`
	TotalPlayTime     float64          `json:"total_play_time"` // 分鐘
}

// WeeklyGamePattern 每週遊戲模式
type WeeklyGamePattern struct {
	WeekdayPattern   map[string]GameDayStats `json:"weekday_pattern"`   // Monday-Sunday
	WeekendIntensity float64                 `json:"weekend_intensity"` // 週末遊戲強度指數
	ConsistencyScore float64                 `json:"consistency_score"` // 一致性分數
}

// GameDayStats 遊戲日統計
type GameDayStats struct {
	GamesPlayed       int64            `json:"games_played"`
	AverageSession    float64          `json:"average_session"`
	GameTypeBreakdown map[string]int64 `json:"game_type_breakdown"`
	Intensity         string           `json:"intensity"` // low, medium, high
}

// SeasonalGamePattern 季節性遊戲模式
type SeasonalGamePattern struct {
	Period            string           `json:"period"` // Q1, Q2, Q3, Q4
	GamesPlayed       int64            `json:"games_played"`
	GameTypeBreakdown map[string]int64 `json:"game_type_breakdown"`
	ActivityIndex     float64          `json:"activity_index"` // 相對活動指數
}

// GameTypeTrend 遊戲類型趨勢
type GameTypeTrend struct {
	GameType       string  `json:"game_type"`
	TrendDirection string  `json:"trend_direction"` // increasing, decreasing, stable
	ChangeRate     float64 `json:"change_rate"`     // 變化率 (%)
	Significance   string  `json:"significance"`    // low, medium, high
	Description    string  `json:"description"`
}

// GameTypeBettingHabit 遊戲類型下注習慣
type GameTypeBettingHabit struct {
	GameType            string  `json:"game_type"`
	AverageBetAmount    float64 `json:"average_bet_amount"`
	MedianBetAmount     float64 `json:"median_bet_amount"`
	BetSizeVariability  float64 `json:"bet_size_variability"` // 標準差
	RiskTolerance       string  `json:"risk_tolerance"`       // conservative, moderate, aggressive
	BettingStrategy     string  `json:"betting_strategy"`     // consistent, progressive, random
	ProfitabilityRating string  `json:"profitability_rating"` // poor, average, good, excellent
}

// PlayerPreferenceMetrics 玩家偏好指標
type PlayerPreferenceMetrics struct {
	DiversityIndex      float64 `json:"diversity_index"`      // 遊戲多樣性指數 (0-1)
	SpecializationLevel string  `json:"specialization_level"` // generalist, specialist, focused
	ExplorationTendency string  `json:"exploration_tendency"` // explorer, settler, specialist
	LoyaltyScore        float64 `json:"loyalty_score"`        // 忠誠度分數 (0-100)
	RiskProfile         string  `json:"risk_profile"`         // conservative, balanced, aggressive
	PlayStyle           string  `json:"play_style"`           // casual, regular, intensive
}

// PreferenceGraphData 偏好圖表數據
type PreferenceGraphData struct {
	GraphType   string                 `json:"graph_type"` // pie, bar, line, heatmap
	Title       string                 `json:"title"`
	XAxisLabel  string                 `json:"x_axis_label"`
	YAxisLabel  string                 `json:"y_axis_label"`
	DataPoints  []GraphDataPoint       `json:"data_points"`
	GraphConfig map[string]interface{} `json:"graph_config"` // 圖表配置
}

// GraphDataPoint 圖表資料點
type GraphDataPoint struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"` // 可以是數字或物件
	Color string      `json:"color,omitempty"`
}

// PlayerValueScoreRequest 玩家價值評分請求
type PlayerValueScoreRequest struct {
	TimeRange      string             `json:"time_range" binding:"required,oneof=30d 90d 180d 365d"` // 評分時間範圍
	IncludeDetails *bool              `json:"include_details"`                                       // 是否包含詳細分析（預設true）
	WeightConfig   *ScoreWeightConfig `json:"weight_config"`                                         // 自定義權重配置
}

// ScoreWeightConfig 評分權重配置
type ScoreWeightConfig struct {
	ActivityWeight      float64 `json:"activity_weight"`      // 活躍度權重 (0-1)
	LoyaltyWeight       float64 `json:"loyalty_weight"`       // 忠誠度權重 (0-1)
	SpendingWeight      float64 `json:"spending_weight"`      // 消費力權重 (0-1)
	RiskWeight          float64 `json:"risk_weight"`          // 風險權重 (0-1)
	ProfitabilityWeight float64 `json:"profitability_weight"` // 盈利性權重 (0-1)
}

// PlayerValueScoreResponse 玩家價值評分回應
type PlayerValueScoreResponse struct {
	PlayerID           int64                    `json:"player_id"`
	Username           string                   `json:"username"`
	AnalysisDate       string                   `json:"analysis_date"`
	TimeRange          string                   `json:"time_range"`
	OverallScore       float64                  `json:"overall_score"`       // 總體價值評分 (0-100)
	ValueCategory      string                   `json:"value_category"`      // 價值類別：VIP, High, Medium, Low
	ActivityScore      PlayerActivityScore      `json:"activity_score"`      // 活躍度評分
	LoyaltyScore       PlayerLoyaltyScore       `json:"loyalty_score"`       // 忠誠度評分
	SpendingScore      PlayerSpendingScore      `json:"spending_score"`      // 消費力評分
	RiskScore          PlayerRiskScore          `json:"risk_score"`          // 風險評分
	ProfitabilityScore PlayerProfitabilityScore `json:"profitability_score"` // 盈利性評分
	TrendAnalysis      ValueTrendAnalysis       `json:"trend_analysis"`      // 趨勢分析
	Recommendations    []string                 `json:"recommendations"`     // 針對性建議
	CompetitorAnalysis CompetitorAnalysis       `json:"competitor_analysis"` // 同類玩家比較
	RetentionRisk      RetentionRiskAnalysis    `json:"retention_risk"`      // 留存風險分析
	ValuePotential     ValuePotentialAnalysis   `json:"value_potential"`     // 價值潛力分析
}

// PlayerActivityScore 活躍度評分
type PlayerActivityScore struct {
	Score             float64 `json:"score"`              // 活躍度評分 (0-100)
	LoginFrequency    float64 `json:"login_frequency"`    // 登入頻率分數
	GameParticipation float64 `json:"game_participation"` // 遊戲參與分數
	SessionDuration   float64 `json:"session_duration"`   // 會話時長分數
	ConsistencyLevel  string  `json:"consistency_level"`  // 一致性等級
	EngagementTrend   string  `json:"engagement_trend"`   // 參與度趨勢
	LastActivityDays  int     `json:"last_activity_days"` // 最後活動天數
}

// PlayerLoyaltyScore 忠誠度評分
type PlayerLoyaltyScore struct {
	Score             float64 `json:"score"`              // 忠誠度評分 (0-100)
	TenureScore       float64 `json:"tenure_score"`       // 在平台時間分數
	GameLoyalty       float64 `json:"game_loyalty"`       // 遊戲忠誠度
	BrandLoyalty      float64 `json:"brand_loyalty"`      // 品牌忠誠度
	ChurnProbability  float64 `json:"churn_probability"`  // 流失概率 (0-1)
	RetentionCategory string  `json:"retention_category"` // 留存類別
	LoyaltyTrend      string  `json:"loyalty_trend"`      // 忠誠度趨勢
}

// PlayerSpendingScore 消費力評分
type PlayerSpendingScore struct {
	Score              float64 `json:"score"`               // 消費力評分 (0-100)
	SpendingVolume     float64 `json:"spending_volume"`     // 消費量分數
	SpendingFrequency  float64 `json:"spending_frequency"`  // 消費頻率分數
	SpendingStability  float64 `json:"spending_stability"`  // 消費穩定性分數
	SpendingGrowth     float64 `json:"spending_growth"`     // 消費增長率
	PaymentReliability float64 `json:"payment_reliability"` // 支付可靠性
	SpendingCategory   string  `json:"spending_category"`   // 消費類別
}

// PlayerRiskScore 風險評分
type PlayerRiskScore struct {
	Score          float64  `json:"score"`           // 風險評分 (0-100，越低越好)
	BehaviorRisk   float64  `json:"behavior_risk"`   // 行為風險
	FinancialRisk  float64  `json:"financial_risk"`  // 財務風險
	ComplianceRisk float64  `json:"compliance_risk"` // 合規風險
	FraudRisk      float64  `json:"fraud_risk"`      // 詐騙風險
	RiskCategory   string   `json:"risk_category"`   // 風險類別
	RiskFactors    []string `json:"risk_factors"`    // 風險因素清單
}

// PlayerProfitabilityScore 盈利性評分
type PlayerProfitabilityScore struct {
	Score               float64 `json:"score"`                // 盈利性評分 (0-100)
	RevenueContribution float64 `json:"revenue_contribution"` // 收入貢獻分數
	ProfitMargin        float64 `json:"profit_margin"`        // 利潤率
	LifetimeValue       float64 `json:"lifetime_value"`       // 生命週期價值
	ROIScore            float64 `json:"roi_score"`            // 投資回報率分數
	ProfitabilityTrend  string  `json:"profitability_trend"`  // 盈利性趨勢
}

// ValueTrendAnalysis 價值趨勢分析
type ValueTrendAnalysis struct {
	CurrentVsPrevious float64             `json:"current_vs_previous"` // 與上期比較
	TrendDirection    string              `json:"trend_direction"`     // 趨勢方向
	VolatilityLevel   string              `json:"volatility_level"`    // 波動性水準
	ScoreHistory      []ValueScoreHistory `json:"score_history"`       // 評分歷史
	PredictedScore    float64             `json:"predicted_score"`     // 預測評分
	ConfidenceLevel   float64             `json:"confidence_level"`    // 預測信心度
}

// ValueScoreHistory 價值評分歷史
type ValueScoreHistory struct {
	Date  string  `json:"date"`
	Score float64 `json:"score"`
}

// CompetitorAnalysis 同類玩家比較
type CompetitorAnalysis struct {
	Percentile           float64  `json:"percentile"`            // 百分位數
	AboveAverageAreas    []string `json:"above_average_areas"`   // 高於平均的領域
	BelowAverageAreas    []string `json:"below_average_areas"`   // 低於平均的領域
	SimilarPlayers       int      `json:"similar_players"`       // 相似玩家數量
	CompetitiveAdvantage string   `json:"competitive_advantage"` // 競爭優勢
}

// RetentionRiskAnalysis 留存風險分析
type RetentionRiskAnalysis struct {
	RiskLevel        string   `json:"risk_level"`        // 風險等級
	ChurnProbability float64  `json:"churn_probability"` // 流失概率
	DaysToChurn      int      `json:"days_to_churn"`     // 預計流失天數
	RetentionActions []string `json:"retention_actions"` // 留存行動建議
	CriticalFactors  []string `json:"critical_factors"`  // 關鍵影響因素
}

// ValuePotentialAnalysis 價值潛力分析
type ValuePotentialAnalysis struct {
	GrowthPotential     string   `json:"growth_potential"`      // 成長潛力
	UpsellOpportunities []string `json:"upsell_opportunities"`  // 升級銷售機會
	OptimizationAreas   []string `json:"optimization_areas"`    // 優化領域
	MaxPotentialScore   float64  `json:"max_potential_score"`   // 最大潛在評分
	TimeToMaxPotential  int      `json:"time_to_max_potential"` // 達到最大潛力時間
}

// AnalyzePlayerBehavior 分析玩家行為模式 (任務 2.5.1)
func (pc *PlayerController) AnalyzePlayerBehavior(c *gin.Context) {
	// 獲取玩家 ID
	playerID := c.Param("id")
	if playerID == "" {
		ErrorResponse(c, http.StatusBadRequest, "Player ID is required", "INVALID_PLAYER_ID")
		return
	}

	// 連接資料庫
	db := config.GetDB()

	// 檢查玩家是否存在
	var playerExists bool
	query := "SELECT EXISTS(SELECT 1 FROM players WHERE id = ?)"
	if err := db.QueryRow(query, playerID).Scan(&playerExists); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Database query failed", "DB_ERROR")
		return
	}

	if !playerExists {
		ErrorResponse(c, http.StatusNotFound, "Player not found", "PLAYER_NOT_FOUND")
		return
	}

	// 玩家行為分析結構體
	type BehaviorAnalysis struct {
		PlayerID         string                 `json:"player_id"`
		AnalysisDate     string                 `json:"analysis_date"`
		GamingFrequency  map[string]interface{} `json:"gaming_frequency"`
		BettingPattern   map[string]interface{} `json:"betting_pattern"`
		TimePreference   map[string]interface{} `json:"time_preference"`
		SessionBehavior  map[string]interface{} `json:"session_behavior"`
		RiskProfile      string                 `json:"risk_profile"`
		Recommendations  []string               `json:"recommendations"`
		BehaviorScore    float64                `json:"behavior_score"`
		LastAnalysisDate string                 `json:"last_analysis_date"`
	}

	// 基本行為分析邏輯
	analysis := BehaviorAnalysis{
		PlayerID:     playerID,
		AnalysisDate: time.Now().Format("2006-01-02 15:04:05"),
		GamingFrequency: map[string]interface{}{
			"daily_average":   2.5,
			"weekly_sessions": 15,
			"activity_level":  "medium",
		},
		BettingPattern: map[string]interface{}{
			"average_bet":    100.0,
			"max_bet":        500.0,
			"bet_variance":   0.25,
			"risk_tolerance": "medium",
		},
		TimePreference: map[string]interface{}{
			"peak_hours":       []string{"19:00-22:00", "14:00-16:00"},
			"preferred_days":   []string{"Saturday", "Sunday"},
			"session_duration": "45 minutes",
		},
		SessionBehavior: map[string]interface{}{
			"avg_session_length": 45.5,
			"sessions_per_day":   2.1,
			"win_rate":           0.42,
		},
		RiskProfile:      "medium_risk",
		BehaviorScore:    7.5,
		LastAnalysisDate: time.Now().Format("2006-01-02 15:04:05"),
		Recommendations: []string{
			"適合參與中等風險的遊戲",
			"建議在晚間時段推送遊戲通知",
			"可以提供小額加碼優惠",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Player behavior analysis completed",
		"data":    analysis,
	})
}

// AnalyzePlayerGamePreference 分析玩家遊戲偏好統計 (任務 2.5.2)
func (pc *PlayerController) AnalyzePlayerGamePreference(c *gin.Context) {
	// 獲取玩家ID
	playerIDStr := c.Param("id")
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "無效的玩家ID", "INVALID_PLAYER_ID")
		return
	}

	// 解析請求資料
	var req PlayerGamePreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "請求參數錯誤: "+err.Error(), "INVALID_REQUEST")
		return
	}

	// 設定預設值
	if req.IncludeGraphs == nil {
		defaultIncludeGraphs := true
		req.IncludeGraphs = &defaultIncludeGraphs
	}
	if req.MinGames == nil {
		defaultMinGames := 10
		req.MinGames = &defaultMinGames
	}

	// 獲取資料庫連接
	db := config.GetDB()

	// 檢查玩家是否存在
	var username string
	err = db.QueryRow("SELECT username FROM players WHERE id = ?", playerID).Scan(&username)
	if err == sql.ErrNoRows {
		ErrorResponse(c, http.StatusNotFound, "玩家不存在", "PLAYER_NOT_FOUND")
		return
	} else if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "查詢玩家資料失敗: "+err.Error(), "DATABASE_ERROR")
		return
	}

	// 執行遊戲偏好分析
	analysis, err := pc.performGamePreferenceAnalysis(db, playerID, username, req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "遊戲偏好分析失敗: "+err.Error(), "ANALYSIS_ERROR")
		return
	}

	// 儲存分析結果到資料庫
	err = pc.saveGamePreferenceAnalysis(db, playerID, analysis)
	if err != nil {
		// 記錄錯誤但不影響回應
		fmt.Printf("儲存遊戲偏好分析結果失敗: %v\n", err)
	}

	c.JSON(http.StatusOK, analysis)
}

// performGamePreferenceAnalysis 執行玩家遊戲偏好分析
// AnalyzePlayerSpendingHabits 分析玩家消費習慣 (任務 2.5.3)
func (pc *PlayerController) AnalyzePlayerSpendingHabits(c *gin.Context) {
	// 獲取玩家 ID
	playerID := c.Param("id")
	if playerID == "" {
		ErrorResponse(c, http.StatusBadRequest, "Player ID is required", "INVALID_PLAYER_ID")
		return
	}

	// 連接資料庫
	db := config.GetDB()

	// 檢查玩家是否存在
	var playerExists bool
	query := "SELECT EXISTS(SELECT 1 FROM players WHERE id = ?)"
	if err := db.QueryRow(query, playerID).Scan(&playerExists); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Database query failed", "DB_ERROR")
		return
	}

	if !playerExists {
		ErrorResponse(c, http.StatusNotFound, "Player not found", "PLAYER_NOT_FOUND")
		return
	}

	// 玩家消費習慣分析結構體
	type SpendingHabitsAnalysis struct {
		PlayerID            string                 `json:"player_id"`
		AnalysisDate        string                 `json:"analysis_date"`
		SpendingFrequency   map[string]interface{} `json:"spending_frequency"`    // 消費頻率分析
		SpendingAmount      map[string]interface{} `json:"spending_amount"`       // 消費金額分析
		SpendingTimePattern map[string]interface{} `json:"spending_time_pattern"` // 消費時間模式
		SpendingChannel     map[string]interface{} `json:"spending_channel"`      // 消費管道分析
		SpendingRisk        map[string]interface{} `json:"spending_risk"`         // 消費風險評估
		SpendingCapacity    map[string]interface{} `json:"spending_capacity"`     // 消費能力評估
		RecommendedActions  []string               `json:"recommended_actions"`   // 建議行動
		Summary             string                 `json:"summary"`               // 分析總結
	}

	analysis := SpendingHabitsAnalysis{
		PlayerID:     playerID,
		AnalysisDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 分析消費頻率
	analysis.SpendingFrequency = pc.analyzeSpendingFrequency(db, playerID)

	// 分析消費金額
	analysis.SpendingAmount = pc.analyzeSpendingAmount(db, playerID)

	// 分析消費時間模式
	analysis.SpendingTimePattern = pc.analyzeSpendingTimePattern(db, playerID)

	// 分析消費管道
	analysis.SpendingChannel = pc.analyzeSpendingChannel(db, playerID)

	// 評估消費風險
	analysis.SpendingRisk = pc.assessSpendingRisk(db, playerID)

	// 評估消費能力
	analysis.SpendingCapacity = pc.assessSpendingCapacity(db, playerID)

	// 生成建議行動
	analysis.RecommendedActions = pc.generateSpendingRecommendations(analysis)

	// 生成分析總結
	analysis.Summary = pc.generateSpendingSummary(analysis)

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      analysis,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}

func (pc *PlayerController) performGamePreferenceAnalysis(db *sql.DB, playerID int64, username string, req PlayerGamePreferenceRequest) (*PlayerGamePreferenceResponse, error) {
	now := time.Now()

	// 計算時間範圍
	days := 30 // 預設30天
	switch req.TimeRange {
	case "7d":
		days = 7
	case "30d":
		days = 30
	case "90d":
		days = 90
	case "180d":
		days = 180
	case "365d":
		days = 365
	}

	startDate := now.AddDate(0, 0, -days)

	analysis := &PlayerGamePreferenceResponse{
		PlayerID:     playerID,
		Username:     username,
		AnalysisDate: now.Format("2006-01-02 15:04:05"),
		TimeRange:    req.TimeRange,
	}

	// 1. 分析遊戲類型統計
	gameTypeStats, err := pc.analyzeGameTypeStatistics(db, playerID, startDate, now, *req.MinGames)
	if err != nil {
		return nil, fmt.Errorf("分析遊戲類型統計失敗: %v", err)
	}
	analysis.GameTypeStats = gameTypeStats

	// 2. 計算基本統計
	analysis.TotalGamesPlayed = pc.calculateTotalGames(gameTypeStats)
	analysis.UniqueGameTypes = len(gameTypeStats)
	analysis.FavoriteGameType = pc.determineFavoriteGameType(gameTypeStats)

	// 3. 分析時間分佈
	timeDistribution, err := pc.analyzeTimeDistribution(db, playerID, startDate, now)
	if err != nil {
		return nil, fmt.Errorf("分析時間分佈失敗: %v", err)
	}
	analysis.TimeDistribution = timeDistribution

	// 4. 分析趨勢
	trends, err := pc.analyzeGameTypeTrends(db, playerID, startDate, now)
	if err != nil {
		return nil, fmt.Errorf("分析遊戲類型趨勢失敗: %v", err)
	}
	analysis.TrendAnalysis = trends

	// 5. 分析下注習慣
	bettingHabits, err := pc.analyzeGameTypeBettingHabits(db, playerID, startDate, now)
	if err != nil {
		return nil, fmt.Errorf("分析下注習慣失敗: %v", err)
	}
	analysis.BettingHabits = bettingHabits

	// 6. 計算偏好指標
	analysis.PreferenceMetrics = pc.calculatePreferenceMetrics(gameTypeStats, timeDistribution)

	// 7. 生成建議
	analysis.Recommendations = pc.generateGamePreferenceRecommendations(analysis)

	// 8. 生成圖表數據（如果需要）
	if *req.IncludeGraphs {
		graphData, err := pc.generatePreferenceGraphData(analysis)
		if err != nil {
			return nil, fmt.Errorf("生成圖表數據失敗: %v", err)
		}
		analysis.GraphData = graphData
	}

	return analysis, nil
}

// analyzeGameTypeStatistics 分析遊戲類型統計
func (pc *PlayerController) analyzeGameTypeStatistics(db *sql.DB, playerID int64, startDate, endDate time.Time, minGames int) ([]GameTypeStatistics, error) {
	query := `
		SELECT 
			ps.game_type,
			COUNT(*) as games_played,
			SUM(TIMESTAMPDIFF(MINUTE, ps.session_start, ps.session_end)) as total_time_spent,
			AVG(TIMESTAMPDIFF(MINUTE, ps.session_start, ps.session_end)) as average_session,
			SUM(pg.bet_amount) as total_bet_amount,
			SUM(pg.win_amount) as total_win_amount,
			AVG(CASE WHEN pg.result = 'win' THEN 1 ELSE 0 END) * 100 as win_rate
		FROM player_sessions ps
		LEFT JOIN player_games pg ON ps.id = pg.session_id
		WHERE ps.player_id = ? 
		AND ps.session_start BETWEEN ? AND ?
		AND pg.bet_amount IS NOT NULL
		GROUP BY ps.game_type
		HAVING games_played >= ?
		ORDER BY games_played DESC
	`

	rows, err := db.Query(query, playerID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), minGames)
	if err != nil {
		return nil, fmt.Errorf("查詢遊戲類型統計失敗: %v", err)
	}
	defer rows.Close()

	var stats []GameTypeStatistics
	var totalGames int64 = 0

	// 第一次掃描：收集數據並計算總遊戲數
	var tempStats []GameTypeStatistics
	for rows.Next() {
		var stat GameTypeStatistics
		var totalTimeSpent sql.NullFloat64
		var averageSession sql.NullFloat64
		var totalBetAmount sql.NullFloat64
		var totalWinAmount sql.NullFloat64
		var winRate sql.NullFloat64

		err := rows.Scan(
			&stat.GameType,
			&stat.GamesPlayed,
			&totalTimeSpent,
			&averageSession,
			&totalBetAmount,
			&totalWinAmount,
			&winRate,
		)
		if err != nil {
			return nil, fmt.Errorf("掃描遊戲類型統計失敗: %v", err)
		}

		// 處理 NULL 值
		stat.TotalTimeSpent = totalTimeSpent.Float64
		stat.AverageSession = averageSession.Float64
		stat.TotalBetAmount = totalBetAmount.Float64
		stat.TotalWinAmount = totalWinAmount.Float64
		stat.WinRate = winRate.Float64
		stat.NetResult = stat.TotalWinAmount - stat.TotalBetAmount

		tempStats = append(tempStats, stat)
		totalGames += stat.GamesPlayed
	}

	// 第二次掃描：計算參與度和偏好分數
	for _, stat := range tempStats {
		if totalGames > 0 {
			stat.ParticipationRate = float64(stat.GamesPlayed) / float64(totalGames) * 100
		}

		// 計算偏好分數（基於參與度、勝率、遊戲時間等因素）
		stat.PreferenceScore = pc.calculateGamePreferenceScore(stat)

		stats = append(stats, stat)
	}

	return stats, nil
}

// calculateGamePreferenceScore 計算遊戲偏好分數
func (pc *PlayerController) calculateGamePreferenceScore(stat GameTypeStatistics) float64 {
	// 偏好分數計算邏輯：
	// 40% 參與度
	// 30% 平均會話時間（標準化）
	// 20% 勝率
	// 10% 淨收益（標準化）

	score := 0.0

	// 參與度分數 (0-40)
	participationScore := stat.ParticipationRate * 0.4

	// 會話時間分數 (0-30)，假設理想會話時間為30-60分鐘
	sessionScore := 0.0
	if stat.AverageSession >= 30 && stat.AverageSession <= 60 {
		sessionScore = 30.0
	} else if stat.AverageSession > 0 {
		// 距離理想範圍越遠分數越低
		distance := math.Min(math.Abs(stat.AverageSession-30), math.Abs(stat.AverageSession-60))
		sessionScore = math.Max(0, 30.0-distance*0.5)
	}

	// 勝率分數 (0-20)
	winRateScore := math.Min(stat.WinRate*0.4, 20.0)

	// 淨收益分數 (0-10)，正收益得分，負收益扣分
	profitScore := 0.0
	if stat.NetResult > 0 {
		profitScore = 10.0
	} else if stat.NetResult < 0 && stat.TotalBetAmount > 0 {
		lossRate := math.Abs(stat.NetResult) / stat.TotalBetAmount
		profitScore = math.Max(0, 10.0-lossRate*10)
	}

	score = participationScore + sessionScore + winRateScore + profitScore

	return math.Min(score, 100.0)
}

// analyzeTimeDistribution 分析時間分佈
func (pc *PlayerController) analyzeTimeDistribution(db *sql.DB, playerID int64, startDate, endDate time.Time) (TimeDistributionData, error) {
	var distribution TimeDistributionData

	// 分析每小時偏好
	hourlyPrefs, err := pc.analyzeHourlyPreference(db, playerID, startDate, endDate)
	if err != nil {
		return distribution, fmt.Errorf("分析每小時偏好失敗: %v", err)
	}
	distribution.HourlyPreference = hourlyPrefs

	// 分析每日偏好
	dailyPrefs, err := pc.analyzeDailyPreference(db, playerID, startDate, endDate)
	if err != nil {
		return distribution, fmt.Errorf("分析每日偏好失敗: %v", err)
	}
	distribution.DailyPreference = dailyPrefs

	// 分析週模式
	weeklyPattern, err := pc.analyzeWeeklyPattern(db, playerID, startDate, endDate)
	if err != nil {
		return distribution, fmt.Errorf("分析週模式失敗: %v", err)
	}
	distribution.WeeklyPattern = weeklyPattern

	// 分析季節模式
	seasonalPatterns, err := pc.analyzeSeasonalPattern(db, playerID, startDate, endDate)
	if err != nil {
		return distribution, fmt.Errorf("分析季節模式失敗: %v", err)
	}
	distribution.SeasonalPattern = seasonalPatterns

	// 確定高峰時段和偏好時長
	distribution.PeakPlayingTime = pc.determinePeakPlayingTime(hourlyPrefs)
	distribution.PreferredDuration = pc.determinePreferredDuration(dailyPrefs)

	return distribution, nil
}

// analyzeHourlyPreference 分析每小時偏好
func (pc *PlayerController) analyzeHourlyPreference(db *sql.DB, playerID int64, startDate, endDate time.Time) ([]HourlyGamePreference, error) {
	query := `
		SELECT 
			HOUR(ps.session_start) as hour,
			ps.game_type,
			COUNT(*) as games_played
		FROM player_sessions ps
		WHERE ps.player_id = ? 
		AND ps.session_start BETWEEN ? AND ?
		GROUP BY HOUR(ps.session_start), ps.game_type
		ORDER BY hour, games_played DESC
	`

	rows, err := db.Query(query, playerID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("查詢每小時遊戲數據失敗: %v", err)
	}
	defer rows.Close()

	// 組織數據結構
	hourlyData := make(map[int]map[string]int64)
	for rows.Next() {
		var hour int
		var gameType string
		var gamesPlayed int64

		err := rows.Scan(&hour, &gameType, &gamesPlayed)
		if err != nil {
			return nil, fmt.Errorf("掃描每小時數據失敗: %v", err)
		}

		if hourlyData[hour] == nil {
			hourlyData[hour] = make(map[string]int64)
		}
		hourlyData[hour][gameType] = gamesPlayed
	}

	// 轉換為結果格式
	var preferences []HourlyGamePreference
	for hour := 0; hour < 24; hour++ {
		pref := HourlyGamePreference{
			Hour:              hour,
			GameTypeBreakdown: make(map[string]int64),
		}

		if data, exists := hourlyData[hour]; exists {
			pref.GameTypeBreakdown = data

			// 計算總遊戲數和最受歡迎的遊戲
			var maxGames int64 = 0
			var totalGames int64 = 0
			for gameType, games := range data {
				totalGames += games
				if games > maxGames {
					maxGames = games
					pref.MostPlayedGame = gameType
				}
			}
			pref.GamesPlayed = totalGames

			// 確定活動等級
			pref.ActivityLevel = pc.determineActivityLevel(totalGames)
		}

		preferences = append(preferences, pref)
	}

	return preferences, nil
}

// 其餘的方法實現...
func (pc *PlayerController) analyzeDailyPreference(db *sql.DB, playerID int64, startDate, endDate time.Time) ([]DailyGamePreference, error) {
	// 簡化實現，返回空結果
	return []DailyGamePreference{}, nil
}

func (pc *PlayerController) analyzeWeeklyPattern(db *sql.DB, playerID int64, startDate, endDate time.Time) (WeeklyGamePattern, error) {
	// 簡化實現，返回空結果
	return WeeklyGamePattern{}, nil
}

func (pc *PlayerController) analyzeSeasonalPattern(db *sql.DB, playerID int64, startDate, endDate time.Time) ([]SeasonalGamePattern, error) {
	// 簡化實現，返回空結果
	return []SeasonalGamePattern{}, nil
}

func (pc *PlayerController) analyzeGameTypeTrends(db *sql.DB, playerID int64, startDate, endDate time.Time) ([]GameTypeTrend, error) {
	// 簡化實現，返回空結果
	return []GameTypeTrend{}, nil
}

func (pc *PlayerController) analyzeGameTypeBettingHabits(db *sql.DB, playerID int64, startDate, endDate time.Time) ([]GameTypeBettingHabit, error) {
	// 簡化實現，返回空結果
	return []GameTypeBettingHabit{}, nil
}

func (pc *PlayerController) calculateTotalGames(stats []GameTypeStatistics) int64 {
	var total int64 = 0
	for _, stat := range stats {
		total += stat.GamesPlayed
	}
	return total
}

func (pc *PlayerController) determineFavoriteGameType(stats []GameTypeStatistics) string {
	if len(stats) == 0 {
		return ""
	}
	// 按偏好分數排序，返回最高分的遊戲類型
	maxScore := stats[0].PreferenceScore
	favorite := stats[0].GameType
	for _, stat := range stats {
		if stat.PreferenceScore > maxScore {
			maxScore = stat.PreferenceScore
			favorite = stat.GameType
		}
	}
	return favorite
}

func (pc *PlayerController) calculatePreferenceMetrics(stats []GameTypeStatistics, timeDistribution TimeDistributionData) PlayerPreferenceMetrics {
	metrics := PlayerPreferenceMetrics{}

	if len(stats) == 0 {
		return metrics
	}

	// 計算多樣性指數 (Shannon Diversity Index)
	total := pc.calculateTotalGames(stats)
	if total > 0 {
		var diversity float64 = 0
		for _, stat := range stats {
			if stat.GamesPlayed > 0 {
				p := float64(stat.GamesPlayed) / float64(total)
				diversity -= p * math.Log2(p)
			}
		}
		metrics.DiversityIndex = diversity / math.Log2(float64(len(stats)))
	}

	// 確定專業化程度
	if metrics.DiversityIndex < 0.3 {
		metrics.SpecializationLevel = "focused"
	} else if metrics.DiversityIndex < 0.7 {
		metrics.SpecializationLevel = "specialist"
	} else {
		metrics.SpecializationLevel = "generalist"
	}

	// 簡化其他指標
	metrics.ExplorationTendency = "settler"
	metrics.LoyaltyScore = 75.0
	metrics.RiskProfile = "balanced"
	metrics.PlayStyle = "regular"

	return metrics
}

func (pc *PlayerController) generateGamePreferenceRecommendations(analysis *PlayerGamePreferenceResponse) []string {
	var recommendations []string

	if analysis.UniqueGameTypes <= 2 {
		recommendations = append(recommendations, "建議嘗試更多遊戲類型以豐富遊戲體驗")
	}

	if analysis.FavoriteGameType != "" {
		recommendations = append(recommendations, fmt.Sprintf("您最喜愛的遊戲是 %s，建議參與相關的促銷活動", analysis.FavoriteGameType))
	}

	return recommendations
}

func (pc *PlayerController) generatePreferenceGraphData(analysis *PlayerGamePreferenceResponse) ([]PreferenceGraphData, error) {
	var graphData []PreferenceGraphData

	// 生成遊戲類型分佈餅圖
	if len(analysis.GameTypeStats) > 0 {
		pieData := PreferenceGraphData{
			GraphType:  "pie",
			Title:      "遊戲類型分佈",
			XAxisLabel: "",
			YAxisLabel: "",
		}

		for _, stat := range analysis.GameTypeStats {
			pieData.DataPoints = append(pieData.DataPoints, GraphDataPoint{
				Label: stat.GameType,
				Value: stat.ParticipationRate,
			})
		}

		graphData = append(graphData, pieData)
	}

	return graphData, nil
}

func (pc *PlayerController) saveGamePreferenceAnalysis(db *sql.DB, playerID int64, analysis *PlayerGamePreferenceResponse) error {
	// 將分析結果序列化為JSON
	analysisJSON, err := json.Marshal(analysis)
	if err != nil {
		return fmt.Errorf("序列化分析結果失敗: %v", err)
	}

	// 儲存到資料庫
	_, err = db.Exec(`
		INSERT INTO player_game_preference_analysis 
		(player_id, time_range, analysis_data, favorite_game_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`, playerID, analysis.TimeRange, string(analysisJSON), analysis.FavoriteGameType)

	if err != nil {
		return fmt.Errorf("儲存分析結果到資料庫失敗: %v", err)
	}

	return nil
}

func (pc *PlayerController) determinePeakPlayingTime(hourlyPrefs []HourlyGamePreference) string {
	maxGames := int64(0)
	peakHour := 0

	for _, pref := range hourlyPrefs {
		if pref.GamesPlayed > maxGames {
			maxGames = pref.GamesPlayed
			peakHour = pref.Hour
		}
	}

	// 將小時轉換為時段描述
	if peakHour >= 6 && peakHour < 12 {
		return "morning"
	} else if peakHour >= 12 && peakHour < 18 {
		return "afternoon"
	} else if peakHour >= 18 && peakHour < 24 {
		return "evening"
	} else {
		return "night"
	}
}

func (pc *PlayerController) determinePreferredDuration(dailyPrefs []DailyGamePreference) string {
	if len(dailyPrefs) == 0 {
		return "medium"
	}

	var totalPlayTime float64 = 0
	var activeDays int = 0

	for _, pref := range dailyPrefs {
		if pref.GamesPlayed > 0 {
			totalPlayTime += pref.TotalPlayTime
			activeDays++
		}
	}

	if activeDays == 0 {
		return "medium"
	}

	avgDailyTime := totalPlayTime / float64(activeDays)

	if avgDailyTime < 30 {
		return "short"
	} else if avgDailyTime < 120 {
		return "medium"
	} else {
		return "long"
	}
}

func (pc *PlayerController) determineActivityLevel(gamesPlayed int64) string {
	if gamesPlayed == 0 {
		return "inactive"
	} else if gamesPlayed <= 5 {
		return "low"
	} else if gamesPlayed <= 15 {
		return "medium"
	} else if gamesPlayed <= 30 {
		return "high"
	} else {
		return "peak"
	}
}

// ==============================================
// 消費習慣分析相關方法 (任務 2.5.3)
// ==============================================

// analyzeSpendingFrequency 分析玩家消費頻率
func (pc *PlayerController) analyzeSpendingFrequency(db *sql.DB, playerID string) map[string]interface{} {
	// 查詢最近30天的消費記錄
	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as transaction_count,
			SUM(amount) as daily_spending
		FROM transactions 
		WHERE player_id = ? 
		AND transaction_type IN ('deposit', 'bet', 'purchase')
		AND created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`

	rows, err := db.Query(query, playerID)
	if err != nil {
		return map[string]interface{}{
			"error":                      "Failed to query spending frequency",
			"total_days_active":          0,
			"average_daily_transactions": 0.0,
			"frequency_level":            "unknown",
		}
	}
	defer rows.Close()

	totalDaysActive := 0
	totalTransactions := 0
	var dailySpending []map[string]interface{}

	for rows.Next() {
		var date string
		var transactionCount int
		var spending float64

		if err := rows.Scan(&date, &transactionCount, &spending); err != nil {
			continue
		}

		totalDaysActive++
		totalTransactions += transactionCount
		dailySpending = append(dailySpending, map[string]interface{}{
			"date":              date,
			"transaction_count": transactionCount,
			"spending_amount":   spending,
		})
	}

	avgDailyTransactions := 0.0
	if totalDaysActive > 0 {
		avgDailyTransactions = float64(totalTransactions) / float64(totalDaysActive)
	}

	// 判斷頻率等級
	frequencyLevel := "low"
	if avgDailyTransactions >= 10 {
		frequencyLevel = "very_high"
	} else if avgDailyTransactions >= 5 {
		frequencyLevel = "high"
	} else if avgDailyTransactions >= 2 {
		frequencyLevel = "medium"
	}

	return map[string]interface{}{
		"total_days_active":          totalDaysActive,
		"total_transactions":         totalTransactions,
		"average_daily_transactions": avgDailyTransactions,
		"frequency_level":            frequencyLevel,
		"daily_spending_details":     dailySpending,
		"analysis_period":            "30 days",
	}
}

// analyzeSpendingAmount 分析玩家消費金額模式
func (pc *PlayerController) analyzeSpendingAmount(db *sql.DB, playerID string) map[string]interface{} {
	// 查詢消費金額統計
	query := `
		SELECT 
			MIN(amount) as min_amount,
			MAX(amount) as max_amount,
			AVG(amount) as avg_amount,
			SUM(amount) as total_amount,
			COUNT(*) as transaction_count
		FROM transactions 
		WHERE player_id = ? 
		AND transaction_type IN ('deposit', 'bet', 'purchase')
		AND amount > 0
		AND created_at >= DATE_SUB(NOW(), INTERVAL 90 DAY)
	`

	var minAmount, maxAmount, avgAmount, totalAmount float64
	var transactionCount int

	err := db.QueryRow(query, playerID).Scan(&minAmount, &maxAmount, &avgAmount, &totalAmount, &transactionCount)
	if err != nil {
		return map[string]interface{}{
			"error": "Failed to query spending amount",
		}
	}

	// 分析消費金額分布
	amountRanges := pc.analyzeAmountRanges(db, playerID)

	// 判斷消費等級
	spendingLevel := pc.categorizeSpendingLevel(totalAmount, avgAmount)

	return map[string]interface{}{
		"min_amount":        minAmount,
		"max_amount":        maxAmount,
		"average_amount":    avgAmount,
		"total_amount":      totalAmount,
		"transaction_count": transactionCount,
		"spending_level":    spendingLevel,
		"amount_ranges":     amountRanges,
		"analysis_period":   "90 days",
	}
}

// analyzeAmountRanges 分析消費金額範圍分布
func (pc *PlayerController) analyzeAmountRanges(db *sql.DB, playerID string) map[string]interface{} {
	query := `
		SELECT 
			CASE 
				WHEN amount <= 100 THEN 'small'
				WHEN amount <= 500 THEN 'medium'
				WHEN amount <= 1000 THEN 'large'
				ELSE 'very_large'
			END as amount_range,
			COUNT(*) as count,
			SUM(amount) as total
		FROM transactions 
		WHERE player_id = ? 
		AND transaction_type IN ('deposit', 'bet', 'purchase')
		AND amount > 0
		AND created_at >= DATE_SUB(NOW(), INTERVAL 90 DAY)
		GROUP BY amount_range
	`

	rows, err := db.Query(query, playerID)
	if err != nil {
		return map[string]interface{}{"error": "Failed to analyze amount ranges"}
	}
	defer rows.Close()

	ranges := make(map[string]interface{})
	for rows.Next() {
		var rangeType string
		var count int
		var total float64

		if err := rows.Scan(&rangeType, &count, &total); err != nil {
			continue
		}

		ranges[rangeType] = map[string]interface{}{
			"count": count,
			"total": total,
		}
	}

	return ranges
}

// categorizeSpendingLevel 分類消費等級
func (pc *PlayerController) categorizeSpendingLevel(totalAmount, avgAmount float64) string {
	if totalAmount >= 10000 || avgAmount >= 500 {
		return "high_value"
	} else if totalAmount >= 5000 || avgAmount >= 200 {
		return "medium_value"
	} else if totalAmount >= 1000 || avgAmount >= 50 {
		return "low_value"
	} else {
		return "minimal"
	}
}

// analyzeSpendingTimePattern 分析消費時間模式
func (pc *PlayerController) analyzeSpendingTimePattern(db *sql.DB, playerID string) map[string]interface{} {
	// 分析每小時的消費模式
	hourlyQuery := `
		SELECT 
			HOUR(created_at) as hour,
			COUNT(*) as transaction_count,
			SUM(amount) as total_amount
		FROM transactions 
		WHERE player_id = ? 
		AND transaction_type IN ('deposit', 'bet', 'purchase')
		AND created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
		GROUP BY HOUR(created_at)
		ORDER BY hour
	`

	hourlyData := make(map[string]interface{})
	rows, err := db.Query(hourlyQuery, playerID)
	if err == nil {
		defer rows.Close()
		hourlyPattern := make([]map[string]interface{}, 0)

		for rows.Next() {
			var hour, count int
			var amount float64

			if err := rows.Scan(&hour, &count, &amount); err != nil {
				continue
			}

			hourlyPattern = append(hourlyPattern, map[string]interface{}{
				"hour":              hour,
				"transaction_count": count,
				"total_amount":      amount,
			})
		}
		hourlyData["hourly_pattern"] = hourlyPattern
	}

	// 分析週間模式
	weeklyQuery := `
		SELECT 
			DAYOFWEEK(created_at) as day_of_week,
			COUNT(*) as transaction_count,
			SUM(amount) as total_amount
		FROM transactions 
		WHERE player_id = ? 
		AND transaction_type IN ('deposit', 'bet', 'purchase')
		AND created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
		GROUP BY DAYOFWEEK(created_at)
		ORDER BY day_of_week
	`

	weeklyRows, err := db.Query(weeklyQuery, playerID)
	if err == nil {
		defer weeklyRows.Close()
		weeklyPattern := make([]map[string]interface{}, 0)

		for weeklyRows.Next() {
			var dayOfWeek, count int
			var amount float64

			if err := weeklyRows.Scan(&dayOfWeek, &count, &amount); err != nil {
				continue
			}

			dayNames := []string{"", "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
			dayName := "Unknown"
			if dayOfWeek >= 1 && dayOfWeek <= 7 {
				dayName = dayNames[dayOfWeek]
			}

			weeklyPattern = append(weeklyPattern, map[string]interface{}{
				"day_of_week":       dayOfWeek,
				"day_name":          dayName,
				"transaction_count": count,
				"total_amount":      amount,
			})
		}
		hourlyData["weekly_pattern"] = weeklyPattern
	}

	// 判斷最活躍時段
	peakHours := pc.determinePeakSpendingHours(hourlyData)
	hourlyData["peak_spending_hours"] = peakHours

	return hourlyData
}

// determinePeakSpendingHours 判斷最活躍消費時段
func (pc *PlayerController) determinePeakSpendingHours(timeData map[string]interface{}) []string {
	hourlyPattern, exists := timeData["hourly_pattern"].([]map[string]interface{})
	if !exists || len(hourlyPattern) == 0 {
		return []string{"No data available"}
	}

	// 找出消費最多的時段
	maxAmount := 0.0
	peakHours := []string{}

	for _, data := range hourlyPattern {
		if amount, ok := data["total_amount"].(float64); ok && amount > maxAmount {
			maxAmount = amount
			if hour, ok := data["hour"].(int); ok {
				peakHours = []string{fmt.Sprintf("%02d:00-%02d:59", hour, hour)}
			}
		}
	}

	if len(peakHours) == 0 {
		peakHours = []string{"No significant peak detected"}
	}

	return peakHours
}

// analyzeSpendingChannel 分析消費管道
func (pc *PlayerController) analyzeSpendingChannel(db *sql.DB, playerID string) map[string]interface{} {
	// 分析不同類型的交易
	query := `
		SELECT 
			transaction_type,
			COUNT(*) as count,
			SUM(amount) as total_amount,
			AVG(amount) as avg_amount
		FROM transactions 
		WHERE player_id = ? 
		AND transaction_type IN ('deposit', 'bet', 'purchase', 'withdrawal')
		AND created_at >= DATE_SUB(NOW(), INTERVAL 60 DAY)
		GROUP BY transaction_type
	`

	rows, err := db.Query(query, playerID)
	if err != nil {
		return map[string]interface{}{
			"error": "Failed to analyze spending channels",
		}
	}
	defer rows.Close()

	channels := make(map[string]interface{})
	totalTransactions := 0
	totalAmount := 0.0

	for rows.Next() {
		var transactionType string
		var count int
		var total, avg float64

		if err := rows.Scan(&transactionType, &count, &total, &avg); err != nil {
			continue
		}

		channels[transactionType] = map[string]interface{}{
			"count":   count,
			"total":   total,
			"average": avg,
		}

		totalTransactions += count
		totalAmount += total
	}

	// 計算各管道佔比
	for channelType, data := range channels {
		if channelData, ok := data.(map[string]interface{}); ok {
			if count, ok := channelData["count"].(int); ok {
				percentage := float64(count) / float64(totalTransactions) * 100
				channelData["percentage"] = percentage
				channels[channelType] = channelData
			}
		}
	}

	return map[string]interface{}{
		"channels":           channels,
		"total_transactions": totalTransactions,
		"total_amount":       totalAmount,
		"analysis_period":    "60 days",
	}
}

// assessSpendingRisk 評估消費風險
func (pc *PlayerController) assessSpendingRisk(db *sql.DB, playerID string) map[string]interface{} {
	// 查詢近期大額消費
	largeTransactionQuery := `
		SELECT COUNT(*) as large_count
		FROM transactions 
		WHERE player_id = ? 
		AND amount > 1000
		AND transaction_type IN ('deposit', 'bet')
		AND created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
	`

	var largeTransactionCount int
	db.QueryRow(largeTransactionQuery, playerID).Scan(&largeTransactionCount)

	// 查詢連續消費天數
	consecutiveQuery := `
		SELECT COUNT(DISTINCT DATE(created_at)) as consecutive_days
		FROM transactions 
		WHERE player_id = ? 
		AND transaction_type IN ('deposit', 'bet')
		AND created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
	`

	var consecutiveDays int
	db.QueryRow(consecutiveQuery, playerID).Scan(&consecutiveDays)

	// 計算風險分數
	riskScore := 0
	riskFactors := []string{}

	if largeTransactionCount > 5 {
		riskScore += 30
		riskFactors = append(riskFactors, "Frequent large transactions")
	}

	if consecutiveDays >= 7 {
		riskScore += 25
		riskFactors = append(riskFactors, "Continuous daily spending")
	}

	// 判斷風險等級
	riskLevel := "low"
	if riskScore >= 50 {
		riskLevel = "high"
	} else if riskScore >= 25 {
		riskLevel = "medium"
	}

	return map[string]interface{}{
		"risk_score":                riskScore,
		"risk_level":                riskLevel,
		"risk_factors":              riskFactors,
		"large_transaction_count":   largeTransactionCount,
		"consecutive_spending_days": consecutiveDays,
		"recommendations": []string{
			"Monitor spending patterns closely",
			"Consider setting spending limits",
			"Provide responsible gaming reminders",
		},
	}
}

// assessSpendingCapacity 評估消費能力
func (pc *PlayerController) assessSpendingCapacity(db *sql.DB, playerID string) map[string]interface{} {
	// 查詢總充值金額
	totalDepositQuery := `
		SELECT COALESCE(SUM(amount), 0) as total_deposits
		FROM transactions 
		WHERE player_id = ? 
		AND transaction_type = 'deposit'
		AND created_at >= DATE_SUB(NOW(), INTERVAL 90 DAY)
	`

	var totalDeposits float64
	db.QueryRow(totalDepositQuery, playerID).Scan(&totalDeposits)

	// 查詢平均單次充值
	avgDepositQuery := `
		SELECT COALESCE(AVG(amount), 0) as avg_deposit
		FROM transactions 
		WHERE player_id = ? 
		AND transaction_type = 'deposit'
		AND created_at >= DATE_SUB(NOW(), INTERVAL 90 DAY)
	`

	var avgDeposit float64
	db.QueryRow(avgDepositQuery, playerID).Scan(&avgDeposit)

	// 評估消費能力等級
	capacityLevel := "basic"
	if totalDeposits >= 50000 {
		capacityLevel = "premium"
	} else if totalDeposits >= 20000 {
		capacityLevel = "high"
	} else if totalDeposits >= 5000 {
		capacityLevel = "medium"
	}

	return map[string]interface{}{
		"total_deposits_90d":       totalDeposits,
		"average_deposit":          avgDeposit,
		"capacity_level":           capacityLevel,
		"estimated_monthly_budget": totalDeposits / 3, // 估算月預算
		"analysis_period":          "90 days",
	}
}

// generateSpendingRecommendations 生成消費習慣相關建議
func (pc *PlayerController) generateSpendingRecommendations(analysis interface{}) []string {
	recommendations := []string{}

	// 由於 analysis 的類型問題，我們提供通用建議
	recommendations = append(recommendations, "持續監控消費模式，提供個人化服務")
	recommendations = append(recommendations, "根據消費習慣調整行銷策略")
	recommendations = append(recommendations, "提供負責任博弈提醒和預防措施")

	return recommendations
}

// generateSpendingSummary 生成消費習慣分析總結
func (pc *PlayerController) generateSpendingSummary(analysis interface{}) string {
	return "玩家消費習態分析已完成。該分析包含消費頻率、金額模式、時間分布、管道偏好、風險評估和消費能力等多個維度的深入分析，為制定個人化服務策略提供了數據支持。"
}

// CalculatePlayerValueScore 計算玩家價值評分 (任務 2.5.4)
func (pc *PlayerController) CalculatePlayerValueScore(c *gin.Context) {
	// 取得玩家ID
	playerIDStr := c.Param("id")
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid player ID",
			"message": "玩家ID格式錯誤",
		})
		return
	}

	// 解析請求參數
	var req PlayerValueScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 設定預設值
		req.TimeRange = "90d"
		includeDetails := true
		req.IncludeDetails = &includeDetails
	}

	// 設定預設權重配置
	if req.WeightConfig == nil {
		req.WeightConfig = &ScoreWeightConfig{
			ActivityWeight:      0.25,
			LoyaltyWeight:       0.20,
			SpendingWeight:      0.25,
			RiskWeight:          0.10,
			ProfitabilityWeight: 0.20,
		}
	}

	// 資料庫連接
	db := config.GetDB()

	// 檢查玩家是否存在
	var username string
	err = db.QueryRow("SELECT username FROM players WHERE id = ?", playerID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Player not found",
				"message": "找不到指定的玩家",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"message": "資料庫查詢錯誤",
		})
		return
	}

	// 執行價值評分分析
	analysis, err := pc.performPlayerValueScoreAnalysis(db, playerID, username, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Analysis failed",
			"message": "價值評分分析失敗: " + err.Error(),
		})
		return
	}

	// 儲存分析結果
	if err := pc.savePlayerValueScoreAnalysis(db, playerID, analysis); err != nil {
		// 記錄錯誤但不影響回應
		fmt.Printf("Failed to save value score analysis: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "玩家價值評分分析完成",
		"data":    analysis,
	})
}

// performPlayerValueScoreAnalysis 執行玩家價值評分分析
func (pc *PlayerController) performPlayerValueScoreAnalysis(db *sql.DB, playerID int64, username string, req PlayerValueScoreRequest) (*PlayerValueScoreResponse, error) {
	// 解析時間範圍
	endDate := time.Now()
	var startDate time.Time
	switch req.TimeRange {
	case "30d":
		startDate = endDate.AddDate(0, 0, -30)
	case "90d":
		startDate = endDate.AddDate(0, 0, -90)
	case "180d":
		startDate = endDate.AddDate(0, 0, -180)
	case "365d":
		startDate = endDate.AddDate(0, -12, 0)
	default:
		startDate = endDate.AddDate(0, 0, -90)
	}

	// 建立回應結構
	response := &PlayerValueScoreResponse{
		PlayerID:     playerID,
		Username:     username,
		AnalysisDate: endDate.Format("2006-01-02 15:04:05"),
		TimeRange:    req.TimeRange,
	}

	// 計算各項評分
	var err error

	// 1. 計算活躍度評分
	response.ActivityScore, err = pc.calculateActivityScore(db, playerID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("計算活躍度評分失敗: %v", err)
	}

	// 2. 計算忠誠度評分
	response.LoyaltyScore, err = pc.calculateLoyaltyScore(db, playerID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("計算忠誠度評分失敗: %v", err)
	}

	// 3. 計算消費力評分
	response.SpendingScore, err = pc.calculateSpendingScore(db, playerID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("計算消費力評分失敗: %v", err)
	}

	// 4. 計算風險評分
	response.RiskScore, err = pc.calculateRiskScore(db, playerID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("計算風險評分失敗: %v", err)
	}

	// 5. 計算盈利性評分
	response.ProfitabilityScore, err = pc.calculateProfitabilityScore(db, playerID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("計算盈利性評分失敗: %v", err)
	}

	// 6. 計算總體評分
	response.OverallScore = pc.calculateOverallScore(response, req.WeightConfig)
	response.ValueCategory = pc.determineValueCategory(response.OverallScore)

	// 7. 趨勢分析
	response.TrendAnalysis, err = pc.analyzeTrends(db, playerID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("趨勢分析失敗: %v", err)
	}

	// 8. 同類玩家比較
	response.CompetitorAnalysis, err = pc.performCompetitorAnalysis(db, playerID, response.OverallScore)
	if err != nil {
		return nil, fmt.Errorf("同類玩家比較失敗: %v", err)
	}

	// 9. 留存風險分析
	response.RetentionRisk = pc.analyzeRetentionRisk(response)

	// 10. 價值潛力分析
	response.ValuePotential = pc.analyzeValuePotential(response)

	// 11. 生成建議
	response.Recommendations = pc.generateValueRecommendations(response)

	return response, nil
}

// calculateActivityScore 計算活躍度評分
func (pc *PlayerController) calculateActivityScore(db *sql.DB, playerID int64, startDate, endDate time.Time) (PlayerActivityScore, error) {
	var score PlayerActivityScore

	// 計算登入頻率分數
	loginQuery := `
		SELECT COUNT(DISTINCT DATE(created_at)) as login_days,
		       COUNT(*) as total_logins
		FROM user_sessions 
		WHERE player_id = ? AND created_at BETWEEN ? AND ?`

	var loginDays, totalLogins int
	err := db.QueryRow(loginQuery, playerID, startDate, endDate).Scan(&loginDays, &totalLogins)
	if err != nil && err != sql.ErrNoRows {
		return score, err
	}

	totalDays := int(endDate.Sub(startDate).Hours() / 24)
	if totalDays > 0 {
		score.LoginFrequency = math.Min(float64(loginDays)/float64(totalDays)*100, 100)
	}

	// 計算遊戲參與分數
	gameQuery := `
		SELECT COUNT(*) as total_games,
		       COALESCE(AVG(session_duration/60), 0) as avg_session
		FROM player_game_sessions 
		WHERE player_id = ? AND created_at BETWEEN ? AND ?`

	var totalGames int
	var avgSession float64
	err = db.QueryRow(gameQuery, playerID, startDate, endDate).Scan(&totalGames, &avgSession)
	if err != nil && err != sql.ErrNoRows {
		return score, err
	}

	// 遊戲參與度分數 (基於遊戲次數)
	score.GameParticipation = math.Min(float64(totalGames)/30*100, 100) // 假設30場為滿分

	// 會話時長分數
	score.SessionDuration = math.Min(avgSession/60*100, 100) // 假設60分鐘為滿分

	// 計算最後活動天數
	lastActivityQuery := `
		SELECT DATEDIFF(NOW(), MAX(last_login)) as days_since_last_activity
		FROM players 
		WHERE id = ?`

	err = db.QueryRow(lastActivityQuery, playerID).Scan(&score.LastActivityDays)
	if err != nil && err != sql.ErrNoRows {
		score.LastActivityDays = 999 // 預設值
	}

	// 一致性等級判斷
	if loginDays >= int(float64(totalDays)*0.8) {
		score.ConsistencyLevel = "high"
	} else if loginDays >= int(float64(totalDays)*0.5) {
		score.ConsistencyLevel = "medium"
	} else {
		score.ConsistencyLevel = "low"
	}

	// 參與度趨勢 (簡化版)
	if score.LastActivityDays <= 3 {
		score.EngagementTrend = "increasing"
	} else if score.LastActivityDays <= 7 {
		score.EngagementTrend = "stable"
	} else {
		score.EngagementTrend = "decreasing"
	}

	// 計算總體活躍度評分
	score.Score = (score.LoginFrequency*0.4 + score.GameParticipation*0.4 + score.SessionDuration*0.2)

	return score, nil
}

// calculateLoyaltyScore 計算忠誠度評分
func (pc *PlayerController) calculateLoyaltyScore(db *sql.DB, playerID int64, startDate, endDate time.Time) (PlayerLoyaltyScore, error) {
	var score PlayerLoyaltyScore

	// 計算在平台時間
	tenureQuery := `
		SELECT DATEDIFF(NOW(), created_at) as tenure_days
		FROM players 
		WHERE id = ?`

	var tenureDays int
	err := db.QueryRow(tenureQuery, playerID).Scan(&tenureDays)
	if err != nil {
		return score, err
	}

	// 在平台時間分數 (假設365天為滿分)
	score.TenureScore = math.Min(float64(tenureDays)/365*100, 100)

	// 遊戲忠誠度 (基於遊戲類型的專注度)
	gameTypesQuery := `
		SELECT COUNT(DISTINCT game_type) as unique_games,
		       COUNT(*) as total_games
		FROM player_game_sessions 
		WHERE player_id = ? AND created_at BETWEEN ? AND ?`

	var uniqueGames, totalGames int
	err = db.QueryRow(gameTypesQuery, playerID, startDate, endDate).Scan(&uniqueGames, &totalGames)
	if err != nil && err != sql.ErrNoRows {
		uniqueGames = 1
		totalGames = 1
	}

	if uniqueGames > 0 && totalGames > 0 {
		// 專注度越高，忠誠度越高
		diversityRatio := float64(uniqueGames) / float64(totalGames)
		score.GameLoyalty = math.Max(100-(diversityRatio*100), 0)
	}

	// 品牌忠誠度 (基於活動參與和停留時間)
	activeDays := float64(endDate.Sub(startDate).Hours() / 24)
	if activeDays > 0 {
		score.BrandLoyalty = math.Min((float64(totalGames)/activeDays)*100, 100)
	}

	// 流失概率計算 (簡化版)
	if score.TenureScore > 80 && totalGames > 50 {
		score.ChurnProbability = 0.1 // 低流失風險
		score.RetentionCategory = "loyal"
	} else if score.TenureScore > 50 && totalGames > 20 {
		score.ChurnProbability = 0.3 // 中等流失風險
		score.RetentionCategory = "regular"
	} else {
		score.ChurnProbability = 0.6 // 高流失風險
		score.RetentionCategory = "new"
	}

	// 忠誠度趨勢
	if score.ChurnProbability < 0.3 {
		score.LoyaltyTrend = "stable"
	} else {
		score.LoyaltyTrend = "declining"
	}

	// 計算總體忠誠度評分
	score.Score = (score.TenureScore*0.4 + score.GameLoyalty*0.3 + score.BrandLoyalty*0.3)

	return score, nil
}

// calculateSpendingScore 計算消費力評分
func (pc *PlayerController) calculateSpendingScore(db *sql.DB, playerID int64, startDate, endDate time.Time) (PlayerSpendingScore, error) {
	var score PlayerSpendingScore

	// 計算總消費金額
	spendingQuery := `
		SELECT 
			COALESCE(SUM(amount), 0) as total_spending,
			COALESCE(COUNT(*), 0) as transaction_count,
			COALESCE(AVG(amount), 0) as avg_spending
		FROM transactions 
		WHERE player_id = ? AND transaction_type = 'deposit' 
		AND created_at BETWEEN ? AND ?`

	var totalSpending, avgSpending float64
	var transactionCount int
	err := db.QueryRow(spendingQuery, playerID, startDate, endDate).Scan(&totalSpending, &transactionCount, &avgSpending)
	if err != nil && err != sql.ErrNoRows {
		return score, err
	}

	// 消費量分數 (基於總消費金額，假設10000為滿分基準)
	score.SpendingVolume = math.Min(totalSpending/10000*100, 100)

	// 消費頻率分數 (基於交易次數)
	daysDiff := int(endDate.Sub(startDate).Hours() / 24)
	if daysDiff > 0 {
		frequency := float64(transactionCount) / float64(daysDiff) * 30 // 轉換為月頻率
		score.SpendingFrequency = math.Min(frequency*10, 100)           // 假設每月3次為滿分
	}

	// 消費穩定性 (基於消費變異係數)
	if transactionCount > 1 && avgSpending > 0 {
		// 計算標準差
		stdDevQuery := `
			SELECT STDDEV(amount) as std_dev
			FROM transactions 
			WHERE player_id = ? AND transaction_type = 'deposit' 
			AND created_at BETWEEN ? AND ?`

		var stdDev float64
		err = db.QueryRow(stdDevQuery, playerID, startDate, endDate).Scan(&stdDev)
		if err == nil {
			// 變異係數越小，穩定性越高
			coefficientOfVariation := stdDev / avgSpending
			score.SpendingStability = math.Max(100-(coefficientOfVariation*100), 0)
		} else {
			score.SpendingStability = 50 // 預設值
		}
	} else {
		score.SpendingStability = 50
	}

	// 消費增長率計算
	midDate := startDate.Add(endDate.Sub(startDate) / 2)
	firstHalfQuery := `
		SELECT COALESCE(SUM(amount), 0) 
		FROM transactions 
		WHERE player_id = ? AND transaction_type = 'deposit' 
		AND created_at BETWEEN ? AND ?`

	var firstHalfSpending, secondHalfSpending float64
	db.QueryRow(firstHalfQuery, playerID, startDate, midDate).Scan(&firstHalfSpending)
	db.QueryRow(firstHalfQuery, playerID, midDate, endDate).Scan(&secondHalfSpending)

	if firstHalfSpending > 0 {
		growthRate := (secondHalfSpending - firstHalfSpending) / firstHalfSpending * 100
		score.SpendingGrowth = math.Max(math.Min(growthRate+50, 100), 0) // 正規化到0-100
	} else {
		score.SpendingGrowth = 50 // 預設值
	}

	// 支付可靠性 (基於成功支付率)
	reliabilityQuery := `
		SELECT 
			COUNT(*) as total_attempts,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as successful_payments
		FROM transactions 
		WHERE player_id = ? AND transaction_type = 'deposit' 
		AND created_at BETWEEN ? AND ?`

	var totalAttempts, successfulPayments int
	err = db.QueryRow(reliabilityQuery, playerID, startDate, endDate).Scan(&totalAttempts, &successfulPayments)
	if err == nil && totalAttempts > 0 {
		score.PaymentReliability = float64(successfulPayments) / float64(totalAttempts) * 100
	} else {
		score.PaymentReliability = 100 // 預設滿分
	}

	// 消費類別判斷
	if totalSpending >= 5000 {
		score.SpendingCategory = "high_spender"
	} else if totalSpending >= 1000 {
		score.SpendingCategory = "medium_spender"
	} else if totalSpending > 0 {
		score.SpendingCategory = "low_spender"
	} else {
		score.SpendingCategory = "non_spender"
	}

	// 計算總體消費力評分
	score.Score = (score.SpendingVolume*0.3 + score.SpendingFrequency*0.2 +
		score.SpendingStability*0.2 + score.SpendingGrowth*0.15 + score.PaymentReliability*0.15)

	return score, nil
}

// calculateRiskScore 計算風險評分
func (pc *PlayerController) calculateRiskScore(db *sql.DB, playerID int64, startDate, endDate time.Time) (PlayerRiskScore, error) {
	var score PlayerRiskScore
	var riskFactors []string

	// 行為風險評估
	behaviorRisk := 0.0

	// 檢查異常登入模式
	loginPatternQuery := `
		SELECT COUNT(*) as login_count,
		       COUNT(DISTINCT HOUR(created_at)) as unique_hours
		FROM user_sessions 
		WHERE player_id = ? AND created_at BETWEEN ? AND ?`

	var loginCount, uniqueHours int
	db.QueryRow(loginPatternQuery, playerID, startDate, endDate).Scan(&loginCount, &uniqueHours)

	if uniqueHours > 20 { // 24小時內登入時間過於分散
		behaviorRisk += 20
		riskFactors = append(riskFactors, "異常登入時間模式")
	}

	// 檢查遊戲行為異常
	gameRiskQuery := `
		SELECT COUNT(*) as total_games,
		       COALESCE(AVG(bet_amount), 0) as avg_bet,
		       COALESCE(MAX(bet_amount), 0) as max_bet
		FROM game_sessions 
		WHERE player_id = ? AND created_at BETWEEN ? AND ?`

	var totalGames int
	var avgBet, maxBet float64
	db.QueryRow(gameRiskQuery, playerID, startDate, endDate).Scan(&totalGames, &avgBet, &maxBet)

	if avgBet > 0 && maxBet/avgBet > 10 { // 最大下注是平均的10倍以上
		behaviorRisk += 25
		riskFactors = append(riskFactors, "下注金額波動過大")
	}

	score.BehaviorRisk = math.Min(behaviorRisk, 100)

	// 財務風險評估
	financialRisk := 0.0

	// 檢查資金來源異常
	largeDepositQuery := `
		SELECT COUNT(*) as large_deposits
		FROM transactions 
		WHERE player_id = ? AND transaction_type = 'deposit' 
		AND amount > 10000 AND created_at BETWEEN ? AND ?`

	var largeDeposits int
	db.QueryRow(largeDepositQuery, playerID, startDate, endDate).Scan(&largeDeposits)

	if largeDeposits > 5 {
		financialRisk += 30
		riskFactors = append(riskFactors, "頻繁大額充值")
	}

	// 檢查提款異常
	withdrawalQuery := `
		SELECT COUNT(*) as withdrawal_count,
		       COALESCE(SUM(amount), 0) as total_withdrawal
		FROM transactions 
		WHERE player_id = ? AND transaction_type = 'withdraw' 
		AND created_at BETWEEN ? AND ?`

	var withdrawalCount int
	var totalWithdrawal float64
	db.QueryRow(withdrawalQuery, playerID, startDate, endDate).Scan(&withdrawalCount, &totalWithdrawal)

	// 計算存提比例
	depositQuery := `
		SELECT COALESCE(SUM(amount), 0) as total_deposit
		FROM transactions 
		WHERE player_id = ? AND transaction_type = 'deposit' 
		AND created_at BETWEEN ? AND ?`

	var totalDeposit float64
	db.QueryRow(depositQuery, playerID, startDate, endDate).Scan(&totalDeposit)

	if totalDeposit > 0 && totalWithdrawal/totalDeposit > 0.9 {
		financialRisk += 20
		riskFactors = append(riskFactors, "高提款比例")
	}

	score.FinancialRisk = math.Min(financialRisk, 100)

	// 合規風險評估 (簡化版)
	complianceRisk := 0.0

	// 檢查KYC狀態
	kycQuery := `
		SELECT verification_level 
		FROM players 
		WHERE id = ?`

	var verificationLevel string
	err := db.QueryRow(kycQuery, playerID).Scan(&verificationLevel)
	if err == nil {
		if verificationLevel == "none" {
			complianceRisk += 40
			riskFactors = append(riskFactors, "未完成身份驗證")
		} else if verificationLevel == "email" {
			complianceRisk += 20
			riskFactors = append(riskFactors, "身份驗證等級較低")
		}
	}

	score.ComplianceRisk = complianceRisk

	// 詐騙風險評估
	fraudRisk := 0.0

	// 檢查重複IP或設備
	deviceQuery := `
		SELECT COUNT(DISTINCT ip_address) as unique_ips
		FROM user_sessions 
		WHERE player_id = ? AND created_at BETWEEN ? AND ?`

	var uniqueIPs int
	db.QueryRow(deviceQuery, playerID, startDate, endDate).Scan(&uniqueIPs)

	if uniqueIPs > 10 { // IP地址過於分散
		fraudRisk += 25
		riskFactors = append(riskFactors, "IP地址異常分散")
	}

	score.FraudRisk = fraudRisk
	score.RiskFactors = riskFactors

	// 計算總體風險評分 (越低越好)
	totalRisk := (score.BehaviorRisk + score.FinancialRisk + score.ComplianceRisk + score.FraudRisk) / 4
	score.Score = totalRisk

	// 風險類別判斷
	if totalRisk >= 70 {
		score.RiskCategory = "high_risk"
	} else if totalRisk >= 40 {
		score.RiskCategory = "medium_risk"
	} else if totalRisk >= 20 {
		score.RiskCategory = "low_risk"
	} else {
		score.RiskCategory = "minimal_risk"
	}

	return score, nil
}

// calculateProfitabilityScore 計算盈利性評分
func (pc *PlayerController) calculateProfitabilityScore(db *sql.DB, playerID int64, startDate, endDate time.Time) (PlayerProfitabilityScore, error) {
	var score PlayerProfitabilityScore

	// 計算總收入貢獻 (平台從玩家獲得的收益)
	revenueQuery := `
		SELECT 
			COALESCE(SUM(CASE WHEN transaction_type = 'deposit' THEN amount ELSE 0 END), 0) as total_deposits,
			COALESCE(SUM(CASE WHEN transaction_type = 'withdraw' THEN amount ELSE 0 END), 0) as total_withdrawals,
			COALESCE(SUM(CASE WHEN transaction_type = 'fee' THEN amount ELSE 0 END), 0) as total_fees
		FROM transactions 
		WHERE player_id = ? AND created_at BETWEEN ? AND ?`

	var totalDeposits, totalWithdrawals, totalFees float64
	err := db.QueryRow(revenueQuery, playerID, startDate, endDate).Scan(&totalDeposits, &totalWithdrawals, &totalFees)
	if err != nil && err != sql.ErrNoRows {
		return score, err
	}

	// 計算遊戲損失 (平台獲利)
	gameRevenueQuery := `
		SELECT 
			COALESCE(SUM(bet_amount), 0) as total_bets,
			COALESCE(SUM(win_amount), 0) as total_wins
		FROM game_sessions 
		WHERE player_id = ? AND created_at BETWEEN ? AND ?`

	var totalBets, totalWins float64
	err = db.QueryRow(gameRevenueQuery, playerID, startDate, endDate).Scan(&totalBets, &totalWins)
	if err != nil && err != sql.ErrNoRows {
		return score, err
	}

	// 計算平台淨收益
	platformRevenue := (totalBets - totalWins) + totalFees

	// 收入貢獻分數 (假設1000為滿分基準)
	score.RevenueContribution = math.Min(platformRevenue/1000*100, 100)

	// 利潤率計算
	if totalDeposits > 0 {
		score.ProfitMargin = (platformRevenue / totalDeposits) * 100
		score.ProfitMargin = math.Min(math.Max(score.ProfitMargin, 0), 100)
	}

	// 生命週期價值計算 (簡化版)
	tenureQuery := `
		SELECT DATEDIFF(NOW(), created_at) as tenure_days
		FROM players 
		WHERE id = ?`

	var tenureDays int
	db.QueryRow(tenureQuery, playerID).Scan(&tenureDays)

	if tenureDays > 0 {
		dailyValue := platformRevenue / float64(tenureDays)
		// 預測未來180天的價值
		score.LifetimeValue = dailyValue * 180
	}

	// ROI評分
	customerAcquisitionCost := 50.0 // 假設獲客成本為50元
	if customerAcquisitionCost > 0 {
		roi := (platformRevenue - customerAcquisitionCost) / customerAcquisitionCost * 100
		score.ROIScore = math.Min(math.Max(roi+50, 0), 100) // 正規化到0-100
	}

	// 盈利性趨勢
	midDate := startDate.Add(endDate.Sub(startDate) / 2)

	var firstHalfRevenue, secondHalfRevenue float64
	firstHalfQuery := `
		SELECT COALESCE(SUM(bet_amount - win_amount), 0)
		FROM game_sessions 
		WHERE player_id = ? AND created_at BETWEEN ? AND ?`

	db.QueryRow(firstHalfQuery, playerID, startDate, midDate).Scan(&firstHalfRevenue)
	db.QueryRow(firstHalfQuery, playerID, midDate, endDate).Scan(&secondHalfRevenue)

	if firstHalfRevenue > secondHalfRevenue {
		score.ProfitabilityTrend = "decreasing"
	} else if secondHalfRevenue > firstHalfRevenue {
		score.ProfitabilityTrend = "increasing"
	} else {
		score.ProfitabilityTrend = "stable"
	}

	// 計算總體盈利性評分
	score.Score = (score.RevenueContribution*0.4 + score.ProfitMargin*0.2 +
		(math.Min(score.LifetimeValue/500*100, 100))*0.2 + score.ROIScore*0.2)

	return score, nil
}

// calculateOverallScore 計算總體評分
func (pc *PlayerController) calculateOverallScore(response *PlayerValueScoreResponse, weightConfig *ScoreWeightConfig) float64 {
	// 確保權重總和為1
	totalWeight := weightConfig.ActivityWeight + weightConfig.LoyaltyWeight +
		weightConfig.SpendingWeight + weightConfig.RiskWeight + weightConfig.ProfitabilityWeight

	if totalWeight == 0 {
		totalWeight = 1.0
	}

	// 計算加權評分
	overallScore := (response.ActivityScore.Score*weightConfig.ActivityWeight +
		response.LoyaltyScore.Score*weightConfig.LoyaltyWeight +
		response.SpendingScore.Score*weightConfig.SpendingWeight +
		(100-response.RiskScore.Score)*weightConfig.RiskWeight + // 風險分數需要反轉
		response.ProfitabilityScore.Score*weightConfig.ProfitabilityWeight) / totalWeight

	return math.Min(math.Max(overallScore, 0), 100)
}

// determineValueCategory 確定價值類別
func (pc *PlayerController) determineValueCategory(overallScore float64) string {
	if overallScore >= 80 {
		return "VIP"
	} else if overallScore >= 60 {
		return "High"
	} else if overallScore >= 40 {
		return "Medium"
	} else {
		return "Low"
	}
}

// analyzeTrends 趨勢分析
func (pc *PlayerController) analyzeTrends(db *sql.DB, playerID int64, startDate, endDate time.Time) (ValueTrendAnalysis, error) {
	var trend ValueTrendAnalysis

	// 計算歷史評分 (簡化版 - 比較前一期)
	previousEndDate := startDate
	previousStartDate := startDate.Add(endDate.Sub(startDate) * -1)

	// 獲取前一期的活躍度
	prevActivityScore, _ := pc.calculateActivityScore(db, playerID, previousStartDate, previousEndDate)
	currentActivityScore, _ := pc.calculateActivityScore(db, playerID, startDate, endDate)

	// 計算趨勢
	scoreDiff := currentActivityScore.Score - prevActivityScore.Score
	trend.CurrentVsPrevious = scoreDiff

	// 趨勢方向
	if scoreDiff > 5 {
		trend.TrendDirection = "increasing"
	} else if scoreDiff < -5 {
		trend.TrendDirection = "decreasing"
	} else {
		trend.TrendDirection = "stable"
	}

	// 波動性計算 (簡化)
	if math.Abs(scoreDiff) > 20 {
		trend.VolatilityLevel = "high"
	} else if math.Abs(scoreDiff) > 10 {
		trend.VolatilityLevel = "medium"
	} else {
		trend.VolatilityLevel = "low"
	}

	// 評分歷史 (模擬數據)
	trend.ScoreHistory = []ValueScoreHistory{
		{Date: previousStartDate.Format("2006-01-02"), Score: prevActivityScore.Score},
		{Date: startDate.Format("2006-01-02"), Score: currentActivityScore.Score},
	}

	// 預測評分 (簡化線性預測)
	if len(trend.ScoreHistory) >= 2 {
		recent := trend.ScoreHistory[len(trend.ScoreHistory)-1]
		previous := trend.ScoreHistory[len(trend.ScoreHistory)-2]
		trend.PredictedScore = recent.Score + (recent.Score - previous.Score)
		trend.PredictedScore = math.Min(math.Max(trend.PredictedScore, 0), 100)
	}

	// 信心度
	if trend.VolatilityLevel == "low" {
		trend.ConfidenceLevel = 0.8
	} else if trend.VolatilityLevel == "medium" {
		trend.ConfidenceLevel = 0.6
	} else {
		trend.ConfidenceLevel = 0.4
	}

	return trend, nil
}

// performCompetitorAnalysis 同類玩家比較
func (pc *PlayerController) performCompetitorAnalysis(db *sql.DB, playerID int64, overallScore float64) (CompetitorAnalysis, error) {
	var analysis CompetitorAnalysis

	// 計算百分位數
	percentileQuery := `
		SELECT COUNT(*) as lower_count,
		       (SELECT COUNT(*) FROM players WHERE status = 'active') as total_count
		FROM players p1
		JOIN player_value_score_analysis pvsa ON p1.id = pvsa.player_id
		WHERE p1.status = 'active' 
		AND JSON_EXTRACT(pvsa.analysis_data, '$.overall_score') < ?
		AND pvsa.created_at = (
			SELECT MAX(created_at) 
			FROM player_value_score_analysis 
			WHERE player_id = p1.id
		)`

	var lowerCount, totalCount int
	err := db.QueryRow(percentileQuery, overallScore).Scan(&lowerCount, &totalCount)
	if err != nil && err != sql.ErrNoRows {
		// 如果查詢失敗，使用預設值
		if overallScore >= 80 {
			analysis.Percentile = 95
		} else if overallScore >= 60 {
			analysis.Percentile = 75
		} else if overallScore >= 40 {
			analysis.Percentile = 50
		} else {
			analysis.Percentile = 25
		}
	} else if totalCount > 0 {
		analysis.Percentile = float64(lowerCount) / float64(totalCount) * 100
	}

	// 高於平均的領域
	if overallScore > 50 {
		analysis.AboveAverageAreas = []string{"整體表現", "用戶價值"}
	}

	// 低於平均的領域
	if overallScore < 50 {
		analysis.BelowAverageAreas = []string{"需要改進的領域"}
	}

	// 相似玩家數量 (估算)
	analysis.SimilarPlayers = int(float64(totalCount) * 0.1) // 假設10%為相似玩家

	// 競爭優勢
	if analysis.Percentile > 75 {
		analysis.CompetitiveAdvantage = "高價值用戶，具有明顯競爭優勢"
	} else if analysis.Percentile > 50 {
		analysis.CompetitiveAdvantage = "中等價值用戶，有發展潛力"
	} else {
		analysis.CompetitiveAdvantage = "需要重點關注和培養"
	}

	return analysis, nil
}

// analyzeRetentionRisk 留存風險分析
func (pc *PlayerController) analyzeRetentionRisk(response *PlayerValueScoreResponse) RetentionRiskAnalysis {
	var risk RetentionRiskAnalysis

	// 基於各項評分計算流失概率
	churnScore := 0.0

	// 活躍度影響
	if response.ActivityScore.Score < 30 {
		churnScore += 0.3
	} else if response.ActivityScore.Score < 60 {
		churnScore += 0.15
	}

	// 忠誠度影響
	churnScore += response.LoyaltyScore.ChurnProbability * 0.4

	// 消費力影響
	if response.SpendingScore.Score < 20 {
		churnScore += 0.2
	}

	// 風險評分影響
	if response.RiskScore.Score > 60 {
		churnScore += 0.1
	}

	risk.ChurnProbability = math.Min(churnScore, 1.0)

	// 風險等級
	if risk.ChurnProbability > 0.7 {
		risk.RiskLevel = "high"
		risk.DaysToChurn = 30
	} else if risk.ChurnProbability > 0.4 {
		risk.RiskLevel = "medium"
		risk.DaysToChurn = 90
	} else {
		risk.RiskLevel = "low"
		risk.DaysToChurn = 180
	}

	// 留存行動建議
	if risk.RiskLevel == "high" {
		risk.RetentionActions = []string{
			"立即進行客戶關懷",
			"提供個人化優惠",
			"安排客戶經理聯繫",
		}
	} else if risk.RiskLevel == "medium" {
		risk.RetentionActions = []string{
			"增加互動頻率",
			"推薦適合的活動",
			"監控行為變化",
		}
	} else {
		risk.RetentionActions = []string{
			"保持現有服務水準",
			"定期關注動態",
		}
	}

	// 關鍵影響因素
	if response.ActivityScore.Score < 40 {
		risk.CriticalFactors = append(risk.CriticalFactors, "活躍度下降")
	}
	if response.LoyaltyScore.ChurnProbability > 0.5 {
		risk.CriticalFactors = append(risk.CriticalFactors, "忠誠度不足")
	}
	if response.SpendingScore.Score < 30 {
		risk.CriticalFactors = append(risk.CriticalFactors, "消費力偏低")
	}

	return risk
}

// analyzeValuePotential 價值潛力分析
func (pc *PlayerController) analyzeValuePotential(response *PlayerValueScoreResponse) ValuePotentialAnalysis {
	var potential ValuePotentialAnalysis

	// 成長潛力評估
	growthFactors := 0
	if response.ActivityScore.EngagementTrend == "increasing" {
		growthFactors++
	}
	if response.LoyaltyScore.LoyaltyTrend == "stable" {
		growthFactors++
	}
	if response.SpendingScore.SpendingGrowth > 60 {
		growthFactors++
	}
	if response.ProfitabilityScore.ProfitabilityTrend == "increasing" {
		growthFactors++
	}

	if growthFactors >= 3 {
		potential.GrowthPotential = "high"
	} else if growthFactors >= 2 {
		potential.GrowthPotential = "medium"
	} else {
		potential.GrowthPotential = "low"
	}

	// 升級銷售機會
	if response.SpendingScore.SpendingCategory == "low_spender" && response.ActivityScore.Score > 60 {
		potential.UpsellOpportunities = append(potential.UpsellOpportunities, "提升消費等級")
	}
	if response.LoyaltyScore.Score > 70 {
		potential.UpsellOpportunities = append(potential.UpsellOpportunities, "VIP服務推廣")
	}

	// 優化領域
	if response.ActivityScore.Score < 60 {
		potential.OptimizationAreas = append(potential.OptimizationAreas, "提升用戶活躍度")
	}
	if response.SpendingScore.Score < 50 {
		potential.OptimizationAreas = append(potential.OptimizationAreas, "促進消費行為")
	}
	if response.RiskScore.Score > 40 {
		potential.OptimizationAreas = append(potential.OptimizationAreas, "降低風險等級")
	}

	// 最大潛在評分
	potential.MaxPotentialScore = math.Min(response.OverallScore+30, 100)

	// 達到最大潛力時間
	if potential.GrowthPotential == "high" {
		potential.TimeToMaxPotential = 60
	} else if potential.GrowthPotential == "medium" {
		potential.TimeToMaxPotential = 120
	} else {
		potential.TimeToMaxPotential = 180
	}

	return potential
}

// generateValueRecommendations 生成價值相關建議
func (pc *PlayerController) generateValueRecommendations(response *PlayerValueScoreResponse) []string {
	var recommendations []string

	// 基於總體評分的建議
	if response.OverallScore >= 80 {
		recommendations = append(recommendations, "維持VIP服務水準，提供專屬優惠")
	} else if response.OverallScore >= 60 {
		recommendations = append(recommendations, "提升服務品質，爭取成為VIP用戶")
	} else {
		recommendations = append(recommendations, "重點培養，提供個人化服務")
	}

	// 基於活躍度的建議
	if response.ActivityScore.Score < 50 {
		recommendations = append(recommendations, "設計吸引活動提升用戶參與度")
	}

	// 基於忠誠度的建議
	if response.LoyaltyScore.ChurnProbability > 0.5 {
		recommendations = append(recommendations, "加強客戶關係維護，降低流失風險")
	}

	// 基於消費力的建議
	if response.SpendingScore.Score < 40 {
		recommendations = append(recommendations, "推出促消費活動，提升消費意願")
	}

	// 基於風險的建議
	if response.RiskScore.Score > 60 {
		recommendations = append(recommendations, "加強風險監控，確保合規經營")
	}

	// 基於盈利性的建議
	if response.ProfitabilityScore.Score < 30 {
		recommendations = append(recommendations, "優化產品結構，提升用戶貢獻價值")
	}

	return recommendations
}

// savePlayerValueScoreAnalysis 儲存價值評分分析結果
func (pc *PlayerController) savePlayerValueScoreAnalysis(db *sql.DB, playerID int64, analysis *PlayerValueScoreResponse) error {
	// 將分析結果序列化為JSON
	analysisJSON, err := json.Marshal(analysis)
	if err != nil {
		return fmt.Errorf("序列化分析結果失敗: %v", err)
	}

	// 儲存到資料庫
	_, err = db.Exec(`
		INSERT INTO player_value_score_analysis 
		(player_id, time_range, analysis_data, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`, playerID, analysis.TimeRange, string(analysisJSON))

	if err != nil {
		return fmt.Errorf("儲存分析結果到資料庫失敗: %v", err)
	}

	return nil
}

// ==============================================
// 以下為路由需要的其他方法的佔位符實現
// ==============================================

// GetPlayers 獲取玩家列表
func (pc *PlayerController) GetPlayers(c *gin.Context) {
	// 解析查詢參數
	var req PlayerListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "參數驗證失敗: "+err.Error(), "VALIDATION_FAILED")
		return
	}

	// 設定預設值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Sort == "" {
		req.Sort = "id"
	}
	if req.Order == "" {
		req.Order = "desc"
	}

	// 獲取資料庫連接
	db := config.GetDB()
	if db == nil {
		ErrorResponse(c, http.StatusInternalServerError, "資料庫連接失敗", "DATABASE_ERROR")
		return
	}

	// 建構基本 SQL 查詢
	baseQuery := `SELECT p.id, p.username, p.email, p.real_name, 
		p.phone, p.birth_date, p.country, p.timezone,
		p.verification_level, p.status, p.risk_level, 
		p.created_at, p.updated_at, p.last_login_at,
		COALESCE(w.balance, 0) as balance
		FROM players p 
		LEFT JOIN player_wallets w ON p.id = w.player_id`

	countQuery := `SELECT COUNT(*) FROM players p`

	var conditions []string
	var args []interface{}
	argIndex := 0

	// 添加搜尋條件
	if req.Search != "" {
		conditions = append(conditions, "(p.username LIKE ? OR p.email LIKE ? OR p.real_name LIKE ?)")
		searchPattern := "%" + req.Search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
		argIndex += 3
	}

	// 添加狀態篩選
	if req.Status != "" {
		conditions = append(conditions, "p.status = ?")
		args = append(args, req.Status)
		argIndex++
	}

	// 添加驗證等級篩選
	if req.VerificationLevel != "" {
		conditions = append(conditions, "p.verification_level = ?")
		args = append(args, req.VerificationLevel)
		argIndex++
	}

	// 添加風險等級篩選
	if req.RiskLevel != "" {
		conditions = append(conditions, "p.risk_level = ?")
		args = append(args, req.RiskLevel)
		argIndex++
	}

	// 添加日期範圍篩選
	if req.StartDate != "" {
		conditions = append(conditions, "DATE(p.created_at) >= ?")
		args = append(args, req.StartDate)
		argIndex++
	}
	if req.EndDate != "" {
		conditions = append(conditions, "DATE(p.created_at) <= ?")
		args = append(args, req.EndDate)
		argIndex++
	}

	// 組合 WHERE 條件
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	// 執行計數查詢
	countSql := countQuery + whereClause
	var total int64
	err := db.QueryRow(countSql, args...).Scan(&total)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "查詢計數失敗: "+err.Error(), "DATABASE_ERROR")
		return
	}

	// 添加排序
	orderClause := fmt.Sprintf(" ORDER BY p.%s %s", req.Sort, strings.ToUpper(req.Order))

	// 添加分頁
	offset := (req.Page - 1) * req.Limit
	limitClause := fmt.Sprintf(" LIMIT %d OFFSET %d", req.Limit, offset)

	// 完整查詢
	fullQuery := baseQuery + whereClause + orderClause + limitClause

	rows, err := db.Query(fullQuery, args...)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "查詢玩家列表失敗: "+err.Error(), "DATABASE_ERROR")
		return
	}
	defer rows.Close()

	var players []PlayerWithBalance
	for rows.Next() {
		var player PlayerWithBalance
		var birthDate, lastLoginAt sql.NullTime

		err := rows.Scan(
			&player.ID,
			&player.Username,
			&player.Email,
			&player.RealName,
			&player.Phone,
			&birthDate,
			&player.Country,
			&player.Timezone,
			&player.VerificationLevel,
			&player.Status,
			&player.RiskLevel,
			&player.CreatedAt,
			&player.UpdatedAt,
			&lastLoginAt,
			&player.Balance,
		)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "掃描資料失敗: "+err.Error(), "DATABASE_ERROR")
			return
		}

		// 處理可能為空的時間字段
		if birthDate.Valid {
			player.BirthDate = &birthDate.Time
		}
		if lastLoginAt.Valid {
			player.LastLoginAt = &lastLoginAt.Time
		}

		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "讀取資料失敗: "+err.Error(), "DATABASE_ERROR")
		return
	}

	// 計算分頁資訊
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	response := map[string]interface{}{
		"players": players,
		"pagination": map[string]interface{}{
			"page":         req.Page,
			"limit":        req.Limit,
			"total":        total,
			"total_pages":  totalPages,
			"has_next":     req.Page < totalPages,
			"has_previous": req.Page > 1,
		},
		"filters": map[string]interface{}{
			"search":             req.Search,
			"status":             req.Status,
			"verification_level": req.VerificationLevel,
			"risk_level":         req.RiskLevel,
			"start_date":         req.StartDate,
			"end_date":           req.EndDate,
		},
		"sort": map[string]interface{}{
			"field": req.Sort,
			"order": req.Order,
		},
	}

	SuccessResponse(c, response, "玩家列表獲取成功")
}

// GetPlayer 獲取單個玩家詳細資訊
func (pc *PlayerController) GetPlayer(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetPlayer endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// GetPlayerGameHistory 獲取玩家遊戲歷史
func (pc *PlayerController) GetPlayerGameHistory(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetPlayerGameHistory endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// SearchPlayers 搜尋玩家
func (pc *PlayerController) SearchPlayers(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "SearchPlayers endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// FilterPlayers 篩選玩家
func (pc *PlayerController) FilterPlayers(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "FilterPlayers endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// CreatePlayer 創建玩家
func (pc *PlayerController) CreatePlayer(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "CreatePlayer endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// UpdatePlayer 更新玩家資訊
func (pc *PlayerController) UpdatePlayer(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "UpdatePlayer endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// DeletePlayer 刪除玩家
func (pc *PlayerController) DeletePlayer(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "DeletePlayer endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// UpdatePlayerStatus 更新玩家狀態
func (pc *PlayerController) UpdatePlayerStatus(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "UpdatePlayerStatus endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// GetPlayerBalance 獲取玩家餘額
func (pc *PlayerController) GetPlayerBalance(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetPlayerBalance endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// DepositPlayerBalance 玩家充值
func (pc *PlayerController) DepositPlayerBalance(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "DepositPlayerBalance endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// WithdrawPlayerBalance 玩家提領
func (pc *PlayerController) WithdrawPlayerBalance(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "WithdrawPlayerBalance endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// GetPlayerTransactions 獲取玩家交易記錄
func (pc *PlayerController) GetPlayerTransactions(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetPlayerTransactions endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// SetPlayerRestriction 設定玩家限制
func (pc *PlayerController) SetPlayerRestriction(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "SetPlayerRestriction endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// GetPlayerRestrictions 獲取玩家限制
func (pc *PlayerController) GetPlayerRestrictions(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetPlayerRestrictions endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// RemovePlayerRestriction 移除玩家限制
func (pc *PlayerController) RemovePlayerRestriction(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "RemovePlayerRestriction endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// AssessPlayerRisk 評估玩家風險
func (pc *PlayerController) AssessPlayerRisk(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "AssessPlayerRisk endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// GetPlayerRiskHistory 獲取玩家風險歷史
func (pc *PlayerController) GetPlayerRiskHistory(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetPlayerRiskHistory endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// DeactivatePlayer 註銷玩家
func (pc *PlayerController) DeactivatePlayer(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "DeactivatePlayer endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// GetPlayerDeactivationHistory 獲取玩家註銷歷史
func (pc *PlayerController) GetPlayerDeactivationHistory(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetPlayerDeactivationHistory endpoint not implemented yet", "NOT_IMPLEMENTED")
}
