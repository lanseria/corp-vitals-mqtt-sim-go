// internal/model/device.go
package model

import "sync"

// DeviceState 设备的实际状态数据
type DeviceState struct {
	DeviceType  string `json:"deviceType"`  // 设备类型
	DeviceID    string `json:"deviceId"`    // 设备编号
	AlarmStatus string `json:"alarmStatus"` // 报警状态："无" 或者 "用户SOS报警"
	Battery     int    `json:"battery"`     // 剩余电量 (0-100)
	HeartRate   int    `json:"heartRate"`   // 心率
	BloodOxygen int    `json:"bloodOxygen"` // 血氧 (0-100)
}

// Device 并发安全的设备包装器
type Device struct {
	mu    sync.RWMutex
	state DeviceState
}

func NewDevice(id, deviceType string) *Device {
	return &Device{
		state: DeviceState{
			DeviceID:    id,
			DeviceType:  deviceType,
			AlarmStatus: "无",
			Battery:     100, // 初始电量100%
			HeartRate:   75,  // 初始心率 75 bpm
			BloodOxygen: 98,  // 初始血氧 98%
		},
	}
}

// Snapshot 获取设备的当前状态快照
func (d *Device) Snapshot() DeviceState {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state
}

// Update 更新设备属性 (仅更新传入非零/非空的值)
func (d *Device) Update(updates DeviceState) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if updates.AlarmStatus != "" {
		d.state.AlarmStatus = updates.AlarmStatus
	}
	if updates.Battery > 0 {
		d.state.Battery = updates.Battery
	}
	if updates.HeartRate > 0 {
		d.state.HeartRate = updates.HeartRate
	}
	if updates.BloodOxygen > 0 {
		d.state.BloodOxygen = updates.BloodOxygen
	}
}
