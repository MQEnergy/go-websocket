package response

import (
	"github.com/MQEnergy/go-websocket/utils/code"
	"github.com/gorilla/websocket"
)

// responseData 响应结构体
type responseData struct {
	Code   code.Code   `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Params interface{} `json:"params"` // 自定义参数
}

// WsJson 返回给客户端的信息
func WsJson(conn *websocket.Conn, code code.Code, message string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code,
		Msg:    message,
		Params: params,
		Data:   data,
	})
}

// WsSuccessJson 返回客户端连接成功
func WsSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.Success,
		Msg:    code.Success.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsFailedJson 返回客户端连接失败
func WsFailedJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.Failed,
		Msg:    code.Failed.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsClientFailedJson 返回客户端主动断连
func WsClientFailedJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.ClientFailed,
		Msg:    code.ClientFailed.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsClientNotExistJson 返回客户端不存在
func WsClientNotExistJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.ClientNotExist,
		Msg:    code.ClientNotExist.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsClientCloseSuccessJson 返回客户端关闭成功
func WsClientCloseSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.ClientCloseSuccess,
		Msg:    code.ClientCloseSuccess.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsClientCloseFailedJson 返回客户端关闭失败
func WsClientCloseFailedJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.ClientCloseFailed,
		Msg:    code.ClientCloseFailed.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsReadMsgErrJson 返回读取消息体失败
func WsReadMsgErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.ReadMsgErr,
		Msg:    code.ReadMsgErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsReadMsgSuccessJson 返回读取消息体成功
func WsReadMsgSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.ReadMsgSuccess,
		Msg:    code.ReadMsgSuccess.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsSendMsgErrJson 返回发送消息体失败
func WsSendMsgErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.SendMsgErr,
		Msg:    code.SendMsgErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsSendMsgSuccessJson 返回发送消息体成功
func WsSendMsgSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.SendMsgSuccess,
		Msg:    code.SendMsgSuccess.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsHeartbeatErrJson 返回心跳检测失败
func WsHeartbeatErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.HeartbeatErr,
		Msg:    code.HeartbeatErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsSystemErrJson 返回系统不能为空
func WsSystemErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.SystemErr,
		Msg:    code.SystemErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsBindGroupSuccessJson 返回绑定群组成功
func WsBindGroupSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.BindGroupSuccess,
		Msg:    code.BindGroupSuccess.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsBindGroupErrJson 返回绑定群组失败
func WsBindGroupErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.BindGroupErr,
		Msg:    code.BindGroupErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsUnAuthedJson 返回用户未认证
func WsUnAuthedJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.UnAuthed,
		Msg:    code.UnAuthed.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsInternalErrJson 返回服务器内部错误
func WsInternalErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.InternalErr,
		Msg:    code.InternalErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsRequestMethodErrJson 返回请求方式错误
func WsRequestMethodErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.RequestMethodErr,
		Msg:    code.RequestMethodErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WsRequestParamErrJson 返回请求参数错误
func WsRequestParamErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code.RequestParamErr,
		Msg:    code.RequestParamErr.Msg(),
		Params: params,
		Data:   data,
	})
}
