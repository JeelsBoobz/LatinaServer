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
		serverCC     = helper.GetIpInfo().CountryCode
	)

	supabase.Connect().DB.From("proxies").Select("*").Eq("vpn", "shadowsocks").Eq("region", "Asia").Execute(&proxies)

	for _, proxy := range proxies {
		if len(outboundsMap[proxy.CountryCode]) < 3 && serverCC != proxy.CountryCode {
			outboundsMap[proxy.CountryCode] = append(outboundsMap[proxy.CountryCode], option.Outbound{
				Tag:  proxy.Remark,
				Type: proxy.VPN,
				ShadowsocksOptions: option.ShadowsocksOutboundOptions{
					ServerOptions: option.ServerOptions{
						Server:     proxy.Server,
						ServerPort: uint16(proxy.ServerPort),
					},
					Method:   proxy.Method,
					Password: proxy.Password,
				},
			})
		}
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
