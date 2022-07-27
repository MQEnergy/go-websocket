package response

import (
	"github.com/MQEnergy/go-websocket/utils/code"
	"github.com/MQEnergy/go-websocket/utils/log"
	"github.com/gorilla/websocket"
)

// responseData 响应结构体
type responseData struct {
	SystemId  string      `json:"system_id"`
	ClientId  string      `json:"client_id"`
	MessageId string      `json:"message_id"`
	Code      code.Code   `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	Params    interface{} `json:"params"` // 自定义参数
}

// WsJson 返回给客户端的信息
func WsJson(conn *websocket.Conn, systemId, clientId, messageId string, code code.Code, message string, data, params interface{}) error {
	log.WriteLog(systemId, clientId, messageId, data, code, message, 4)
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code,
		Msg:       message,
		Data:      data,
		Params:    params,
	})
}
