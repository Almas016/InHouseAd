package handlers

import (
	"InHouseAd/internal/app"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	wc *app.WebsiteChecker
}

func NewHandler(wc *app.WebsiteChecker) *Handler {
	return &Handler{wc: wc}
}

func (h *Handler) AccessTime(c *fiber.Ctx) error {
	url := c.Params("url")
	url = "http://" + url
	accessTime, err := h.wc.GetAccessTime(url)
	if err != nil {
		c.Status(404)
	}
	return c.JSON(accessTime)
}

func (h *Handler) MinAccessTime(c *fiber.Ctx) error {
	minWeb := h.wc.GetMinAccessTime()
	return c.JSON(minWeb)
}

func (h *Handler) MaxAccessTime(c *fiber.Ctx) error {
	maxWeb := h.wc.GetMaxAccessTime()
	return c.JSON(maxWeb)
}
