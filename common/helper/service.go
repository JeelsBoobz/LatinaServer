package helper

import (
	"fmt"
	"os/exec"
	"strings"
)

func ReloadService(names ...string) {
	for _, name := range names {
		fmt.Println("Reloading", name, "...")
		out, err := exec.Command("systemctl", "reload", name).Output()
		if strings.HasSuffix(string(out), "cannot reload.") {
			exec.Command("systemctl", "restart", name)
		} else if err != nil {
			panic(err)
		}

		fmt.Println(name, "successfully reloaded !")
	}
}
