package log

import (
	"github.com/sirupsen/logrus"
	"gowebsocket/utils/code"
)

var (
	Logger *logrus.Logger
)

// WriteLog 写日志
func WriteLog(systemId, clientId, messageId string, data interface{}, code code.Code, msg string, level logrus.Level) {
	Logger.WithFields(logrus.Fields{
		"client_id":  clientId,
		"system_id":  systemId,
		"message_id": messageId,
		"code":       code,
		"data":       data,
	}).Log(level, msg)
}
