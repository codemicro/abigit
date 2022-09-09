package static

import (
	"embed"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"
)

//go:embed static/*
var Static embed.FS

func NewHandler() fiber.Handler {
	return filesystem.New(filesystem.Config{
		Root:       http.FS(Static),
		PathPrefix: "static",
	})
}
