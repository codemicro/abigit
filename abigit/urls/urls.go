package urls

import (
	"fmt"
	"github.com/codemicro/abigit/abigit/config"
	"strings"
)

const (
	Index = "/"

	Auth = "/auth"

	AuthLogin = Auth + "/login"

	AuthOIDC         = Auth + "/oidc"
	AuthOIDCInbound  = AuthOIDC + "/inbound"
	AuthOIDCOutbound = AuthOIDC + "/outbound"

	CreateRepository           = "/create"
	CreateRepositoryValidation = "/create/validate"

	Repositories          = "/~"
	RepositoryByName      = Repositories + "/:slug"
	ServeRepositoryByName = Repositories + "/:slug"

	RepositoryTabs      = RepositoryByName + "/tab"
	RepositoryTabReadme = RepositoryTabs + "/?tab=readme"
	RepositoryTabTree   = RepositoryTabs + "/?tab=tree"
	RepositoryTabRefs   = RepositoryTabs + "/?tab=refs"
)

func Make(template string, replacements ...interface{}) string {
	return config.HTTP.ExternalURL + MakeRelative(template, replacements...)
}

func MakeRelative(template string, replacements ...any) string {
	spt := strings.Split(template, "/")
	for i, part := range spt {
		if len(part) == 0 {
			continue
		}
		if part[0] == ':' {
			spt[i] = "%s"
		}
	}
	return fmt.Sprintf(strings.Join(spt, "/"), replacements...)
}
