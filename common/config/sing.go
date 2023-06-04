package config

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/LalatinaHub/LatinaServer/common/db"
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
	options := ReadSingConfig()

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
