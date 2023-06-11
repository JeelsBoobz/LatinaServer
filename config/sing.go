package config

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func ReadSingConfig() option.Options {
	body, err := os.ReadFile("/usr/local/etc/latinaserver/config.json")
	if err != nil {
		panic(err)
	}

	var options option.Options
	err = options.UnmarshalJSON(body)
	if err != nil {
		panic(err)
	}

	return options
}

func WriteSingConfig() option.Options {
	premiumList := db.GetPremiumList()
	relayOutbounds := relay.GetRelayOutbounds()
	options := ReadSingConfig()
	options.Experimental = &option.ExperimentalOptions{
		ClashAPI: &option.ClashAPIOptions{
			ExternalController: "127.0.0.1:9090",
			ExternalUI:         "/usr/local/latinaserver/dashboard/",
		},
		V2RayAPI: &option.V2RayAPIOptions{
			Listen: CS.V2rayAPIAddress,
			Stats: &option.V2RayStatsServiceOptions{
				Enabled:   true,
				Inbounds:  []string{},
				Outbounds: []string{"direct"},
				Users:     []string{},
			},
		},
	}

	var inbounds []option.Inbound
	for i, inbound := range options.Inbounds {
		var (
			port = 52000 + i
		)

		switch inbound.Type {
		case C.TypeTrojan:
			inbound.TrojanOptions.ListenPort = uint16(port)
			inbound.TrojanOptions.Users = []option.TrojanUser{}

			for _, user := range premiumList[C.TypeTrojan] {
				inbound.TrojanOptions.Users = append(inbound.TrojanOptions.Users, option.TrojanUser{
					Name:     strconv.Itoa(int(user.Id)),
					Password: user.Password,
				})
			}

			if inbound.TrojanOptions.Transport != nil {
				switch inbound.TrojanOptions.Transport.Type {
				case C.V2RayTransportTypeWebsocket:
					inbound.TrojanOptions.Transport.WebsocketOptions.Path = "/" + inbound.Type
				case C.V2RayTransportTypeGRPC:
					inbound.TrojanOptions.Transport.GRPCOptions.ServiceName = inbound.Type
				}
			}
		case C.TypeVMess:
			inbound.VMessOptions.ListenPort = uint16(port)
			inbound.VMessOptions.Users = []option.VMessUser{}

			for _, user := range premiumList[C.TypeVMess] {
				inbound.VMessOptions.Users = append(inbound.VMessOptions.Users, option.VMessUser{
					Name: strconv.Itoa(int(user.Id)),
					UUID: user.Password,
				})
			}

			if inbound.VMessOptions.Transport != nil {
				switch inbound.VMessOptions.Transport.Type {
				case C.V2RayTransportTypeWebsocket:
					inbound.VMessOptions.Transport.WebsocketOptions.Path = "/" + inbound.Type
				case C.V2RayTransportTypeGRPC:
					inbound.VMessOptions.Transport.GRPCOptions.ServiceName = inbound.Type
				}
			}
		case C.TypeVLESS:
			inbound.VLESSOptions.ListenPort = uint16(port)
			inbound.VLESSOptions.Users = []option.VLESSUser{}

			for _, user := range premiumList[C.TypeVLESS] {
				inbound.VLESSOptions.Users = append(inbound.VLESSOptions.Users, option.VLESSUser{
					Name: strconv.Itoa(int(user.Id)),
					UUID: user.Password,
				})
			}

			if inbound.VLESSOptions.Transport != nil {
				switch inbound.VLESSOptions.Transport.Type {
				case C.V2RayTransportTypeWebsocket:
					inbound.VLESSOptions.Transport.WebsocketOptions.Path = "/" + inbound.Type
				case C.V2RayTransportTypeGRPC:
					inbound.VLESSOptions.Transport.GRPCOptions.ServiceName = inbound.Type
				}
			}
		}

		inbounds = append(inbounds, inbound)
		options.Experimental.V2RayAPI.Stats.Inbounds = append(options.Experimental.V2RayAPI.Stats.Inbounds, inbound.Tag)
	}

	for _, list := range premiumList {
		for _, user := range list {
			options.Experimental.V2RayAPI.Stats.Users = append(options.Experimental.V2RayAPI.Stats.Users, strconv.Itoa(int(user.Id)))
		}
	}

	options.Inbounds = inbounds

	options.Outbounds = []option.Outbound{
		{
			Type: C.TypeDirect,
			Tag:  "direct",
		},
		{
			Type: C.TypeDNS,
			Tag:  "dns-out",
		},
	}
	options.Outbounds = append(options.Outbounds, relayOutbounds...)

	options.Route.Final = "direct"
	options.Route.Rules = []option.Rule{
		{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				Protocol: option.Listable[string]{"dns"},
				Outbound: "dns-out",
			},
		},
		{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				IPCIDR:   option.Listable[string]{"1.1.1.1", "8.8.8.8"},
				Outbound: "direct",
			},
		},
		{
			Type: C.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				Port:     option.Listable[uint16]{53},
				Outbound: "direct",
			},
		},
	}
	for _, outbound := range relayOutbounds {
		if len(outbound.Tag) < 5 {
			rule := option.Rule{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					AuthUser: []string{},
					Outbound: outbound.Tag,
				},
			}

			for _, premium := range premiumList {
				for _, user := range premium {
					if user.CC == outbound.Tag {
						rule.DefaultOptions.AuthUser = append(rule.DefaultOptions.AuthUser, strconv.Itoa(int(user.Id)))
					}
				}
			}

			if len(rule.DefaultOptions.AuthUser) > 0 {
				options.Route.Rules = append(options.Route.Rules, rule)
			}
		}
	}

	// Write new config
	f, err := os.Create("/usr/local/etc/latinaserver/config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := json.MarshalIndent(options, "", "\t")
	if err != nil {
		panic(err)
	}
	f.WriteString(string(b))

	return options
}
