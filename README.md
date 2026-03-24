# MQTT Service Template

基于 Go 语言的 MQTT 服务模板，提供开箱即用的 MQTT 客户端封装。

## 项目结构

```
.
├── cmd/app/main.go      # 入口文件
├── internal/
│   ├── model/           # 数据模型
│   └── mqtt/            # MQTT 客户端封装
├── pkg/config/          # 配置管理
├── .env.example         # 环境变量示例
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

## 快速开始

### 1. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，配置 MQTT Broker 地址等信息
```

### 2. 本地运行

```bash
go run ./cmd/app
```

### 3. Docker 运行

```bash
docker-compose up -d --build
```

## MQTT Topic 规范

| Topic Pattern | 用途 |
|---------------|------|
| `devices/{deviceId}/telemetry` | 遥测数据上报 |
| `devices/{deviceId}/heartbeat` | 心跳 |
| `devices/{deviceId}/commands` | 接收命令 |

## 扩展开发

1. **添加命令处理**: 在 `App.ProcessCommand` 中添加新的命令类型处理
2. **自定义 Topic**: 扩展 `mqtt.Client` 添加新的发布/订阅方法
3. **数据模型**: 在 `internal/model/` 中定义业务数据结构

## 依赖

- [paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang) - MQTT 客户端库
- [godotenv](https://github.com/joho/godotenv) - 环境变量加载
