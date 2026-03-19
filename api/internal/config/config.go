package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	MinIO    MinIOConfig
	Asynq    AsynqConfig
	Midtrans MidtransConfig
}

type AppConfig struct {
	Name    string
	Env     string
	Version string
}

type HTTPConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	URL string
	TTL time.Duration
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type AsynqConfig struct {
	RedisURL    string
	Concurrency int
}

type MidtransConfig struct {
	ServerKey string
	ClientKey string
	Env       string
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	v.SetDefault("APP_NAME", "bmt-saas")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_VERSION", "1.0.0")
	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("HTTP_READ_TIMEOUT", "30s")
	v.SetDefault("HTTP_WRITE_TIMEOUT", "30s")
	v.SetDefault("HTTP_IDLE_TIMEOUT", "120s")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)
	v.SetDefault("DB_CONN_MAX_LIFETIME", "5m")
	v.SetDefault("REDIS_TTL", "5m")
	v.SetDefault("JWT_ACCESS_EXPIRY", "15m")
	v.SetDefault("JWT_REFRESH_EXPIRY", "168h")
	v.SetDefault("ASYNQ_CONCURRENCY", 10)
	v.SetDefault("MINIO_USE_SSL", false)

	_ = v.ReadInConfig()

	return &Config{
		App: AppConfig{
			Name:    v.GetString("APP_NAME"),
			Env:     v.GetString("APP_ENV"),
			Version: v.GetString("APP_VERSION"),
		},
		HTTP: HTTPConfig{
			Port:         v.GetString("HTTP_PORT"),
			ReadTimeout:  v.GetDuration("HTTP_READ_TIMEOUT"),
			WriteTimeout: v.GetDuration("HTTP_WRITE_TIMEOUT"),
			IdleTimeout:  v.GetDuration("HTTP_IDLE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			URL:             v.GetString("DATABASE_URL"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: v.GetDuration("DB_CONN_MAX_LIFETIME"),
		},
		Redis: RedisConfig{
			URL: v.GetString("REDIS_URL"),
			TTL: v.GetDuration("REDIS_TTL"),
		},
		JWT: JWTConfig{
			AccessSecret:  v.GetString("JWT_ACCESS_SECRET"),
			RefreshSecret: v.GetString("JWT_REFRESH_SECRET"),
			AccessExpiry:  v.GetDuration("JWT_ACCESS_EXPIRY"),
			RefreshExpiry: v.GetDuration("JWT_REFRESH_EXPIRY"),
		},
		MinIO: MinIOConfig{
			Endpoint:  v.GetString("MINIO_ENDPOINT"),
			AccessKey: v.GetString("MINIO_ACCESS_KEY"),
			SecretKey: v.GetString("MINIO_SECRET_KEY"),
			Bucket:    v.GetString("MINIO_BUCKET"),
			UseSSL:    v.GetBool("MINIO_USE_SSL"),
		},
		Asynq: AsynqConfig{
			RedisURL:    v.GetString("REDIS_URL"),
			Concurrency: v.GetInt("ASYNQ_CONCURRENCY"),
		},
		Midtrans: MidtransConfig{
			ServerKey: v.GetString("MIDTRANS_SERVER_KEY"),
			ClientKey: v.GetString("MIDTRANS_CLIENT_KEY"),
			Env:       v.GetString("MIDTRANS_ENV"),
		},
	}, nil
}
