package go_websocket

import "sync"

type Hub struct {
	Clients       map[*Client]bool    // 全部客户端列表 {*Client1: bool, *Client2: bool...}
	SystemClients map[string][]string // 全部系统列表 {"systemId1": []string{"clientId1", "clientId2"...}, "systemId2": []string{"clientId3", "clientId4"...}}
	Groups        map[string][]string // 全部群组列表  {"systemId1:groupName": []string{"clientId1", "clientId2"...}}

	ClientRegister   chan *Client // 客户端连接处理
	ClientUnregister chan *Client // 客户端断开连接处理
	ClientLock       sync.RWMutex // 客户端列表读写锁
	Broadcast        chan []byte  // 来自客户端的入站消息
}

// NewHub 实例化
func NewHub() *Hub {
	return &Hub{
		Clients:          make(map[*Client]bool),
		Groups:           make(map[string][]string, 1000),
		SystemClients:    make(map[string][]string, 1000),
		ClientRegister:   make(chan *Client),
		ClientUnregister: make(chan *Client),
		Broadcast:        make(chan []byte),
	}
}

func (m *Hub) Run() {
	for {
		select {
		case client := <-m.ClientRegister:
			m.handleClientRegister(client)

		case client := <-m.ClientUnregister:
			m.handleClientUnregister(client)
			close(client.send)

		case message := <-m.Broadcast:
			for client := range m.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					m.handleClientUnregister(client)
				}
			}
		}
	}
}

// handleClientRegister 客户端连接处理
func (m *Hub) handleClientRegister(client *Client) {
	m.ClientLock.Lock()
	m.SystemClients[client.SystemId] = append(m.SystemClients[client.SystemId], client.ClientId)
	m.Clients[client] = true
	m.ClientLock.Unlock()
}

// handleClientUnregister 客户端断开连接处理
func (m *Hub) handleClientUnregister(client *Client) {
	m.ClientLock.Lock()
	if _, ok := m.Clients[client]; ok {
		delete(m.Clients, client)
	}
	for index, clientId := range m.SystemClients[client.SystemId] {
		if clientId == client.ClientId {
			m.SystemClients[client.SystemId] = append(m.SystemClients[client.SystemId][:index], m.SystemClients[client.SystemId][index+1:]...)
			break
		}
	}
	m.ClientLock.Unlock()
}
