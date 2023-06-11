package relay

import (
	supabase "github.com/LalatinaHub/LatinaServer/db"
	"github.com/LalatinaHub/LatinaServer/helper"
	db "github.com/LalatinaHub/LatinaSub-go/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func GetRelayOutbounds() []option.Outbound {
	var (
		proxies      []db.DBScheme
		outbounds    = []option.Outbound{}
		outboundsMap = map[string][]option.Outbound{}
	)

	supabase.Connect().DB.From("proxies").Select("*").Eq("vpn", "vmess").Eq("conn_mode", "cdn").Eq("tls", "0").Eq("transport", "ws").Neq("country_code", helper.GetIpInfo().CountryCode).Execute(&proxies)

	for _, proxy := range proxies {
		outboundsMap[proxy.CountryCode] = append(outboundsMap[proxy.CountryCode], option.Outbound{
			Tag:  proxy.Remark,
			Type: proxy.VPN,
			VMessOptions: option.VMessOutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     "172.67.73.39",
					ServerPort: uint16(proxy.ServerPort),
				},
				UUID:     proxy.UUID,
				Security: proxy.Security,
				AlterId:  proxy.AlterId,
				Transport: &option.V2RayTransportOptions{
					Type: C.V2RayTransportTypeWebsocket,
					WebsocketOptions: option.V2RayWebsocketOptions{
						Path: proxy.Path,
						Headers: map[string]option.Listable[string]{
							"Host": {proxy.Server},
						},
					},
				},
			},
		})
	}

	for cc, out := range outboundsMap {
		urltest := option.Outbound{
			Tag:  cc,
			Type: C.TypeURLTest,
			URLTestOptions: option.URLTestOutboundOptions{
				Outbounds: []string{},
			},
		}

		for _, outbound := range out {
			urltest.URLTestOptions.Outbounds = append(urltest.URLTestOptions.Outbounds, outbound.Tag)
		}
		outbounds = append(outbounds, urltest)
		outbounds = append(outbounds, out...)
	}

	return outbounds
}
