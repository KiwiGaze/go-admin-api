package common

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetClientIP(c *gin.Context) string {
	ip := c.Request.Header.Get("X-Forwarded-For")
	if ip == "" || strings.Contains(ip, "127.0.0.1") {
		// If empty or a local address, try getting it from X-Real-IP.
		ip = c.Request.Header.Get("X-Real-IP")
	}
	if ip == "" {
		// If it is still empty, use RemoteIP.
		ip = c.RemoteIP()
	}
	if ip == "" || ip == "127.0.0.1" {
		// If it is still empty or a local address, use ClientIP.
		ip = c.ClientIP()
	}
	if ip == "" {
		// Fall back to the local address as a last resort.
		ip = "127.0.0.1"
	}
	return ip
}
