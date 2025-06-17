package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config 應用程式配置結構
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	JWT      JWTConfig      `json:"jwt"`
	Security SecurityConfig `json:"security"`
	Game     GameConfig     `json:"game"`
}

// ServerConfig 伺服器配置
type ServerConfig struct {
	Port         string        `json:"port"`
	Host         string        `json:"host"`
	Mode         string        `json:"mode"` // debug, release, test
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// DatabaseConfig 資料庫配置
type DatabaseConfig struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	Username     string        `json:"username"`
	Password     string        `json:"password"`
	Database     string        `json:"database"`
	MaxOpenConns int           `json:"max_open_conns"`
	MaxIdleConns int           `json:"max_idle_conns"`
	MaxLifetime  time.Duration `json:"max_lifetime"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	PoolSize int    `json:"pool_size"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string        `json:"secret"`
	ExpireTime time.Duration `json:"expire_time"`
	Issuer     string        `json:"issuer"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	PasswordMinLength  int           `json:"password_min_length"`
	MaxLoginAttempts   int           `json:"max_login_attempts"`
	LockoutDuration    time.Duration `json:"lockout_duration"`
	AllowedOrigins     []string      `json:"allowed_origins"`
	RateLimitPerMinute int           `json:"rate_limit_per_minute"`
	SessionTimeout     time.Duration `json:"session_timeout"`
}

// GameConfig 遊戲配置
type GameConfig struct {
	DefaultCurrency       string  `json:"default_currency"`
	MinBetAmount          float64 `json:"min_bet_amount"`
	MaxBetAmount          float64 `json:"max_bet_amount"`
	HouseEdge             float64 `json:"house_edge"`
	MaxPlayersPerTable    int     `json:"max_players_per_table"`
	SessionTimeoutMinutes int     `json:"session_timeout_minutes"`
}

// 全域配置實例
var AppConfig *Config

// LoadConfig 載入應用程式配置
func LoadConfig() (*Config, error) {
	// 載入 .env 檔案
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Mode:         getEnv("GIN_MODE", "debug"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "33061"),
			Username:     getEnv("DB_USERNAME", "root"),
			Password:     getEnv("DB_PASSWORD", "rootpassword"),
			Database:     getEnv("DB_DATABASE", "nexus_gaming"),
			MaxOpenConns: getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getIntEnv("DB_MAX_IDLE_CONNS", 5),
			MaxLifetime:  getDurationEnv("DB_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "63791"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
			PoolSize: getIntEnv("REDIS_POOL_SIZE", 10),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "nexus-gaming-secret-key-change-in-production"),
			ExpireTime: getDurationEnv("JWT_EXPIRE_TIME", 24*time.Hour),
			Issuer:     getEnv("JWT_ISSUER", "nexus-gaming"),
		},
		Security: SecurityConfig{
			PasswordMinLength:  getIntEnv("PASSWORD_MIN_LENGTH", 8),
			MaxLoginAttempts:   getIntEnv("MAX_LOGIN_ATTEMPTS", 5),
			LockoutDuration:    getDurationEnv("LOCKOUT_DURATION", 15*time.Minute),
			AllowedOrigins:     getStringSliceEnv("ALLOWED_ORIGINS", []string{"*"}),
			RateLimitPerMinute: getIntEnv("RATE_LIMIT_PER_MINUTE", 60),
			SessionTimeout:     getDurationEnv("SESSION_TIMEOUT", 30*time.Minute),
		},
		Game: GameConfig{
			DefaultCurrency:       getEnv("DEFAULT_CURRENCY", "TWD"),
			MinBetAmount:          getFloatEnv("MIN_BET_AMOUNT", 1.0),
			MaxBetAmount:          getFloatEnv("MAX_BET_AMOUNT", 10000.0),
			HouseEdge:             getFloatEnv("HOUSE_EDGE", 0.025),
			MaxPlayersPerTable:    getIntEnv("MAX_PLAYERS_PER_TABLE", 6),
			SessionTimeoutMinutes: getIntEnv("GAME_SESSION_TIMEOUT", 30),
		},
	}

	// 設定全域配置
	AppConfig = config

	return config, nil
}

// GetDSN 取得資料庫連線字串
func (db *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db.Username, db.Password, db.Host, db.Port, db.Database)
}

// GetRedisAddr 取得 Redis 連線地址
func (r *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

// 輔助函式：從環境變數取得字串值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 輔助函式：從環境變數取得整數值
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// 輔助函式：從環境變數取得浮點數值
func getFloatEnv(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// 輔助函式：從環境變數取得時間長度值
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// 輔助函式：從環境變數取得字串陣列值
func getStringSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// 簡單實作，以逗號分隔
		// 可以根據需求使用更複雜的解析邏輯
		return []string{value}
	}
	return defaultValue
}

// InitConfig 初始化配置（main.go 調用的函數別名）
func InitConfig() error {
	_, err := LoadConfig()
	return err
}

// GetJWTSecret 獲取 JWT 密鑰
func GetJWTSecret() string {
	if AppConfig != nil {
		return AppConfig.JWT.Secret
	}
	// 如果配置未初始化，從環境變數取得
	return getEnv("JWT_SECRET", "nexus-gaming-secret-key-change-in-production")
}
