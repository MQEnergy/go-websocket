package response

import (
	"github.com/gorilla/websocket"
	"gowebsocket/utils/code"
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
