package go_websocket

import (
	"encoding/binary"
	"errors"
	"net"
)

// convertToIntIP 转换ip为int
func convertToIntIP(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

// GetLocalIpToInt 获取本机IP转成int
func GetLocalIpToInt() (uint32, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return 0, err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return convertToIntIP(ipnet.IP), nil
			}
		}
	}
	return 0, errors.New("can not find the client ip address")
}
