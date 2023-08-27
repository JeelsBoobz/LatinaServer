package web

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/LalatinaHub/LatinaServer/config"
	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/gin-gonic/gin"
)

var (
	password     = os.Getenv("PASSWORD")
	realityRegex = regexp.MustCompile("reality")
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
		case "/reality":
			var (
				singConfig = config.ReadSingConfig()
				text       = []string{}
			)

			text = append(text, "REALITY SERVER INFORMATION")
			text = append(text, "--------------------------")
			text = append(text, "VPN Type : VLESS")
			text = append(text, "Reality Public Key : "+config.RealityPublicKey)
			text = append(text, "Reality ShortID : "+config.RealityShortID[0])
			text = append(text, "")

			for i, inbound := range singConfig.Inbounds {
				if realityRegex.MatchString(inbound.Tag) {
					tag := strings.Split(inbound.Tag, "-")
					text = append(text, fmt.Sprintf("%s : %d", tag[2], 52000+i))
				}
			}
			c.String(http.StatusOK, strings.Join(text, "\n"))
		default:
			if proxy, err := reverse(c, "http://fool.azurewebsites.net/get"); err == nil {
				proxy.ServeHTTP(c.Writer, c.Request)
			}
		}
	})

	return r
}
