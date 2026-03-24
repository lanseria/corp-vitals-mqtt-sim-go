// cmd/app/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"corp-vitals-sim-go/internal/device"
	"corp-vitals-sim-go/internal/webhook"
	"corp-vitals-sim-go/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	// 1. 初始化 Webhook 客户端
	webhookClient := webhook.NewClient(cfg)

	// 2. 初始化设备管理器并启动仿真流
	deviceManager := device.NewManager(webhookClient)
	deviceManager.StartSimulation()

	// 3. 启动 HTTP 调试接口
	mux := http.NewServeMux()
	mux.HandleFunc("GET /devices", deviceManager.HandleListDevices)
	mux.HandleFunc("POST /devices/{id}", deviceManager.HandleUpdateDevice)

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: mux,
	}

	go func() {
		log.Printf("INFO: HTTP Debug API started on port %s", cfg.HTTPPort)
		log.Printf("INFO: Pushing data to Webhook URL: %s", cfg.WebhookURL)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("FATAL: HTTP server error: %v", err)
		}
	}()

	// 4. 优雅关停
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("INFO: Shutting down...")
	httpServer.Close()
	log.Println("INFO: Server stopped.")
}
