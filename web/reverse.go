package web

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/LalatinaHub/LatinaServer/config"
	"github.com/gin-gonic/gin"
	C "github.com/sagernet/sing-box/constant"
)

func proxy(c *gin.Context) {
	var (
		config        = config.ReadSingConfig()
		port   uint16 = 50000
		path          = c.Param("proxyPath")
	)

	switch c.Request.Header.Get("Upgrade") {
	case "websocket":
		for _, inbound := range config.Inbounds {
			if "/"+inbound.Type == path {
				switch inbound.Type {
				case C.TypeTrojan:
					port = inbound.TrojanOptions.ListenPort
				case C.TypeVMess:
					port = inbound.VMessOptions.ListenPort
				case C.TypeVLESS:
					port = inbound.VLESSOptions.ListenPort
				}
			}
		}
	}

	remote, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = path
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func ReverseProxy() http.Handler {
	r := gin.New()
	r.GET("/*proxyPath", proxy)

	return r
}
