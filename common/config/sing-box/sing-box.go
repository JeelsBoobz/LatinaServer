package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LalatinaHub/LatinaServer/common/config/nginx"
	"github.com/LalatinaHub/LatinaServer/common/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

var (
	domain = os.Getenv("DOMAIN")
)

func ReadAndWriteConfig() option.Options {
	premiumList := db.GetPremiumList()
	body, err := os.ReadFile("/usr/local/etc/latinaserver/config.json")
	if err != nil {
		panic(err)
	}

	var options option.Options
	err = options.UnmarshalJSON(body)
	if err != nil {
		panic(err)
	}

	var inbounds []option.Inbound
	for i, inbound := range options.Inbounds {
		var (
			port       = 52000 + i
			serverName = fmt.Sprintf("%s.%s", inbound.Tag, domain)
		)

		switch inbound.Type {
		case C.TypeTrojan:
			inbound.TrojanOptions.ListenPort = uint16(port)
			inbound.TrojanOptions.Users = []option.TrojanUser{}

			for _, user := range premiumList[C.TypeTrojan] {
				inbound.TrojanOptions.Users = append(inbound.TrojanOptions.Users, option.TrojanUser{
					Name:     user.Name,
					Password: user.Password,
				})
			}

			if inbound.TrojanOptions.TLS != nil {
				inbound.TrojanOptions.TLS.ServerName = serverName
			}
		case C.TypeVMess:
			inbound.VMessOptions.ListenPort = uint16(port)
			inbound.VMessOptions.Users = []option.VMessUser{}

			for _, user := range premiumList[C.TypeVMess] {
				inbound.VMessOptions.Users = append(inbound.VMessOptions.Users, option.VMessUser{
					Name: user.Name,
					UUID: user.Password,
				})
			}

			if inbound.VMessOptions.TLS != nil {
				inbound.VMessOptions.TLS.ServerName = serverName
			}
		case C.TypeVLESS:
			inbound.VLESSOptions.ListenPort = uint16(port)
			inbound.VLESSOptions.Users = []option.VLESSUser{}

			for _, user := range premiumList[C.TypeVLESS] {
				inbound.VLESSOptions.Users = append(inbound.VLESSOptions.Users, option.VLESSUser{
					Name: user.Name,
					UUID: user.Password,
				})
			}

			if inbound.VLESSOptions.TLS != nil {
				inbound.VLESSOptions.TLS.ServerName = serverName
			}
		}

		inbounds = append(inbounds, inbound)
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

	// Generate nginx configuration
	nginx.GenerateReverseProxy(options)

	return options
}
