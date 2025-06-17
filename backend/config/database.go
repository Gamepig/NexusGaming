package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// DB MySQL 資料庫連線實例
	DB *sql.DB
	// RedisClient Redis 連線實例
	RedisClient *redis.Client
)

// InitDatabase 初始化資料庫連線
func InitDatabase(config *Config) error {
	// 初始化 MySQL
	if err := initMySQL(&config.Database); err != nil {
		return fmt.Errorf("failed to init MySQL: %w", err)
	}

	// 初始化 Redis
	if err := initRedis(&config.Redis); err != nil {
		return fmt.Errorf("failed to init Redis: %w", err)
	}

	log.Println("Database connections initialized successfully")
	return nil
}

// initMySQL 初始化 MySQL 連線
func initMySQL(config *DatabaseConfig) error {
	dsn := config.GetDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	// 設定連線池參數
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.MaxLifetime)

	// 測試連線
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping MySQL: %w", err)
	}

	DB = db
	log.Printf("MySQL connected successfully to %s:%s/%s", config.Host, config.Port, config.Database)
	return nil
}

// initRedis 初始化 Redis 連線
func initRedis(config *RedisConfig) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddr(),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	// 測試連線
	ctx := rdb.Context()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	RedisClient = rdb
	log.Printf("Redis connected successfully to %s", config.GetRedisAddr())
	return nil
}

// CloseDatabase 關閉資料庫連線
func CloseDatabase() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Error closing MySQL connection: %v", err)
		} else {
			log.Println("MySQL connection closed")
		}
	}

	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		} else {
			log.Println("Redis connection closed")
		}
	}
}

// GetDB 取得資料庫連線
func GetDB() *sql.DB {
	return DB
}

// GetRedis 取得 Redis 連線
func GetRedis() *redis.Client {
	return RedisClient
}

// CheckDatabaseHealth 檢查資料庫健康狀態
func CheckDatabaseHealth() error {
	if DB == nil {
		return fmt.Errorf("MySQL connection is nil")
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("MySQL ping failed: %w", err)
	}

	if RedisClient == nil {
		return fmt.Errorf("Redis connection is nil")
	}

	ctx := RedisClient.Context()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis ping failed: %w", err)
	}

	return nil
}

// ConnectDatabase 建立資料庫連線（main.go 調用的函數別名）
func ConnectDatabase() error {
	if AppConfig == nil {
		return fmt.Errorf("config not initialized, call InitConfig() first")
	}
	return InitDatabase(AppConfig)
}
