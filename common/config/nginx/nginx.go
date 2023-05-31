package nginx

import (
	"fmt"
	"os"
	"strings"

	"github.com/sagernet/sing-box/option"
)

func GenerateReverseProxy(config option.Options) {
	var (
		upstreamList []string = []string{}
		domain       string   = os.Getenv("DOMAIN")
		ports        []uint16 = []uint16{443, 8443}
	)

	for _, inbound := range config.Inbounds {
		upstreamList = append(upstreamList, inbound.Tag)
	}

	result := []string{"stream {"}
	result = append(result, "\tmap $ssl_preread_server_name $singbox {")

	for _, upstream := range upstreamList {
		result = append(result, fmt.Sprintf("\t\t%s.%s %s;", upstream, domain, upstream))
	}
	result = append(result, "\t}")

	for i, upstream := range upstreamList {
		result = append(result, fmt.Sprintf("\tupstream %s {", upstream))
		result = append(result, fmt.Sprintf("\t\tserver 127.0.0.1:%d;", 52000+i))
		result = append(result, "\t}")
	}

	result = append(result, "\tserver {")
	for _, port := range ports {
		result = append(result, fmt.Sprintf("\t\tlisten %d reuseport;", port))
		result = append(result, fmt.Sprintf("\t\tlisten [::]:%d reuseport;", port))
	}
	result = append(result, "\t\tproxy_pass $singbox;")
	result = append(result, "\t\tssl_preread on;")
	result = append(result, "\t\tproxy_protocol on;")
	result = append(result, "\t}")
	result = append(result, "}")

	// Write nginx config
	f, err := os.Create("~/stream.conf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(strings.Join(result[:], "\n"))
}
