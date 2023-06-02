package web

import (
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func WebServer() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	r.Use(static.Serve("/", static.LocalFile("/usr/local/openresty/nginx/html", false)))

	return r
}
