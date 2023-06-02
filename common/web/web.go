package web

import (
	"net/http"
	"os"

	"github.com/LalatinaHub/LatinaServer/common/helper"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	C "github.com/sagernet/sing-box/constant"
)

var (
	password = os.Getenv("PASSWORD")
	domain   = os.Getenv("DOMAIN")
)

func WebServer() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/reload", func(c *gin.Context) {
		if c.Query("pass") == password {
			helper.ReloadService([]string{CS.ServiceSingBox, CS.ServiceOpenresty}...)
			c.Status(http.StatusOK)
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})

	r.GET("/info", func(c *gin.Context) {
		if c.Query("pass") == password {
			c.JSON(http.StatusOK, gin.H{
				"domain": domain,
				"ip":     helper.GetOutboundIP(),
				"ports": map[string]uint16{
					"tls":  443,
					"ntls": 80,
				},
				"networks": []string{C.V2RayTransportTypeWebsocket, ""},
			})
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	})

	r.Use(static.Serve("/", static.LocalFile("/usr/local/openresty/nginx/html", false)))

	return r
}
