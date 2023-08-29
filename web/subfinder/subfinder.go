package subfinder

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	httpxRunner "github.com/projectdiscovery/httpx/runner"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

func SubfinderHandler(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": fmt.Sprintf("Failed to create subfinder runner: %v", err),
		})
		return
	}

	output := &bytes.Buffer{}
	if err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), domain, []io.Writer{output}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": fmt.Sprintf("Failed to enumerate domain: %v", err),
		})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": fmt.Sprintf("Failed to validate httpx options: %v", err),
		})
		return
	}

	hxRunner, err := httpxRunner.New(&httpxOpts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": fmt.Sprintf("Failed to start httpx runner: %v", err),
		})
		return
	}
	defer hxRunner.Close()

	hxRunner.RunEnumeration()

	c.JSON(http.StatusOK, gin.H{
		"error":  false,
		"result": result,
	})
}
