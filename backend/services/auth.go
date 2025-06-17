package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"nexus-gaming-backend/config"
	"nexus-gaming-backend/models"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService 身份驗證服務
type AuthService struct {
	JWTSecret string
	DB        *sql.DB
}

// NewAuthService 建立新的身份驗證服務
func NewAuthService() *AuthService {
	return &AuthService{
		JWTSecret: config.GetJWTSecret(),
		DB:        config.GetDB(),
	}
}

// JWTClaims JWT 聲明結構
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	if user == nil {
		return "", errors.New("使用者資料不能為空")
	}

	// 獲取角色名稱
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}

	// 設定 Token 有效期限為 24 小時
	expirationTime := time.Now().Add(24 * time.Hour)

	// 建立 JWT 聲明
	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "nexus-gaming",
			Subject:   fmt.Sprintf("user-%d", user.ID),
		},
	}

	// 建立 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密鑰簽名
	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("無法生成 Token: %v", err)
	}

	return tokenString, nil
}

// ValidateToken 驗證 JWT Token
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, errors.New("Token 不能為空")
	}

	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 驗證簽名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("無效的簽名方法: %v", token.Header["alg"])
		}
		return []byte(s.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Token 解析失敗: %v", err)
	}

	// 檢查 Token 是否有效
	if !token.Valid {
		return nil, errors.New("Token 無效")
	}

	// 提取聲明
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("無法提取 Token 聲明")
	}

	// 檢查 Token 是否過期
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("Token 已過期")
	}

	return claims, nil
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("原 Token 驗證失敗: %v", err)
	}

	// 檢查 Token 是否還有足夠時間刷新（至少還有 1 小時有效期）
	if time.Until(claims.ExpiresAt.Time) < time.Hour {
		return "", errors.New("Token 剩餘時間不足，無法刷新")
	}

	// 從資料庫重新獲取使用者資訊
	user, err := s.getUserByID(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("無法獲取使用者資訊: %v", err)
	}

	return s.GenerateToken(user)
}

// AuthenticateUser 使用者身份驗證
func (s *AuthService) AuthenticateUser(username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("使用者名稱和密碼不能為空")
	}

	// 從資料庫查詢使用者（包含角色資訊）
	user, err := s.getUserByUsername(username)
	if err != nil {
		return nil, errors.New("使用者名稱或密碼錯誤")
	}

	if user.Status != "active" {
		return nil, errors.New("帳號已被停用")
	}

	// 驗證密碼
	if !user.CheckPassword(password) {
		return nil, errors.New("使用者名稱或密碼錯誤")
	}

	// 更新最後登入時間
	if err := s.updateLastLogin(user.ID); err != nil {
		// 記錄錯誤但不影響登入流程
		fmt.Printf("Warning: failed to update last login time: %v\n", err)
	}

	return user, nil
}

// GetUserFromToken 從 Token 獲取使用者資訊
func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 從資料庫獲取完整的使用者資訊
	user, err := s.getUserByID(claims.UserID)
	if err != nil {
		return nil, errors.New("使用者不存在或已被停用")
	}

	if user.Status != "active" {
		return nil, errors.New("使用者已被停用")
	}

	return user, nil
}

// CheckPermission 檢查使用者權限
func (s *AuthService) CheckPermission(user *models.User, requiredRole string) bool {
	if user == nil || user.Role == nil {
		return false
	}

	// 定義角色層級
	roleHierarchy := map[string]int{
		"super_admin": 4,
		"admin":       3,
		"agent":       2,
		"dealer":      1,
	}

	userLevel, userExists := roleHierarchy[user.Role.Name]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	if !userExists || !requiredExists {
		return false
	}

	// 檢查使用者角色是否有足夠權限
	return userLevel >= requiredLevel
}

// 私有方法：根據 ID 獲取使用者
func (s *AuthService) getUserByID(id int) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password, u.role_id, u.status, u.last_login_at, u.created_at, u.updated_at,
		       r.id, r.name, r.description, r.permissions, r.created_at, r.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = ?
	`

	row := s.DB.QueryRow(query, id)

	user := &models.User{Role: &models.Role{}}
	var _, roleCreatedAt, roleUpdatedAt sql.NullTime
	var roleIDNull, roleNameNull, roleDescNull, rolePermNull sql.NullString

	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.RoleID, &user.Status, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		&roleIDNull, &roleNameNull, &roleDescNull, &rolePermNull, &roleCreatedAt, &roleUpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("使用者不存在")
		}
		return nil, err
	}

	// 處理角色資訊
	if roleNameNull.Valid {
		user.Role.Name = roleNameNull.String
		user.Role.Description = roleDescNull.String
		// 處理權限（假設以逗號分隔）
		if rolePermNull.Valid && rolePermNull.String != "" {
			// 這裡需要根據實際的權限存儲格式來解析
			user.Role.Permissions = []string{rolePermNull.String}
		}
	} else {
		user.Role = nil
	}

	return user, nil
}

// 私有方法：根據使用者名稱獲取使用者
func (s *AuthService) getUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password, u.role_id, u.status, u.last_login_at, u.created_at, u.updated_at,
		       r.id, r.name, r.description, r.permissions, r.created_at, r.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.username = ?
	`

	row := s.DB.QueryRow(query, username)

	user := &models.User{Role: &models.Role{}}
	var _, roleCreatedAt, roleUpdatedAt sql.NullTime
	var roleIDNull, roleNameNull, roleDescNull, rolePermNull sql.NullString

	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.RoleID, &user.Status, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		&roleIDNull, &roleNameNull, &roleDescNull, &rolePermNull, &roleCreatedAt, &roleUpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("使用者不存在")
		}
		return nil, err
	}

	// 處理角色資訊
	if roleNameNull.Valid {
		user.Role.Name = roleNameNull.String
		user.Role.Description = roleDescNull.String
		// 處理權限（假設以逗號分隔）
		if rolePermNull.Valid && rolePermNull.String != "" {
			// 這裡需要根據實際的權限存儲格式來解析
			user.Role.Permissions = []string{rolePermNull.String}
		}
	} else {
		user.Role = nil
	}

	return user, nil
}

// 私有方法：更新最後登入時間
func (s *AuthService) updateLastLogin(userID int) error {
	query := "UPDATE users SET last_login_at = ? WHERE id = ?"
	_, err := s.DB.Exec(query, time.Now(), userID)
	return err
}
