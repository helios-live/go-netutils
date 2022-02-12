package netutils

import (
	"net"
	"regexp"
	"strconv"
)

const (
	// TypeHTTP is used for HTTP upstream proxies
	TypeHTTP ProxyType = "HTTP"

	// TypeSocks5 is used for Socks5 upstream proxies
	TypeSocks5 ProxyType = "Socks5"
)

// ProxyType is a string either HTTP or Socks5
type ProxyType string

// ProxyInfo forwards represents a http proxy
type ProxyInfo struct {
	User       string
	Pass       string
	Host       string
	Port       int
	Type       ProxyType
	connection net.Conn
}

// ReadProxy returns a pointer to a new instance of ProxyInfo parsed from proxy
// @param proxy string format user:pass@ip:port
func ReadProxy(proxy string) *ProxyInfo {

	r, _ := regexp.Compile("^(?P<type>(http|socks5)://){0,1}((?P<user>[^:@]+)(:(?P<pass>[^:@]+)){0,}@){0,1}(?P<host>[0-9A-Za-z\\.\\-]+):(?P<port>[0-9]+)")
	res := r.FindStringSubmatch(proxy)
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = res[i]
		}
	}
	var tp ProxyType
	if v, ok := result["type"]; !ok {
		tp = "HTTP"
	} else {
		if v == "socks5://" {
			tp = TypeSocks5
		} else {
			tp = TypeHTTP
		}
	}

	port, _ := strconv.Atoi(result["port"])
	return &ProxyInfo{
		User: result["user"],
		Pass: result["pass"],
		Host: result["host"],
		Port: port,
		Type: tp,
	}
}
