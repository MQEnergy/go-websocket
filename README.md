# go-websocket
websocket库

## client_manager.go方法和函数
```
SetClientToList             添加客户端到列表
SetSystemClientToList       添加系统ID和客户端到列表
SetClientToGroupList        添加客户端到分组

GetAllClient                获取所有的客户端
GetAllClientCount           获取所有客户端数量
GetClientByList             通过客户端列表获取*Client
GetSystemClientList         获取指定系统的客户端列表
GetGroupClientList          获取本地分组的成员

RemoveAllClient             删除所有客户端 (包含 ClientList SystemClientList GroupList)
RemoveGroupClient           删除分组里的客户端
RemoveClientByList          从列表删除*Client
RemoveSystemClientByList    删除系统里的客户端
CloseClient                 关闭客户端 (发送关闭client到ClientDisConnect通道中)
```
