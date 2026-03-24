// internal/mqtt/client.go
package mqtt

import (
	"fmt"
	"log"
	"time"

	"corp-vitals-mqtt-sim-go/internal/model"
	"corp-vitals-mqtt-sim-go/pkg/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// CommandProcessor 命令处理器接口
type CommandProcessor interface {
	ProcessCommand(cmd model.CommandRequest)
}

// Client MQTT 客户端封装
type Client struct {
	client mqtt.Client
	cfg    *config.Config
	cmdSvc CommandProcessor
}

func NewClient(cfg *config.Config, cmdSvc CommandProcessor) *Client {
	return &Client{cfg: cfg, cmdSvc: cmdSvc}
}

func (c *Client) Connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.cfg.MQTTURL)
	opts.SetClientID(c.cfg.MQTTClientID)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	opts.SetOnConnectHandler(c.onConnectHandler)
	opts.SetConnectionLostHandler(c.connectionLostHandler)

	log.Printf("INFO: Connecting to MQTT broker: %s (ClientID: %s)", c.cfg.MQTTURL, c.cfg.MQTTClientID)

	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("could not connect to MQTT broker: %w", token.Error())
	}
	return nil
}

func (c *Client) onConnectHandler(client mqtt.Client) {
	log.Println("INFO: Connected to MQTT broker.")
	// 如果需要全局订阅命令，可在此添加：go c.subscribeToCommands("devices/+/commands")
}

func (c *Client) connectionLostHandler(client mqtt.Client, err error) {
	log.Printf("WARN: MQTT connection lost: %v. Reconnecting...", err)
}

func (c *Client) Publish(topic string, payload []byte) {
	if c.client == nil || !c.client.IsConnected() {
		return
	}
	token := c.client.Publish(topic, 1, false, payload)
	go func() {
		if token.WaitTimeout(3*time.Second) && token.Error() != nil {
			log.Printf("ERROR: Failed to publish to %s: %v", topic, token.Error())
		}
	}()
}

// PublishTelemetry 发布具体某个设备的遥测数据
func (c *Client) PublishTelemetry(deviceID string, payload []byte) {
	topic := fmt.Sprintf("devices/%s/telemetry", deviceID)
	c.Publish(topic, payload)
}

func (c *Client) Disconnect() {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
		log.Println("INFO: Disconnected from MQTT broker.")
	}
}
