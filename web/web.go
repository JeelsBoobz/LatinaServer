package web

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/LalatinaHub/LatinaServer/config"
	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/gin-gonic/gin"
)

var (
	password = os.Getenv("PASSWORD")
)

func reverse(c *gin.Context, target string) (*httputil.ReverseProxy, error) {
	remote, err := url.Parse(target)
	if err != nil {
		fmt.Println(err)
		return &httputil.ReverseProxy{}, err
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
	}

	return proxy, err
}

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
		case "/relay":
			c.JSON(http.StatusOK, relay.Relays)
		case "/port":
			var (
				singConfig = config.ReadSingConfig()
				portPair   = []string{}
			)

			for _, inbound := range singConfig.Inbounds {
				portPair = append(portPair, fmt.Sprintf("%s : %d", inbound.Tag, 52000+len(portPair)))
			}
			c.String(http.StatusOK, strings.Join(portPair, "\n"))
		default:
			if proxy, err := reverse(c, "http://fool.azurewebsites.net/get"); err == nil {
				proxy.ServeHTTP(c.Writer, c.Request)
			}
		}
	})

	return r
}
