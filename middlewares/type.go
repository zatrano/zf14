package middlewares

import (
	"davet.link/configs/sessionconfig"
	"davet.link/models"
	"davet.link/pkg/flashmessages"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
)

func TypeMiddleware(requiredType models.UserType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := sessionconfig.SessionStart(c)
		if err != nil {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Oturum başlatılamadı")
			return c.Redirect("/auth/login")
		}

		userID, err := sessionconfig.GetUserIDFromSession(sess)
		if err != nil {
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Yetkili oturum bulunamadı")
			return c.Redirect("/auth/login")
		}

		authService := services.NewAuthService()
		user, err := authService.GetUserProfile(userID)
		if err != nil {
			_ = sess.Destroy()
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı bulunamadı")
			return c.Redirect("/auth/login")
		}

		if user.Type != requiredType {
			_ = sess.Destroy()
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Bu sayfaya erişim izniniz yok")
			return c.Redirect("/auth/login")
		}

		return c.Next()
	}
}
