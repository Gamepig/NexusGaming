package models

import (
	"time"
)

// Review 評論結構體
type Review struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	User         *User     `json:"user,omitempty"`
	GameID       int       `json:"game_id" db:"game_id"`
	Game         *Game     `json:"game,omitempty"`
	Rating       int       `json:"rating" db:"rating"` // 1-5 星評分
	Title        string    `json:"title" db:"title"`
	Content      string    `json:"content" db:"content"`
	PlayedHours  int       `json:"played_hours" db:"played_hours"`
	Recommend    bool      `json:"recommend" db:"recommend"`
	HelpfulVotes int       `json:"helpful_votes" db:"helpful_votes"`
	Status       string    `json:"status" db:"status"`           // published, draft, flagged, deleted
	IsVerified   bool      `json:"is_verified" db:"is_verified"` // 驗證購買
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ReviewVote 評論投票結構體
type ReviewVote struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	ReviewID  int       `json:"review_id" db:"review_id"`
	VoteType  string    `json:"vote_type" db:"vote_type"` // helpful, not_helpful
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ReviewReport 評論檢舉結構體
type ReviewReport struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	ReviewID  int       `json:"review_id" db:"review_id"`
	Reason    string    `json:"reason" db:"reason"` // spam, inappropriate, fake, other
	Comment   string    `json:"comment" db:"comment"`
	Status    string    `json:"status" db:"status"` // pending, reviewed, dismissed
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// GameRating 遊戲評分統計結構體
type GameRating struct {
	GameID          int       `json:"game_id" db:"game_id"`
	AverageRating   float64   `json:"average_rating" db:"average_rating"`
	TotalReviews    int       `json:"total_reviews" db:"total_reviews"`
	FiveStars       int       `json:"five_stars" db:"five_stars"`
	FourStars       int       `json:"four_stars" db:"four_stars"`
	ThreeStars      int       `json:"three_stars" db:"three_stars"`
	TwoStars        int       `json:"two_stars" db:"two_stars"`
	OneStar         int       `json:"one_star" db:"one_star"`
	RecommendedRate float64   `json:"recommended_rate" db:"recommended_rate"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// ReviewRepository 評論資料存取介面
type ReviewRepository interface {
	Create(review *Review) error
	GetByID(id int) (*Review, error)
	GetByGameID(gameID int, offset, limit int, sortBy string) ([]*Review, error)
	GetByUserID(userID int, offset, limit int) ([]*Review, error)
	Update(review *Review) error
	Delete(id int) error
	UpdateStatus(id int, status string) error
	List(offset, limit int, filters ReviewFilters) ([]*Review, error)
	Count(filters ReviewFilters) (int, error)
	HasUserReviewed(userID, gameID int) (bool, error)
	GetUserReview(userID, gameID int) (*Review, error)
	GetFeaturedReviews(gameID int, limit int) ([]*Review, error)
	UpdateHelpfulVotes(reviewID int, votes int) error
}

// ReviewVoteRepository 評論投票資料存取介面
type ReviewVoteRepository interface {
	Create(vote *ReviewVote) error
	GetByUserAndReview(userID, reviewID int) (*ReviewVote, error)
	Update(vote *ReviewVote) error
	Delete(userID, reviewID int) error
	CountByReview(reviewID int) (helpful int, notHelpful int, err error)
}

// ReviewReportRepository 評論檢舉資料存取介面
type ReviewReportRepository interface {
	Create(report *ReviewReport) error
	GetByID(id int) (*ReviewReport, error)
	List(offset, limit int, status string) ([]*ReviewReport, error)
	UpdateStatus(id int, status string) error
	Count(status string) (int, error)
}

// GameRatingRepository 遊戲評分統計資料存取介面
type GameRatingRepository interface {
	GetByGameID(gameID int) (*GameRating, error)
	UpdateRating(gameID int) error
	GetTopRatedGames(limit int) ([]*GameRating, error)
}

// ReviewFilters 評論查詢過濾器
type ReviewFilters struct {
	GameID     int    `json:"game_id"`
	UserID     int    `json:"user_id"`
	MinRating  int    `json:"min_rating"`
	MaxRating  int    `json:"max_rating"`
	Status     string `json:"status"`
	Recommend  *bool  `json:"recommend"`
	IsVerified *bool  `json:"is_verified"`
	MinHours   int    `json:"min_hours"`
	SortBy     string `json:"sort_by"`    // created_at, rating, helpful_votes
	SortOrder  string `json:"sort_order"` // asc, desc
}

// TableName 返回評論表名
func (r *Review) TableName() string {
	return "reviews"
}

// TableName 返回評論投票表名
func (rv *ReviewVote) TableName() string {
	return "review_votes"
}

// TableName 返回評論檢舉表名
func (rr *ReviewReport) TableName() string {
	return "review_reports"
}

// TableName 返回遊戲評分表名
func (gr *GameRating) TableName() string {
	return "game_ratings"
}

// IsPublished 檢查評論是否已發布
func (r *Review) IsPublished() bool {
	return r.Status == "published"
}

// CanEdit 檢查評論是否可以編輯
func (r *Review) CanEdit() bool {
	return r.Status == "published" || r.Status == "draft"
}

// IsFlagged 檢查評論是否被檢舉
func (r *Review) IsFlagged() bool {
	return r.Status == "flagged"
}

// IsPositive 檢查評論是否為正面評價
func (r *Review) IsPositive() bool {
	return r.Rating >= 4 && r.Recommend
}

// GetRatingText 取得評分文字描述
func (r *Review) GetRatingText() string {
	switch r.Rating {
	case 5:
		return "excellent"
	case 4:
		return "good"
	case 3:
		return "average"
	case 2:
		return "poor"
	case 1:
		return "terrible"
	default:
		return "unknown"
	}
}

// IsHelpful 檢查是否為有用的投票
func (rv *ReviewVote) IsHelpful() bool {
	return rv.VoteType == "helpful"
}

// IsPending 檢查檢舉是否待處理
func (rr *ReviewReport) IsPending() bool {
	return rr.Status == "pending"
}

// GetRatingDistribution 取得評分分佈百分比
func (gr *GameRating) GetRatingDistribution() map[int]float64 {
	if gr.TotalReviews == 0 {
		return map[int]float64{
			5: 0, 4: 0, 3: 0, 2: 0, 1: 0,
		}
	}

	total := float64(gr.TotalReviews)
	return map[int]float64{
		5: float64(gr.FiveStars) / total * 100,
		4: float64(gr.FourStars) / total * 100,
		3: float64(gr.ThreeStars) / total * 100,
		2: float64(gr.TwoStars) / total * 100,
		1: float64(gr.OneStar) / total * 100,
	}
}

// GetQualityScore 取得品質分數（基於評分和推薦率）
func (gr *GameRating) GetQualityScore() float64 {
	if gr.TotalReviews == 0 {
		return 0
	}

	// 結合平均評分和推薦率的綜合分數
	ratingScore := gr.AverageRating / 5.0 * 100 // 0-100
	recommendScore := gr.RecommendedRate        // 0-100

	// 加權平均：評分權重 70%，推薦率權重 30%
	return (ratingScore * 0.7) + (recommendScore * 0.3)
}

// ReviewQueryBuilder 評論查詢建構器
type ReviewQueryBuilder struct {
	query string
	args  []interface{}
	joins []string
}

// NewReviewQueryBuilder 建立新的評論查詢建構器
func NewReviewQueryBuilder() *ReviewQueryBuilder {
	return &ReviewQueryBuilder{
		query: "SELECT r.* FROM reviews r",
		args:  make([]interface{}, 0),
		joins: make([]string, 0),
	}
}

// WithUser 加入使用者關聯查詢
func (qb *ReviewQueryBuilder) WithUser() *ReviewQueryBuilder {
	qb.joins = append(qb.joins, "LEFT JOIN users u ON r.user_id = u.id")
	return qb
}

// WithGame 加入遊戲關聯查詢
func (qb *ReviewQueryBuilder) WithGame() *ReviewQueryBuilder {
	qb.joins = append(qb.joins, "LEFT JOIN games g ON r.game_id = g.id")
	return qb
}

// WhereGameID 依遊戲 ID 過濾
func (qb *ReviewQueryBuilder) WhereGameID(gameID int) *ReviewQueryBuilder {
	if gameID > 0 {
		qb.addWhere("r.game_id = ?", gameID)
	}
	return qb
}

// WhereUserID 依使用者 ID 過濾
func (qb *ReviewQueryBuilder) WhereUserID(userID int) *ReviewQueryBuilder {
	if userID > 0 {
		qb.addWhere("r.user_id = ?", userID)
	}
	return qb
}

// WhereStatus 依狀態過濾
func (qb *ReviewQueryBuilder) WhereStatus(status string) *ReviewQueryBuilder {
	if status != "" {
		qb.addWhere("r.status = ?", status)
	}
	return qb
}

// WhereRating 依評分過濾
func (qb *ReviewQueryBuilder) WhereRating(minRating, maxRating int) *ReviewQueryBuilder {
	if minRating > 0 {
		qb.addWhere("r.rating >= ?", minRating)
	}
	if maxRating > 0 && maxRating >= minRating {
		qb.addWhere("r.rating <= ?", maxRating)
	}
	return qb
}

// WhereRecommend 依推薦狀態過濾
func (qb *ReviewQueryBuilder) WhereRecommend(recommend *bool) *ReviewQueryBuilder {
	if recommend != nil {
		qb.addWhere("r.recommend = ?", *recommend)
	}
	return qb
}

// WhereVerified 依驗證狀態過濾
func (qb *ReviewQueryBuilder) WhereVerified(verified *bool) *ReviewQueryBuilder {
	if verified != nil {
		qb.addWhere("r.is_verified = ?", *verified)
	}
	return qb
}

// OrderBy 排序
func (qb *ReviewQueryBuilder) OrderBy(column, direction string) *ReviewQueryBuilder {
	if direction != "ASC" && direction != "DESC" {
		direction = "DESC"
	}
	qb.query += " ORDER BY " + column + " " + direction
	return qb
}

// Limit 設定分頁
func (qb *ReviewQueryBuilder) Limit(offset, limit int) *ReviewQueryBuilder {
	qb.query += " LIMIT ? OFFSET ?"
	qb.args = append(qb.args, limit, offset)
	return qb
}

// Build 建構查詢
func (qb *ReviewQueryBuilder) Build() (string, []interface{}) {
	// 加入 JOIN 子句
	for _, join := range qb.joins {
		qb.query += " " + join
	}
	return qb.query, qb.args
}

// addWhere 新增 WHERE 條件
func (qb *ReviewQueryBuilder) addWhere(condition string, args ...interface{}) {
	if len(qb.args) == 0 {
		qb.query += " WHERE " + condition
	} else {
		qb.query += " AND " + condition
	}
	qb.args = append(qb.args, args...)
}
