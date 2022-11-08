# go-websocket
基于gorilla/websocket封装的websocket库，实现基于系统纬度的消息推送，基于群组纬度的消息推送，基于单个和多个客户端消息推送。

[![GoDoc](https://godoc.org/github.com/MQEnergy/go-websocket/?status.svg)](https://pkg.go.dev/github.com/MQEnergy/go-websocket)
[![Go Report Card](https://goreportcard.com/badge/github.com/MQEnergy/go-websocket)](https://goreportcard.com/report/github.com/MQEnergy/go-websocket)
[![codebeat badge](https://codebeat.co/badges/063ec0b6-5059-4b1b-92c0-4f750438faa8)](https://codebeat.co/projects/github-com-mqenergy-go-websocket-main)
[![GitHub license](https://img.shields.io/github/license/MQEnergy/go-websocket)](https://github.com/MQEnergy/go-websocket/blob/main/LICENSE)

## 一、目录结构
```
├── LICENSE
├── README.md
├── client.go           // 客户端
├── client_hub.go       // 客户端集线器
├── code.go             // 状态码
├── example             // 案例
│   └── ws.go
├── go.mod
├── go.sum
├── log.go              // 日志
├── node.go             // 节点（用于在分布式系统生成基于节点的客户端连接ID）
├── response.go         // 客户端发送消息
└── server.go           // 服务

```
## 二、在项目中安装使用
```go
go get -u github.com/MQEnergy/go-websocket
```
## 三、运行example
### 1、开启服务
```go
go run examples/ws.go
```
```
服务器启动成功，端口号 :9991 
```
代表启动成功

### 2、案例
具体查看example目录

#### 1）连接ws并加群组
system_id为系统ID（不必填 不填默认当前节点ip的int值）
group_id为群组ID（不必填 不填连接不加群组 注意：群组id为全局唯一ID 不然可能会出现不同系统的相同群组都推送消息）

请求
```
ws://127.0.0.1:9991/ws?system_id=123&group_id=test
```
可选多种返回方式 如： Text，Json，Binary（二进制方式）
返回如下json示例：
```
{
    "code": 0,
    "msg": "客户端连接成功",
    "data": {
        "client_id": "1589962851152388096",
        "group_id": "test",
        "system_id": "123"
    },
    "params": null
}
```

#### 2）全局广播消息群发
请求
```
http://127.0.0.1:9991/push_to_system?system_id=123&data={"hello":"world"}
```
返回
```
{
    "msg": "系统消息发送成功",
}
```

#### 3）单个系统消息群发
请求
```
http://127.0.0.1:9991/push_to_system?system_id=123&data={"hello":"world"}
```
返回
```
{
    "msg": "系统消息发送成功",
}
```

#### 4）推送消息到群组
请求
```
http://127.0.0.1:9991/push_to_group?system_id=123&group_id=test&data={"hello":"world1"}
```
返回
```
{
    "msg": "群组消息发送成功",
}
```

#### 5）单个客户端消息发送
请求
```
http://127.0.0.1:9991/push_to_client?client_id=123&data={"hello":"world"}
```
返回
```
{
    "msg": "客户端消息发送成功",
}
```