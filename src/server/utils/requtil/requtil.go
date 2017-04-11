package requtil

import (
	"net/http"
	"strings"
)

/**
 *  获取用户ip地址
 */
func GetIp(req *http.Request) string {
	ip := req.RemoteAddr
	if "" == ip {
		ip = req.Header.Get("x-forwarded-for") //
	}
	if "" == ip{
		ip = req.Header.Get("X-Real-IP") //
	}
	if "" != ip{
		ips := strings.Split(ip, ",")
		ip = ips[0]
		return strings.Split(ip,":")[0]
	}
	return ""
}
