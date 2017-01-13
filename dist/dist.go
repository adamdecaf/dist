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
func (a Address) Equal(o Address) bool {
	return a.IP.Equal(o.IP) && (a.Port == o.Port)
}
func (a Address) Valid() bool {
	return (len(a.IP) > 0 && len(a.IP.String()) > 0) && a.Port > 0
}
