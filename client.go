package go_websocket

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/websocket"
	"strconv"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 5012

	readBufferSize  = 1024 // 读缓冲区大小
	writeBufferSize = 1024 // 写缓冲区大小
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	ClientId string `json:"client_id"` // 客户端连接ID
	GroupId  string `json:"group_id"`  // 群组id
	SystemId string `json:"system_id"` // 系统ID 为分布式做准备的
	Conn     *websocket.Conn
	Send     chan []byte
	hub      *Hub
}

// GenerateUuid 生成唯一ID
func GenerateUuid(node *snowflake.Node) string {
	if node == nil {
		var err error
		node, err = snowflake.NewNode(1)
		if err != nil {
			return ""
		}
	}
	id := node.Generate()
	return strconv.FormatInt(id.Int64(), 10)
}
