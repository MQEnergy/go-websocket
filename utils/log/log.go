package log

import (
	"github.com/MQEnergy/go-websocket/utils/code"
	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
)

// TraceLog 写日志
func TraceLog(code code.Code, params, data, err interface{}, level logrus.Level) {
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
		"code":   code.HeartbeatErr,
		"err":    err,
		"params": params,
		"data":   data,
	}).Log(level, code.HeartbeatErr.Msg())
}

// TraceClientCloseFailedLog 客户端关闭失败消息
func TraceClientCloseFailedLog(params, data, err interface{}, level logrus.Level) {
	Logger.WithFields(logrus.Fields{
		"code":   code.ClientCloseFailed,
		"err":    err,
		"params": params,
		"data":   data,
	}).Log(level, code.ClientCloseFailed.Msg())
}

// TraceSuccessLog 客户端连接成功消息
func TraceSuccessLog(params, data interface{}, level logrus.Level) {
	Logger.WithFields(logrus.Fields{
		"code":   code.Success,
		"params": params,
		"data":   data,
	}).Log(level, code.Success.Msg())
}

// TraceReadMsgSuccessLog 读取消息体成功消息
func TraceReadMsgSuccessLog(params, data interface{}, level logrus.Level) {
	Logger.WithFields(logrus.Fields{
		"code":   code.ReadMsgSuccess,
		"params": params,
		"data":   data,
	}).Log(level, code.ReadMsgSuccess.Msg())
}

// TraceSendMsgErrLog 发送消息体失败
func TraceSendMsgErrLog(params, data, err interface{}, level logrus.Level) {
	Logger.WithFields(logrus.Fields{
		"code":   code.SendMsgErr,
		"err":    err,
		"params": params,
		"data":   data,
	}).Log(level, code.SendMsgErr.Msg())
}
