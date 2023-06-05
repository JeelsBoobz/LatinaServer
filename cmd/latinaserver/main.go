package main

import (
	"fmt"
	"time"

	"github.com/LalatinaHub/LatinaServer/config"
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
	for _, user := range config.ReadSingConfig().Experimental.V2RayAPI.Stats.Users {
		db.UpdatePremiumQuota(user)
	}
}

func main() {
	fmt.Println("Service started !")
	s := gocron.NewScheduler(loc)

	s.Every(1).Day().At("00:00").Tag("hot-reload").Do(HotReload)
	s.Every(5).Minutes().Tag("update-quota").Do(UpdateUsersQuota)

	HotReload()
	s.StartAsync()

	web.StartWebService()
}
