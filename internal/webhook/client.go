// internal/webhook/client.go
package webhook

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"corp-vitals-sim-go/pkg/config"
)

// Client Webhook 推送客户端
type Client struct {
	cfg        *config.Config
	httpClient *http.Client
}

// NewClient 初始化 Webhook 客户端
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // 设置 5 秒超时，避免阻塞
		},
	}
}

// SendData 发送数据到 Webhook
func (c *Client) SendData(payload []byte) {
	if c.cfg.WebhookURL == "" {
		return
	}

	req, err := http.NewRequest("POST", c.cfg.WebhookURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("ERROR: Failed to create webhook request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("ERROR: Webhook push failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("WARN: Webhook returned non-200 status: %d", resp.StatusCode)
	}
}
