package go_websocket

import (
	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
)

// TraceLog 写日志
func TraceLog(code Code, params, data, err interface{}, level logrus.Level) {
	Logger.WithFields(logrus.Fields{
		"code":   code,
		"err":    err,
		"params": params,
		"data":   data,
	}).Log(level, code.Msg())
}

// TraceHeartbeatErrdLog 心跳检测失败消息
func TraceHeartbeatErrdLog(params, data, err interface{}, level logrus.Level) {
	TraceLog(HeartbeatErr, params, data, err, level)
}

// TraceClientCloseFailedLog 客户端关闭失败消息
func TraceClientCloseFailedLog(params, data, err interface{}, level logrus.Level) {
	TraceLog(ClientCloseFailed, params, data, err, level)
}

// TraceClientCloseSuccessLog 客户端关闭成功消息
func TraceClientCloseSuccessLog(params, data, err interface{}, level logrus.Level) {
	TraceLog(ClientCloseSuccess, params, data, err, level)
}

// TraceSuccessLog 客户端连接成功消息
func TraceSuccessLog(params, data interface{}, level logrus.Level) {
	TraceLog(Success, params, data, nil, level)
}

// TraceReadMsgSuccessLog 读取消息体成功消息
func TraceReadMsgSuccessLog(params, data interface{}, level logrus.Level) {
	TraceLog(ReadMsgSuccess, params, data, nil, level)
}

// TraceSendMsgErrLog 发送消息体失败
func TraceSendMsgErrLog(params, data, err interface{}, level logrus.Level) {
	TraceLog(SendMsgErr, params, data, err, level)
}
