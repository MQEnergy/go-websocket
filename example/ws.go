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
		go_websocket.WsServer(hub, writer, request)
	})

	// 推送到系统
	http.HandleFunc("/push_to_system", func(writer http.ResponseWriter, request *http.Request) {
		systemId := request.FormValue("system_id")
		data := request.FormValue("data")
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		if systemId == "" || data == "" {
			writer.Write([]byte("{\"msg\":\"参数错误\"}"))
			return
		}
		hub.Broadcast <- []byte(data)
		return
	})
	log.Println("服务启动成功。端口号 :9991")
	if err := http.ListenAndServe(":9991", nil); err != nil {
		log.Println("ListenAndServe: ", err)
	}
}
