package helper

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func ReloadService(names ...string) {
	for _, name := range names {
		fmt.Println("Reloading", name, "...")
		out, err := exec.Command("systemctl", "reload", name).Output()
		var ignoredText = []string{
			"reload.",
			"inactive.",
		}
		if m, _ := regexp.MatchString(regexp.MustCompile(strings.Join(ignoredText[:], "|")).String(), string(out)); m {
			exec.Command("systemctl", "start", name)
		} else if err != nil {
			panic(err)
		}

		fmt.Println(name, "successfully reloaded !")
	}
}
