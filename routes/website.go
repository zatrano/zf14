package routes

import (
	handlers "davet.link/handlers/website"

	"github.com/gofiber/fiber/v2"
)

func registerWebsiteRoutes(app *fiber.App) {
	app.Get("/", handlers.NewWebsiteHandler().ShowHomePage)
	app.Get("/kullanim-sartlari", handlers.NewWebsiteHandler().ShowTermsOfUsePage)
}
