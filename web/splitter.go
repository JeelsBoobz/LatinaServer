package web

import (
	"fmt"
	"os"

	"github.com/LalatinaHub/LatinaServer/config"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	splitterConfig "github.com/LalatinaHub/tls-splitter/config/raw"
	C "github.com/sagernet/sing-box/constant"
)

var (
	splitterConfig_ = splitterConfig.RawConfig{
		Listen:             "0.0.0.0:443",
		RedirectHttps:      "",
		InboundBufferSize:  4,
		OutboundBufferSize: 32,
		VHosts: []splitterConfig.RawVHost{
			{
				Name:          os.Getenv("DOMAIN"),
				TlsOffloading: true,
				ManagedCert:   true,
				Alpn:          "http/1.1",
				Protocols:     "tls12,tls13",
				Http: splitterConfig.RawHttpHandler{
					Handler: "proxyPass",
					Args:    CS.ReverseProxyAddress,
				},
				Trojan: splitterConfig.RawHandler{
					Handler: "proxyPass",
					Args:    "",
				},
				Default: splitterConfig.RawHandler{
					Handler: "proxyPass",
					Args:    "",
				},
			},
		},
	}
)

func ReadSplitterConfig() splitterConfig.RawConfig {
	for _, inbound := range config.ReadSingConfig().Inbounds {
		if inbound.Type == C.TypeTrojan {
			trojan := inbound.TrojanOptions
			if trojan.Transport == nil {
				continue
			}

			switch trojan.Transport.Type {
			case C.V2RayTransportTypeWebsocket:
				splitterConfig_.VHosts[0].Trojan.Args = fmt.Sprintf("127.0.0.1:%d", int(trojan.ListenPort))
			}
		} else if inbound.Type == C.TypeVMess {
			vmess := inbound.VMessOptions
			if vmess.Transport == nil {
				continue
			}

			switch vmess.Transport.Type {
			case C.V2RayTransportTypeWebsocket:
				splitterConfig_.VHosts[0].Default.Args = fmt.Sprintf("127.0.0.1:%d", int(vmess.ListenPort))
			}
		}
	}

	fmt.Println(splitterConfig_)
	return splitterConfig_
}
