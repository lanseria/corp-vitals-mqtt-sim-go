// internal/device/manager.go
package device

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"corp-vitals-sim-go/internal/model"
	"corp-vitals-sim-go/internal/webhook"
)

type Manager struct {
	devices map[string]*model.Device
	webhook *webhook.Client
}

func NewManager(whClient *webhook.Client) *Manager {
	m := &Manager{
		devices: make(map[string]*model.Device),
		webhook: whClient,
	}

	// 初始化随机种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 5; i++ {
		id := fmt.Sprintf("DEV-VITAL-%03d", i)
		device := model.NewDevice(id, "VitalBand-X1")

		// 为每个设备设置随机属性
		randomState := model.DeviceState{
			Battery:     60 + r.Intn(41),     // 60-100
			HeartRate:   60 + r.Intn(61),     // 60-120
			BloodOxygen: 94 + r.Intn(7),      // 94-100
			AlarmStatus: r.Intn(10) == 0,     // 10% 概率触发报警
		}
		device.Update(randomState)

		m.devices[id] = device
	}
	return m
}

func (m *Manager) StartSimulation() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			for deviceID, d := range m.devices {
				state := d.Snapshot()
				payload := map[string]interface{}{
					"deviceId":  deviceID,
					"timestamp": time.Now().UnixMilli(),
					"vitals":    state,
				}

				if data, err := json.Marshal(payload); err == nil {
					// 启动协程推送数据，防止阻塞定时器
					go m.webhook.SendData(data)
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
