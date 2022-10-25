package go_websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

// responseData 响应结构体
type responseData struct {
	Code   Code        `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Params interface{} `json:"params"` // 自定义参数
}
type MsgType int

const (
	Text MsgType = 1
	Json         = iota + Text
	Binary
)

// WriteMessage 返回给客户端的信息
func WriteMessage(conn *websocket.Conn, code Code, message string, data, params interface{}, msgType MsgType) error {
	r := responseData{
		Code:   code,
		Msg:    message,
		Params: params,
		Data:   data,
	}
	switch msgType {
	case Text:
		marshal, _ := json.Marshal(r)
		return conn.WriteMessage(1, marshal)
	case Binary:
		marshal, _ := json.Marshal(r)
		return conn.WriteMessage(2, marshal)
	case Json:
		return conn.WriteJSON(r)
	}
	return nil
}

// WriteJson 返回给客户端的信息
func WriteJson(conn *websocket.Conn, code Code, message string, data, params interface{}) error {
	return WriteMessage(conn, code, message, data, params, Json)
}

// WriteSuccessJson 返回客户端连接成功
func WriteSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, Success, Success.Msg(), data, params)
}

// WriteFailedJson 返回客户端连接失败
func WriteFailedJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, Failed, Failed.Msg(), data, params)
}

// WriteClientFailedJson 返回客户端主动断连
func WriteClientFailedJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, ClientFailed, ClientFailed.Msg(), data, params)
}

// WriteClientNotExistJson 返回客户端不存在
func WriteClientNotExistJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, ClientNotExist, ClientNotExist.Msg(), data, params)
}

// WriteClientCloseSuccessJson 返回客户端关闭成功
func WriteClientCloseSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, ClientCloseSuccess, ClientCloseSuccess.Msg(), data, params)
}

// WriteClientCloseFailedJson 返回客户端关闭失败
func WriteClientCloseFailedJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, ClientCloseFailed, ClientCloseFailed.Msg(), data, params)
}

// WriteReadMsgErrJson 返回读取消息体失败
func WriteReadMsgErrJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, ReadMsgErr, ReadMsgErr.Msg(), data, params)
}

// WriteReadMsgSuccessJson 返回读取消息体成功
func WriteReadMsgSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, ReadMsgSuccess, ReadMsgSuccess.Msg(), data, params)
}

// WriteSendMsgErrJson 返回发送消息体失败
func WriteSendMsgErrJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, SendMsgErr, SendMsgErr.Msg(), data, params)
}

// WriteSendMsgSuccessJson 返回发送消息体成功
func WriteSendMsgSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, SendMsgSuccess, SendMsgSuccess.Msg(), data, params)
}

// WriteHeartbeatErrJson 返回心跳检测失败
func WriteHeartbeatErrJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, HeartbeatErr, HeartbeatErr.Msg(), data, params)
}

// WriteBindGroupSuccessJson 返回绑定群组成功
func WriteBindGroupSuccessJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, BindGroupSuccess, BindGroupSuccess.Msg(), data, params)
}

// WriteRequestParamErrJson 返回请求参数错误
func WriteRequestParamErrJson(conn *websocket.Conn, data, params interface{}) error {
	return WriteJson(conn, RequestParamErr, RequestParamErr.Msg(), data, params)
}
