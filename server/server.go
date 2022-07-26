package server

import (
	"encoding/json"
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
		Data      *string
	}

	// Sender 发送者结构体
	Sender struct {
		SystemId  string
		ClientId  string
		MessageId string
		GroupName string
		Code      code.Code
		Msg       string
		Data      *string
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

// Run 执行
func Run() {
	for {
		select {
		case client := <-Manager.ClientConnect:
			// 客户端连接处理
			Manager.ClientConnectHandler(client)

		case disClient := <-Manager.ClientDisConnect:
			// 客户端断连处理
			Manager.ClientDisConnectHandler(disClient)
			//// 客户端下线通知
			marshal, _ := json.Marshal(map[string]string{
				"clientId": disClient.ClientId,
			})
			data := string(marshal)
			if len(disClient.GroupList) > 0 {
				for _, groupName := range disClient.GroupList {
					//发送下线通知
					SendMessageToLocalGroup(&Sender{
						SystemId:  disClient.SystemId,
						ClientId:  disClient.ClientId,
						MessageId: client.GenerateUuid(32),
						GroupName: groupName,
						Code:      code.ClientFailed,
						Msg:       code.ClientFailed.Msg(),
						Data:      &data,
					})
				}
			}
		}
	}
}

// PushToClientMsgChan 发送消息体到通道
func PushToClientMsgChan(clientId, messageId string, code code.Code, msg string, data *string) {
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
	log.WriteLog(sender.SystemId, sender.ClientId, sender.Data, sender.Code, sender.Msg, 4)
	PushToClientMsgChan(sender.ClientId, sender.MessageId, sender.Code, sender.Msg, sender.Data)
}

// MessagePushListener 监听并发送给客户端消息
func MessagePushListener() {
	for {
		clientInfo := <-ToClientMsgChan
		if client, err := Manager.GetClientByList(clientInfo.ClientId); err == nil && client != nil {
			if err := response.WsJson(client.Conn, client.SystemId, client.ClientId, clientInfo.MessageId, clientInfo.Code, clientInfo.Msg, clientInfo.Data, nil); err != nil {
				Manager.ClientDisConnect <- client
				log.WriteLog(client.SystemId, client.ClientId, clientInfo, code.ClientNotExist, "客户端异常离线 "+err.Error(), 4)
			}
		} else {
			log.WriteLog("", clientInfo.ClientId, clientInfo, code.ClientNotExist, code.ClientNotExist.Msg(), 4)
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
			//发送心跳
			allClient := Manager.GetAllClient()
			for clientId, client := range allClient {
				if err := client.Conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
					Manager.ClientDisConnect <- client
					log.WriteLog(client.SystemId, client.ClientId, map[string]interface{}{"clientId": clientId, "clientCount": Manager.GetAllClientCount()}, code.HeartbeatErr, "心跳检测失败 "+err.Error(), 4)
				}
			}
		}
	}
}
