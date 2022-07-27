package main

import (
	"encoding/json"
	"fmt"
	"github.com/MQEnergy/go-websocket/client"
	"github.com/MQEnergy/go-websocket/connect"
	"github.com/MQEnergy/go-websocket/server"
	"github.com/MQEnergy/go-websocket/utils/code"
	"github.com/bwmarrin/snowflake"
	"net/http"
)

func main() {
	// 监听消息发送
	go server.MessagePushListener()

	// 启动websocket
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		connect.NewConn(writer, request, writer.Header()).OnHandshake()
		return
	})

	// 推送消息到单个客户端
	http.HandleFunc("/push_to_client", func(writer http.ResponseWriter, request *http.Request) {
		clientId := request.FormValue("client_id")
		systemId := request.FormValue("system_id")
		data := request.FormValue("data")
		messageId := client.GenerateUuid(32, nil)
		sender := &server.Sender{
			ClientId:  clientId,
			SystemId:  systemId,
			MessageId: messageId,
			Code:      code.SendMsgSuccess,
			Msg:       code.SendMsgSuccess.Msg(),
			Data:      &data,
		}
		server.SendMessageToClient(sender)
		msg, _ := json.Marshal(sender)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.Write(msg)
		return
	})
	// 推送消息到多个客户端
	http.HandleFunc("/push_to_clients", func(writer http.ResponseWriter, request *http.Request) {
		clientIds := request.FormValue("client_ids")
		systemId := request.FormValue("system_id")
		data := request.FormValue("data")
		var clientIdList []string
		if err := json.Unmarshal([]byte(clientIds), &clientIdList); err != nil {
			writer.Header().Set("Content-Type", "application/json; charset=utf-8")
			writer.Write([]byte(err.Error()))
			return
		}
		var senderList = make([]*server.Sender, 0)
		node, _ := snowflake.NewNode(1)
		for _, clientId := range clientIdList {
			messageId := client.GenerateUuid(32, node)
			sender := &server.Sender{
				ClientId:  clientId,
				SystemId:  systemId,
				MessageId: messageId,
				Code:      code.SendMsgSuccess,
				Msg:       code.SendMsgSuccess.Msg(),
				Data:      &data,
			}
			server.SendMessageToClient(sender)
			senderList = append(senderList, sender)
		}
		msg, _ := json.Marshal(senderList)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.Write(msg)
		return
	})
	// 绑定到群组
	http.HandleFunc("/bind_to_group", func(writer http.ResponseWriter, request *http.Request) {

	})
	// 推送消息到群组
	http.HandleFunc("/push_to_group", func(writer http.ResponseWriter, request *http.Request) {

	})
	// 关闭连接
	http.HandleFunc("/close", func(writer http.ResponseWriter, request *http.Request) {
		connect.NewConn(writer, request, writer.Header()).OnClose()
		return
	})

	fmt.Println("服务器启动成功，端口号 :9991 \n")
	if err := http.ListenAndServe(":9991", nil); err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
