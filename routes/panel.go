package routes

import (
	handlers "davet.link/handlers/panel"
	"davet.link/middlewares"
	"davet.link/models"

	"github.com/gofiber/fiber/v2"
)

func registerPanelRoutes(app *fiber.App) {
	panelGroup := app.Group("/panel")
	panelGroup.Use(
		middlewares.AuthMiddleware,
		middlewares.StatusMiddleware,
		middlewares.TypeMiddleware(models.Panel),
		middlewares.VerifiedMiddleware,
	)

	panelGroup.Get("/home", handlers.PanelHomeHandler)

	panelInvitationHandler := handlers.NewPanelInvitationHandler()
	panelGroup.Get("/invitations", panelInvitationHandler.ListPanelInvitations)
	panelGroup.Get("/invitations/create", panelInvitationHandler.ShowCreatePanelInvitation)
	panelGroup.Post("/invitations/create", panelInvitationHandler.CreatePanelInvitation)
	panelGroup.Get("/invitations/update/:id", panelInvitationHandler.ShowUpdatePanelInvitation)
	panelGroup.Post("/invitations/update/:id", panelInvitationHandler.UpdatePanelInvitation)
	panelGroup.Delete("/invitations/delete/:id", panelInvitationHandler.DeletePanelInvitation)

	panelCardHandler := handlers.NewPanelCardHandler()
	panelGroup.Get("/cards", panelCardHandler.ListPanelCards)
	panelGroup.Get("/cards/create", panelCardHandler.ShowCreatePanelCard)
	panelGroup.Post("/cards/create", panelCardHandler.CreatePanelCard)
	panelGroup.Get("/cards/update/:id", panelCardHandler.ShowUpdatePanelCard)
	panelGroup.Post("/cards/update/:id", panelCardHandler.UpdatePanelCard)
	panelGroup.Delete("/cards/delete/:id", panelCardHandler.DeletePanelCard)
}
