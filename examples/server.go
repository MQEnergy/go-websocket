package main

import (
	"encoding/json"
	"fmt"
	"github.com/MQEnergy/go-websocket/connect"
	"github.com/MQEnergy/go-websocket/server"
	"github.com/MQEnergy/go-websocket/utils/code"
	"net/http"
)

func main() {
	// 启动websocket
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		conn := connect.NewConn(writer, request, nil)
		conn.OnHandshake()
		// 监听消息推送
		conn.OnPush()
	})

	// 推送消息到单个客户端
	http.HandleFunc("/push_to_client", func(writer http.ResponseWriter, request *http.Request) {
		clientId := request.FormValue("client_id")
		systemId := request.FormValue("system_id")
		data := request.FormValue("data")
		server.SendMessageToClient(&server.Sender{
			ClientId: clientId,
			SystemId: systemId,
			Code:     code.SendMsgSuccess,
			Msg:      code.SendMsgSuccess.Msg(),
			Data:     &data,
		})
		msg, _ := json.Marshal(server.Sender{
			ClientId: clientId,
			SystemId: systemId,
			Code:     code.SendMsgSuccess,
			Msg:      code.SendMsgSuccess.Msg(),
			Data:     &data,
		})
		writer.Write(msg)

	})
	// 推送消息到多个客户端
	http.HandleFunc("/push_to_clients", func(writer http.ResponseWriter, request *http.Request) {

	})
	// 绑定到群组
	http.HandleFunc("/bind_to_group", func(writer http.ResponseWriter, request *http.Request) {

	})
	// 推送消息到群组
	http.HandleFunc("/push_to_group", func(writer http.ResponseWriter, request *http.Request) {

	})
	// 关闭连接
	http.HandleFunc("/close", func(writer http.ResponseWriter, request *http.Request) {
		conn := connect.NewConn(writer, request, nil)
		conn.OnClose()
	})

	fmt.Println("服务器启动成功，端口号 :9991 \n")
	if err := http.ListenAndServe(":9991", nil); err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
