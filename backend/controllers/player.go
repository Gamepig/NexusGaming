package controllers

import (
	"database/sql"
	"fmt"
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
	Page              int    `form:"page" binding:"min=1"`                                                             // 頁碼，從1開始
	Limit             int    `form:"limit" binding:"min=1,max=100"`                                                    // 每頁數量，最大100
	Sort              string `form:"sort" binding:"oneof=id username email created_at updated_at total_bet total_win"` // 排序字段
	Order             string `form:"order" binding:"oneof=asc desc"`                                                   // 排序順序
	Search            string `form:"search"`                                                                           // 搜尋關鍵字（姓名、用戶名、郵箱）
	Status            string `form:"status" binding:"omitempty,oneof=active inactive suspended deleted"`               // 狀態篩選
	StartDate         string `form:"start_date"`                                                                       // 註冊開始日期 (YYYY-MM-DD)
	EndDate           string `form:"end_date"`                                                                         // 註冊結束日期 (YYYY-MM-DD)
	MinBalance        string `form:"min_balance"`                                                                      // 最小餘額
	MaxBalance        string `form:"max_balance"`                                                                      // 最大餘額
	VerificationLevel string `form:"verification_level" binding:"omitempty,oneof=none email phone identity"`           // 驗證等級
	RiskLevel         string `form:"risk_level" binding:"omitempty,oneof=low medium high"`                             // 風險等級
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

// PlayerListResponse 玩家列表回應
type PlayerListResponse struct {
	Players    []PlayerWithBalance `json:"players"`     // 玩家列表
	Total      int64               `json:"total"`       // 總數量
	Page       int                 `json:"page"`        // 當前頁碼
	Limit      int                 `json:"limit"`       // 每頁數量
	TotalPages int                 `json:"total_pages"` // 總頁數
	Sort       string              `json:"sort"`        // 排序字段
	Order      string              `json:"order"`       // 排序順序
}

// PlayerGameHistoryRequest 玩家遊戲歷史查詢請求
type PlayerGameHistoryRequest struct {
	Page      int    `form:"page" binding:"min=1"`                                                  // 頁碼，從1開始
	Limit     int    `form:"limit" binding:"min=1,max=100"`                                         // 每頁數量，最大100
	Sort      string `form:"sort" binding:"omitempty,oneof=played_at bet_amount win_amount result"` // 排序字段
	Order     string `form:"order" binding:"omitempty,oneof=asc desc"`                              // 排序順序
	GameType  string `form:"game_type" binding:"omitempty,oneof=texas_holdem stud_poker baccarat"`  // 遊戲類型篩選
	Result    string `form:"result" binding:"omitempty,oneof=win loss draw"`                        // 結果篩選
	StartDate string `form:"start_date" binding:"omitempty" example:"2024-01-01"`                   // 開始日期 (YYYY-MM-DD)
	EndDate   string `form:"end_date" binding:"omitempty" example:"2024-12-31"`                     // 結束日期 (YYYY-MM-DD)
	MinBet    string `form:"min_bet" binding:"omitempty,numeric"`                                   // 最小下注金額
	MaxBet    string `form:"max_bet" binding:"omitempty,numeric"`                                   // 最大下注金額
}

// GameHistoryDetail 遊戲歷史詳細資訊
type GameHistoryDetail struct {
	ID          int64     `json:"id"`           // 遊戲參與記錄ID
	SessionID   int64     `json:"session_id"`   // 遊戲場次ID
	GameType    string    `json:"game_type"`    // 遊戲類型
	GameName    string    `json:"game_name"`    // 遊戲名稱
	BetAmount   float64   `json:"bet_amount"`   // 下注金額
	WinAmount   float64   `json:"win_amount"`   // 獲勝金額
	NetAmount   float64   `json:"net_amount"`   // 淨輸贏金額 (win_amount - bet_amount)
	Result      string    `json:"result"`       // 結果 (win/loss/draw)
	PlayedAt    time.Time `json:"played_at"`    // 遊戲時間
	Duration    int       `json:"duration"`     // 遊戲時長（秒）
	GameDetails string    `json:"game_details"` // 遊戲詳細資訊（JSON格式）
}

// PlayerGameHistoryResponse 玩家遊戲歷史回應
type PlayerGameHistoryResponse struct {
	Games      []GameHistoryDetail `json:"games"`       // 遊戲歷史列表
	Total      int64               `json:"total"`       // 總數量
	Page       int                 `json:"page"`        // 當前頁碼
	Limit      int                 `json:"limit"`       // 每頁數量
	TotalPages int                 `json:"total_pages"` // 總頁數
	Sort       string              `json:"sort"`        // 排序字段
	Order      string              `json:"order"`       // 排序順序
	Summary    GameHistorySummary  `json:"summary"`     // 統計摘要
}

// GameHistorySummary 遊戲歷史統計摘要
type GameHistorySummary struct {
	TotalGames  int     `json:"total_games"`  // 總遊戲次數
	TotalBet    float64 `json:"total_bet"`    // 總下注金額
	TotalWin    float64 `json:"total_win"`    // 總獲勝金額
	NetProfit   float64 `json:"net_profit"`   // 淨利潤
	WinRate     float64 `json:"win_rate"`     // 勝率
	AverageBet  float64 `json:"average_bet"`  // 平均下注金額
	BiggestWin  float64 `json:"biggest_win"`  // 最大單次獲勝
	BiggestLoss float64 `json:"biggest_loss"` // 最大單次損失
}

// GetPlayers 取得玩家列表
// @Summary 取得玩家列表
// @Description 取得玩家列表，支援分頁、排序、搜尋、篩選功能
// @Tags 玩家管理
// @Accept json
// @Produce json
// @Param page query int false "頁碼" default(1) minimum(1)
// @Param limit query int false "每頁數量" default(20) minimum(1) maximum(100)
// @Param sort query string false "排序字段" Enums(id,username,email,created_at,updated_at,total_bet,total_win) default(id)
// @Param order query string false "排序順序" Enums(asc,desc) default(desc)
// @Param search query string false "搜尋關鍵字（姓名、用戶名、郵箱）"
// @Param status query string false "狀態篩選" Enums(active,inactive,suspended,deleted)
// @Param start_date query string false "註冊開始日期 (YYYY-MM-DD)"
// @Param end_date query string false "註冊結束日期 (YYYY-MM-DD)"
// @Param min_balance query string false "最小餘額"
// @Param max_balance query string false "最大餘額"
// @Param verification_level query string false "驗證等級" Enums(none,email,phone,identity)
// @Param risk_level query string false "風險等級" Enums(low,medium,high)
// @Security BearerAuth
// @Success 200 {object} APIResponse{data=PlayerListResponse}
// @Failure 400 {object} APIResponse
// @Failure 401 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/v1/players [get]
func (pc *PlayerController) GetPlayers(c *gin.Context) {
	var req PlayerListRequest

	// 設定預設值
	req.Page = 1
	req.Limit = 20
	req.Sort = "id"
	req.Order = "desc"

	// 綁定查詢參數
	if err := c.ShouldBindQuery(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters: "+err.Error(), "INVALID_PARAMS")
		return
	}

	// 驗證日期格式
	if req.StartDate != "" {
		if _, err := time.Parse("2006-01-02", req.StartDate); err != nil {
			ErrorResponse(c, http.StatusBadRequest, "Invalid start_date format, use YYYY-MM-DD", "INVALID_DATE_FORMAT")
			return
		}
	}
	if req.EndDate != "" {
		if _, err := time.Parse("2006-01-02", req.EndDate); err != nil {
			ErrorResponse(c, http.StatusBadRequest, "Invalid end_date format, use YYYY-MM-DD", "INVALID_DATE_FORMAT")
			return
		}
	}

	// 驗證餘額格式
	var minBalance, maxBalance *float64
	if req.MinBalance != "" {
		if val, err := strconv.ParseFloat(req.MinBalance, 64); err != nil {
			ErrorResponse(c, http.StatusBadRequest, "Invalid min_balance format", "INVALID_BALANCE_FORMAT")
			return
		} else {
			minBalance = &val
		}
	}
	if req.MaxBalance != "" {
		if val, err := strconv.ParseFloat(req.MaxBalance, 64); err != nil {
			ErrorResponse(c, http.StatusBadRequest, "Invalid max_balance format", "INVALID_BALANCE_FORMAT")
			return
		} else {
			maxBalance = &val
		}
	}

	// 建立查詢
	db := config.GetDB()
	if db == nil {
		ErrorResponse(c, http.StatusInternalServerError, "Database connection not available", "DB_CONNECTION_ERROR")
		return
	}

	// 建立基本查詢
	baseQuery := `FROM players p 
	LEFT JOIN player_wallets pw ON p.id = pw.player_id`

	whereClause := []string{"1=1"}
	args := []interface{}{}

	// 新增搜尋條件
	if req.Search != "" {
		whereClause = append(whereClause, "(p.username LIKE ? OR p.email LIKE ? OR p.real_name LIKE ?)")
		searchPattern := "%" + req.Search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	// 新增狀態篩選
	if req.Status != "" {
		whereClause = append(whereClause, "p.status = ?")
		args = append(args, req.Status)
	}

	// 新增日期範圍篩選
	if req.StartDate != "" {
		whereClause = append(whereClause, "DATE(p.created_at) >= ?")
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		whereClause = append(whereClause, "DATE(p.created_at) <= ?")
		args = append(args, req.EndDate)
	}

	// 新增餘額範圍篩選
	if minBalance != nil {
		whereClause = append(whereClause, "pw.balance >= ?")
		args = append(args, *minBalance)
	}
	if maxBalance != nil {
		whereClause = append(whereClause, "pw.balance <= ?")
		args = append(args, *maxBalance)
	}

	// 新增驗證等級篩選
	if req.VerificationLevel != "" {
		whereClause = append(whereClause, "p.verification_level = ?")
		args = append(args, req.VerificationLevel)
	}

	// 新增風險等級篩選
	if req.RiskLevel != "" {
		whereClause = append(whereClause, "p.risk_level = ?")
		args = append(args, req.RiskLevel)
	}

	whereSQL := " WHERE " + strings.Join(whereClause, " AND ")

	// 計算總數
	countQuery := "SELECT COUNT(DISTINCT p.id) " + baseQuery + whereSQL
	var total int64
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to count players: "+err.Error(), "DB_QUERY_ERROR")
		return
	}

	// 計算總頁數
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	// 建立查詢語句
	selectFields := `p.id, p.username, p.email, p.real_name, p.phone, p.status, 
	p.verification_level, p.risk_level, p.total_bet, p.total_win, 
	p.last_login_at, p.created_at, p.updated_at,
	COALESCE(pw.balance, 0) as balance`

	orderBy := fmt.Sprintf(" ORDER BY %s %s", req.Sort, strings.ToUpper(req.Order))
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d", req.Limit, (req.Page-1)*req.Limit)

	query := "SELECT " + selectFields + " " + baseQuery + whereSQL + orderBy + limit

	// 執行查詢
	rows, err := db.Query(query, args...)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to query players: "+err.Error(), "DB_QUERY_ERROR")
		return
	}
	defer rows.Close()

	// 解析結果
	var players []PlayerWithBalance
	for rows.Next() {
		var player models.Player
		var balance float64
		var lastLoginAt sql.NullTime

		err := rows.Scan(
			&player.ID, &player.Username, &player.Email, &player.RealName, &player.Phone,
			&player.Status, &player.VerificationLevel, &player.RiskLevel,
			&player.TotalBet, &player.TotalWin, &lastLoginAt, &player.CreatedAt, &player.UpdatedAt,
			&balance,
		)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, "Failed to scan player data: "+err.Error(), "DATA_SCAN_ERROR")
			return
		}

		if lastLoginAt.Valid {
			player.LastLoginAt = &lastLoginAt.Time
		}

		// 組合玩家資料和餘額
		playerWithBalance := PlayerWithBalance{
			Player:  player,
			Balance: balance,
		}

		players = append(players, playerWithBalance)
	}

	if err = rows.Err(); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Row iteration error: "+err.Error(), "DB_ITERATION_ERROR")
		return
	}

	// 準備回應
	response := PlayerListResponse{
		Players:    players,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
		Sort:       req.Sort,
		Order:      req.Order,
	}

	SuccessResponse(c, response, "Players retrieved successfully")
}

// GetPlayer 取得特定玩家資訊
// @Summary 取得玩家詳細資訊
// @Description 根據玩家ID取得詳細的玩家資訊，包括基本資料、錢包狀態、標籤、限制等
// @Tags 玩家管理
// @Accept json
// @Produce json
// @Param id path int true "玩家ID"
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Success 200 {object} APIResponse{data=PlayerDetailResponse} "成功取得玩家詳細資訊"
// @Failure 400 {object} APIResponse "請求參數錯誤"
// @Failure 401 {object} APIResponse "未授權"
// @Failure 404 {object} APIResponse "玩家不存在"
// @Failure 500 {object} APIResponse "伺服器內部錯誤"
// @Router /api/v1/players/{id} [get]
func (pc *PlayerController) GetPlayer(c *gin.Context) {
	// 解析玩家ID
	playerIDStr := c.Param("id")
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid player ID format", "INVALID_PLAYER_ID")
		return
	}

	db := config.GetDB()

	// 查詢玩家基本資訊和錢包餘額
	var player models.Player
	var balance float64
	var lastLoginAt sql.NullTime

	playerQuery := `
		SELECT p.id, p.username, p.email, p.real_name, p.phone, p.status, 
		       p.verification_level, p.risk_level, p.total_bet, p.total_win, 
		       p.last_login_at, p.created_at, p.updated_at,
		       COALESCE(pw.balance, 0) as balance
		FROM players p
		LEFT JOIN player_wallets pw ON p.id = pw.player_id
		WHERE p.id = ?`

	err = db.QueryRow(playerQuery, playerID).Scan(
		&player.ID, &player.Username, &player.Email, &player.RealName, &player.Phone,
		&player.Status, &player.VerificationLevel, &player.RiskLevel,
		&player.TotalBet, &player.TotalWin, &lastLoginAt, &player.CreatedAt, &player.UpdatedAt,
		&balance,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			ErrorResponse(c, http.StatusNotFound, "Player not found", "PLAYER_NOT_FOUND")
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Failed to query player: "+err.Error(), "DATABASE_ERROR")
		return
	}

	if lastLoginAt.Valid {
		player.LastLoginAt = &lastLoginAt.Time
	}

	// 建立基本回應物件
	playerWithBalance := PlayerWithBalance{
		Player:  player,
		Balance: balance,
	}

	// 查詢玩家標籤
	tags, err := pc.getPlayerTags(db, playerID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to query player tags: "+err.Error(), "DATABASE_ERROR")
		return
	}

	// 查詢玩家限制
	restrictions, err := pc.getPlayerRestrictions(db, playerID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to query player restrictions: "+err.Error(), "DATABASE_ERROR")
		return
	}

	// 計算玩家統計數據
	statistics, err := pc.getPlayerStatistics(db, playerID, player.CreatedAt)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to calculate player statistics: "+err.Error(), "DATABASE_ERROR")
		return
	}

	// 查詢最近遊戲記錄
	recentGames, err := pc.getRecentGames(db, playerID, 10) // 最近10場遊戲
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to query recent games: "+err.Error(), "DATABASE_ERROR")
		return
	}

	// 建立詳細回應
	response := PlayerDetailResponse{
		PlayerWithBalance: playerWithBalance,
		Tags:              tags,
		Restrictions:      restrictions,
		Statistics:        statistics,
		RecentGames:       recentGames,
	}

	SuccessResponse(c, response, "Player details retrieved successfully")
}

// GetPlayerGameHistory 取得玩家遊戲歷史
// @Summary 取得玩家遊戲歷史記錄
// @Description 查詢特定玩家的詳細遊戲歷史，支援分頁、排序、篩選和統計摘要
// @Tags 玩家管理
// @Accept json
// @Produce json
// @Param id path int true "玩家ID"
// @Param page query int false "頁碼（預設1）" default(1)
// @Param limit query int false "每頁數量（預設20，最大100）" default(20)
// @Param sort query string false "排序字段" Enums(played_at,bet_amount,win_amount,result) default(played_at)
// @Param order query string false "排序順序" Enums(asc,desc) default(desc)
// @Param game_type query string false "遊戲類型篩選" Enums(texas_holdem,stud_poker,baccarat)
// @Param result query string false "結果篩選" Enums(win,loss,draw)
// @Param start_date query string false "開始日期 (YYYY-MM-DD)" format(date)
// @Param end_date query string false "結束日期 (YYYY-MM-DD)" format(date)
// @Param min_bet query number false "最小下注金額"
// @Param max_bet query number false "最大下注金額"
// @Param Authorization header string true "Bearer token" default(Bearer )
// @Success 200 {object} APIResponse{data=PlayerGameHistoryResponse} "成功取得玩家遊戲歷史"
// @Failure 400 {object} APIResponse "請求參數錯誤"
// @Failure 401 {object} APIResponse "未授權"
// @Failure 404 {object} APIResponse "玩家不存在"
// @Failure 500 {object} APIResponse "伺服器內部錯誤"
// @Router /api/v1/players/{id}/games [get]
func (pc *PlayerController) GetPlayerGameHistory(c *gin.Context) {
	// 解析玩家ID
	playerIDStr := c.Param("id")
	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid player ID format", "INVALID_PLAYER_ID")
		return
	}

	// 解析查詢參數
	var req PlayerGameHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters: "+err.Error(), "INVALID_PARAMETERS")
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
		req.Sort = "played_at"
	}
	if req.Order == "" {
		req.Order = "desc"
	}

	// 取得資料庫連接
	db := config.GetDB()
	if db == nil {
		ErrorResponse(c, http.StatusInternalServerError, "Database connection is not available", "DATABASE_ERROR")
		return
	}

	// 檢查玩家是否存在
	var playerExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM players WHERE id = ?)", playerID).Scan(&playerExists)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to check player existence: "+err.Error(), "DATABASE_ERROR")
		return
	}
	if !playerExists {
		ErrorResponse(c, http.StatusNotFound, "Player not found", "PLAYER_NOT_FOUND")
		return
	}

	// 查詢遊戲歷史
	games, total, err := pc.getPlayerGameHistoryData(db, playerID, req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to query game history: "+err.Error(), "DATABASE_ERROR")
		return
	}

	// 計算統計摘要
	summary, err := pc.getGameHistorySummary(db, playerID, req)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to calculate summary: "+err.Error(), "DATABASE_ERROR")
		return
	}

	// 計算總頁數
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	// 建立回應
	response := PlayerGameHistoryResponse{
		Games:      games,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
		Sort:       req.Sort,
		Order:      req.Order,
		Summary:    summary,
	}

	SuccessResponse(c, response, "Game history retrieved successfully")
}

// SearchPlayers 搜尋玩家 (2.2.3 的預備實現)
func (pc *PlayerController) SearchPlayers(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "SearchPlayers endpoint will be implemented in task 2.2.3", "NOT_IMPLEMENTED")
}

// FilterPlayers 篩選玩家 (2.2.4 的預備實現)
func (pc *PlayerController) FilterPlayers(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "FilterPlayers endpoint will be implemented in task 2.2.4", "NOT_IMPLEMENTED")
}

// getPlayerTags 查詢玩家標籤
func (pc *PlayerController) getPlayerTags(db *sql.DB, playerID int64) ([]PlayerTag, error) {
	query := `
		SELECT pt.id, pt.name, pt.tag_type 
		FROM player_tags pt
		INNER JOIN player_tag_relations ptr ON pt.id = ptr.tag_id
		WHERE ptr.player_id = ?`

	rows, err := db.Query(query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []PlayerTag
	for rows.Next() {
		var tag PlayerTag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Type)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if tags == nil {
		tags = []PlayerTag{} // 返回空陣列而非 nil
	}
	return tags, nil
}

// getPlayerRestrictions 查詢玩家限制
func (pc *PlayerController) getPlayerRestrictions(db *sql.DB, playerID int64) ([]PlayerRestriction, error) {
	query := `
		SELECT id, restriction_type, value, is_active, expires_at
		FROM player_restrictions
		WHERE player_id = ? AND (expires_at IS NULL OR expires_at > NOW())`

	rows, err := db.Query(query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restrictions []PlayerRestriction
	for rows.Next() {
		var restriction PlayerRestriction
		var expiresAt sql.NullTime

		err := rows.Scan(&restriction.ID, &restriction.RestrictionType, &restriction.Value,
			&restriction.IsActive, &expiresAt)
		if err != nil {
			return nil, err
		}

		if expiresAt.Valid {
			restriction.ExpiresAt = &expiresAt.Time
		}

		restrictions = append(restrictions, restriction)
	}

	if restrictions == nil {
		restrictions = []PlayerRestriction{} // 返回空陣列而非 nil
	}
	return restrictions, nil
}

// getPlayerStatistics 計算玩家統計數據
func (pc *PlayerController) getPlayerStatistics(db *sql.DB, playerID int64, createdAt time.Time) (PlayerStatistics, error) {
	var stats PlayerStatistics

	// 計算註冊天數
	stats.DaysRegistered = int(time.Since(createdAt).Hours() / 24)

	// 查詢遊戲統計數據
	gameStatsQuery := `
		SELECT 
			COUNT(*) as total_games,
			COALESCE(AVG(bet_amount), 0) as average_bet,
			COALESCE(MAX(win_amount), 0) as biggest_win,
			COALESCE(MIN(CASE WHEN win_amount < bet_amount THEN win_amount - bet_amount END), 0) as biggest_loss,
			COALESCE(MAX(played_at), '1970-01-01') as last_activity
		FROM game_participations 
		WHERE player_id = ?`

	var lastActivity time.Time
	err := db.QueryRow(gameStatsQuery, playerID).Scan(
		&stats.TotalGames, &stats.AverageBet, &stats.BiggestWin, &stats.BiggestLoss, &lastActivity)
	if err != nil {
		return stats, err
	}

	// 計算最後活動天數
	if !lastActivity.IsZero() && lastActivity.Year() > 1970 {
		stats.LastActivityDays = int(time.Since(lastActivity).Hours() / 24)
	} else {
		stats.LastActivityDays = stats.DaysRegistered // 如果沒有遊戲記錄，使用註冊天數
	}

	// 計算勝率
	if stats.TotalGames > 0 {
		winQuery := `
			SELECT COUNT(*) 
			FROM game_participations 
			WHERE player_id = ? AND win_amount > bet_amount`

		var winCount int
		err = db.QueryRow(winQuery, playerID).Scan(&winCount)
		if err != nil {
			return stats, err
		}
		stats.WinRate = float64(winCount) / float64(stats.TotalGames) * 100
	}

	return stats, nil
}

// getRecentGames 查詢最近遊戲記錄
func (pc *PlayerController) getRecentGames(db *sql.DB, playerID int64, limit int) ([]GameParticipation, error) {
	query := `
		SELECT gp.id, gs.game_type, gp.bet_amount, gp.win_amount, 
		       CASE 
		           WHEN gp.win_amount > gp.bet_amount THEN 'win'
		           WHEN gp.win_amount = gp.bet_amount THEN 'draw'
		           ELSE 'loss'
		       END as result,
		       gp.played_at
		FROM game_participations gp
		INNER JOIN game_sessions gs ON gp.session_id = gs.id
		WHERE gp.player_id = ?
		ORDER BY gp.played_at DESC
		LIMIT ?`

	rows, err := db.Query(query, playerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []GameParticipation
	for rows.Next() {
		var game GameParticipation
		err := rows.Scan(&game.ID, &game.GameType, &game.BetAmount, &game.WinAmount,
			&game.Result, &game.PlayedAt)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	if games == nil {
		games = []GameParticipation{} // 返回空陣列而非 nil
	}
	return games, nil
}

// getPlayerGameHistoryData 查詢玩家遊戲歷史數據
func (pc *PlayerController) getPlayerGameHistoryData(db *sql.DB, playerID int64, req PlayerGameHistoryRequest) ([]GameHistoryDetail, int64, error) {
	// 建立基本查詢和計數查詢
	baseQuery := `
		FROM game_participations gp
		INNER JOIN game_sessions gs ON gp.session_id = gs.id
		WHERE gp.player_id = ?`

	countQuery := "SELECT COUNT(*) " + baseQuery

	// 建立完整查詢
	selectQuery := `
		SELECT gp.id, gp.session_id, gs.game_type, 
			   CASE gs.game_type 
				   WHEN 'texas_holdem' THEN '德州撲克'
				   WHEN 'stud_poker' THEN '梭哈撲克' 
				   WHEN 'baccarat' THEN '百家樂'
				   ELSE gs.game_type 
			   END as game_name,
			   gp.bet_amount, gp.win_amount, 
			   (gp.win_amount - gp.bet_amount) as net_amount,
			   gp.result, gp.played_at, gp.duration, gp.game_details` + baseQuery

	// 建立 WHERE 條件和參數
	var whereConditions []string
	var queryParams []interface{}
	var countParams []interface{}

	// 基本參數
	queryParams = append(queryParams, playerID)
	countParams = append(countParams, playerID)

	// 遊戲類型篩選
	if req.GameType != "" {
		whereConditions = append(whereConditions, "gs.game_type = ?")
		queryParams = append(queryParams, req.GameType)
		countParams = append(countParams, req.GameType)
	}

	// 結果篩選
	if req.Result != "" {
		whereConditions = append(whereConditions, "gp.result = ?")
		queryParams = append(queryParams, req.Result)
		countParams = append(countParams, req.Result)
	}

	// 日期範圍篩選
	if req.StartDate != "" {
		whereConditions = append(whereConditions, "DATE(gp.played_at) >= ?")
		queryParams = append(queryParams, req.StartDate)
		countParams = append(countParams, req.StartDate)
	}
	if req.EndDate != "" {
		whereConditions = append(whereConditions, "DATE(gp.played_at) <= ?")
		queryParams = append(queryParams, req.EndDate)
		countParams = append(countParams, req.EndDate)
	}

	// 下注金額範圍篩選
	if req.MinBet != "" {
		whereConditions = append(whereConditions, "gp.bet_amount >= ?")
		queryParams = append(queryParams, req.MinBet)
		countParams = append(countParams, req.MinBet)
	}
	if req.MaxBet != "" {
		whereConditions = append(whereConditions, "gp.bet_amount <= ?")
		queryParams = append(queryParams, req.MaxBet)
		countParams = append(countParams, req.MaxBet)
	}

	// 添加額外的 WHERE 條件
	if len(whereConditions) > 0 {
		whereClause := " AND " + strings.Join(whereConditions, " AND ")
		selectQuery += whereClause
		countQuery += whereClause
	}

	// 查詢總數
	var total int64
	err := db.QueryRow(countQuery, countParams...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 添加排序和分頁
	orderClause := fmt.Sprintf(" ORDER BY gp.%s %s", req.Sort, strings.ToUpper(req.Order))
	selectQuery += orderClause

	offset := (req.Page - 1) * req.Limit
	limitClause := fmt.Sprintf(" LIMIT %d OFFSET %d", req.Limit, offset)
	selectQuery += limitClause

	// 執行查詢
	rows, err := db.Query(selectQuery, queryParams...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 掃描結果
	var games []GameHistoryDetail
	for rows.Next() {
		var game GameHistoryDetail
		err := rows.Scan(
			&game.ID, &game.SessionID, &game.GameType, &game.GameName,
			&game.BetAmount, &game.WinAmount, &game.NetAmount,
			&game.Result, &game.PlayedAt, &game.Duration, &game.GameDetails,
		)
		if err != nil {
			return nil, 0, err
		}
		games = append(games, game)
	}

	if games == nil {
		games = []GameHistoryDetail{} // 返回空陣列而非 nil
	}

	return games, total, nil
}

// getGameHistorySummary 計算玩家遊戲歷史統計摘要
func (pc *PlayerController) getGameHistorySummary(db *sql.DB, playerID int64, req PlayerGameHistoryRequest) (GameHistorySummary, error) {
	// 建立基本查詢
	baseQuery := `
		FROM game_participations gp
		INNER JOIN game_sessions gs ON gp.session_id = gs.id
		WHERE gp.player_id = ?`

	// 建立 WHERE 條件和參數（與主查詢相同的篩選條件）
	var whereConditions []string
	var params []interface{}
	params = append(params, playerID)

	// 應用與主查詢相同的篩選條件
	if req.GameType != "" {
		whereConditions = append(whereConditions, "gs.game_type = ?")
		params = append(params, req.GameType)
	}
	if req.Result != "" {
		whereConditions = append(whereConditions, "gp.result = ?")
		params = append(params, req.Result)
	}
	if req.StartDate != "" {
		whereConditions = append(whereConditions, "DATE(gp.played_at) >= ?")
		params = append(params, req.StartDate)
	}
	if req.EndDate != "" {
		whereConditions = append(whereConditions, "DATE(gp.played_at) <= ?")
		params = append(params, req.EndDate)
	}
	if req.MinBet != "" {
		whereConditions = append(whereConditions, "gp.bet_amount >= ?")
		params = append(params, req.MinBet)
	}
	if req.MaxBet != "" {
		whereConditions = append(whereConditions, "gp.bet_amount <= ?")
		params = append(params, req.MaxBet)
	}

	// 添加 WHERE 條件
	if len(whereConditions) > 0 {
		baseQuery += " AND " + strings.Join(whereConditions, " AND ")
	}

	// 統計查詢
	summaryQuery := `
		SELECT 
			COUNT(*) as total_games,
			COALESCE(SUM(gp.bet_amount), 0) as total_bet,
			COALESCE(SUM(gp.win_amount), 0) as total_win,
			COALESCE(SUM(gp.win_amount - gp.bet_amount), 0) as net_profit,
			COALESCE(AVG(gp.bet_amount), 0) as average_bet,
			COALESCE(MAX(CASE WHEN gp.result = 'win' THEN gp.win_amount - gp.bet_amount ELSE 0 END), 0) as biggest_win,
			COALESCE(MIN(CASE WHEN gp.result = 'loss' THEN gp.win_amount - gp.bet_amount ELSE 0 END), 0) as biggest_loss,
			COALESCE(COUNT(CASE WHEN gp.result = 'win' THEN 1 END), 0) as win_count` + baseQuery

	var summary GameHistorySummary
	var winCount int

	err := db.QueryRow(summaryQuery, params...).Scan(
		&summary.TotalGames, &summary.TotalBet, &summary.TotalWin,
		&summary.NetProfit, &summary.AverageBet, &summary.BiggestWin,
		&summary.BiggestLoss, &winCount,
	)
	if err != nil {
		return summary, err
	}

	// 計算勝率
	if summary.TotalGames > 0 {
		summary.WinRate = float64(winCount) / float64(summary.TotalGames) * 100
	}

	return summary, nil
}
