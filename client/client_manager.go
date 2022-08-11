package client

import (
	"errors"
	"github.com/MQEnergy/go-websocket/utils/code"
	"github.com/MQEnergy/go-websocket/utils/response"
	"sync"
)

type Manager struct {
	ClientList       map[string]*Client  // 全部客户端列表 {"clientId1": *Client, "clientId2": *Client...}
	GroupList        map[string][]string // 全部群组列表  {"systemId:groupName": []string{"clientId1", "clientId2"...}}
	SystemClientList map[string][]string // 全部系统列表 {"systemId1": []string{"clientId1", "clientId2"...}, "systemId2": []string{"clientId3", "clientId4"...}}

	ClientConnect    chan *Client // 客户端连接处理
	ClientDisConnect chan *Client // 客户端断开连接处理

	ClientListLock       sync.RWMutex // 客户端列表读写锁
	GroupListLock        sync.RWMutex // 群组列表读写锁
	SystemClientListLock sync.RWMutex // 系统列表读写锁
}

// NewManager 实例化
func NewManager() *Manager {
	return &Manager{
		ClientList:       make(map[string]*Client),
		GroupList:        make(map[string][]string, 10000),
		SystemClientList: make(map[string][]string, 10000),
		ClientConnect:    make(chan *Client, 10000),
		ClientDisConnect: make(chan *Client, 10000),
	}
}

// ClientConnectHandler 客户端连接handler
func (m *Manager) ClientConnectHandler(client *Client) error {
	// 建立连接事件
	m.SetClientToList(client)
	// 发送给客户端连接成功
	if err := response.WsJson(client.Conn, client.SystemId, client.ClientId, "", code.Success, code.Success.Msg(), "", ""); err != nil {
		m.ClientDisConnect <- client
		return err
	}
	return nil
}

// ClientDisConnectHandler 客户端断连handler
func (m *Manager) ClientDisConnectHandler(client *Client) error {
	// 断开连接事件
	if err := client.Conn.Close(); err != nil {
		return err
	}
	// 删除客户端
	m.RemoveAllClient(client)
	// 清除当前客户端
	client = nil
	return nil
}

// SetClientToList 添加客户端到列表
func (m *Manager) SetClientToList(client *Client) {
	m.ClientListLock.Lock()
	defer m.ClientListLock.Unlock()
	m.ClientList[client.ClientId] = client
}

// SetSystemClientToList 添加系统ID和客户端到列表
func (m *Manager) SetSystemClientToList(client *Client) {
	m.SystemClientListLock.Lock()
	defer m.SystemClientListLock.Unlock()
	m.SystemClientList[client.SystemId] = append(m.SystemClientList[client.SystemId], client.ClientId)
}

// SetClientToGroupList 添加客户端到分组
func (m *Manager) SetClientToGroupList(groupName string, client *Client) error {
	//判断之前是否有添加过
	groupFlag := true
	for _, groupValue := range client.GroupList {
		if groupValue == groupName {
			groupFlag = false
			break
		}
	}
	// 组名添加
	if groupFlag {
		client.GroupList = append(client.GroupList, groupName)
	}
	groupKey := client.SystemId + ":" + groupName
	// 判断客户端是否已添加至分组
	flag := true
	groupClientList := m.GetGroupClientList(groupKey)
	for _, clientId := range groupClientList {
		if clientId == client.ClientId {
			// 群组中存在当前客户端
			flag = false
			break
		}
	}
	if !flag {
		return errors.New("请勿重复添加到同一群组")
	}
	m.GroupListLock.Lock()
	defer m.GroupListLock.Unlock()
	m.GroupList[groupKey] = append(m.GroupList[groupKey], client.ClientId)
	return nil
}

// GetAllClient 获取所有的客户端
func (m *Manager) GetAllClient() map[string]*Client {
	m.ClientListLock.RLock()
	defer m.ClientListLock.RUnlock()
	return m.ClientList
}

// GetAllClientCount 获取所有客户端数量
func (m *Manager) GetAllClientCount() int {
	m.ClientListLock.RLock()
	defer m.ClientListLock.RUnlock()
	return len(m.ClientList)
}

// GetClientByList 通过客户端列表获取*Client
func (m *Manager) GetClientByList(clientId string) (*Client, error) {
	m.ClientListLock.RLock()
	defer m.ClientListLock.RUnlock()
	client, ok := m.ClientList[clientId]
	if !ok {
		return nil, errors.New(code.ClientNotExist.Msg())
	}
	return client, nil
}

// GetSystemClientList 获取指定系统的客户端列表
func (m *Manager) GetSystemClientList(systemId string) []string {
	m.SystemClientListLock.RLock()
	defer m.SystemClientListLock.RUnlock()
	return m.SystemClientList[systemId]
}

// GetGroupClientList 获取本地分组的成员
func (m *Manager) GetGroupClientList(groupKey string) []string {
	m.GroupListLock.RLock()
	defer m.GroupListLock.RUnlock()
	return m.GroupList[groupKey]
}

// RemoveAllClient 删除所有客户端
func (m *Manager) RemoveAllClient(client *Client) {
	// 删除 *Client
	m.RemoveClientByList(client.ClientId)
	// 删除所在的分组
	if len(client.GroupList) > 0 {
		for _, groupName := range client.GroupList {
			groupKey := client.SystemId + ":" + groupName
			m.RemoveGroupClient(groupKey, client.ClientId)
		}
	}
	// 删除系统里的客户端
	m.RemoveSystemClientByList(client)
}

// RemoveGroupClient 删除分组里的客户端
func (m *Manager) RemoveGroupClient(groupKey, clientId string) {
	m.GroupListLock.Lock()
	defer m.GroupListLock.Unlock()

	for index, _clientId := range m.GroupList[groupKey] {
		if _clientId == clientId {
			m.GroupList[groupKey] = append(m.GroupList[groupKey][:index], m.GroupList[groupKey][index+1:]...)
		}
	}
}

// RemoveClientByList 从列表删除*Client
func (m *Manager) RemoveClientByList(clientId string) {
	m.ClientListLock.Lock()
	defer m.ClientListLock.Unlock()
	delete(m.ClientList, clientId)
}

// RemoveSystemClientByList 删除系统里的客户端
func (m *Manager) RemoveSystemClientByList(client *Client) {
	m.SystemClientListLock.Lock()
	defer m.SystemClientListLock.Unlock()

	for index, clientId := range m.SystemClientList[client.SystemId] {
		if clientId == client.ClientId {
			m.SystemClientList[client.SystemId] = append(m.SystemClientList[client.SystemId][:index], m.SystemClientList[client.SystemId][index+1:]...)
		}
	}
}

// CloseClient 关闭客户端
func (m *Manager) CloseClient(clientId, systemId string) error {
	conn, err := m.GetClientByList(clientId)
	if err == nil && conn != nil {
		if conn.SystemId != systemId {
			return errors.New(code.RequestParamErr.Msg())
		}
		m.ClientDisConnect <- conn
	}
	return nil
}
