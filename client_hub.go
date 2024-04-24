package go_websocket

import (
	"errors"
	"sync"
)

type BroadcastChan struct {
	Name string `json:"name"`
	Msg  []byte `json:"msg"`
}

type Hub struct {
	Clients       map[*Client]bool     // 全部客户端列表 {*Client1: bool, *Client2: bool...}
	SystemClients map[string][]*Client // 全部系统列表 {"systemId1": []*Clients{*Client1, *Client2...}, "systemId2": []*Clients{*Client1, *Client2...}}
	GroupClients  map[string][]*Client // 全部群组列表 {"groupId": []*Clients{*Client1, *Client2...}}

	ClientRegister   chan *Client        // 客户端连接处理
	ClientUnregister chan *Client        // 客户端断开连接处理
	ClientLock       sync.RWMutex        // 客户端列表读写锁
	Broadcast        chan []byte         // 来自广播的入站消息
	SystemBroadcast  chan *BroadcastChan // 来自系统的入站消息 {Name:"systemId", Msg:"msg"}
	GroupBroadcast   chan *BroadcastChan // 来自群组的入站消息 {Name:"groupId", Msg:"msg"}
	ClientBroadcast  chan *BroadcastChan // 来自客户端的入站消息 {Name:"clientId", Msg:"msg"}
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
		SystemBroadcast:  make(chan *BroadcastChan, 1000),
		GroupBroadcast:   make(chan *BroadcastChan, 1000),
		ClientBroadcast:  make(chan *BroadcastChan, 1000),
	}
}

// Run run chan listener
func (m *Hub) Run() {
	for {
		select {
		case client := <-m.ClientRegister:
			m.handleClientRegister(client)

		case client := <-m.ClientUnregister:
			m.handleClientUnregister(client)
			close(client.send)

		// 全局广播
		case message := <-m.Broadcast:
			m.AllBroadcastHandle(message)
		// 系统广播
		case systems := <-m.SystemBroadcast:
			m.SystemBroadcastHandle(systems.Name, systems.Msg)
		// 群组广播
		case groups := <-m.GroupBroadcast:
			m.GroupBroadcastHandle(groups.Name, groups.Msg)
		// 客户端推送
		case clients := <-m.ClientBroadcast:
			m.ClientBroadcastHandle(clients.Name, clients.Msg)
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
	clients, ok := m.GroupClients[client.GroupId]
	if ok {
		for index, _client := range clients {
			if _client.ClientId == client.ClientId {
				m.GroupClients[client.GroupId] = append(m.GroupClients[client.GroupId][:index], m.GroupClients[client.GroupId][index+1:]...)
			}
		}
	}
	m.ClientLock.Unlock()
}

// AllBroadcastHandle 全局广播
func (m *Hub) AllBroadcastHandle(msg []byte) {
	for client := range m.Clients {
		select {
		case client.send <- msg:
		default:
			close(client.send)
			m.handleClientUnregister(client)
		}
	}
}

// SystemBroadcastHandle 系统广播处理
func (m *Hub) SystemBroadcastHandle(systemId string, msg []byte) {
	clients, err := m.GetSystemClients(systemId)
	if err != nil {
		m.RemoveSystem(systemId)
	}
	for _, client := range clients {
		select {
		case client.send <- msg:
		default:
			close(client.send)
			m.handleClientUnregister(client)
		}
	}
}

// GroupBroadcastHandle 群组消息通道处理
func (m *Hub) GroupBroadcastHandle(groupId string, msg []byte) {
	clients, err := m.GetGroupClients(groupId)
	if err != nil {
		m.RemoveGroup(groupId)
	}
	for _, client := range clients {
		select {
		case client.send <- msg:
		default:
			close(client.send)
			m.handleClientUnregister(client)
		}
	}
}

// ClientBroadcastHandle 单客户端通道处理
func (m *Hub) ClientBroadcastHandle(clientId string, msg []byte) {
	var _client *Client
	for client := range m.Clients {
		if client.ClientId == clientId {
			_client = client
			break
		}
	}
	if _client != nil {
		select {
		case _client.send <- msg:
			break
		default:
			close(_client.send)
			m.handleClientUnregister(_client)
		}
	}
}

// SetClientToGroups 添加客户端到分组
func (m *Hub) SetClientToGroups(groupId string, client *Client) bool {
	clients, ok := m.GroupClients[groupId]
	if !ok {
		return false
	}
	for _, _client := range clients {
		if _client.ClientId == client.ClientId {
			return false
		}
	}
	m.ClientLock.Lock()
	m.GroupClients[groupId] = append(m.GroupClients[groupId], client)
	m.ClientLock.Unlock()
	return true
}

// GetSystemClients 获取系统的客户端列表
func (m *Hub) GetSystemClients(name string) ([]*Client, error) {
	clients, ok := m.SystemClients[name]
	if !ok {
		return []*Client{}, errors.New("group does not exist")
	}
	return clients, nil
}

// GetGroupClients 获取群组的客户端列表
func (m *Hub) GetGroupClients(name string) ([]*Client, error) {
	clients, ok := m.GroupClients[name]
	if !ok {
		return []*Client{}, errors.New("group does not exist")
	}
	return clients, nil
}

// RemoveSystem 删除system和系统中的client
func (m *Hub) RemoveSystem(name string) {
	delete(m.SystemClients, name)
}

// RemoveGroup 删除group和群组中的client
func (m *Hub) RemoveGroup(name string) {
	delete(m.GroupClients, name)
}

// RemoveClientByGroup 从群组删除客户端
func (m *Hub) RemoveClientByGroup(client *Client) error {
	m.ClientLock.Lock()
	clients, ok := m.GroupClients[client.GroupId]
	if !ok {
		return errors.New("group does not exist")
	}
	for index, _client := range clients {
		if _client.ClientId == client.ClientId {
			m.GroupClients[client.GroupId] = append(m.GroupClients[client.GroupId][:index], m.GroupClients[client.GroupId][index+1:]...)
		}
	}
	m.ClientLock.Unlock()
	return nil
}
