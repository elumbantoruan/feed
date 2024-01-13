package main

import (
	"github.com/elumbantoruan/feed/pkg/web"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Routes
	var handler web.Handler
	app.Get("/", handler.RenderFeedsRoute)
	app.Post("/update", handler.UpdateFeedRoute)

	app.Listen(":5000")
}
