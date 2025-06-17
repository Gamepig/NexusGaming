package main

import (
	"fmt"
	"log"
	"nexus-gaming-backend/config"
	"nexus-gaming-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化資料庫連接
	config.ConnectDatabase()

	// 初始化 Gin 路由器
	r := gin.Default()

	// 設置路由
	routes.SetupRoutes(r)
	routes.SetupAPIV2Routes(r)

	// 啟動服務器
	fmt.Println("Nexus Gaming Backend Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}
