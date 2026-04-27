package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	AWSServices AWSConfig
	JWT         JWTConfig
	Uploads     UploadConfig
	SMTP        SMTPConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Mode     string
}

type AWSConfig struct {
	Region         string
	Key            string
	KeyId          string
	S3Bucket       string
	S3Endpoint     string
	EventQueueName string
}

type JWTConfig struct {
	JWTSecret                 string
	JWTTokenExpiration        time.Duration
	JWTRefreshTokenExpiration time.Duration
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type UploadConfig struct {
	Path       string
	UploadSize int64

	// Upload provider can be s3 or locally
	UploadProvider string
}

func LoadEnv() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiresAt, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "1h"))
	refreshTokenExpiresAt, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRES_IN", "72h"))
	maxUpload, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "1025"))

	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},

		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Name:     getEnv("DB_NAME", "armory_db"),
			Port:     getEnv("DB_PORT", "5432"),
			Password: getEnv("DB_PASSWORD", "password"),
			User:     getEnv("DB_USER", "postgres"),
			Mode:     getEnv("DB_SSL_MODE", "disable"),
		},
		AWSServices: AWSConfig{
			Region:         getEnv("AWS_REGION", "us-east-1"),
			Key:            getEnv("AWS_SECRET_ACCESS_KEY", "test"),
			KeyId:          getEnv("AWS_ACCESS_KEY_ID", "test"),
			S3Bucket:       getEnv("AWS_S3_BUCKET", "armory-db-bucket"),
			S3Endpoint:     getEnv("AWS_S3_ENDPOINT", "http://localhost:4566"),
			EventQueueName: getEnv("EVENT_QUEUE_NAME", "armory-events"),
		},
		JWT: JWTConfig{
			JWTSecret:                 getEnv("JWT_SCRET", "your_jwt_secret"),
			JWTTokenExpiration:        jwtExpiresAt,
			JWTRefreshTokenExpiration: refreshTokenExpiresAt,
		},
		Uploads: UploadConfig{
			Path:           getEnv("UPLOAD_DIR", "/tmp/uploads"),
			UploadSize:     maxUpload,
			UploadProvider: getEnv("UPLOAD_PROVIDER", "local"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     smtpPort,
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@shop.com"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}
