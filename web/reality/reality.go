package reality

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/LalatinaHub/LatinaServer/config"
	"github.com/gin-gonic/gin"
)

var (
	domain       = os.Getenv("DOMAIN")
	realityRegex = regexp.MustCompile("reality")
)

func RealityHandler(c *gin.Context) {
	var (
		singConfig = config.ReadSingConfig()
		text       = []string{}
	)

	text = append(text, "REALITY SERVER INFORMATION")
	text = append(text, "--------------------------")
	text = append(text, "VPN Type : VLESS")
	text = append(text, "Reality Public Key : "+config.RealityPublicKey)
	text = append(text, "Reality ShortID : "+config.RealityShortID[0])
	text = append(text, "")
	text = append(text, "Example :")
	text = append(text, fmt.Sprintf("vless://00000000-0000-0000-0000-000000000000@%s:52005/?type=tcp&encryption=none&flow=&sni=meet.google.com&allowInsecure=1&fp=random&security=reality&pbk=%s&sid=%s#REALITY", domain, config.RealityPublicKey, config.RealityShortID[0]))
	text = append(text, "")
	text = append(text, "")
	text = append(text, "SNI AND PORT BINDING")
	text = append(text, "--------------------")

	for i, inbound := range singConfig.Inbounds {
		if realityRegex.MatchString(inbound.Tag) {
			tag := strings.Split(inbound.Tag, "-")
			text = append(text, fmt.Sprintf("%s : %d", strings.Join(tag[2:], "-"), 52000+i))
		}
	}

	c.String(http.StatusOK, strings.Join(text, "\n"))
}
