package controllers

import (
	"net/http"
	"strings"

	"nexus-gaming-backend/models"
	"nexus-gaming-backend/services"

	"github.com/gin-gonic/gin"
)

// AuthController 身份驗證控制器
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController 建立新的身份驗證控制器
func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// LoginRequest 登入請求結構
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse 登入回應結構
type LoginResponse struct {
	Token     string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User      *models.User `json:"user"`
	ExpiresIn int64        `json:"expires_in" example:"86400"`
	TokenType string       `json:"token_type" example:"Bearer"`
}

// Login 使用者登入
// @Summary 使用者登入
// @Description 使用帳號密碼進行身份驗證，成功後返回 JWT Token
// @Tags 身份驗證
// @Accept json
// @Produce json
// @Param login body LoginRequest true "登入資訊"
// @Success 200 {object} APIResponse{data=LoginResponse} "登入成功"
// @Failure 400 {object} APIResponse "請求參數錯誤"
// @Failure 401 {object} APIResponse "帳號或密碼錯誤"
// @Failure 500 {object} APIResponse "伺服器內部錯誤"
// @Router /api/auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "請求參數錯誤",
			Data:    gin.H{"error": err.Error()},
		})
		return
	}

	// 驗證使用者身份
	user, err := ac.authService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "帳號或密碼錯誤",
			Data:    gin.H{"error": err.Error()},
		})
		return
	}

	// 生成 JWT Token
	token, err := ac.authService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Token 生成失敗",
			Data:    gin.H{"error": err.Error()},
		})
		return
	}

	// 清除敏感資訊
	user.Password = ""

	// 返回登入結果
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "登入成功",
		Data: LoginResponse{
			Token:     token,
			User:      user,
			ExpiresIn: 86400, // 24 小時
			TokenType: "Bearer",
		},
	})
}

// LogoutRequest 登出請求結構
type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

// Logout 使用者登出
// @Summary 使用者登出
// @Description 使用者登出，使 Token 失效（未來可實現 Token 黑名單）
// @Tags 身份驗證
// @Accept json
// @Produce json
// @Param logout body LogoutRequest true "登出資訊"
// @Security BearerAuth
// @Success 200 {object} APIResponse "登出成功"
// @Failure 400 {object} APIResponse "請求參數錯誤"
// @Failure 401 {object} APIResponse "Token 無效"
// @Failure 500 {object} APIResponse "伺服器內部錯誤"
// @Router /api/auth/logout [post]
func (ac *AuthController) Logout(c *gin.Context) {
	// 從 Header 或請求體中獲取 Token
	var token string

	// 優先從 Authorization Header 獲取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		// 從請求體獲取
		var req LogoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "請求參數錯誤",
				Data:    gin.H{"error": err.Error()},
			})
			return
		}
		token = req.Token
	}

	// 驗證 Token 有效性
	_, err := ac.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Token 無效",
			Data:    gin.H{"error": err.Error()},
		})
		return
	}

	// TODO: 實現 Token 黑名單機制
	// 目前暫時返回成功，未來可以將 Token 加入黑名單

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "登出成功",
		Data:    gin.H{"message": "已成功登出"},
	})
}

// RefreshRequest Token 刷新請求結構
type RefreshRequest struct {
	Token string `json:"token" binding:"required"`
}

// RefreshResponse Token 刷新回應結構
type RefreshResponse struct {
	Token     string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn int64  `json:"expires_in" example:"86400"`
	TokenType string `json:"token_type" example:"Bearer"`
}

// RefreshToken 刷新 Token
// @Summary 刷新 Token
// @Description 使用舊 Token 刷新獲取新的 Token
// @Tags 身份驗證
// @Accept json
// @Produce json
// @Param refresh body RefreshRequest true "刷新資訊"
// @Success 200 {object} APIResponse{data=RefreshResponse} "刷新成功"
// @Failure 400 {object} APIResponse "請求參數錯誤"
// @Failure 401 {object} APIResponse "Token 無效或已過期"
// @Failure 500 {object} APIResponse "伺服器內部錯誤"
// @Router /api/auth/refresh [post]
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "請求參數錯誤",
			Data:    gin.H{"error": err.Error()},
		})
		return
	}

	// 刷新 Token
	newToken, err := ac.authService.RefreshToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Token 刷新失敗",
			Data:    gin.H{"error": err.Error()},
		})
		return
	}

	// 返回新 Token
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Token 刷新成功",
		Data: RefreshResponse{
			Token:     newToken,
			ExpiresIn: 86400, // 24 小時
			TokenType: "Bearer",
		},
	})
}

// GetProfile 獲取使用者資訊
// @Summary 獲取當前使用者資訊
// @Description 根據 Token 獲取當前登入使用者的詳細資訊
// @Tags 身份驗證
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} APIResponse{data=models.User} "獲取成功"
// @Failure 401 {object} APIResponse "Token 無效"
// @Failure 500 {object} APIResponse "伺服器內部錯誤"
// @Router /api/auth/profile [get]
func (ac *AuthController) GetProfile(c *gin.Context) {
	// 從 Header 獲取 Token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "未提供有效的授權 Token",
			Data:    nil,
		})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// 從 Token 獲取使用者資訊
	user, err := ac.authService.GetUserFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Token 無效",
			Data:    gin.H{"error": err.Error()},
		})
		return
	}

	// 清除敏感資訊
	user.Password = ""

	// 返回使用者資訊
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "獲取使用者資訊成功",
		Data:    user,
	})
}

// AuthMiddleware 身份驗證中介軟體
func (ac *AuthController) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Authorization header 獲取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "缺少認證 Token",
				Data:    gin.H{"error": "MISSING_TOKEN"},
			})
			c.Abort()
			return
		}

		// 檢查 Bearer 格式
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "Token 格式錯誤，應為 'Bearer <token>'",
				Data:    gin.H{"error": "INVALID_TOKEN_FORMAT"},
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// 驗證 token
		claims, err := ac.authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "Token 驗證失敗",
				Data:    gin.H{"error": err.Error()},
			})
			c.Abort()
			return
		}

		// 設置使用者資訊到 context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminPermissionMiddleware 管理員權限檢查中介軟體
func (ac *AuthController) AdminPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 檢查使用者角色
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, APIResponse{
				Success: false,
				Message: "未找到使用者角色資訊",
				Data:    gin.H{"error": "MISSING_USER_ROLE"},
			})
			c.Abort()
			return
		}

		// 檢查是否為管理員
		if userRole != "admin" && userRole != "super_admin" {
			c.JSON(http.StatusForbidden, APIResponse{
				Success: false,
				Message: "需要管理員權限",
				Data:    gin.H{"error": "INSUFFICIENT_PERMISSION"},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
