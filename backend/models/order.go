package models

import (
	"time"
)

// Order 訂單結構體
type Order struct {
	ID            int         `json:"id" db:"id"`
	UserID        int         `json:"user_id" db:"user_id"`
	User          *User       `json:"user,omitempty"`
	OrderNumber   string      `json:"order_number" db:"order_number"`
	Status        string      `json:"status" db:"status"` // pending, paid, cancelled, refunded
	TotalAmount   float64     `json:"total_amount" db:"total_amount"`
	PaymentMethod string      `json:"payment_method" db:"payment_method"`
	PaymentStatus string      `json:"payment_status" db:"payment_status"` // pending, completed, failed
	PaymentID     string      `json:"payment_id" db:"payment_id"`
	BillingInfo   BillingInfo `json:"billing_info"`
	Items         []OrderItem `json:"items,omitempty"`
	Notes         string      `json:"notes" db:"notes"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	PaidAt        *time.Time  `json:"paid_at" db:"paid_at"`
}

// OrderItem 訂單項目結構體
type OrderItem struct {
	ID       int     `json:"id" db:"id"`
	OrderID  int     `json:"order_id" db:"order_id"`
	GameID   int     `json:"game_id" db:"game_id"`
	Game     *Game   `json:"game,omitempty"`
	Price    float64 `json:"price" db:"price"`
	Discount float64 `json:"discount" db:"discount"`
	Quantity int     `json:"quantity" db:"quantity"`
}

// BillingInfo 帳單資訊結構體
type BillingInfo struct {
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Address   string `json:"address" db:"address"`
	City      string `json:"city" db:"city"`
	State     string `json:"state" db:"state"`
	ZipCode   string `json:"zip_code" db:"zip_code"`
	Country   string `json:"country" db:"country"`
}

// Cart 購物車結構體
type Cart struct {
	ID        int        `json:"id" db:"id"`
	UserID    int        `json:"user_id" db:"user_id"`
	Items     []CartItem `json:"items,omitempty"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// CartItem 購物車項目結構體
type CartItem struct {
	ID       int       `json:"id" db:"id"`
	CartID   int       `json:"cart_id" db:"cart_id"`
	GameID   int       `json:"game_id" db:"game_id"`
	Game     *Game     `json:"game,omitempty"`
	Quantity int       `json:"quantity" db:"quantity"`
	AddedAt  time.Time `json:"added_at" db:"added_at"`
}

// Wishlist 願望清單結構體
type Wishlist struct {
	ID        int            `json:"id" db:"id"`
	UserID    int            `json:"user_id" db:"user_id"`
	Items     []WishlistItem `json:"items,omitempty"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// WishlistItem 願望清單項目結構體
type WishlistItem struct {
	ID         int       `json:"id" db:"id"`
	WishlistID int       `json:"wishlist_id" db:"wishlist_id"`
	GameID     int       `json:"game_id" db:"game_id"`
	Game       *Game     `json:"game,omitempty"`
	AddedAt    time.Time `json:"added_at" db:"added_at"`
}

// OrderRepository 訂單資料存取介面
type OrderRepository interface {
	Create(order *Order) error
	GetByID(id int) (*Order, error)
	GetByOrderNumber(orderNumber string) (*Order, error)
	GetByUserID(userID int, offset, limit int) ([]*Order, error)
	Update(order *Order) error
	UpdateStatus(id int, status string) error
	UpdatePaymentStatus(id int, status, paymentID string) error
	List(offset, limit int, filters OrderFilters) ([]*Order, error)
	Count(filters OrderFilters) (int, error)
	GetRevenue(startDate, endDate time.Time) (float64, error)
	GetOrderStats() (*OrderStats, error)
}

// CartRepository 購物車資料存取介面
type CartRepository interface {
	GetByUserID(userID int) (*Cart, error)
	AddItem(userID, gameID, quantity int) error
	UpdateItem(userID, gameID, quantity int) error
	RemoveItem(userID, gameID int) error
	Clear(userID int) error
	GetItemCount(userID int) (int, error)
	GetTotal(userID int) (float64, error)
}

// WishlistRepository 願望清單資料存取介面
type WishlistRepository interface {
	GetByUserID(userID int) (*Wishlist, error)
	AddItem(userID, gameID int) error
	RemoveItem(userID, gameID int) error
	HasItem(userID, gameID int) (bool, error)
	GetItemCount(userID int) (int, error)
}

// OrderFilters 訂單查詢過濾器
type OrderFilters struct {
	UserID        int        `json:"user_id"`
	Status        string     `json:"status"`
	PaymentStatus string     `json:"payment_status"`
	PaymentMethod string     `json:"payment_method"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	MinAmount     *float64   `json:"min_amount"`
	MaxAmount     *float64   `json:"max_amount"`
	SortBy        string     `json:"sort_by"`    // created_at, total_amount
	SortOrder     string     `json:"sort_order"` // asc, desc
}

// OrderStats 訂單統計結構體
type OrderStats struct {
	TotalOrders       int     `json:"total_orders"`
	TotalRevenue      float64 `json:"total_revenue"`
	PendingOrders     int     `json:"pending_orders"`
	CompletedOrders   int     `json:"completed_orders"`
	CancelledOrders   int     `json:"cancelled_orders"`
	AverageOrderValue float64 `json:"average_order_value"`
}

// TableName 返回訂單表名
func (o *Order) TableName() string {
	return "orders"
}

// TableName 返回訂單項目表名
func (oi *OrderItem) TableName() string {
	return "order_items"
}

// TableName 返回購物車表名
func (c *Cart) TableName() string {
	return "carts"
}

// TableName 返回購物車項目表名
func (ci *CartItem) TableName() string {
	return "cart_items"
}

// TableName 返回願望清單表名
func (w *Wishlist) TableName() string {
	return "wishlists"
}

// TableName 返回願望清單項目表名
func (wi *WishlistItem) TableName() string {
	return "wishlist_items"
}

// IsPending 檢查訂單是否待處理
func (o *Order) IsPending() bool {
	return o.Status == "pending"
}

// IsPaid 檢查訂單是否已付款
func (o *Order) IsPaid() bool {
	return o.Status == "paid" && o.PaymentStatus == "completed"
}

// CanCancel 檢查訂單是否可以取消
func (o *Order) CanCancel() bool {
	return o.Status == "pending" || (o.Status == "paid" && o.PaymentStatus == "pending")
}

// CanRefund 檢查訂單是否可以退款
func (o *Order) CanRefund() bool {
	return o.Status == "paid" && o.PaymentStatus == "completed"
}

// GetSubtotal 取得訂單項目小計
func (oi *OrderItem) GetSubtotal() float64 {
	subtotal := oi.Price * float64(oi.Quantity)
	if oi.Discount > 0 {
		subtotal -= oi.Discount
	}
	return subtotal
}

// GetTotalItems 取得購物車總項目數
func (c *Cart) GetTotalItems() int {
	total := 0
	for _, item := range c.Items {
		total += item.Quantity
	}
	return total
}

// GetTotal 取得購物車總金額
func (c *Cart) GetTotal() float64 {
	total := 0.0
	for _, item := range c.Items {
		if item.Game != nil {
			total += item.Game.Price * float64(item.Quantity)
		}
	}
	return total
}

// HasGame 檢查購物車是否包含指定遊戲
func (c *Cart) HasGame(gameID int) bool {
	for _, item := range c.Items {
		if item.GameID == gameID {
			return true
		}
	}
	return false
}

// HasGame 檢查願望清單是否包含指定遊戲
func (w *Wishlist) HasGame(gameID int) bool {
	for _, item := range w.Items {
		if item.GameID == gameID {
			return true
		}
	}
	return false
}

// OrderQueryBuilder 訂單查詢建構器
type OrderQueryBuilder struct {
	query string
	args  []interface{}
	joins []string
}

// NewOrderQueryBuilder 建立新的訂單查詢建構器
func NewOrderQueryBuilder() *OrderQueryBuilder {
	return &OrderQueryBuilder{
		query: "SELECT o.* FROM orders o",
		args:  make([]interface{}, 0),
		joins: make([]string, 0),
	}
}

// WithUser 加入使用者關聯查詢
func (qb *OrderQueryBuilder) WithUser() *OrderQueryBuilder {
	qb.joins = append(qb.joins, "LEFT JOIN users u ON o.user_id = u.id")
	return qb
}

// WhereUserID 依使用者 ID 過濾
func (qb *OrderQueryBuilder) WhereUserID(userID int) *OrderQueryBuilder {
	if userID > 0 {
		qb.addWhere("o.user_id = ?", userID)
	}
	return qb
}

// WhereStatus 依狀態過濾
func (qb *OrderQueryBuilder) WhereStatus(status string) *OrderQueryBuilder {
	if status != "" {
		qb.addWhere("o.status = ?", status)
	}
	return qb
}

// WherePaymentStatus 依付款狀態過濾
func (qb *OrderQueryBuilder) WherePaymentStatus(status string) *OrderQueryBuilder {
	if status != "" {
		qb.addWhere("o.payment_status = ?", status)
	}
	return qb
}

// WhereDateRange 依日期範圍過濾
func (qb *OrderQueryBuilder) WhereDateRange(startDate, endDate *time.Time) *OrderQueryBuilder {
	if startDate != nil {
		qb.addWhere("o.created_at >= ?", *startDate)
	}
	if endDate != nil {
		qb.addWhere("o.created_at <= ?", *endDate)
	}
	return qb
}

// WhereAmountRange 依金額範圍過濾
func (qb *OrderQueryBuilder) WhereAmountRange(minAmount, maxAmount *float64) *OrderQueryBuilder {
	if minAmount != nil {
		qb.addWhere("o.total_amount >= ?", *minAmount)
	}
	if maxAmount != nil {
		qb.addWhere("o.total_amount <= ?", *maxAmount)
	}
	return qb
}

// OrderBy 排序
func (qb *OrderQueryBuilder) OrderBy(column, direction string) *OrderQueryBuilder {
	if direction != "ASC" && direction != "DESC" {
		direction = "DESC"
	}
	qb.query += " ORDER BY " + column + " " + direction
	return qb
}

// Limit 設定分頁
func (qb *OrderQueryBuilder) Limit(offset, limit int) *OrderQueryBuilder {
	qb.query += " LIMIT ? OFFSET ?"
	qb.args = append(qb.args, limit, offset)
	return qb
}

// Build 建構查詢
func (qb *OrderQueryBuilder) Build() (string, []interface{}) {
	// 加入 JOIN 子句
	for _, join := range qb.joins {
		qb.query += " " + join
	}
	return qb.query, qb.args
}

// addWhere 新增 WHERE 條件
func (qb *OrderQueryBuilder) addWhere(condition string, args ...interface{}) {
	if len(qb.args) == 0 {
		qb.query += " WHERE " + condition
	} else {
		qb.query += " AND " + condition
	}
	qb.args = append(qb.args, args...)
}
