package connlib

import (
	"fmt"
)

// IPv4Address Stores IPv4 addresses
// net.IP is not really useful for our purpose since it's not comparable
type IPv4Address [4]byte

func (ip *IPv4Address) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func (ip *IPv4Address) IsUnspecified() bool {
	return ip == nil || *ip == [4]byte{0, 0, 0, 0}
}

type Endpoint struct {
	IP   IPv4Address
	Port uint16
}

func (e Endpoint) IsUnbound() bool {
	return e.IP.IsUnspecified() && e.Port == 0
}

func (e Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.IP.String(), e.Port)
}

type Connection struct {
	Local  Endpoint
	Remote Endpoint
}

type ConnectionList []Connection

func (cl ConnectionList) ToMap() map[Connection]interface{} {
	connMap := make(map[Connection]interface{})

	for _, conn := range cl {
		connMap[conn] = nil
	}

	return connMap
}
