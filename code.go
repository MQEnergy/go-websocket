package go_websocket

type Code int

const (
	Success Code = 0
	Failed  Code = 10001 + iota
	ClientFailed
	ClientNotExist
	ClientCloseSuccess
	ClientCloseFailed
	ReadMsgErr
	ReadMsgSuccess
	SendMsgErr
	SendMsgSuccess
	HeartbeatErr
	SystemErr
	BindGroupSuccess
	BindGroupErr
	UnAuthed
	InternalErr
	RequestMethodErr
	RequestParamErr
)

var CodeMap = map[Code]string{
	Success:            "客户端连接成功",
	Failed:             "客户端连接失败",
	ClientFailed:       "客户端主动断连",
	ClientNotExist:     "客户端不存在",
	ClientCloseSuccess: "客户端关闭成功",
	ClientCloseFailed:  "客户端关闭失败",
	ReadMsgErr:         "读取消息体失败",
	ReadMsgSuccess:     "读取消息体成功",
	SendMsgErr:         "发送消息体失败",
	SendMsgSuccess:     "发送消息体成功",
	HeartbeatErr:       "心跳检测失败",
	SystemErr:          "系统不能为空",
	BindGroupSuccess:   "绑定群组成功",
	BindGroupErr:       "绑定群组失败",
	UnAuthed:           "用户未认证",
	InternalErr:        "服务器内部错误",
	RequestMethodErr:   "请求方式错误",
	RequestParamErr:    "请求参数错误",
}

// Msg 返回错误码对应的说明
func (c Code) Msg() string {
	if v, ok := CodeMap[c]; ok {
		return v
	}
	return ``
}
