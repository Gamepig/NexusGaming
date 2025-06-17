package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 角色權限管理相關
func GetRoles(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetRoles endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func CreateRole(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "CreateRole endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func UpdateRole(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "UpdateRole endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func DeleteRole(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "DeleteRole endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func GetPermissions(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetPermissions endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func CreatePermission(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "CreatePermission endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// 操作日誌相關
func GetOperationLogs(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetOperationLogs endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func GetOperationLog(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetOperationLog endpoint not implemented yet", "NOT_IMPLEMENTED")
}

// 系統設置相關
func GetSystemSettings(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetSystemSettings endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func UpdateSystemSettings(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "UpdateSystemSettings endpoint not implemented yet", "NOT_IMPLEMENTED")
}
