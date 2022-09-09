package endpoints

import (
	"github.com/codemicro/abigit/abigit/urls"
	"github.com/gofiber/fiber/v2"
)

func (e *Endpoints) Index(ctx *fiber.Ctx) error {
	return ctx.Redirect(urls.Make(urls.AuthOIDCOutbound))
}
