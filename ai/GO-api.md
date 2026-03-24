# 体征监测设备模拟器 (corp-vitals-sim) 对接文档

## 1. 概述

本项目用于模拟多台“智能体征监测设备”（如智能手环/胸贴）。设备会持续采集用户的生命体征数据（心率、血氧、电量、报警状态等），并通过 **Webhook (HTTP POST)** 方式实时推送到第三方指定的业务系统接口。

第三方业务系统通过**提供一个标准的 HTTP 接口**来接收这些设备的实时数据。
同时，本项目提供了一组**本地 HTTP 接口**，用于在联调时动态改变设备的体征数据（如模拟突发心脏异常或触发 SOS 报警），以验证第三方系统的预警处理逻辑。

---

## 2. 第三方对接说明（Webhook 方式）

第三方系统作为 **数据接收服务端 (Webhook Receiver)** 接入。

### 2.1 接收配置

- **通信协议**: HTTP / HTTPS
- **请求方法**: `POST`
- **Content-Type**: `application/json`
- **响应要求**: 第三方接口在收到数据后，应立即返回 `200 OK` 状态码。

### 2.2 推送机制定义

模拟器会将采集到的实时遥测数据主动推送到预设的 URL。

- **推送触发**: 每 5 秒推送一次实时数据。
- **重试机制**: 若推送失败（如网络抖动或第三方响应非 200），模拟器通常会进行指数退避重试（根据具体配置而定）。
- **并发处理**: 建议第三方接收端具备异步处理能力，避免因逻辑耗时导致推送链路阻塞。

### 2.3 消息 Payload 数据结构 (JSON)

推送的数据包格式为 JSON：

```json
{
  "deviceId": "DEV-VITAL-001",
  "timestamp": 1711245678123,
  "vitals": {
    "deviceType": "VitalBand-X1",
    "deviceId": "DEV-VITAL-001",
    "alarmStatus": "无",
    "battery": 100,
    "heartRate": 75,
    "bloodOxygen": 98
  }
}
```

**字段说明**:
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `deviceId` | String | 设备编号 |
| `timestamp` | Number | 数据上报时间戳 (毫秒级) |
| `vitals.deviceType` | String | 设备类型名称 |
| `vitals.alarmStatus`| String | 报警状态，枚举值：`"无"` 或 `"用户SOS报警"` |
| `vitals.battery` | Number | 剩余电量，范围 0~100 (%) |
| `vitals.heartRate` | Number | 心率，单位 bpm |
| `vitals.bloodOxygen`| Number | 血氧饱和度，范围 0~100 (%) |

---

## 3. 内部调试控制说明（HTTP 接口）

这部分接口主要用于**开发人员在本地或服务器上控制模拟器**。在与第三方联调时，你可以通过调用这些接口瞬间改变某个设备的体征，让第三方观察他们的系统是否能正确接收和报警。

- **基础 URL**: `http://<模拟器部署IP>:8080`

### 3.1 获取所有模拟设备列表

用于查看当前正在模拟的设备清单及其实时内部状态（这些数据也会定期通过 Webhook 推送）。

- **请求**: `GET /devices`
- **响应示例**:

```json
[
  {
    "deviceType": "VitalBand-X1",
    "deviceId": "DEV-VITAL-001",
    "alarmStatus": "无",
    "battery": 100,
    "heartRate": 75,
    "bloodOxygen": 98
  },
  {
    "deviceType": "VitalBand-X1",
    "deviceId": "DEV-VITAL-002",
    "alarmStatus": "无",
    "battery": 100,
    "heartRate": 75,
    "bloodOxygen": 98
  }
]
```
