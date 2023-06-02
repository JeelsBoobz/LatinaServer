package web

import (
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func StartWebService() {
	web := &http.Server{
		Addr:         ":50000",
		Handler:      WebServer(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	reverseProxy := &http.Server{
		Addr:         ":50001",
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
