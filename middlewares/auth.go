package middlewares

import (
	"context"

	"davet.link/configs/sessionconfig"
	"davet.link/pkg/flashmessages"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	sess, err := sessionconfig.SessionStart(c)
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Oturum başlatılamadı")
		return c.Redirect("/auth/login")
	}

	userID, err := sessionconfig.GetUserIDFromSession(sess)
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Oturum bilgileri geçersiz")
		return c.Redirect("/auth/login")
	}

	authService := services.NewAuthService()
	user, err := authService.GetUserProfile(userID)
	if err != nil {
		_ = sess.Destroy()
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı bulunamadı")
		return c.Redirect("/auth/login")
	}

	ctx := context.WithValue(c.Context(), "user_id", userID)
	ctx = context.WithValue(ctx, "user_type", user.Type)
	ctx = context.WithValue(ctx, "user_account", user.Account)
	c.SetUserContext(ctx)

	c.Locals("userID", userID)
	c.Locals("userType", user.Type)
	c.Locals("userAccount", user.Account)

	return c.Next()
}
