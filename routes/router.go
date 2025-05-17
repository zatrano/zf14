package routes

import (
	"davet.link/configs/sessionconfig"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	app.Use(logger.New())

	sessionStore := sessionconfig.SetupSession()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("session", sessionStore)
		return c.Next()
	})

	registerWebsiteRoutes(app)
	registerAuthRoutes(app)
	registerDashboardRoutes(app)
	registerPanelRoutes(app)
}
