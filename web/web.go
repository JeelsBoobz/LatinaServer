package web

import (
	"net/http"
	"os"

	"github.com/LalatinaHub/LatinaServer/config"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	password = os.Getenv("PASSWORD")
)

func WebServer() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	if password == "" {
		password = "reload"
	}

	r.GET("/"+password, func(c *gin.Context) {
		config.Write()
		helper.ReloadService([]string{CS.ServiceSingBox, CS.ServiceOpenresty}...)
		c.Status(http.StatusOK)
	})

	r.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, helper.GetIpInfo())
	})

	r.Use(static.Serve("/", static.LocalFile("/usr/local/openresty/nginx/html", false)))

	return r
}
