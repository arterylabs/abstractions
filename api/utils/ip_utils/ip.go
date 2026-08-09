package ip_utils

import (
	"net"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetUserIPEcho(c echo.Context) string {
	remoteAddr, _, _ := net.SplitHostPort(c.Request().RemoteAddr)
	forwardedFor := c.Request().Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		forwardedForArray := strings.Split(forwardedFor, ",")
		if len(forwardedForArray) < 1 {
			return remoteAddr
		}
		return forwardedForArray[0]
	}
	return remoteAddr
}
