package main

import (
	"encoding/json"
	"fmt"
	"github.com/MQEnergy/go-websocket"
	"github.com/MQEnergy/go-websocket/client"
	"github.com/MQEnergy/go-websocket/server"
	"github.com/MQEnergy/go-websocket/utils/code"
	"github.com/MQEnergy/go-websocket/utils/log"
	"github.com/MQEnergy/go-websocket/utils/response"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func init() {
	// 日志注入
	log.Logger = logrus.New()
	log.Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
func main() {
	// 监听消息发送
	go server.MessagePushListener()

	// 监听客户端连接或断连
	go server.Run()

	// 心跳检测
	go server.HeartbeatListener()

	// 启动websocket
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		conn := gowebsocket.NewConn(writer, request, writer.Header(), &client.Client{})
		conn.OnHandshake()
		// 开启协程读取信息
		conn.OnMessage(func(c *client.Client, msg []byte) error {
			data := make(map[string]interface{}, 0)
			if err := json.Unmarshal(msg, &data); err != nil {
				response.WsJson(c.Conn, c.SystemId, c.ClientId, "", code.RequestParamErr, code.RequestParamErr.Msg(), nil, nil)
				return err
			}
			response.WsJson(c.Conn, c.SystemId, c.ClientId, "", code.ReadMsgSuccess, code.ReadMsgSuccess.Msg(), data, nil)
			log.WriteLog(c.SystemId, c.ClientId, "", string(msg), code.ReadMsgSuccess, code.ReadMsgSuccess.Msg(), 4)
			return nil
		})
		return
	})

	// 推送消息到单个客户端
	http.HandleFunc("/push_to_client", func(writer http.ResponseWriter, request *http.Request) {
		clientId := request.FormValue("client_id")
		systemId := request.FormValue("system_id")
		data := request.FormValue("data")

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")

		if systemId == "" || systemId == "" || data == "" {
			writer.Write([]byte("{\"msg\":\"参数错误\"}"))
			return
		}
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
		writer.Write(msg)
		return
	})
	// 推送消息到多个客户端
	http.HandleFunc("/push_to_clients", func(writer http.ResponseWriter, request *http.Request) {
		clientIds := request.FormValue("client_ids")
		systemId := request.FormValue("system_id")
		data := request.FormValue("data")
		var clientIdList []string

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")

		if systemId == "" || clientIds == "" || data == "" {
			writer.Write([]byte("{\"msg\":\"参数错误\"}"))
			return
		}
		if err := json.Unmarshal([]byte(clientIds), &clientIdList); err != nil {
			writer.Write([]byte("{\"msg\":\"" + err.Error() + "\"}"))
			return
		}
		senderList := make([]*server.Sender, 0)
		node, _ := snowflake.NewNode(1)
		for _, clientId := range clientIdList {
			messageId := client.GenerateUuid(32, node)
			sender := &server.Sender{
				ClientId:  clientId,
				SystemId:  systemId,
				MessageId: messageId,
				Code:      code.SendMsgSuccess,
				Msg:       code.SendMsgSuccess.Msg(),
				Data:      data,
			}
			server.SendMessageToClient(sender)
			senderList = append(senderList, sender)
		}
		msg, _ := json.Marshal(senderList)
		writer.Write(msg)
		return
	})
	// 绑定到群组
	http.HandleFunc("/bind_to_group", func(writer http.ResponseWriter, request *http.Request) {
		groupName := request.FormValue("group_name")
		clientId := request.FormValue("client_id")
		systemId := request.FormValue("system_id")

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")

		if systemId == "" || groupName == "" || clientId == "" {
			writer.Write([]byte("{\"msg\":\"参数错误\"}"))
			return
		}
		// 绑定操作
		if err := server.Manager.SetClientToGroupList(groupName, &client.Client{ClientId: clientId, SystemId: systemId}); err != nil {
			writer.Write([]byte("{\"msg\":\"" + err.Error() + "\"}"))
			return
		}

		// 发送信息到群组
		sender := &server.Sender{
			SystemId:  systemId,
			ClientId:  clientId,
			MessageId: client.GenerateUuid(32, nil),
			GroupName: groupName,
			Code:      code.BindGroupSuccess,
			Msg:       "客户端id：" + clientId + " " + code.BindGroupSuccess.Msg(),
			Data:      nil,
		}
		//发送系统通知
		server.SendMessageToLocalGroup(sender)

		// 返回
		msg, _ := json.Marshal(sender)
		writer.Write(msg)
		return
	})
	// 推送消息到群组
	http.HandleFunc("/push_to_group", func(writer http.ResponseWriter, request *http.Request) {
		systemId := request.FormValue("system_id")
		groupName := request.FormValue("group_name")
		data := request.FormValue("data")
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")

		if systemId == "" || groupName == "" || data == "" {
			writer.Write([]byte("{\"msg\":\"参数错误\"}"))
			return
		}
		groupKey := systemId + ":" + groupName
		groupClientList := server.Manager.GetGroupClientList(groupKey)
		if len(groupClientList) == 0 {
			writer.Write([]byte("{\"msg\":\"系统对应的群组不存在\"}"))
			return
		}
		// 发送信息到群组
		sender := &server.Sender{
			SystemId:  systemId,
			MessageId: client.GenerateUuid(32, nil),
			GroupName: groupName,
			Code:      code.SendMsgSuccess,
			Msg:       code.SendMsgSuccess.Msg(),
			Data:      data,
		}
		//发送系统通知
		server.SendMessageToLocalGroup(sender)

		sender.ClientId = strings.Join(groupClientList, ",")
		msg, _ := json.Marshal(sender)
		writer.Write(msg)
		return
	})
	// 关闭连接
	http.HandleFunc("/close", func(writer http.ResponseWriter, request *http.Request) {
		systemId := request.FormValue("system_id")
		clientId := request.FormValue("client_id")
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")

		if clientId == "" || systemId == "" {
			writer.Write([]byte("{\"msg\":\"参数错误\"}"))
			return
		}
		if err := gowebsocket.NewConn(writer, request, writer.Header(), &client.Client{ClientId: clientId, SystemId: systemId}).OnClose(); err != nil {
			writer.Write([]byte("{\"msg\":\"" + err.Error() + "\"}"))
			return
		}
		writer.Write([]byte("{\"msg\":\"客户端关闭成功\"}"))
		return
	})

	fmt.Println("服务器启动成功，端口号 :9991")
	if err := http.ListenAndServe(":9991", nil); err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}
