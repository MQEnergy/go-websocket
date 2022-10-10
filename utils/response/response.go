package response

import (
	"github.com/MQEnergy/go-websocket/utils/code"
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

// WsSuccessJson 返回客户端连接成功
func WsSuccessJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.Success,
		Msg:       code.Success.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsFailedJson 返回客户端连接失败
func WsFailedJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.Failed,
		Msg:       code.Failed.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsClientFailedJson 返回客户端主动断连
func WsClientFailedJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.ClientFailed,
		Msg:       code.ClientFailed.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsClientNotExistJson 返回客户端不存在
func WsClientNotExistJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.ClientNotExist,
		Msg:       code.ClientNotExist.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsClientCloseSuccessJson 返回客户端关闭成功
func WsClientCloseSuccessJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.ClientCloseSuccess,
		Msg:       code.ClientCloseSuccess.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsClientCloseFailedJson 返回客户端关闭失败
func WsClientCloseFailedJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.ClientCloseFailed,
		Msg:       code.ClientCloseFailed.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsReadMsgErrJson 返回读取消息体失败
func WsReadMsgErrJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.ReadMsgErr,
		Msg:       code.ReadMsgErr.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsReadMsgSuccessJson 返回读取消息体成功
func WsReadMsgSuccessJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.ReadMsgSuccess,
		Msg:       code.ReadMsgSuccess.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsSendMsgErrJson 返回发送消息体失败
func WsSendMsgErrJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.SendMsgErr,
		Msg:       code.SendMsgErr.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsSendMsgSuccessJson 返回发送消息体成功
func WsSendMsgSuccessJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.SendMsgSuccess,
		Msg:       code.SendMsgSuccess.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsHeartbeatErrJson 返回心跳检测失败
func WsHeartbeatErrJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.HeartbeatErr,
		Msg:       code.HeartbeatErr.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsSystemErrJson 返回系统不能为空
func WsSystemErrJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.SystemErr,
		Msg:       code.SystemErr.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsBindGroupSuccessJson 返回绑定群组成功
func WsBindGroupSuccessJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.BindGroupSuccess,
		Msg:       code.BindGroupSuccess.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsBindGroupErrJson 返回绑定群组失败
func WsBindGroupErrJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.BindGroupErr,
		Msg:       code.BindGroupErr.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsUnAuthedJson 返回用户未认证
func WsUnAuthedJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.UnAuthed,
		Msg:       code.UnAuthed.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsInternalErrJson 返回服务器内部错误
func WsInternalErrJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.InternalErr,
		Msg:       code.InternalErr.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsRequestMethodErrJson 返回请求方式错误
func WsRequestMethodErrJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.RequestMethodErr,
		Msg:       code.RequestMethodErr.Msg(),
		Data:      data,
		Params:    params,
	})
}

// WsRequestParamErrJson 返回请求参数错误
func WsRequestParamErrJson(conn *websocket.Conn, systemId, clientId, messageId string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		SystemId:  systemId,
		ClientId:  clientId,
		MessageId: messageId,
		Code:      code.RequestParamErr,
		Msg:       code.RequestParamErr.Msg(),
		Data:      data,
		Params:    params,
	})
}
