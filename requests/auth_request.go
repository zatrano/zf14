package requests

import (
	"davet.link/pkg/flashmessages"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Account  string `form:"account" validate:"required,min=3"`
	Password string `form:"password" validate:"required,min=6"`
}

func ValidateLoginRequest(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		err := err.(validator.ValidationErrors)[0]
		switch {
		case err.Field() == "Account" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı adı zorunludur")
		case err.Field() == "Password" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifre zorunludur")
		case err.Field() == "Password" && err.Tag() == "min":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifre en az 6 karakter olmalıdır")
		default:
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz giriş bilgileri")
		}
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	c.Locals("loginRequest", req)
	return c.Next()
}

type UpdatePasswordRequest struct {
	CurrentPassword string `form:"current_password" validate:"required,min=6"`
	NewPassword     string `form:"new_password" validate:"required,min=8,nefield=CurrentPassword"`
	ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=NewPassword"`
}

func ValidateUpdatePasswordRequest(c *fiber.Ctx) error {
	var req UpdatePasswordRequest

	if err := c.BodyParser(&req); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı")
		return c.Redirect("/auth/update-password", fiber.StatusSeeOther)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		err := err.(validator.ValidationErrors)[0]
		switch {
		case err.Field() == "CurrentPassword" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Mevcut şifre zorunludur")
		case err.Field() == "CurrentPassword" && err.Tag() == "min":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Mevcut şifre en az 6 karakter olmalıdır")
		case err.Field() == "NewPassword" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Yeni şifre zorunludur")
		case err.Field() == "NewPassword" && err.Tag() == "min":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Yeni şifre en az 8 karakter olmalıdır")
		case err.Field() == "NewPassword" && err.Tag() == "nefield":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Yeni şifre mevcut şifreden farklı olmalıdır")
		case err.Field() == "ConfirmPassword" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifre tekrarı zorunludur")
		case err.Field() == "ConfirmPassword" && err.Tag() == "eqfield":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Yeni şifreler uyuşmuyor")
		default:
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz şifre bilgileri")
		}
		return c.Redirect("/auth/update-password", fiber.StatusSeeOther)
	}

	c.Locals("updatePasswordRequest", req)
	return c.Next()
}

type RegisterRequest struct {
	Name            string `form:"name" validate:"required,min=3"`
	Account         string `form:"account" validate:"required,email"`
	Password        string `form:"password" validate:"required,min=6"`
	ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=Password"`
}

func ValidateRegisterRequest(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı")
		return c.Redirect("/auth/register", fiber.StatusSeeOther)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		err := err.(validator.ValidationErrors)[0]
		switch {
		case err.Field() == "Name" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "İsim zorunludur")
		case err.Field() == "Account" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "E-posta zorunludur")
		case err.Field() == "Account" && err.Tag() == "email":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçerli bir e-posta adresi giriniz")
		case err.Field() == "Password" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifre zorunludur")
		case err.Field() == "Password" && err.Tag() == "min":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifre en az 6 karakter olmalıdır")
		case err.Field() == "ConfirmPassword" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifre tekrarı zorunludur")
		case err.Field() == "ConfirmPassword" && err.Tag() == "eqfield":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifreler eşleşmiyor")
		default:
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz kayıt bilgileri")
		}
		return c.Redirect("/auth/register", fiber.StatusSeeOther)
	}

	c.Locals("registerRequest", req)
	return c.Next()
}

type ForgotPasswordRequest struct {
	Email string `form:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token           string `form:"token" validate:"required"`
	NewPassword     string `form:"new_password" validate:"required,min=8"`
	ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=NewPassword"`
}

func ValidateForgotPasswordRequest(c *fiber.Ctx) error {
	var req ForgotPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		err := err.(validator.ValidationErrors)[0]
		switch {
		case err.Field() == "Email" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "E-posta zorunludur")
		case err.Field() == "Email" && err.Tag() == "email":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçerli bir e-posta adresi giriniz")
		default:
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek")
		}
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}

	c.Locals("forgotPasswordRequest", req)
	return c.Next()
}

func ValidateResetPasswordRequest(c *fiber.Ctx) error {
	var req ResetPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı")
		return c.Redirect("/auth/reset-password", fiber.StatusSeeOther)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		err := err.(validator.ValidationErrors)[0]
		switch {
		case err.Field() == "Token" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Token zorunludur")
		case err.Field() == "NewPassword" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Yeni şifre zorunludur")
		case err.Field() == "NewPassword" && err.Tag() == "min":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Yeni şifre en az 8 karakter olmalıdır")
		case err.Field() == "ConfirmPassword" && err.Tag() == "required":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifre onayı zorunludur")
		case err.Field() == "ConfirmPassword" && err.Tag() == "eqfield":
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifreler eşleşmiyor")
		default:
			_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz bilgiler")
		}
		return c.Redirect("/auth/reset-password", fiber.StatusSeeOther)
	}

	c.Locals("resetPasswordRequest", req)
	return c.Next()
}
