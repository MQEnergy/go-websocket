package server

import (
	"bytes"
	"github.com/MQEnergy/go-websocket/client"
	"github.com/MQEnergy/go-websocket/utils/code"
	"github.com/MQEnergy/go-websocket/utils/log"
	"github.com/MQEnergy/go-websocket/utils/response"
	"github.com/gorilla/websocket"
	"time"
)

type (
	// clientInfo 客户端消息体
	clientInfo struct {
		SystemId  string      `json:"system_id"`
		ClientId  string      `json:"client_id"`
		MessageId string      `json:"message_id"`
		Code      code.Code   `json:"code"`
		Msg       string      `json:"msg"`
		Data      interface{} `json:"data"`
	}

	// Sender 发送者结构体
	Sender struct {
		SystemId  string      `json:"system_id"`
		ClientId  string      `json:"client_id"`
		MessageId string      `json:"message_id"`
		GroupName string      `json:"group_name"`
		Code      code.Code   `json:"code"`
		Msg       string      `json:"msg"`
		Data      interface{} `json:"data"`
	}
)

var (
	ToClientMsgChan chan clientInfo       // 客户端消息channel通道
	Manager         = client.NewManager() // 管理者
	newline         = []byte{'\n'}
	space           = []byte{' '}
	writeWait       = 10 * time.Second    // Time allowed to write a message to the peer.
	pongWait        = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod      = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize  = 8192                // 最大的消息大小

)

func init() {
	ToClientMsgChan = make(chan clientInfo, 10000)
}

// Run 执行客户端连接处理
func Run() {
	for {
		select {
		// 客户端连接处理
		case client := <-Manager.ClientConnect:
			Manager.ClientConnectHandler(client)
		// 客户端断连处理
		case disClient := <-Manager.ClientDisConnect:
			Manager.ClientDisConnectHandler(disClient)
		}
	}
}

// SendMessageToClient 发送消息给客户端
func SendMessageToClient(sender *Sender) {
	ToClientMsgChan <- clientInfo{SystemId: sender.SystemId, ClientId: sender.ClientId, MessageId: sender.MessageId, Code: sender.Code, Msg: sender.Msg, Data: sender.Data}
}

// WriteMessageHandler 监听并发送给客户端消息
func WriteMessageHandler(c *client.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		Manager.ClientDisConnectHandler(c)
	}()
	for {
		select {
		// 接受消息体
		case clientInfo, ok := <-ToClientMsgChan:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				return
			}
			if c.SystemId == clientInfo.SystemId && c.ClientId == clientInfo.ClientId {
				client, err := Manager.GetClientByList(clientInfo.ClientId)
				if err != nil {
					return
				}
				params := map[string]string{
					"system_id":  client.SystemId,
					"client_id":  clientInfo.ClientId,
					"message_id": clientInfo.MessageId,
				}
				// 给客户端发消息
				if err := response.WsJson(client.Conn, clientInfo.Code, clientInfo.Msg, clientInfo.Data, params); err != nil {
					return
				}
				log.TraceLog(clientInfo.Code, params, clientInfo.Data, nil, 4)
			}

		// 定时心跳监测
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.TraceHeartbeatErrdLog(map[string]string{
					"system_id": c.SystemId,
					"client_id": c.ClientId,
				}, nil, err.Error(), 3)
				return
			}
		default:
		}
	}
}

// ReadMessageHandler 读取消息
func ReadMessageHandler(c *client.Client, fn func(client *client.Client, msg []byte) error) {
	c.Conn.SetReadLimit(int64(maxMessageSize))
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(appData string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	defer func() {
		Manager.ClientDisConnectHandler(c)
	}()
	for {
		//接收消息
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			if messageType != -1 && websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.TraceSendMsgErrLog("", "", err.Error(), 2)
			}
			break
		}
		// 回调函数
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if err := fn(c, message); err != nil {
			//Manager.ClientDisConnect <- c
			break
		}
	}
}

// SendMessageToLocalGroup 以群组纬度统一发送消息
func SendMessageToLocalGroup(sender *Sender) {
	if sender.GroupName != "" {
		groupKey := sender.SystemId + ":" + sender.GroupName
		clientIds := Manager.GetGroupClientList(groupKey)
		if len(clientIds) > 0 {
			for _, clientId := range clientIds {
				if _, err := Manager.GetClientByList(clientId); err != nil {
					// 不存在就删除群组中的客户端
					Manager.RemoveGroupClient(groupKey, clientId)
				} else {
					// 通过本服务器发送信息
					sender.ClientId = clientId
					SendMessageToClient(sender)
				}
			}
		}
	}
}

// SendMessageToLocalSystem 以系统纬度统一发送消息
func SendMessageToLocalSystem(sender *Sender) {
	if sender.SystemId != "" {
		clientIds := Manager.GetSystemClientList(sender.SystemId)
		if len(clientIds) > 0 {
			for _, clientId := range clientIds {
				if _, err := Manager.GetClientByList(clientId); err != nil {
					// 不存在就删除系统的客户端
					Manager.RemoveSystemClientByList(&client.Client{ClientId: clientId, SystemId: sender.SystemId})
				} else {
					// 通过本服务器发送信息
					sender.ClientId = clientId
					SendMessageToClient(sender)
				}
			}
		}
	}
}
