// pkg/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 保存应用的所有配置
type Config struct {
	MQTTURL      string // MQTT Broker 地址
	MQTTClientID string // 客户端唯一标识
	HTTPPort     string // 本地 HTTP 服务端口
}

// Load 从环境变量加载配置
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("WARN: No .env file found, using system env vars.")
	}

	cfg := &Config{
		MQTTURL:      getEnv("MQTT_URL", "tcp://localhost:1883"), // 云端 HiveMQ 地址请在 .env 中配置
		MQTTClientID: getEnv("MQTT_CLIENT_ID", "sim-client-001"),
		HTTPPort:     getEnv("HTTP_PORT", "8080"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
