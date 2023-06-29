package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/LalatinaHub/LatinaServer/config"
	"github.com/LalatinaHub/LatinaServer/config/relay"
	CS "github.com/LalatinaHub/LatinaServer/constant"
	"github.com/LalatinaHub/LatinaServer/db"
	"github.com/LalatinaHub/LatinaServer/helper"
	"github.com/LalatinaHub/LatinaServer/web"
	"github.com/go-co-op/gocron"
)

var (
	loc, _ = time.LoadLocation("Asia/Jakarta")
)

func HotReload() {
	config.Write()
	helper.ReloadService([]string{CS.ServiceSingBox, CS.ServiceOpenresty}...)
}

func UpdateUsersQuota() {
	var isAnyExceed bool = false
	for _, user := range config.ReadSingConfig().Experimental.V2RayAPI.Stats.Users {
		if !db.UpdatePremiumQuota(user) {
			isAnyExceed = true
		}

		time.Sleep(500 * time.Millisecond)
	}

	if isAnyExceed {
		HotReload()
	}
}

func main() {
	fmt.Println("Service started !")
	s := gocron.NewScheduler(loc)

	s.Every(1).Day().At("00:00").Tag("hot-reload").Do(HotReload)
	s.Every(30).Minutes().Tag("get-relays").Do(relay.GatherRelays)
	s.Every(5).Minutes().Tag("update-quota").Do(UpdateUsersQuota)
	s.Every(1).Day().At("03:00").Tag("reboot").Do(func() {
		if err := exec.Command("reboot").Run(); err != nil {
			fmt.Println("Failed to reboot the server !", err)
		}
	})

	helper.GetIpInfo()
	HotReload()
	s.StartAsync()

	web.StartWebService()
}
