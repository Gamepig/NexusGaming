package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware 處理跨域請求
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// LoggerMiddleware 記錄請求日誌
func LoggerMiddleware() gin.HandlerFunc {
	return gin.Logger()
}

// ErrorHandlerMiddleware 統一錯誤處理
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 處理錯誤
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
				"code":    "INTERNAL_ERROR",
			})
		}
	}
}

// 身份驗證和權限相關中介軟體已移至 auth.go 檔案
