package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 使用者結構體
type User struct {
	ID          int        `json:"id" db:"id"`
	Username    string     `json:"username" db:"username"`
	Email       string     `json:"email" db:"email"`
	Password    string     `json:"-" db:"password"` // 不在 JSON 中顯示密碼
	RoleID      int        `json:"role_id" db:"role_id"`
	Role        *Role      `json:"role,omitempty"` // 關聯的角色
	Status      string     `json:"status" db:"status"`
	LastLoginAt *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Role 角色結構體
type Role struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Permissions []string  `json:"permissions" db:"permissions"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UserRepository 使用者資料存取介面
type UserRepository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id int) error
	List(offset, limit int, status string) ([]*User, error)
	Count(status string) (int, error)
	UpdateLastLogin(id int) error
}

// RoleRepository 角色資料存取介面
type RoleRepository interface {
	GetByID(id int) (*Role, error)
	GetByName(name string) (*Role, error)
	List() ([]*Role, error)
	HasPermission(roleID int, permission string) (bool, error)
}

// TableName 返回使用者表名
func (u *User) TableName() string {
	return "users"
}

// TableName 返回角色表名
func (r *Role) TableName() string {
	return "roles"
}

// IsActive 檢查使用者是否啟用
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// HasRole 檢查使用者是否有指定角色
func (u *User) HasRole(roleName string) bool {
	if u.Role == nil {
		return false
	}
	return u.Role.Name == roleName
}

// HasPermission 檢查角色是否有指定權限
func (r *Role) HasPermission(permission string) bool {
	for _, perm := range r.Permissions {
		if perm == permission || perm == "*" {
			return true
		}
	}
	return false
}

// UserQueryBuilder 使用者查詢建構器
type UserQueryBuilder struct {
	query  string
	args   []interface{}
	offset int
	limit  int
}

// NewUserQueryBuilder 建立新的查詢建構器
func NewUserQueryBuilder() *UserQueryBuilder {
	return &UserQueryBuilder{
		query: "SELECT u.id, u.username, u.email, u.role_id, u.status, u.last_login_at, u.created_at, u.updated_at FROM users u",
		args:  make([]interface{}, 0),
	}
}

// WithRole 加入角色關聯查詢
func (qb *UserQueryBuilder) WithRole() *UserQueryBuilder {
	qb.query = "SELECT u.id, u.username, u.email, u.role_id, u.status, u.last_login_at, u.created_at, u.updated_at, " +
		"r.id as role_id, r.name as role_name, r.description as role_description, r.permissions as role_permissions " +
		"FROM users u LEFT JOIN roles r ON u.role_id = r.id"
	return qb
}

// WhereStatus 依狀態過濾
func (qb *UserQueryBuilder) WhereStatus(status string) *UserQueryBuilder {
	if status != "" {
		if len(qb.args) == 0 {
			qb.query += " WHERE u.status = ?"
		} else {
			qb.query += " AND u.status = ?"
		}
		qb.args = append(qb.args, status)
	}
	return qb
}

// WhereRole 依角色過濾
func (qb *UserQueryBuilder) WhereRole(roleID int) *UserQueryBuilder {
	if roleID > 0 {
		if len(qb.args) == 0 {
			qb.query += " WHERE u.role_id = ?"
		} else {
			qb.query += " AND u.role_id = ?"
		}
		qb.args = append(qb.args, roleID)
	}
	return qb
}

// OrderBy 排序
func (qb *UserQueryBuilder) OrderBy(column string, direction string) *UserQueryBuilder {
	if direction != "ASC" && direction != "DESC" {
		direction = "ASC"
	}
	qb.query += " ORDER BY " + column + " " + direction
	return qb
}

// Limit 設定分頁
func (qb *UserQueryBuilder) Limit(offset, limit int) *UserQueryBuilder {
	qb.offset = offset
	qb.limit = limit
	qb.query += " LIMIT ? OFFSET ?"
	qb.args = append(qb.args, limit, offset)
	return qb
}

// Build 建構查詢
func (qb *UserQueryBuilder) Build() (string, []interface{}) {
	return qb.query, qb.args
}

// CheckPassword 檢查密碼是否正確
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// HashPassword 加密密碼
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
