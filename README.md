# go-websocket
基于gorilla/websocket封装的websocket库，实现基于系统纬度的消息推送，基于群组纬度的消息推送，基于单个和多个客户端消息推送。

## 目录结构
```
├── LICENSE
├── README.md
├── client
│   ├── client.go             // 客户端
│   └── client_manager.go     // 客户端管理者
├── connect
│   └── connect.go            // websocket连接
├── examples
│   └── server.go             // 服务端案例
├── go.mod
├── go.sum
├── server
│   └── server.go             // 服务端
└── utils
    ├── code
    │   └── code.go           // 状态码
    ├── log
    │   └── log.go            // log
    └── response
        └── response.go       // websocket返回给客户端响应
```
## 在项目中安装使用
```go
go get -u github.com/MQEnergy/go-websocket
```
## 安装依赖
```go
go mod tidy
```
## 运行example
### 一、开启服务
```go
go run examples/server.go
```
```
服务器启动成功，端口号 :9991 
```
代表启动成功

### 二、测试
#### 1、连接websocket
system_id为系统ID（测试随机填写）

请求
```
ws://127.0.0.1:9991/ws?system_id=123
```
返回
```
{
    "system_id": "123",
    "client_id": "1552131603575083008",
    "message_id": "bmn18osnryryy",
    "code": 10009,
    "msg": "发送消息体成功",
    "data": "{\"hello\":\"world1\"}",
    "params": null
}
```

#### 2、推送单个客户端 
请求 post
```
http://127.0.0.1:9991/push_to_client
```
请求参数
```
client_id:1552131603575083008
system_id:123
data:{"hello":"world1"}
```
返回
```
{
    "system_id": "123",
    "client_id": "1552106203776028672",
    "message_id": "bmn1pgxsyyryy",
    "group_name": "",
    "code": 10009,
    "msg": "发送消息体成功",
    "data": "{\"hello\":\"world1\"}"
}
```

#### 3、推送多个客户端 
请求 post
```
http://127.0.0.1:9991/push_to_clients
```
请求参数
```
client_ids:["1552131603575083008","1552131575854927872"]
system_id:123
data:{"hello":"world1"}
```
返回
```
[
    {
        "system_id": "123",
        "client_id": "1552131603575083008",
        "message_id": "bmnwpcepyyryy",
        "group_name": "",
        "code": 10009,
        "msg": "发送消息体成功",
        "data": "{\"hello\":\"world1\"}"
    },
    {
        "system_id": "123",
        "client_id": "1552131575854927872",
        "message_id": "bmnwpcepyyryb",
        "group_name": "",
        "code": 10009,
        "msg": "发送消息体成功",
        "data": "{\"hello\":\"world1\"}"
    }
]
```

## 其他
### client_manager.go方法和函数
```
SetClientToList             添加客户端到列表
SetSystemClientToList       添加系统ID和客户端到列表
SetClientToGroupList        添加客户端到分组

GetAllClient                获取所有的客户端
GetAllClientCount           获取所有客户端数量
GetClientByList             通过客户端列表获取*Client
GetSystemClientList         获取指定系统的客户端列表
GetGroupClientList          获取本地分组的成员

RemoveAllClient             删除所有客户端 (包含 ClientList SystemClientList GroupList)
RemoveGroupClient           删除分组里的客户端
RemoveClientByList          从列表删除*Client
RemoveSystemClientByList    删除系统里的客户端
CloseClient                 关闭客户端 (发送关闭client到ClientDisConnect通道中)
```
### server.go函数
```
Run                         执行连接（断连，连接进行clientList的操作）
PushToClientMsgChan         发送消息体到通道
SendMessageToLocalGroup     以群组纬度统一发送消息
SendMessageToLocalSystem    以系统纬度统一发送消息
SendMessageToClient         发送消息给客户端
MessagePushListener         监听并发送给客户端消息
HeartbeatListener           心跳监听
```