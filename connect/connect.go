package connect

import (
	"errors"
	"github.com/MQEnergy/go-websocket/client"
	"github.com/MQEnergy/go-websocket/server"
	"github.com/MQEnergy/go-websocket/utils/code"
	"github.com/MQEnergy/go-websocket/utils/log"
	"github.com/MQEnergy/go-websocket/utils/response"
	"github.com/gorilla/websocket"
	"net/http"
)

const (
	maxMessageSize  = 8192 // 最大的消息大小
	readBufferSize  = 1024 // 读缓冲区大小
	writeBufferSize = 1024 // 写缓冲区大小
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Connect struct {
	writer   http.ResponseWriter
	request  *http.Request
	header   http.Header
	_client  client.Client   // 客户端
	_manager *client.Manager // 连接管理器
}

type ConnInterface interface {
	OnHandshake() error // OnHandshake 握手
	OnOpen() error      // OnOpen 连接
	OnMessage() error   // OnMessage 接收消息
	OnPush() error      // OnPush 发送消息
	OnClose() error     // OnClose 关闭连接
}

// NewConn 实例化
func NewConn(w http.ResponseWriter, r *http.Request, header http.Header) ConnInterface {
	return &Connect{
		writer:   w,
		request:  r,
		header:   header,
		_manager: server.Manager,
	}
}

// OnHandshake 握手
func (c *Connect) OnHandshake() error {
	// http服务升级为websocket协议
	conn, err := upgrader.Upgrade(c.writer, c.request, c.header)
	if err != nil {
		log.WriteLog(c._client.SystemId, c._client.ClientId, map[string]string{"err": err.Error()}, code.ReadMsgErr, code.ClientFailed.Msg(), 4)
		return err
	}
	conn.SetReadLimit(maxMessageSize)

	c._client.Conn = conn

	//判断系统ID
	systemId := c.request.FormValue("system_id")
	if systemId == "" {
		// 给客户端发送信息
		response.WsJson(c._client.Conn, c._client.SystemId, c._client.ClientId, "", code.Failed, code.Failed.Msg(), []string{}, []string{})
		// 关闭连接
		c.OnClose()
		return errors.New(code.Failed.Msg())
	}
	// 生成客户端ID
	clientId := client.GenerateUuid(0)
	// 实例化新客户端连接
	c._client = *client.NewClient(clientId, systemId, conn)
	// 添加系统ID和客户端到列表
	c._manager.SetSystemClientToList(&c._client)
	// 打开websocket 给客户端发送消息
	c.OnOpen()
	// 心跳检测
	go server.HeartbeatListener()
	return nil
}

// OnOpen 开启websocket
func (c *Connect) OnOpen() error {
	// 开启协程读取信息
	c.OnMessage()
	// 客户端连接事件
	c._manager.ClientConnect <- &c._client

	// 监听客户端连接或断连
	go server.Run()

	return nil
}

// OnMessage 接收消息
func (c *Connect) OnMessage() error {
	go func() {
		for {
			//接收消息
			if messageType, message, err := c._client.Conn.ReadMessage(); err != nil {
				if messageType == -1 &&
					websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
					c._manager.ClientDisConnect <- &c._client
				} else {
					log.WriteLog(c._client.SystemId, c._client.ClientId, "", code.ReadMsgErr, code.ReadMsgErr.Msg()+" err: "+err.Error(), 4)
				}
				return
			} else {
				// 推送给客户端
				response.WsJson(c._client.Conn, c._client.SystemId, c._client.ClientId, "", code.ReadMsgSuccess, code.ReadMsgSuccess.Msg(), string(message), nil)
			}
		}
	}()
	return nil
}

// OnPush 监听消息推送
func (c *Connect) OnPush() error {
	go server.MessagePushListener()
	return nil
}

// OnClose 关闭连接
func (c *Connect) OnClose() error {
	if err := c._manager.CloseClient(c._client.ClientId, c._client.SystemId); err != nil {
		return err
	}
	return nil
}
