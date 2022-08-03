package client

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/websocket"
	"strconv"
)

type Client struct {
	ClientId  string          `json:"client_id"`  // 客户端连接ID
	SystemId  string          `json:"system_id"`  // 系统ID 为分布式做准备的
	Conn      *websocket.Conn `json:"conn"`       // websocket连接
	GroupList []string        `json:"group_list"` // 分组列表 []string{ groupName1, groupName2... }
}

// NewClient 实例化客户端
func NewClient(clientId string, systemId string, conn *websocket.Conn) *Client {
	return &Client{
		ClientId:  clientId,
		SystemId:  systemId,
		Conn:      conn,
		GroupList: make([]string, 0),
	}
}

// GenerateUuid 生成唯一ID
func GenerateUuid(num int, node *snowflake.Node) string {
	if node == nil {
		var err error
		node, err = snowflake.NewNode(1)
		if err != nil {
			return ""
		}
	}
	id := node.Generate()
	switch num {
	case 2:
		return id.Base2()
	case 32:
		return id.Base32()
	case 36:
		return id.Base36()
	case 58:
		return id.Base58()
	case 64:
		return id.Base64()
	default:
		return strconv.FormatInt(id.Int64(), 10)
	}
}
