// Code generated by Neon - DO NOT EDIT
// https://github.com/codemicro/go-neon

package views

import (
	"github.com/codemicro/abigit/abigit/core"
	ntcWRcOenDZxj "html"
	ntcSqmcWketfv "strings"
)

type IndexProps struct {
	AllRepos []*core.RepoOnDisk
}

func Index(ctx *RenderContext, props *IndexProps) string {
	ntcRijqFUuoGg := new(ntcSqmcWketfv.Builder)
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    ")
	ctx.pageTitle = "Welcome to " + ctx.platformName

	_, _ = ntcRijqFUuoGg.WriteString("\n\n    ")
	_, _ = ntcRijqFUuoGg.WriteString(PageBegin(ctx))
	_, _ = ntcRijqFUuoGg.WriteString("\n    ")
	_, _ = ntcRijqFUuoGg.WriteString(Navbar(ctx))
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    <div class=\"container\">\n        <h3>Recently updated</h3>\n\n        <h3>All repositories</h3>\n\n        ")
	if len(props.AllRepos) == 0 {
		_, _ = ntcRijqFUuoGg.WriteString("\n            <p class=\"secondary\">Nothing to see here!</p>\n        ")
	} else {
		_, _ = ntcRijqFUuoGg.WriteString("\n            ")
		for _, repo := range props.AllRepos {
			_, _ = ntcRijqFUuoGg.WriteString("\n            <p>")
			_, _ = ntcRijqFUuoGg.WriteString(ntcWRcOenDZxj.EscapeString(repo.Slug))
			_, _ = ntcRijqFUuoGg.WriteString(" - ")
			_, _ = ntcRijqFUuoGg.WriteString(ntcWRcOenDZxj.EscapeString(repo.Description))
			_, _ = ntcRijqFUuoGg.WriteString(" - ")
			_, _ = ntcRijqFUuoGg.WriteString(ntcWRcOenDZxj.EscapeString(repo.Path))
			_, _ = ntcRijqFUuoGg.WriteString("</p>\n            ")
		}
		_, _ = ntcRijqFUuoGg.WriteString("\n        ")
	}
	_, _ = ntcRijqFUuoGg.WriteString("\n\n    </div>\n\n    ")
	_, _ = ntcRijqFUuoGg.WriteString(PageEnd(ctx))
	_, _ = ntcRijqFUuoGg.WriteString("\n\n")
	return ntcRijqFUuoGg.String()
}
