package web

import (
	"net/http"
	"time"

	CS "github.com/LalatinaHub/LatinaServer/constant"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func StartWebService() {
	web := &http.Server{
		Addr:         CS.WebServerAddress,
		Handler:      WebServer(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	reverseProxy := &http.Server{
		Addr:         CS.ReverseProxyAddress,
		Handler:      ReverseProxy(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return web.ListenAndServe()
	})
	g.Go(func() error {
		return reverseProxy.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
