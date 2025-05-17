package handlers

import (
	"net/http"

	"davet.link/pkg/renderer"

	"github.com/gofiber/fiber/v2"
)

type WebsiteHandler struct {
}

func NewWebsiteHandler() *WebsiteHandler {
	return &WebsiteHandler{}
}

func (h *WebsiteHandler) ShowHomePage(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/home", "layouts/website", mapData, http.StatusOK)
}

func (h *WebsiteHandler) ShowTermsOfUsePage(c *fiber.Ctx) error {
	mapData := fiber.Map{}
	return renderer.Render(c, "website/terms-of-use", "layouts/website", mapData, http.StatusOK)
}
