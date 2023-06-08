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

	g.Go(func() error {
		return web.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
