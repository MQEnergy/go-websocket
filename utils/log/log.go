package log

import (
	"github.com/MQEnergy/go-websocket/utils/code"
	"github.com/sirupsen/logrus"
	"time"
)

// WriteLog 写日志
func WriteLog(systemId, clientId, messageId string, data interface{}, code code.Code, msg string, level logrus.Level) {
	logrus.WithFields(logrus.Fields{
		"clientId":  clientId,
		"systemId":  systemId,
		"MessageId": messageId,
		"code":      code,
		"data":      data,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}).Log(level, msg)
}
