package routes

import (
	"nexus-gaming-backend/controllers"
	"nexus-gaming-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 設置所有路由
func SetupRoutes(r *gin.Engine) {
	// 設置全域中介軟體
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())

	// 健康檢查路由
	r.GET("/health", controllers.HealthCheck)

	// API v1 路由群組
	v1 := r.Group("/api/v1")
	{
		// 身份驗證路由（不需要驗證）
		auth := v1.Group("/auth")
		authController := controllers.NewAuthController()
		{
			auth.POST("/login", authController.Login)
			auth.POST("/logout", authController.Logout)
			auth.POST("/refresh", authController.RefreshToken)
			auth.GET("/profile", authController.GetProfile)
		}

		// 需要身份驗證的路由
		authenticated := v1.Group("/")
		authMiddleware := controllers.NewAuthController().AuthMiddleware()
		authenticated.Use(authMiddleware)
		{
			// 使用者管理路由
			users := authenticated.Group("/users")
			{
				users.GET("/", controllers.GetUsers)
				users.GET("/:id", controllers.GetUser)
				users.POST("/", controllers.CreateUser)
				users.PUT("/:id", controllers.UpdateUser)
				users.DELETE("/:id", controllers.DeleteUser)
				users.PUT("/:id/password", controllers.ChangePassword)
			}

			// 玩家管理路由
			players := authenticated.Group("/players")
			playerController := controllers.NewPlayerController()
			{
				players.GET("/", playerController.GetPlayers)
				players.GET("/:id", playerController.GetPlayer)
				players.GET("/:id/games", playerController.GetPlayerGameHistory) // 新增：玩家遊戲歷史
				players.GET("/search", playerController.SearchPlayers)
				players.GET("/filter", playerController.FilterPlayers)
				players.POST("/", playerController.CreatePlayer)
				players.PUT("/:id", playerController.UpdatePlayer)
				players.DELETE("/:id", playerController.DeletePlayer)
				players.PUT("/:id/status", playerController.UpdatePlayerStatus)

				// 玩家點數管理
				players.GET("/:id/balance", playerController.GetPlayerBalance)
				players.POST("/:id/deposit", playerController.DepositPlayerBalance)
				players.POST("/:id/withdraw", playerController.WithdrawPlayerBalance)
				players.GET("/:id/transactions", playerController.GetPlayerTransactions)

				// 玩家限制管理
				players.POST("/:id/restrictions", playerController.SetPlayerRestriction)
				players.GET("/:id/restrictions", playerController.GetPlayerRestrictions)
				players.DELETE("/:id/restrictions/:restriction_id", playerController.RemovePlayerRestriction)

				// 玩家風險評估
				players.POST("/:id/risk-assessment", playerController.AssessPlayerRisk)
				players.GET("/:id/risk-history", playerController.GetPlayerRiskHistory)

				// 玩家註銷功能
				players.POST("/:id/deactivate", playerController.DeactivatePlayer)
				players.GET("/:id/deactivation-history", playerController.GetPlayerDeactivationHistory)

				// 玩家行為分析
				players.POST("/:id/behavior-analysis", playerController.AnalyzePlayerBehavior)
				players.POST("/:id/game-preference", playerController.AnalyzePlayerGamePreference)
				players.POST("/:id/spending-habits", playerController.AnalyzePlayerSpendingHabits)
			}

			// 遊戲管理路由
			games := authenticated.Group("/games")
			{
				games.GET("/", controllers.GetGames)
				games.GET("/:id", controllers.GetGame)
				games.POST("/", controllers.CreateGame)
				games.PUT("/:id", controllers.UpdateGame)
				games.DELETE("/:id", controllers.DeleteGame)
				games.PUT("/:id/status", controllers.UpdateGameStatus)

				// 遊戲配置管理
				games.GET("/:id/config", controllers.GetGameConfig)
				games.PUT("/:id/config", controllers.UpdateGameConfig)

				// 賠率管理
				games.GET("/:id/odds", controllers.GetGameOdds)
				games.PUT("/:id/odds", controllers.UpdateGameOdds)

				// 遊戲統計
				games.GET("/:id/stats", controllers.GetGameStats)
			}

			// 財務管理路由
			financial := authenticated.Group("/financial")
			{
				// 交易記錄
				financial.GET("/transactions", controllers.GetTransactions)
				financial.GET("/transactions/:id", controllers.GetTransaction)

				// 儲值管理
				financial.GET("/deposits", controllers.GetDeposits)
				financial.POST("/deposits", controllers.CreateDeposit)
				financial.PUT("/deposits/:id/confirm", controllers.ConfirmDeposit)

				// 提領管理
				financial.GET("/withdrawals", controllers.GetWithdrawals)
				financial.POST("/withdrawals", controllers.CreateWithdrawal)
				financial.PUT("/withdrawals/:id/approve", controllers.ApproveWithdrawal)

				// 對帳報表
				financial.GET("/reconciliation/daily", controllers.GetDailyReconciliation)
				financial.GET("/reconciliation/monthly", controllers.GetMonthlyReconciliation)
			}

			// 代理商管理路由
			agents := authenticated.Group("/agents")
			{
				agents.GET("/", controllers.GetAgents)
				agents.GET("/:id", controllers.GetAgent)
				agents.POST("/", controllers.CreateAgent)
				agents.PUT("/:id", controllers.UpdateAgent)
				agents.DELETE("/:id", controllers.DeleteAgent)
				agents.PUT("/:id/status", controllers.UpdateAgentStatus)

				// 經銷商管理
				agents.GET("/:id/dealers", controllers.GetAgentDealers)
				agents.POST("/:id/dealers", controllers.CreateDealer)

				// 分潤管理
				agents.GET("/:id/commission", controllers.GetAgentCommission)
				agents.PUT("/:id/commission", controllers.UpdateAgentCommission)
				agents.GET("/:id/settlements", controllers.GetAgentSettlements)
			}

			// 經銷商管理路由
			dealers := authenticated.Group("/dealers")
			{
				dealers.GET("/", controllers.GetDealers)
				dealers.GET("/:id", controllers.GetDealer)
				dealers.PUT("/:id", controllers.UpdateDealer)
				dealers.DELETE("/:id", controllers.DeleteDealer)
				dealers.PUT("/:id/status", controllers.UpdateDealerStatus)

				// 經銷商分潤
				dealers.GET("/:id/commission", controllers.GetDealerCommission)
				dealers.GET("/:id/settlements", controllers.GetDealerSettlements)

				// 經銷商玩家
				dealers.GET("/:id/players", controllers.GetDealerPlayers)
			}

			// 報表管理路由
			reports := authenticated.Group("/reports")
			{
				// 營運報表
				reports.GET("/dashboard", controllers.GetDashboardData)
				reports.GET("/revenue", controllers.GetRevenueReport)
				reports.GET("/player-analysis", controllers.GetPlayerAnalysisReport)
				reports.GET("/game-performance", controllers.GetGamePerformanceReport)

				// 代理商報表
				reports.GET("/agent-performance", controllers.GetAgentPerformanceReport)
				reports.GET("/commission-summary", controllers.GetCommissionSummaryReport)

				// 自訂報表
				reports.POST("/custom", controllers.GenerateCustomReport)
				reports.GET("/export/:type", controllers.ExportReport)
			}

			// 系統管理路由
			admin := authenticated.Group("/admin")
			adminMiddleware := controllers.NewAuthController().AdminPermissionMiddleware()
			admin.Use(adminMiddleware) // 需要管理員權限
			{
				// 角色權限管理
				admin.GET("/roles", controllers.GetRoles)
				admin.POST("/roles", controllers.CreateRole)
				admin.PUT("/roles/:id", controllers.UpdateRole)
				admin.DELETE("/roles/:id", controllers.DeleteRole)

				// 權限管理
				admin.GET("/permissions", controllers.GetPermissions)
				admin.POST("/permissions", controllers.CreatePermission)

				// 操作日誌
				admin.GET("/logs", controllers.GetOperationLogs)
				admin.GET("/logs/:id", controllers.GetOperationLog)

				// 系統設置
				admin.GET("/settings", controllers.GetSystemSettings)
				admin.PUT("/settings", controllers.UpdateSystemSettings)
			}
		}
	}
}

// SetupAPIV2Routes 設置 API v2 路由（預留未來版本）
func SetupAPIV2Routes(r *gin.Engine) {
	v2 := r.Group("/api/v2")
	{
		v2.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"version": "2.0",
				"status":  "under_development",
			})
		})
	}
}
