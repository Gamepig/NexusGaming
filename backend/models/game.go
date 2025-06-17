package models

import (
	"time"
)

// Game 遊戲結構體
type Game struct {
	ID           int        `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Description  string     `json:"description" db:"description"`
	GenreID      int        `json:"genre_id" db:"genre_id"`
	Genre        *Genre     `json:"genre,omitempty"`
	DeveloperID  int        `json:"developer_id" db:"developer_id"`
	Developer    *Developer `json:"developer,omitempty"`
	Price        float64    `json:"price" db:"price"`
	ReleaseDate  time.Time  `json:"release_date" db:"release_date"`
	Status       string     `json:"status" db:"status"` // released, coming_soon, early_access
	Rating       float64    `json:"rating" db:"rating"`
	RatingCount  int        `json:"rating_count" db:"rating_count"`
	ImageURL     string     `json:"image_url" db:"image_url"`
	TrailerURL   string     `json:"trailer_url" db:"trailer_url"`
	SystemReqs   *SystemReq `json:"system_requirements,omitempty"`
	Tags         []Tag      `json:"tags,omitempty"`
	Screenshots  []string   `json:"screenshots,omitempty"`
	Features     []string   `json:"features,omitempty"`
	Languages    []string   `json:"languages,omitempty"`
	MetadataJSON string     `json:"-" db:"metadata_json"` // 儲存額外的 JSON 元資料
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// Genre 遊戲類型結構體
type Genre struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Developer 開發商結構體
type Developer struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Website     string    `json:"website" db:"website"`
	Country     string    `json:"country" db:"country"`
	LogoURL     string    `json:"logo_url" db:"logo_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SystemReq 系統需求結構體
type SystemReq struct {
	ID              int    `json:"id" db:"id"`
	GameID          int    `json:"game_id" db:"game_id"`
	MinOS           string `json:"min_os" db:"min_os"`
	MinProcessor    string `json:"min_processor" db:"min_processor"`
	MinMemory       string `json:"min_memory" db:"min_memory"`
	MinGraphics     string `json:"min_graphics" db:"min_graphics"`
	MinStorage      string `json:"min_storage" db:"min_storage"`
	RecommendedOS   string `json:"recommended_os" db:"recommended_os"`
	RecommendedProc string `json:"recommended_processor" db:"recommended_processor"`
	RecommendedMem  string `json:"recommended_memory" db:"recommended_memory"`
	RecommendedGfx  string `json:"recommended_graphics" db:"recommended_graphics"`
	RecommendedSto  string `json:"recommended_storage" db:"recommended_storage"`
}

// Tag 標籤結構體
type Tag struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Color       string    `json:"color" db:"color"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// GameRepository 遊戲資料存取介面
type GameRepository interface {
	Create(game *Game) error
	GetByID(id int) (*Game, error)
	Update(game *Game) error
	Delete(id int) error
	List(offset, limit int, filters GameFilters) ([]*Game, error)
	Count(filters GameFilters) (int, error)
	GetByGenre(genreID int, offset, limit int) ([]*Game, error)
	GetByDeveloper(developerID int, offset, limit int) ([]*Game, error)
	Search(query string, offset, limit int) ([]*Game, error)
	GetFeatured(limit int) ([]*Game, error)
	GetNewReleases(limit int) ([]*Game, error)
	GetTopRated(limit int) ([]*Game, error)
	UpdateRating(gameID int, rating float64, count int) error
}

// GenreRepository 遊戲類型資料存取介面
type GenreRepository interface {
	GetByID(id int) (*Genre, error)
	GetByName(name string) (*Genre, error)
	List() ([]*Genre, error)
	Create(genre *Genre) error
	Update(genre *Genre) error
	Delete(id int) error
}

// DeveloperRepository 開發商資料存取介面
type DeveloperRepository interface {
	GetByID(id int) (*Developer, error)
	GetByName(name string) (*Developer, error)
	List(offset, limit int) ([]*Developer, error)
	Create(developer *Developer) error
	Update(developer *Developer) error
	Delete(id int) error
}

// TagRepository 標籤資料存取介面
type TagRepository interface {
	GetByID(id int) (*Tag, error)
	GetByName(name string) (*Tag, error)
	List() ([]*Tag, error)
	GetByGame(gameID int) ([]*Tag, error)
	Create(tag *Tag) error
	Update(tag *Tag) error
	Delete(id int) error
	AddToGame(gameID, tagID int) error
	RemoveFromGame(gameID, tagID int) error
}

// GameFilters 遊戲查詢過濾器
type GameFilters struct {
	GenreID     int      `json:"genre_id"`
	DeveloperID int      `json:"developer_id"`
	MinPrice    *float64 `json:"min_price"`
	MaxPrice    *float64 `json:"max_price"`
	Status      string   `json:"status"`
	MinRating   *float64 `json:"min_rating"`
	Tags        []int    `json:"tags"`
	SortBy      string   `json:"sort_by"`    // name, price, rating, release_date
	SortOrder   string   `json:"sort_order"` // asc, desc
}

// TableName 返回遊戲表名
func (g *Game) TableName() string {
	return "games"
}

// TableName 返回類型表名
func (g *Genre) TableName() string {
	return "genres"
}

// TableName 返回開發商表名
func (d *Developer) TableName() string {
	return "developers"
}

// TableName 返回標籤表名
func (t *Tag) TableName() string {
	return "tags"
}

// IsReleased 檢查遊戲是否已發布
func (g *Game) IsReleased() bool {
	return g.Status == "released" && g.ReleaseDate.Before(time.Now())
}

// IsFree 檢查遊戲是否免費
func (g *Game) IsFree() bool {
	return g.Price == 0
}

// GetDiscountedPrice 取得折扣價格（如果有的話）
func (g *Game) GetDiscountedPrice(discountPercent float64) float64 {
	if discountPercent <= 0 || discountPercent >= 100 {
		return g.Price
	}
	return g.Price * (1 - discountPercent/100)
}

// GameQueryBuilder 遊戲查詢建構器
type GameQueryBuilder struct {
	query string
	args  []interface{}
	joins []string
}

// NewGameQueryBuilder 建立新的遊戲查詢建構器
func NewGameQueryBuilder() *GameQueryBuilder {
	return &GameQueryBuilder{
		query: "SELECT g.* FROM games g",
		args:  make([]interface{}, 0),
		joins: make([]string, 0),
	}
}

// WithGenre 加入類型關聯查詢
func (qb *GameQueryBuilder) WithGenre() *GameQueryBuilder {
	qb.joins = append(qb.joins, "LEFT JOIN genres gen ON g.genre_id = gen.id")
	return qb
}

// WithDeveloper 加入開發商關聯查詢
func (qb *GameQueryBuilder) WithDeveloper() *GameQueryBuilder {
	qb.joins = append(qb.joins, "LEFT JOIN developers dev ON g.developer_id = dev.id")
	return qb
}

// WhereGenre 依類型過濾
func (qb *GameQueryBuilder) WhereGenre(genreID int) *GameQueryBuilder {
	if genreID > 0 {
		qb.addWhere("g.genre_id = ?", genreID)
	}
	return qb
}

// WhereStatus 依狀態過濾
func (qb *GameQueryBuilder) WhereStatus(status string) *GameQueryBuilder {
	if status != "" {
		qb.addWhere("g.status = ?", status)
	}
	return qb
}

// WherePriceRange 依價格範圍過濾
func (qb *GameQueryBuilder) WherePriceRange(minPrice, maxPrice *float64) *GameQueryBuilder {
	if minPrice != nil {
		qb.addWhere("g.price >= ?", *minPrice)
	}
	if maxPrice != nil {
		qb.addWhere("g.price <= ?", *maxPrice)
	}
	return qb
}

// WhereRating 依評分過濾
func (qb *GameQueryBuilder) WhereRating(minRating *float64) *GameQueryBuilder {
	if minRating != nil {
		qb.addWhere("g.rating >= ?", *minRating)
	}
	return qb
}

// Search 搜尋關鍵字
func (qb *GameQueryBuilder) Search(keyword string) *GameQueryBuilder {
	if keyword != "" {
		qb.addWhere("(g.name LIKE ? OR g.description LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	return qb
}

// OrderBy 排序
func (qb *GameQueryBuilder) OrderBy(column, direction string) *GameQueryBuilder {
	if direction != "ASC" && direction != "DESC" {
		direction = "ASC"
	}
	qb.query += " ORDER BY " + column + " " + direction
	return qb
}

// Limit 設定分頁
func (qb *GameQueryBuilder) Limit(offset, limit int) *GameQueryBuilder {
	qb.query += " LIMIT ? OFFSET ?"
	qb.args = append(qb.args, limit, offset)
	return qb
}

// Build 建構查詢
func (qb *GameQueryBuilder) Build() (string, []interface{}) {
	// 加入 JOIN 子句
	for _, join := range qb.joins {
		qb.query += " " + join
	}
	return qb.query, qb.args
}

// addWhere 新增 WHERE 條件
func (qb *GameQueryBuilder) addWhere(condition string, args ...interface{}) {
	if len(qb.args) == 0 {
		qb.query += " WHERE " + condition
	} else {
		qb.query += " AND " + condition
	}
	qb.args = append(qb.args, args...)
}
