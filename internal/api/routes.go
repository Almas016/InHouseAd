package api

import (
	"InHouseAd/internal/api/handlers"
	"github.com/gofiber/fiber/v2"
)

func Routes(h *handlers.Handler) error {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  true,
			"message": "success",
		})
	})
	app.Get("/access/:url", h.AccessTime)
	app.Get("/min", h.MinAccessTime)
	app.Get("/max", h.MaxAccessTime)

	return app.Listen(":3000")
}
