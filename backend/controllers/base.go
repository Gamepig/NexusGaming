package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse 統一的 API 回應格式
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// SuccessResponse 成功回應
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 錯誤回應
func ErrorResponse(c *gin.Context, statusCode int, message string, code string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Code:    code,
	})
}

// HealthCheck 健康檢查
func HealthCheck(c *gin.Context) {
	SuccessResponse(c, gin.H{
		"status":    "ok",
		"timestamp": "2024-12-19T12:00:00Z",
		"service":   "nexus-gaming-backend",
	}, "Service is healthy")
}

// 以下是佔位符函數，將在後續階段實現具體邏輯

// 身份驗證相關（已移至 auth.go）
// func Login(c *gin.Context) - 已在 auth.go 實現
// func Logout(c *gin.Context) - 已在 auth.go 實現
// func RefreshToken(c *gin.Context) - 已在 auth.go 實現

// 使用者管理相關
func GetUsers(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetUsers endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func GetUser(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetUser endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func CreateUser(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "CreateUser endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func UpdateUser(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "UpdateUser endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func DeleteUser(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "DeleteUser endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func ChangePassword(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "ChangePassword endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// 玩家管理相關
func GetPlayers(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetPlayers endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func GetPlayer(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetPlayer endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func CreatePlayer(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "CreatePlayer endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func UpdatePlayer(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "UpdatePlayer endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func DeletePlayer(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "DeletePlayer endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func UpdatePlayerStatus(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "UpdatePlayerStatus endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func GetPlayerBalance(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetPlayerBalance endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func DepositPlayerBalance(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "DepositPlayerBalance endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func WithdrawPlayerBalance(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "WithdrawPlayerBalance endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func GetPlayerTransactions(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetPlayerTransactions endpoint not implemented yet", "NOT_IMPLEMENTED")
}
