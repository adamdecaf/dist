package dist

import (
	"fmt"
	"net"
)

type Address struct {
	IP net.IP `json:"ip"`
	Port int `json:"port"`
}
func (a Address) String() string {
	return fmt.Sprintf("%s:%d", a.IP.String(), a.Port)
}
