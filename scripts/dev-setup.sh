#!/bin/bash

echo "🚀 NexusGaming 開發環境設置"
echo "=============================="

# 檢查 Docker 是否安裝
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安裝，請先安裝 Docker"
    exit 1
fi

# 檢查 Docker Compose 是否安裝
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安裝，請先安裝 Docker Compose"
    exit 1
fi

# 啟動資料庫服務
echo "📦 啟動 Docker 服務..."
docker-compose up -d

# 等待 MySQL 啟動
echo "⏳ 等待 MySQL 服務啟動..."
sleep 15

# 檢查 MySQL 連線
echo "🔍 檢查 MySQL 連線..."
if docker exec nexus-gaming-mysql mysql -u nexus_user -pnexus_password -e "SELECT 1;" &> /dev/null; then
    echo "✅ MySQL 連線成功"
else
    echo "❌ MySQL 連線失敗"
    exit 1
fi

# 檢查 Redis 連線
echo "🔍 檢查 Redis 連線..."
if docker exec nexus-gaming-redis redis-cli ping &> /dev/null; then
    echo "✅ Redis 連線成功"
else
    echo "❌ Redis 連線失敗"
    exit 1
fi

echo ""
echo "🎉 開發環境設置完成！"
echo ""
echo "📋 服務資訊："
echo "   • MySQL:      localhost:3306"
echo "   • Redis:      localhost:6379"
echo "   • PHPMyAdmin: http://localhost:8081"
echo "   • 後端API:    http://localhost:8080"
echo "   • 前端:       http://localhost:3002"
echo ""
echo "📁 資料庫登入資訊："
echo "   • 使用者: nexus_user"
echo "   • 密碼:   nexus_password"
echo "   • 資料庫: nexus_gaming"
echo ""
echo "🔧 下一步："
echo "   1. 複製 env.example 到 .env 並調整設定"
echo "   2. 啟動後端服務: cd backend && go run main.go"
echo "   3. 啟動前端服務: cd frontend && npm run dev" 