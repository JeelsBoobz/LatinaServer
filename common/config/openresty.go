package config

import (
	"os"
	"strconv"
	"strings"

	C "github.com/sagernet/sing-box/constant"
)

func WriteReverseConfig() {
	config := ReadSingConfig()
	r, err := os.ReadFile("./resources/openresty/stream/reverse.conf")
	if err != nil {
		panic(err)
	}

	reverseConfig := string(r)

	for _, inbound := range config.Inbounds {
		switch inbound.Type {
		case C.TypeTrojan:
			trojan := inbound.TrojanOptions
			if trojan.Transport == nil {
				reverseConfig = strings.Replace(reverseConfig, "TROJAN_TCP_PORT", strconv.Itoa(int(trojan.ListenPort)), 1)
			}
		case C.TypeVMess:
			vmess := inbound.VMessOptions
			if vmess.Transport == nil {
				reverseConfig = strings.Replace(reverseConfig, "VMESS_TCP_PORT", strconv.Itoa(int(vmess.ListenPort)), 1)
			}
		}
	}

	f, err := os.Create("./reverse.conf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(reverseConfig)
}
