package server

import (
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
		ClientId  string
		MessageId string
		Code      code.Code
		Msg       string
		Data      interface{}
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
	heartbeatTimer  = 20 * time.Second    // 20秒心跳
	Manager         = client.NewManager() // 管理者
)

func init() {
	ToClientMsgChan = make(chan clientInfo, 1000)
}

// Run 执行客户端连接处理
func Run() {
	for {
		select {
		// 客户端连接处理
		case client := <-Manager.ClientConnect:
			err := Manager.ClientConnectHandler(client)
			if err != nil {
				log.WriteLog(client.SystemId, client.ClientId, "", client, code.ClientFailed, code.ClientFailed.Msg(), 4)
			} else {
				log.WriteLog(client.SystemId, client.ClientId, "", client, code.Success, code.Success.Msg(), 4)
			}
		// 客户端断连处理
		case disClient := <-Manager.ClientDisConnect:
			err := Manager.ClientDisConnectHandler(disClient)
			if err != nil {
				log.WriteLog(disClient.SystemId, disClient.ClientId, "", disClient, code.ClientCloseFailed, code.ClientCloseFailed.Msg(), 4)
			} else {
				log.WriteLog(disClient.SystemId, disClient.ClientId, "", disClient, code.ClientCloseSuccess, code.ClientCloseSuccess.Msg(), 4)
			}
		}
	}
}

// PushToClientMsgChan 发送消息体到通道
func PushToClientMsgChan(clientId, messageId string, code code.Code, msg string, data interface{}) {
	ToClientMsgChan <- clientInfo{ClientId: clientId, MessageId: messageId, Code: code, Msg: msg, Data: data}
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

// SendMessageToClient 发送消息给客户端
func SendMessageToClient(sender *Sender) {
	PushToClientMsgChan(sender.ClientId, sender.MessageId, sender.Code, sender.Msg, sender.Data)
}

// MessagePushListener 监听并发送给客户端消息
func MessagePushListener() {
	for {
		clientInfo := <-ToClientMsgChan
		if client, err := Manager.GetClientByList(clientInfo.ClientId); err == nil && client != nil {
			if err := response.WsJson(client.Conn, client.SystemId, client.ClientId, clientInfo.MessageId, clientInfo.Code, clientInfo.Msg, clientInfo.Data, nil); err != nil {
				log.WriteLog(client.SystemId, client.ClientId, clientInfo.MessageId, clientInfo.Data, code.ClientNotExist, "客户端异常离线 "+err.Error(), 4)
				Manager.ClientDisConnect <- client
			} else {
				log.WriteLog(client.SystemId, client.ClientId, clientInfo.MessageId, clientInfo.Data, code.SendMsgSuccess, code.SendMsgSuccess.Msg(), 4)
			}
		} else {
			log.WriteLog("", clientInfo.ClientId, clientInfo.MessageId, clientInfo.Data, code.ClientNotExist, code.ClientNotExist.Msg(), 4)
		}
	}
}

// HeartbeatListener 心跳监听
func HeartbeatListener() {
	ticker := time.NewTicker(heartbeatTimer)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			allClient := Manager.GetAllClient()
			for clientId, client := range allClient {
				if err := client.Conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
					Manager.ClientDisConnect <- client
					log.WriteLog(client.SystemId, clientId, "", map[string]interface{}{"clientCount": Manager.GetAllClientCount()}, code.HeartbeatErr, "心跳检测失败 "+err.Error(), 4)
				}
			}
		}
	}
}
