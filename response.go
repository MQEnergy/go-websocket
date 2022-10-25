package go_websocket

import (
	"github.com/gorilla/websocket"
)

// responseData 响应结构体
type responseData struct {
	Code   Code        `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Params interface{} `json:"params"` // 自定义参数
}

// WriteJson 返回给客户端的信息
func WriteJson(conn *websocket.Conn, code Code, message string, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   code,
		Msg:    message,
		Params: params,
		Data:   data,
	})
}

// WsSuccessJson 返回客户端连接成功
func WriteSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   Success,
		Msg:    Success.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteFailedJson 返回客户端连接失败
func WriteFailedJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   Failed,
		Msg:    Failed.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteClientFailedJson 返回客户端主动断连
func WriteClientFailedJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   ClientFailed,
		Msg:    ClientFailed.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteClientNotExistJson 返回客户端不存在
func WriteClientNotExistJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   ClientNotExist,
		Msg:    ClientNotExist.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteClientCloseSuccessJson 返回客户端关闭成功
func WriteClientCloseSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   ClientCloseSuccess,
		Msg:    ClientCloseSuccess.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteClientCloseFailedJson 返回客户端关闭失败
func WriteClientCloseFailedJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   ClientCloseFailed,
		Msg:    ClientCloseFailed.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteReadMsgErrJson 返回读取消息体失败
func WriteReadMsgErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   ReadMsgErr,
		Msg:    ReadMsgErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteReadMsgSuccessJson 返回读取消息体成功
func WriteReadMsgSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   ReadMsgSuccess,
		Msg:    ReadMsgSuccess.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteSendMsgErrJson 返回发送消息体失败
func WriteSendMsgErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   SendMsgErr,
		Msg:    SendMsgErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteSendMsgSuccessJson 返回发送消息体成功
func WriteSendMsgSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   SendMsgSuccess,
		Msg:    SendMsgSuccess.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteHeartbeatErrJson 返回心跳检测失败
func WriteHeartbeatErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   HeartbeatErr,
		Msg:    HeartbeatErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteSystemErrJson 返回系统不能为空
func WriteSystemErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   SystemErr,
		Msg:    SystemErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteBindGroupSuccessJson 返回绑定群组成功
func WriteBindGroupSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   BindGroupSuccess,
		Msg:    BindGroupSuccess.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteBindGroupErrJson 返回绑定群组失败
func WriteBindGroupErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   BindGroupErr,
		Msg:    BindGroupErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteUnAuthedJson 返回用户未认证
func WriteUnAuthedJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   UnAuthed,
		Msg:    UnAuthed.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteInternalErrJson 返回服务器内部错误
func WriteInternalErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   InternalErr,
		Msg:    InternalErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteRequestMethodErrJson 返回请求方式错误
func WriteRequestMethodErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   RequestMethodErr,
		Msg:    RequestMethodErr.Msg(),
		Params: params,
		Data:   data,
	})
}

// WriteRequestParamErrJson 返回请求参数错误
func WriteRequestParamErrJson(conn *websocket.Conn, data, params interface{}) error {
	return conn.WriteJSON(responseData{
		Code:   RequestParamErr,
		Msg:    RequestParamErr.Msg(),
		Params: params,
		Data:   data,
	})
}
