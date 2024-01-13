package main

import (
	"net/http"

	"github.com/elumbantoruan/feed/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/heptiolabs/healthcheck"
)

const healthCheckEndpoint = ":8086"

func main() {
	app := fiber.New()

	// Routes
	var handler web.Handler
	app.Get("/", handler.RenderFeedsRoute)
	app.Post("/update", handler.UpdateFeedRoute)

	health := healthcheck.NewHandler()

	go http.ListenAndServe(healthCheckEndpoint, health)

	app.Listen(":5000")
}
