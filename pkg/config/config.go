// pkg/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 保存应用的所有配置
type Config struct {
	WebhookURL string // Webhook 推送地址
	HTTPPort   string // 本地 HTTP 调试服务端口
}

// Load 从环境变量加载配置
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("WARN: No .env file found, using system env vars.")
	}

	return &Config{
		WebhookURL: getEnv("WEBHOOK_URL", "http://localhost:30003/api/webhook/device-data"),
		HTTPPort:   getEnv("HTTP_PORT", "8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
