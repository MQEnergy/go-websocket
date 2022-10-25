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
	Logger.WithFields(logrus.Fields{
		"code":   HeartbeatErr,
		"err":    err,
		"params": params,
		"data":   data,
	}).Log(level, HeartbeatErr.Msg())
}

// TraceClientCloseFailedLog 客户端关闭失败消息
func TraceClientCloseFailedLog(params, data, err interface{}, level logrus.Level) {
	Logger.WithFields(logrus.Fields{
		"code":   ClientCloseFailed,
		"err":    err,
		"params": params,
		"data":   data,
	}).Log(level, ClientCloseFailed.Msg())
}

// TraceSuccessLog 客户端连接成功消息
func TraceSuccessLog(params, data interface{}, level logrus.Level) {
	Logger.WithFields(logrus.Fields{
		"code":   Success,
		"params": params,
		"data":   data,
	}).Log(level, Success.Msg())
}

// TraceReadMsgSuccessLog 读取消息体成功消息
func TraceReadMsgSuccessLog(params, data interface{}, level logrus.Level) {
	Logger.WithFields(logrus.Fields{
		"code":   ReadMsgSuccess,
		"params": params,
		"data":   data,
	}).Log(level, ReadMsgSuccess.Msg())
}

// TraceSendMsgErrLog 发送消息体失败
func TraceSendMsgErrLog(params, data, err interface{}, level logrus.Level) {
	Logger.WithFields(logrus.Fields{
		"code":   SendMsgErr,
		"err":    err,
		"params": params,
		"data":   data,
	}).Log(level, SendMsgErr.Msg())
}
