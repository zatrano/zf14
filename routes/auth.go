package routes

import (
	handlers "davet.link/handlers/auth"
	"davet.link/middlewares"
	"davet.link/requests"

	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(app *fiber.App) {
	authHandler := handlers.NewAuthHandler()

	authGroup := app.Group("/auth")

	authGroup.Get("/login", middlewares.GuestMiddleware, authHandler.ShowLogin)
	authGroup.Post("/login", middlewares.GuestMiddleware, requests.ValidateLoginRequest, authHandler.Login)

	authGroup.Get("/logout", middlewares.AuthMiddleware, authHandler.Logout)
	authGroup.Get("/profile", middlewares.AuthMiddleware, authHandler.Profile)
	authGroup.Post("/profile/update-password", middlewares.AuthMiddleware, requests.ValidateUpdatePasswordRequest, authHandler.UpdatePassword)
	authGroup.Get("/register", middlewares.GuestMiddleware, authHandler.ShowRegister)
	authGroup.Post("/register", middlewares.GuestMiddleware, requests.ValidateRegisterRequest, authHandler.Register)
	authGroup.Get("/forgot-password", middlewares.GuestMiddleware, authHandler.ShowForgotPassword)
	authGroup.Post("/forgot-password", middlewares.GuestMiddleware, requests.ValidateForgotPasswordRequest, authHandler.ForgotPassword)
	authGroup.Get("/reset-password", middlewares.GuestMiddleware, authHandler.ShowResetPassword)
	authGroup.Post("/reset-password", middlewares.GuestMiddleware, requests.ValidateResetPasswordRequest, authHandler.ResetPassword)
	authGroup.Get("/verify-email", middlewares.GuestMiddleware, authHandler.VerifyEmail)
	authGroup.Get("/resend-verification", middlewares.GuestMiddleware, authHandler.ShowResendVerification)
authGroup.Post("/resend-verification", middlewares.GuestMiddleware, authHandler.ResendVerification)
	authGroup.Get("/google/login", handlers.GoogleLogin)
	authGroup.Get("/google/callback", handlers.GoogleCallback)
}
