package subfinder

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	httpxRunner "github.com/projectdiscovery/httpx/runner"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

func SubfinderHandler(c *gin.Context) (gin.H, error) {
	var (
		domain = c.Query("domain")
		result = []httpxRunner.Result{}
	)

	// Subdomain enumeration
	subfinderOpts := &runner.Options{
		Threads:            10,
		Timeout:            30,
		MaxEnumerationTime: 10,
	}

	subfinder, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		return gin.H{
			"error":   true,
			"message": fmt.Sprintf("Failed to create subfinder runner: %v", err),
		}, err
	}

	output := &bytes.Buffer{}
	if err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), domain, []io.Writer{output}); err != nil {
		return gin.H{
			"error":   true,
			"message": fmt.Sprintf("Failed to enumerate domain: %v", err),
		}, err
	}

	// Scan each subdomain
	httpxOpts := httpxRunner.Options{
		Methods:         "GET",
		InputTargetHost: strings.Split(output.String(), "\n"),
		OnResult: func(r httpxRunner.Result) {
			if r.Input != "" {
				result = append(result, r)
			}

		},
	}

	if err = httpxOpts.ValidateOptions(); err != nil {
		return gin.H{
			"error":   true,
			"message": fmt.Sprintf("Failed to validate httpx options: %v", err),
		}, err
	}

	hxRunner, err := httpxRunner.New(&httpxOpts)
	if err != nil {
		return gin.H{
			"error":   true,
			"message": fmt.Sprintf("Failed to start httpx runner: %v", err),
		}, err
	}
	defer hxRunner.Close()

	hxRunner.RunEnumeration()

	return gin.H{
		"error":  false,
		"result": result,
	}, nil
}
