package go_websocket

import (
	"bytes"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
)

var (
	Node     *snowflake.Node
	upgrader = websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func init() {
	localIp, err := GetLocalIpToInt()
	if err != nil {
		panic(err)
	}
	Node, err = snowflake.NewNode(int64(localIp) % 1023)
	if err != nil {
		panic(err)
	}
}

// ReadMessageHandler 将来自 websocket 连接的消息推送到集线器。
func (c *Client) ReadMessageHandler() {
	defer func() {
		c.hub.ClientUnregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(appData string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				TraceClientCloseSuccessLog("", "", err.Error(), 4)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.Broadcast <- message
	}
}

// WriteMessageHandler 将消息从集线器发送到 websocket 连接
func (c *Client) WriteMessageHandler() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			data := make(map[string]interface{}, 0)
			if err := json.Unmarshal(message, &data); err != nil {
				return
			}
			c.Conn.SetWriteDeadline(time.Time{})
			WriteMessage(c.Conn, SendMsgSuccess, SendMsgSuccess.Msg(), data, nil, Binary)

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// WsServer handles websocket requests from the peer.
func WsServer(hub *Hub, w http.ResponseWriter, r *http.Request) (*Client, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	systemId := r.FormValue("system_id")
	if systemId == "" {
		sid, err := GetLocalIpToInt()
		if err != nil {
			return nil, err
		}
		systemId = strconv.Itoa(sid)
	}
	client := &Client{
		SystemId: systemId,
		ClientId: GenerateUuid(Node),
		hub:      hub,
		Conn:     conn,
		send:     make(chan []byte, 256),
	}
	client.hub.ClientRegister <- client
	WriteMessage(conn, Success, Success.Msg(), map[string]string{"system_id": systemId, "client_id": client.ClientId}, nil, Binary)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WriteMessageHandler()
	go client.ReadMessageHandler()

	return client, nil
}
