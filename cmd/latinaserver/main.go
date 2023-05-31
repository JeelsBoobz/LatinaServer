package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	config "github.com/LalatinaHub/LatinaServer/common/config/sing-box"
	"github.com/go-co-op/gocron"
)

var (
	loc, _ = time.LoadLocation("Asia/Jakarta")
)

func reloadSingBox() {
	fmt.Println("Reloading sing-box...")
	processes, err := exec.Command("ps", "-e", "-o", "pid,comm").Output()
	if err != nil {
		panic(err)
	}

	for _, process := range strings.Split(string(processes), "\n") {
		fields := strings.Fields(process)
		if len(fields) > 1 && fields[1] == "sing-box" {
			pid, err := strconv.Atoi(fields[0])
			if err != nil {
				panic(err)
			}

			_, err = exec.Command("kill", "-HUP", strconv.Itoa(pid)).Output()
			if err != nil {
				panic(err)
			}

			fmt.Println("Sing-box reloaded !")
		}
	}
}

func HotReload() {
	config.ReadAndWriteConfig()
	reloadSingBox()
}

func main() {
	fmt.Println("Service started !")
	s := gocron.NewScheduler(loc)

	s.Every(1).Day().At("00:00").Tag("hot-reload").Do(HotReload)

	HotReload()
	s.StartBlocking()
}
