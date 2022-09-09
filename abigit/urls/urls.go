package urls

import (
	"fmt"
	"github.com/codemicro/abigit/abigit/config"
	"strings"
)

const (
	Index = "/"

	Auth             = "/auth"
	AuthOIDC         = Auth + "/oidc"
	AuthOIDCInbound  = AuthOIDC + "/inbound"
	AuthOIDCOutbound = AuthOIDC + "/outbound"
)

func Make(template string, replacements ...interface{}) string {
	spt := strings.Split(template, "/")
	for i, part := range spt {
		if len(part) == 0 {
			continue
		}
		if part[0] == ':' {
			spt[i] = "%s"
		}
	}
	return config.HTTP.ExternalURL + fmt.Sprintf(strings.Join(spt, "/"), replacements...)
}
