package web

import (
	"net/http"
	"os"

	"github.com/LalatinaHub/LatinaServer/config"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/arl/statsviz"
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

	r.GET("/*filepath", func(c *gin.Context) {
		switch c.Param("filepath") {
		case "/ws":
			statsviz.Ws(c.Writer, c.Request)
		case "/" + password:
			config.Write()
			helper.ReloadService([]string{CS.ServiceSingBox, CS.ServiceOpenresty}...)
			c.Status(http.StatusOK)
		case "/info":
			c.JSON(http.StatusOK, helper.GetIpInfo())
		default:
			statsviz.IndexAtRoot("/").ServeHTTP(c.Writer, c.Request)
		}
	})

	return r
}
