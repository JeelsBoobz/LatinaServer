package helper

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/LalatinaHub/LatinaSub-go/ipapi"
)

var (
	ipinfo ipapi.Ipapi
)

func GetIpInfo() ipapi.Ipapi {
	var (
		buf = new(strings.Builder)
	)
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := httpClient.Get("http://ipinfo.io/json")
	if err != nil {
		return ipinfo
	}
	defer resp.Body.Close()

	io.Copy(buf, resp.Body)
	if resp.StatusCode == 200 {
		ipinfo = ipapi.Parse(buf.String())
		return ipinfo
	}

	return ipinfo
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
