// internal/device/manager.go
package device

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"corp-vitals-mqtt-sim-go/internal/model"
	"corp-vitals-mqtt-sim-go/internal/mqtt"
)

// Manager 管理多台设备并提供 HTTP 控制接口
type Manager struct {
	devices map[string]*model.Device
	mqtt    *mqtt.Client
}

// NewManager 初始化并创建 5 个模拟设备
func NewManager(mqttClient *mqtt.Client) *Manager {
	m := &Manager{
		devices: make(map[string]*model.Device),
		mqtt:    mqttClient,
	}

	// 模拟初始化 5 个设备
	for i := 1; i <= 5; i++ {
		id := fmt.Sprintf("DEV-VITAL-%03d", i)
		m.devices[id] = model.NewDevice(id, "VitalBand-X1")
	}

	return m
}

// StartSimulation 启动周期性上报数据
func (m *Manager) StartSimulation() {
	ticker := time.NewTicker(5 * time.Second) // 每 5 秒上报一次
	go func() {
		for range ticker.C {
			for deviceID, d := range m.devices {
				state := d.Snapshot()

				// 封装符合云端接收格式的 Payload
				payload := map[string]interface{}{
					"deviceId":  deviceID,
					"timestamp": time.Now().UnixMilli(),
					"vitals":    state,
				}

				if data, err := json.Marshal(payload); err == nil {
					m.mqtt.PublishTelemetry(deviceID, data)
				}
			}
		}
	}()
}

// --- HTTP API 处理函数 ---

// HandleListDevices GET /devices
func (m *Manager) HandleListDevices(w http.ResponseWriter, r *http.Request) {
	var list []model.DeviceState
	for _, d := range m.devices {
		list = append(list, d.Snapshot())
	}

	// 排序保证输出顺序固定
	sort.Slice(list, func(i, j int) bool {
		return list[i].DeviceID < list[j].DeviceID
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

// HandleUpdateDevice POST /devices/{id}
func (m *Manager) HandleUpdateDevice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") // Go 1.22+ 路由特性
	d, ok := m.devices[id]
	if !ok {
		http.Error(w, `{"error": "Device not found"}`, http.StatusNotFound)
		return
	}

	var updates model.DeviceState
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, `{"error": "Invalid JSON payload"}`, http.StatusBadRequest)
		return
	}

	// 应用更新 (这里你可以动态传入诸如 心率120、SOS报警 等进行调试)
	d.Update(updates)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d.Snapshot())
}
