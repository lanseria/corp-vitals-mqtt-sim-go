// cmd/app/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"corp-vitals-mqtt-sim-go/internal/device"
	"corp-vitals-mqtt-sim-go/internal/mqtt"
	"corp-vitals-mqtt-sim-go/pkg/config"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	// 2. 初始化 MQTT 客户端
	mqttClient := mqtt.NewClient(cfg, nil) // 目前模拟器只上报数据，不需要主动接受命令可传 nil

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

	// 4. 初始化设备管理器并启动 5 个设备的仿真流
	deviceManager := device.NewManager(mqttClient)
	deviceManager.StartSimulation()

	// 5. 启动 HTTP 接口，用于查看与调试修改设备状态
	mux := http.NewServeMux()
	mux.HandleFunc("GET /devices", deviceManager.HandleListDevices)
	mux.HandleFunc("POST /devices/{id}", deviceManager.HandleUpdateDevice)

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: mux,
	}

	go func() {
		log.Printf("INFO: HTTP Debug API started on port %s", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: HTTP server error: %v", err)
		}
	}()

	// 6. 优雅关停
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("INFO: Shutting down...")
	mqttClient.Disconnect()
	httpServer.Close()
	log.Println("INFO: Server stopped.")
}
