package main

import (
	"fmt"
	"time"

	"github.com/LalatinaHub/LatinaServer/common/config"
	"github.com/LalatinaHub/LatinaServer/common/helper"
	"github.com/go-co-op/gocron"
)

var (
	loc, _   = time.LoadLocation("Asia/Jakarta")
	services = []string{"sing-box", "openresty"}
)

func HotReload() {
	config.Write()
	helper.ReloadService(services...)
}

func main() {
	fmt.Println("Service started !")
	s := gocron.NewScheduler(loc)

	s.Every(1).Day().At("00:00").Tag("hot-reload").Do(HotReload)

	HotReload()
	s.StartBlocking()
}
