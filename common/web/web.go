package web

import (
	"net/http"
	"os"

	"github.com/LalatinaHub/LatinaServer/common/helper"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	password = os.Getenv("PASSWORD")
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

	r.Use(static.Serve("/", static.LocalFile("/usr/local/openresty/nginx/html", false)))

	return r
}
