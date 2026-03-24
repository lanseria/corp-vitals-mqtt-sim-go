package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 保存应用的所有配置
type Config struct {
	// MQTT Broker 配置
	MQTTURL      string // MQTT Broker 地址, 例: tcp://localhost:1883
	MQTTClientID string // 客户端唯一标识
	DeviceID     string // 设备唯一标识, 用于构建 MQTT Topic
}

// Load 从环境变量加载配置
func Load() (*Config, error) {
	// 尝试加载 .env 文件 (本地开发用)
	if err := godotenv.Load(); err != nil {
		log.Println("WARN: No .env file found, using system env vars.")
	}

	cfg := &Config{
		MQTTURL:      getEnv("MQTT_URL", "tcp://localhost:1883"),
		MQTTClientID: getEnv("MQTT_CLIENT_ID", "mqtt-client-001"),
		DeviceID:     getEnv("DEVICE_ID", "device-001"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
