package middlewares

import (
	"davet.link/configs/sessionconfig"
	"davet.link/pkg/flashmessages"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
)

func StatusMiddleware(c *fiber.Ctx) error {
	sess, err := sessionconfig.SessionStart(c)
	if err != nil {
		return c.Redirect("/auth/login")
	}

	userID, err := sessionconfig.GetUserIDFromSession(sess)
	if err != nil {
		return c.Redirect("/auth/login")
	}

	authService := services.NewAuthService()
	user, err := authService.GetUserProfile(userID)
	if err != nil {
		_ = sess.Destroy()
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı bulunamadı")
		return c.Redirect("/auth/login")
	}

	if !user.Status {
		_ = sess.Destroy()
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı durumu geçersiz")
		return c.Redirect("/auth/login")
	}

	return c.Next()
}
