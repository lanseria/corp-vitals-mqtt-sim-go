package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mqtt-service-template/internal/model"
	"mqtt-service-template/internal/mqtt"
	"mqtt-service-template/pkg/config"
)

// App 应用结构体, 实现 CommandProcessor 接口
type App struct {
	mqttClient *mqtt.Client
}

// ProcessCommand 处理接收到的 MQTT 命令
func (a *App) ProcessCommand(cmd model.CommandRequest) {
	log.Printf("INFO: Processing command: %+v", cmd)

	// 在这里添加你的命令处理逻辑
	switch cmd.Type {
	case "ECHO":
		a.handleEchoCommand(cmd)
	default:
		log.Printf("WARN: Unknown command type: %s", cmd.Type)
	}
}

func (a *App) handleEchoCommand(cmd model.CommandRequest) {
	response := map[string]interface{}{
		"type":      "ECHO_RESPONSE",
		"timestamp": time.Now().Unix(),
		"payload":   cmd.Payload,
	}
	if data, err := json.Marshal(response); err == nil {
		a.mqttClient.PublishTelemetry(data)
	}
}

func main() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	// 2. 初始化应用和 MQTT 客户端
	app := &App{}
	mqttClient := mqtt.NewClient(cfg, app)
	app.mqttClient = mqttClient

	// 3. 连接 MQTT (带重试)
	go func() {
		for {
			if err := mqttClient.Connect(); err != nil {
				log.Printf("ERROR: MQTT connect failed, retry in 10s: %v", err)
				time.Sleep(10 * time.Second)
			} else {
				break
			}
		}
	}()

	// 4. 启动心跳定时器 (示例)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			heartbeat := map[string]interface{}{
				"deviceId":  cfg.DeviceID,
				"timestamp": time.Now().Unix(),
				"status":    "online",
			}
			if data, err := json.Marshal(heartbeat); err == nil {
				mqttClient.PublishHeartbeat(data)
			}
		}
	}()

	log.Printf("INFO: Service started (DeviceID: %s)", cfg.DeviceID)

	// 5. 优雅关停
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("INFO: Shutting down...")
	mqttClient.Disconnect()
	log.Println("INFO: Server stopped.")
}
