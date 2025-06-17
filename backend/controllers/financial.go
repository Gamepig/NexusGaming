package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 財務管理相關
func GetTransactions(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetTransactions endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func GetTransaction(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetTransaction endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func GetDeposits(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetDeposits endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func CreateDeposit(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "CreateDeposit endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func ConfirmDeposit(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "ConfirmDeposit endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func GetWithdrawals(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetWithdrawals endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func CreateWithdrawal(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "CreateWithdrawal endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func ApproveWithdrawal(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "ApproveWithdrawal endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func GetDailyReconciliation(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetDailyReconciliation endpoint not implemented yet", "NOT_IMPLEMENTED")
}

func GetMonthlyReconciliation(c *gin.Context) {
	ErrorResponse(c, http.StatusNotImplemented, "GetMonthlyReconciliation endpoint not implemented yet", "NOT_IMPLEMENTED")
}
