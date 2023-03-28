package utils

import (
	"net"
	"net/http"
	"strings"
)

// 获取客户端请求IP
func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip
	}
	ip = r.Header.Get("X-Forwarded-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	if net.ParseIP(ip) != nil {
		return ip
	}
	return ""
}

// ExternalIP获取外部IP.
func ExternalIP() (res []string) {
	inters, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, inter := range inters {
		if !strings.HasPrefix(inter.Name, "lo") {
			addresses, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addresses {
				if ipNet, ok := addr.(*net.IPNet); ok {
					if ipNet.IP.IsLoopback() || ipNet.IP.IsLinkLocalMulticast() || ipNet.IP.IsLinkLocalUnicast() {
						continue
					}
					if ip4 := ipNet.IP.To4(); ip4 != nil {
						switch true {
						case ip4[0] == 10:
							continue
						case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
							continue
						case ip4[0] == 192 && ip4[1] == 168:
							continue
						default:
							res = append(res, ipNet.IP.String())
						}
					}
				}
			}
		}
	}
	return
}

// InternalIP获取内部IP.
func InternalIP() string {
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range inters {
		if !isUp(inter.Flags) {
			continue
		}
		if !strings.HasPrefix(inter.Name, "lo") {
			addresses, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addresses {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						return ipNet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

// isUp Interface is up
func isUp(v net.Flags) bool {
	return v&net.FlagUp == net.FlagUp
}
