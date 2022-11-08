package main

import (
	"github.com/MQEnergy/go-websocket"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

var (
	Node *snowflake.Node
)

func init() {
	// 日志注入
	go_websocket.Logger = logrus.New()
	go_websocket.Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	localIp, err := go_websocket.GetLocalIpToInt()
	if err != nil {
		panic(err)
	}
	Node, err = snowflake.NewNode(int64(localIp) % 1023)
	if err != nil {
		panic(err)
	}
}

func main() {
	hub := go_websocket.NewHub()
	go hub.Run()

	// ws连接
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		_, err := go_websocket.WsServer(hub, writer, request, go_websocket.Binary)
		if err != nil {
			return
		}
	})

	// 推送到所有连接的客户端
	http.HandleFunc("/push_to_all", func(writer http.ResponseWriter, request *http.Request) {
		data := request.FormValue("data")
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		if data == "" {
			writer.Write([]byte("{\"msg\":\"参数错误\"}"))
			return
		}
		hub.Broadcast <- []byte(data)
		writer.Write([]byte("{\"msg\":\"全局消息发送成功\"}"))
		return
	})

	// 推送到所在系统的客户端
	http.HandleFunc("/push_to_system", func(writer http.ResponseWriter, request *http.Request) {
		systemId := request.FormValue("system_id")
		data := request.FormValue("data")
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		if systemId == "" || data == "" {
			writer.Write([]byte("{\"msg\":\"参数错误\"}"))
			return
		}
		hub.SystemBroadcast <- &go_websocket.BroadcastChan{
			Name: systemId,
			Msg:  []byte(data),
		}
		writer.Write([]byte("{\"msg\":\"系统消息发送成功\"}"))
		return
	})

	// 推送到群组
	http.HandleFunc("/push_to_group", func(writer http.ResponseWriter, request *http.Request) {
		groupId := request.FormValue("group_id")
		data := request.FormValue("data")
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		if groupId == "" || data == "" {
			writer.Write([]byte("{\"msg\":\"参数错误\"}"))
			return
		}
		hub.GroupBroadcast <- &go_websocket.BroadcastChan{
			Name: groupId,
			Msg:  []byte(data),
		}
		//hub.GroupBroadcastHandle(groupId, []byte(data))
		writer.Write([]byte("{\"msg\":\"群组消息发送成功\"}"))
		return
	})

	// 推送到单个客户端
	http.HandleFunc("/push_to_client", func(writer http.ResponseWriter, request *http.Request) {
		clientId := request.FormValue("client_id")
		data := request.FormValue("data")
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		if clientId == "" || data == "" {
			writer.Write([]byte("{\"msg\":\"参数错误\"}"))
			return
		}
		hub.ClientBroadcast <- &go_websocket.BroadcastChan{
			Name: clientId,
			Msg:  []byte(data),
		}
		writer.Write([]byte("{\"msg\":\"客户端消息发送成功\"}"))
		return
	})

	log.Println("服务启动成功。端口号 :9991")
	if err := http.ListenAndServe(":9991", nil); err != nil {
		log.Println("ListenAndServe: ", err)
	}
}
