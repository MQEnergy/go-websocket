package gowebsocket

import (
	"github.com/MQEnergy/go-websocket/client"
	"github.com/MQEnergy/go-websocket/server"
	"github.com/MQEnergy/go-websocket/utils/ip"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

const (
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
	Node    *snowflake.Node
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func init() {
	localIp, err := ip.GetLocalIpToInt()
	if err != nil {
		panic(err)
	}
	Node, err = snowflake.NewNode(int64(localIp) % 1023)
	if err != nil {
		panic(err)
	}
}

type Connect struct {
	writer   http.ResponseWriter
	request  *http.Request
	header   http.Header
	_client  *client.Client  // 客户端
	_manager *client.Manager // 连接管理器
}

type ConnInterface interface {
	OnHandshake(fn func(client *client.Client) error) error           // OnHandshake 握手
	OnOpen() error                                                    // OnOpen 连接
	OnMessage(fn func(client *client.Client, msg []byte) error) error // OnMessage 接收消息
	OnClose(fn func(client *client.Client) error) error               // OnClose 关闭连接
}

// New 实例化
func New(w http.ResponseWriter, r *http.Request, header http.Header, c *client.Client) ConnInterface {
	return &Connect{
		writer:   w,
		request:  r,
		header:   header,
		_client:  c,
		_manager: server.Manager,
	}
}

// OnHandshake 握手
func (c *Connect) OnHandshake(fn func(client *client.Client) error) error {
	conn, err := upgrader.Upgrade(c.writer, c.request, c.header)
	if err != nil {
		return err
	}
	systemId := c.request.FormValue("system_id")
	if systemId == "" {
		sid, err := ip.GetLocalIpToInt()
		if err != nil {
			return err
		}
		systemId = strconv.Itoa(sid)
	}
	clientId := client.GenerateUuid(0, Node)
	c._client = client.NewClient(clientId, systemId, conn)

	if err := c.OnOpen(); err != nil {
		return err
	}
	if err := fn(c._client); err != nil {
		if err := c._manager.ClientDisConnectHandler(c._client); err != nil {
			return err
		}
		return err
	}
	return nil
}

// OnOpen 开启websocket
func (c *Connect) OnOpen() error {
	c._manager.SetSystemClientToList(c._client)
	if err := c._manager.ClientConnectHandler(c._client); err != nil {
		return err
	}
	//c._manager.ClientConnect <- c._client
	return nil
}

// OnMessage 接收消息
func (c *Connect) OnMessage(fn func(client *client.Client, msg []byte) error) error {
	// 读取消息协程处理
	go server.ReadMessageHandler(c._client, fn)
	// 写入消息协程处理
	go server.WriteMessageHandler(c._client)
	return nil
}

// OnClose 关闭连接
func (c *Connect) OnClose(fn func(c *client.Client) error) error {
	if err := c._manager.ClientDisConnectHandler(c._client); err != nil {
		return err
	}
	// 回调函数
	if err := fn(c._client); err != nil {
		return err
	}
	return nil
}
