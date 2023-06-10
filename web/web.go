package web

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/LalatinaHub/LatinaServer/config"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/helper"
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

	r.GET("/*path", func(c *gin.Context) {
		switch c.Param("path") {
		case "/" + password:
			config.Write()
			helper.ReloadService([]string{CS.ServiceSingBox, CS.ServiceOpenresty}...)
			c.Status(http.StatusOK)
		case "/info":
			c.JSON(http.StatusOK, helper.GetIpInfo())
		default:
			remote, err := url.Parse("http://fool.azurewebsites.net/get")
			if err != nil {
				panic(err)
			}

			proxy := httputil.NewSingleHostReverseProxy(remote)
			proxy.Director = func(req *http.Request) {
				req.Header = c.Request.Header
				req.Host = remote.Host
				req.URL.Scheme = remote.Scheme
				req.URL.Host = remote.Host
				req.URL.Path = "/get"
			}
			proxy.ServeHTTP(c.Writer, c.Request)
		}
	})

	return r
}
