package go_websocket

import (
	"errors"
	"sync"
)

type Hub struct {
	Clients       map[*Client]bool     // 全部客户端列表 {*Client1: bool, *Client2: bool...}
	SystemClients map[string][]*Client // 全部系统列表 {"systemId1": []*Clients{*Client1, *Client2...}, "systemId2": []*Clients{*Client1, *Client2...}}
	GroupClients  map[string][]*Client // 全部群组列表 {"groupName": []*Clients{*Client1, *Client2...}}

	ClientRegister   chan *Client           // 客户端连接处理
	ClientUnregister chan *Client           // 客户端断开连接处理
	ClientLock       sync.RWMutex           // 客户端列表读写锁
	Broadcast        chan []byte            // 来自客户端的入站消息
	GroupBroadcast   chan map[string][]byte // 来自群组的入站消息 {groupName:[]byte}
}

// NewHub 实例化
func NewHub() *Hub {
	return &Hub{
		Clients:          make(map[*Client]bool),
		GroupClients:     make(map[string][]*Client, 1000),
		SystemClients:    make(map[string][]*Client, 1000),
		ClientRegister:   make(chan *Client),
		ClientUnregister: make(chan *Client),
		Broadcast:        make(chan []byte),
		GroupBroadcast:   make(chan map[string][]byte, 1000),
	}
}

// Run
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
		case groups := <-m.GroupBroadcast:
			for gname, message := range groups {
				clients, err := m.GetGroupClients(gname)
				if err != nil {
					m.RemoveGroup(gname)
					break
				}
				for _, client := range clients {
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
}

// handleClientRegister 客户端连接处理
func (m *Hub) handleClientRegister(client *Client) {
	m.ClientLock.Lock()
	m.SystemClients[client.SystemId] = append(m.SystemClients[client.SystemId], client)
	if client.GroupId != "" {
		m.GroupClients[client.GroupId] = append(m.GroupClients[client.GroupId], client)
	}
	m.Clients[client] = true
	m.ClientLock.Unlock()
}

// handleClientUnregister 客户端断开连接处理
func (m *Hub) handleClientUnregister(client *Client) {
	m.ClientLock.Lock()
	if _, ok := m.Clients[client]; ok {
		delete(m.Clients, client)
	}
	for index, _client := range m.SystemClients[client.SystemId] {
		if _client.ClientId == client.ClientId {
			m.SystemClients[client.SystemId] = append(m.SystemClients[client.SystemId][:index], m.SystemClients[client.SystemId][index+1:]...)
			break
		}
	}
	m.ClientLock.Unlock()
}

// SetClientToGroups 添加客户端到分组
func (m *Hub) SetClientToGroups(groupName string, client *Client) bool {
	clients, ok := m.GroupClients[groupName]
	if !ok {
		return false
	}
	for _, _client := range clients {
		if _client.ClientId == client.ClientId {
			return false
		}
	}
	m.ClientLock.Lock()
	m.GroupClients[groupName] = append(m.GroupClients[groupName], client)
	m.ClientLock.Unlock()
	return true
}

// GetGroupClients 获取群组的客户端列表
func (m *Hub) GetGroupClients(name string) ([]*Client, error) {
	clients, ok := m.GroupClients[name]
	if !ok {
		return []*Client{}, errors.New("group name is not exist")
	}
	return clients, nil
}

// RemoveGroup 删除group和群组中的client
func (m *Hub) RemoveGroup(name string) {
	delete(m.GroupClients, name)
}
