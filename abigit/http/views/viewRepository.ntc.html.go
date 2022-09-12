// Code generated by Neon - DO NOT EDIT
// https://github.com/codemicro/go-neon

package views

import (
	"fmt"
	"github.com/codemicro/abigit/abigit/core"
	"github.com/codemicro/abigit/abigit/urls"
	ntcWzvIICEQKo "html"
	"path/filepath"
	ntctyoxNCvLze "strings"
)

type ViewRepositoryProps struct {
	Repo    *core.RepoOnDisk
	IsEmpty bool
}

func ViewRepository(ctx *RenderContext, props *ViewRepositoryProps) string {
	ntcRijqFUuoGg := new(ntctyoxNCvLze.Builder)
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    ")
	ctx.pageTitle = props.Repo.Slug + " | " + ctx.platformName

	_, _ = ntcRijqFUuoGg.WriteString("\n\n    ")
	_, _ = ntcRijqFUuoGg.WriteString(PageBegin(ctx))
	_, _ = ntcRijqFUuoGg.WriteString("\n    ")
	_, _ = ntcRijqFUuoGg.WriteString(Navbar(ctx))
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    <div class=\"container\">\n        <h1>")
	_, _ = ntcRijqFUuoGg.WriteString(ntcWzvIICEQKo.EscapeString(props.Repo.Slug))
	_, _ = ntcRijqFUuoGg.WriteString("</h1>\n        <p class=\"secondary\">")
	_, _ = ntcRijqFUuoGg.WriteString(ntcWzvIICEQKo.EscapeString(props.Repo.Description))
	_, _ = ntcRijqFUuoGg.WriteString("</p>\n        <p>Size on disk: ")
	_, _ = ntcRijqFUuoGg.WriteString(ntcWzvIICEQKo.EscapeString(formatFileSize(props.Repo.Size)))
	_, _ = ntcRijqFUuoGg.WriteString("</p>\n\n<!--        \\{{ if props.IsEmpty \\}}-->\n<!--            <div class=\"card full-width\">-->\n<!--                <span>This repository is uninitialised. Try creating and pushing a commit, then try again.</span>-->\n<!--            </div>-->\n<!--        \\{{ else \\}}-->\n            <div id=\"tabs\" class=\"tabs\"\n                 hx-get=\"")
	_, _ = ntcRijqFUuoGg.WriteString(ntcWzvIICEQKo.EscapeString(urls.Make(urls.RepositoryTabs, props.Repo.Slug)))
	_, _ = ntcRijqFUuoGg.WriteString("\"\n                 hx-trigger=\"load\"\n                 hx-target=\"#tabs\"\n                 hx-swap=\"innerHTML\"\n            ></div>\n<!--        \\{{ endif \\}}-->\n    </div>\n\n    ")
	_, _ = ntcRijqFUuoGg.WriteString(PageEnd(ctx))
	_, _ = ntcRijqFUuoGg.WriteString("\n\n")
	return ntcRijqFUuoGg.String()
}

const (
	TabSelectorReadme = 1 << iota
	TabSelectorShowTree
	TabSelectorShowRefs
	TabSelectorClone
)

type RepositoryTabProps struct {
	Repo       *core.RepoOnDisk
	ShowTabs   int
	CurrentTab int

	Readme struct {
		Content string
	}

	Clone struct {
		SSHUser        string
		SSHHost        string
		SSHStoragePath string
	}
}

func RepositoryTabs(ctx *RenderContext, props *RepositoryTabProps) string {
	ntcRijqFUuoGg := new(ntctyoxNCvLze.Builder)
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    ")
	selectedTab := props.CurrentTab

	if props.ShowTabs&selectedTab == 0 {
		for selectedTab&0b1 == 0 {
			if selectedTab == 0 {
				break
			}

			selectedTab = selectedTab >> 1
		}
	}

	_, _ = ntcRijqFUuoGg.WriteString("\n\n    <div class=\"tab-list\">\n\n        ")
	if props.ShowTabs&TabSelectorReadme != 0 {
		_, _ = ntcRijqFUuoGg.WriteString("\n            <a hx-get=\"")
		_, _ = ntcRijqFUuoGg.WriteString(ntcWzvIICEQKo.EscapeString(urls.Make(urls.RepositoryTabReadme, props.Repo.Slug)))
		_, _ = ntcRijqFUuoGg.WriteString("\"\n                ")
		if selectedTab == TabSelectorReadme {
			_, _ = ntcRijqFUuoGg.WriteString("\n                     class=\"selected\"\n                ")
		}
		_, _ = ntcRijqFUuoGg.WriteString("\n            >README</a>\n        ")
	}
	_, _ = ntcRijqFUuoGg.WriteString("\n\n        ")
	if props.ShowTabs&TabSelectorShowTree != 0 {
		_, _ = ntcRijqFUuoGg.WriteString("\n            <a hx-get=\"")
		_, _ = ntcRijqFUuoGg.WriteString(ntcWzvIICEQKo.EscapeString(urls.Make(urls.RepositoryTabTree, props.Repo.Slug)))
		_, _ = ntcRijqFUuoGg.WriteString("\"\n                ")
		if selectedTab == TabSelectorShowTree {
			_, _ = ntcRijqFUuoGg.WriteString("\n                    class=\"selected\"\n                ")
		}
		_, _ = ntcRijqFUuoGg.WriteString("\n            >Tree</a>\n        ")
	}
	_, _ = ntcRijqFUuoGg.WriteString("\n\n        ")
	if props.ShowTabs&TabSelectorShowRefs != 0 {
		_, _ = ntcRijqFUuoGg.WriteString("\n            <a hx-get=\"")
		_, _ = ntcRijqFUuoGg.WriteString(ntcWzvIICEQKo.EscapeString(urls.Make(urls.RepositoryTabRefs, props.Repo.Slug)))
		_, _ = ntcRijqFUuoGg.WriteString("\"\n                ")
		if selectedTab == TabSelectorShowRefs {
			_, _ = ntcRijqFUuoGg.WriteString("\n                    class=\"selected\"\n                ")
		}
		_, _ = ntcRijqFUuoGg.WriteString("\n            >Refs</a>\n        ")
	}
	_, _ = ntcRijqFUuoGg.WriteString("\n\n        ")
	if props.ShowTabs&TabSelectorClone != 0 {
		_, _ = ntcRijqFUuoGg.WriteString("\n            <a hx-get=\"")
		_, _ = ntcRijqFUuoGg.WriteString(ntcWzvIICEQKo.EscapeString(urls.Make(urls.RepositoryTabClone, props.Repo.Slug)))
		_, _ = ntcRijqFUuoGg.WriteString("\"\n                ")
		if selectedTab == TabSelectorClone {
			_, _ = ntcRijqFUuoGg.WriteString("\n                    class=\"selected\"\n                ")
		}
		_, _ = ntcRijqFUuoGg.WriteString("\n            >Clone</a>\n        ")
	}
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    </div>\n\n    <div class=\"tab-content\">\n        ")
	if selectedTab == TabSelectorReadme {
		_, _ = ntcRijqFUuoGg.WriteString("\n            ")
		_, _ = ntcRijqFUuoGg.WriteString(repositoryTabReadme(ctx, props))
		_, _ = ntcRijqFUuoGg.WriteString("\n        ")
	} else if selectedTab == TabSelectorShowTree {
		_, _ = ntcRijqFUuoGg.WriteString("\n            ")
		_, _ = ntcRijqFUuoGg.WriteString(repositoryTabTree(ctx, props))
		_, _ = ntcRijqFUuoGg.WriteString("\n        ")
	} else if selectedTab == TabSelectorShowRefs {
		_, _ = ntcRijqFUuoGg.WriteString("\n            ")
		_, _ = ntcRijqFUuoGg.WriteString(repositoryTabRefs(ctx, props))
		_, _ = ntcRijqFUuoGg.WriteString("\n        ")
	} else if selectedTab == TabSelectorClone {
		_, _ = ntcRijqFUuoGg.WriteString("\n            ")
		_, _ = ntcRijqFUuoGg.WriteString(repositoryTabClone(ctx, props))
		_, _ = ntcRijqFUuoGg.WriteString("\n        ")
	}
	_, _ = ntcRijqFUuoGg.WriteString("\n    </div>\n\n")
	return ntcRijqFUuoGg.String()
}
func repositoryTabReadme(ctx *RenderContext, props *RepositoryTabProps) string {
	ntcRijqFUuoGg := new(ntctyoxNCvLze.Builder)
	_, _ = ntcRijqFUuoGg.WriteString("\n    <div class=\"readme-content\">")
	_, _ = ntcRijqFUuoGg.WriteString(props.Readme.Content)
	_, _ = ntcRijqFUuoGg.WriteString("</div>\n")
	return ntcRijqFUuoGg.String()
}
func repositoryTabTree(ctx *RenderContext, props *RepositoryTabProps) string {
	ntcRijqFUuoGg := new(ntctyoxNCvLze.Builder)
	_, _ = ntcRijqFUuoGg.WriteString("\n    <p>TODO: tree here</p>\n")
	return ntcRijqFUuoGg.String()
}
func repositoryTabRefs(ctx *RenderContext, props *RepositoryTabProps) string {
	ntcRijqFUuoGg := new(ntctyoxNCvLze.Builder)
	_, _ = ntcRijqFUuoGg.WriteString("\n    <p>TODO: refs here</p>\n")
	return ntcRijqFUuoGg.String()
}
func repositoryTabClone(ctx *RenderContext, props *RepositoryTabProps) string {
	ntcRijqFUuoGg := new(ntctyoxNCvLze.Builder)
	_, _ = ntcRijqFUuoGg.WriteString("\n    <p>Clone with:</p>\n    <ul>\n        <li>SSH (read/write): ")
	_, _ = ntcRijqFUuoGg.WriteString(ntcWzvIICEQKo.EscapeString(fmt.Sprintf("%s@%s:%s", props.Clone.SSHUser, props.Clone.SSHHost, filepath.Join(props.Clone.SSHStoragePath, props.Repo.Slug))))
	_, _ = ntcRijqFUuoGg.WriteString("</li>\n        <li>HTTP (read-only): ")
	_, _ = ntcRijqFUuoGg.WriteString(ntcWzvIICEQKo.EscapeString(urls.Make(urls.ServeRepositoryByName, props.Repo.Slug+".git")))
	_, _ = ntcRijqFUuoGg.WriteString("</li>\n    </ul>\n\n    <p>// TODO: Remember about securing private repositories.</p>\n")
	return ntcRijqFUuoGg.String()
}
