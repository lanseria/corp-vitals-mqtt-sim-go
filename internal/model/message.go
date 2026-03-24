package model

// CommandRequest 通用命令请求结构
type CommandRequest struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// TelemetryPayload 遥测数据结构示例
type TelemetryPayload struct {
	DeviceID  string                 `json:"deviceId"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}
