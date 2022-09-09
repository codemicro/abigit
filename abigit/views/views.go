package views

import (
	"github.com/codemicro/abigit/abigit/config"
	"github.com/gofiber/fiber/v2"
)

//go:generate neontc --extension ntc.html

type RenderContext struct {
	pageTitle   string
	externalURL string
}

func NewRenderContext() *RenderContext {
	rctx := new(RenderContext)
	rctx.externalURL = config.HTTP.ExternalURL
	return rctx
}

func SendPage(ctx *fiber.Ctx, content string) error {
	ctx.Type("html")
	return ctx.SendString(content)
}
