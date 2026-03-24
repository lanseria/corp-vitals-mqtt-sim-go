package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"mqtt-service-template/internal/model"
	"mqtt-service-template/pkg/config"

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

// NewClient 创建 MQTT 客户端
func NewClient(cfg *config.Config, cmdSvc CommandProcessor) *Client {
	return &Client{cfg: cfg, cmdSvc: cmdSvc}
}

// Connect 连接到 MQTT Broker
func (c *Client) Connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.cfg.MQTTURL)
	opts.SetClientID(c.cfg.MQTTClientID)

	// 启用自动重连
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)

	// 设置连接事件处理器
	opts.SetOnConnectHandler(c.onConnectHandler)
	opts.SetConnectionLostHandler(c.connectionLostHandler)

	log.Printf("INFO: Connecting to MQTT broker: %s (ClientID: %s)",
		c.cfg.MQTTURL, c.cfg.MQTTClientID)

	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("could not connect to MQTT broker: %w", token.Error())
	}
	return nil
}

// onConnectHandler 连接成功回调
func (c *Client) onConnectHandler(client mqtt.Client) {
	log.Println("INFO: Connected to MQTT broker.")
	go c.subscribeToCommands()
}

// connectionLostHandler 连接丢失回调
func (c *Client) connectionLostHandler(client mqtt.Client, err error) {
	log.Printf("WARN: MQTT connection lost: %v. Reconnecting...", err)
}

// Publish 发布消息到指定 Topic
func (c *Client) Publish(topic string, payload []byte) {
	if c.client == nil || !c.client.IsConnected() {
		log.Println("WARN: MQTT client not connected, skipping publish.")
		return
	}

	token := c.client.Publish(topic, 1, false, payload)
	go func() {
		if token.WaitTimeout(3*time.Second) && token.Error() != nil {
			log.Printf("ERROR: Failed to publish to %s: %v", topic, token.Error())
		}
	}()
}

// PublishTelemetry 发布遥测数据
func (c *Client) PublishTelemetry(payload []byte) {
	topic := fmt.Sprintf("devices/%s/telemetry", c.cfg.DeviceID)
	c.Publish(topic, payload)
}

// PublishHeartbeat 发布心跳
func (c *Client) PublishHeartbeat(payload []byte) {
	topic := fmt.Sprintf("devices/%s/heartbeat", c.cfg.DeviceID)
	c.Publish(topic, payload)
}

// Disconnect 断开连接
func (c *Client) Disconnect() {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
		log.Println("INFO: Disconnected from MQTT broker.")
	}
}

// subscribeToCommands 订阅命令 Topic
func (c *Client) subscribeToCommands() {
	topic := fmt.Sprintf("devices/%s/commands", c.cfg.DeviceID)
	if token := c.client.Subscribe(topic, 1, c.commandHandler); token.Wait() && token.Error() != nil {
		log.Printf("ERROR: Failed to subscribe to %s: %v", topic, token.Error())
	} else {
		log.Printf("INFO: Subscribed to commands topic: %s", topic)
	}
}

// commandHandler 命令消息处理
func (c *Client) commandHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("INFO: Received command on %s: %s", msg.Topic(), string(msg.Payload()))

	if c.cmdSvc == nil {
		return
	}

	var cmd model.CommandRequest
	if err := json.Unmarshal(msg.Payload(), &cmd); err != nil {
		log.Printf("ERROR: Failed to unmarshal command: %v", err)
		return
	}

	c.cmdSvc.ProcessCommand(cmd)
}
