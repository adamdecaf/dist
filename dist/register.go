package dist

import (
	"log"
	"fmt"
	"net"
	"net/http"
)

const (
	DirAddress = "127.0.0.1:8080"
)

func freePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func outboundIP() *net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range addrs {
		// find the first non-loopback ip
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return &ip.IP
			}
		}
	}
	return nil
}

func registerWithDir(addr Address) error {
	resp, err := http.Get(fmt.Sprintf("http://%s/register?address=%s&port=%d", DirAddress, addr.IP.String(), addr.Port))
	if err != nil {
		return err
	}
	if resp.StatusCode - 200 > 99 {
		return fmt.Errorf("status code != 2xx, got=%d", resp.StatusCode)
	}
	log.Println(resp.StatusCode)
	return nil
}

// Check err != nil on this!
func Register() (Address, error) {
	ip := outboundIP()
	port := freePort()
	if ip == nil || (port < 0 || port > 65535) {
		log.Fatalf("error getting worker address, ip=%s, port=%d", ip.String(), port)
	}
	addr := Address{
		IP: *ip,
		Port: port,
	}
	log.Println("register self")
	err := registerWithDir(addr)
	log.Println(err)
	return addr, err
}
