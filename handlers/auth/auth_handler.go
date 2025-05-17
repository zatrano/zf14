package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"

	"davet.link/configs/logconfig"
	"davet.link/configs/sessionconfig"
	"davet.link/models"
	"davet.link/pkg/flashmessages"
	"davet.link/pkg/renderer"
	"davet.link/requests"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service services.IAuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{service: services.NewAuthService()}
}

func (h *AuthHandler) handleError(c *fiber.Ctx, err error, userID uint, account string, action string) error {
	var errMsg string
	flashKey := flashmessages.FlashErrorKey
	redirectTarget := "/auth/login"
	logoutUser := false

	switch err {
	case services.ErrInvalidCredentials:
		errMsg = "Kullanıcı adı veya şifre hatalı."
	case services.ErrUserInactive:
		errMsg = "Hesabınız aktif değil. Lütfen yöneticinizle iletişime geçin."
	case services.ErrUserNotFound:
		errMsg = "Kullanıcı bulunamadı, lütfen tekrar giriş yapın."
		logoutUser = true
		logconfig.Log.Warn(action+": Kullanıcı bulunamadı", zap.Uint("user_id", userID))
	case services.ErrCurrentPasswordIncorrect:
		errMsg = "Mevcut şifreniz hatalı."
		redirectTarget = "/auth/profile"
	case services.ErrPasswordTooShort, services.ErrPasswordSameAsOld:
		errMsg = err.Error()
		redirectTarget = "/auth/profile"
	default:
		errMsg = "İşlem sırasında bir sorun oluştu. Lütfen tekrar deneyin."
		logconfig.Log.Error(action+": Beklenmeyen hata",
			zap.Uint("user_id", userID),
			zap.String("account", account),
			zap.Error(err))
	}

	if logoutUser {
		h.destroySession(c)
	}

	_ = flashmessages.SetFlashMessage(c, flashKey, errMsg)
	return c.Redirect(redirectTarget, fiber.StatusSeeOther)
}

func (h *AuthHandler) getSessionUser(c *fiber.Ctx) (uint, error) {
	if userID, ok := c.Locals("userID").(uint); ok {
		return userID, nil
	}

	sess, err := sessionconfig.SessionStart(c)
	if err != nil {
		return 0, err
	}

	userIDValue := sess.Get("user_id")
	switch v := userIDValue.(type) {
	case uint:
		return v, nil
	case int:
		return uint(v), nil
	case float64:
		return uint(v), nil
	default:
		return 0, fiber.ErrUnauthorized
	}
}

func (h *AuthHandler) destroySession(c *fiber.Ctx) {
	sess, err := sessionconfig.SessionStart(c)
	if err != nil {
		logconfig.Log.Warn("Oturum yok edilemedi (zaten yok olabilir)", zap.Error(err))
		return
	}
	if err := sess.Destroy(); err != nil {
		logconfig.Log.Error("Oturum yok edilemedi", zap.Error(err))
	}
}

func (h *AuthHandler) ShowLogin(c *fiber.Ctx) error {
	mapData := fiber.Map{
		"Title": "Giriş",
	}
	return renderer.Render(c, "auth/login", "layouts/auth", mapData, http.StatusOK)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req, ok := c.Locals("loginRequest").(requests.LoginRequest)
	if !ok {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	user, err := h.service.Authenticate(req.Account, req.Password)
	if err != nil {
		return h.handleError(c, err, 0, req.Account, "Login")
	}

	sess, err := sessionconfig.SessionStart(c)
	if err != nil {
		logconfig.Log.Error("Oturum başlatılamadı",
			zap.Uint("user_id", user.ID),
			zap.String("account", user.Account),
			zap.Error(err))
		return h.handleError(c, fiber.ErrInternalServerError, user.ID, user.Account, "Login")
	}

	sess.Set("user_id", user.ID)
	sess.Set("user_type", string(user.Type))
	if err := sess.Save(); err != nil {
		logconfig.Log.Error("Oturum kaydedilemedi",
			zap.Uint("user_id", user.ID),
			zap.String("account", user.Account),
			zap.Error(err))
		return h.handleError(c, fiber.ErrInternalServerError, user.ID, user.Account, "Login")
	}

	switch user.Type {
	case models.Panel:
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Başarıyla giriş yapıldı")
		return c.Redirect("/panel/home", fiber.StatusFound)
	case models.Dashboard:
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Başarıyla giriş yapıldı")
		return c.Redirect("/dashboard/home", fiber.StatusFound)
	default:
		_ = sess.Destroy()
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz kullanıcı tipi")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	userID, err := h.getSessionUser(c)
	if err != nil {
		logconfig.Log.Warn("Profil: Geçersiz oturum", zap.Error(err))
		h.destroySession(c)
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz oturum, lütfen tekrar giriş yapın.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	user, err := h.service.GetUserProfile(userID)
	if err != nil {
		return h.handleError(c, err, userID, "", "Profil")
	}

	mapData := fiber.Map{
		"Title": "Profilim",
		"User":  user,
	}
	return renderer.Render(c, "auth/profile", "layouts/auth", mapData, http.StatusOK)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	h.destroySession(c)
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Başarıyla çıkış yapıldı.")
	return c.Redirect("/auth/login", fiber.StatusFound)
}

func (h *AuthHandler) UpdatePassword(c *fiber.Ctx) error {
	userID, err := h.getSessionUser(c)
	if err != nil {
		h.destroySession(c)
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz oturum bilgisi, lütfen tekrar giriş yapın.")
		return c.Redirect("/auth/login", fiber.StatusSeeOther)
	}

	req, ok := c.Locals("updatePasswordRequest").(requests.UpdatePasswordRequest)
	if !ok {
		logconfig.SLog.Warn("Parola güncelleme: Geçersiz istek formatı")
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek formatı.")
		return c.Redirect("/auth/profile", fiber.StatusSeeOther)
	}

	if err := h.service.UpdatePassword(c.UserContext(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		return h.handleError(c, err, userID, "", "Parola Güncelleme")
	}

	h.destroySession(c)
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Şifre başarıyla güncellendi. Lütfen yeni şifrenizle tekrar giriş yapın.")
	return c.Redirect("/auth/login", fiber.StatusFound)
}

func (h *AuthHandler) ShowRegister(c *fiber.Ctx) error {
	mapData := fiber.Map{
		"Title": "Kayıt Ol",
	}
	return renderer.Render(c, "auth/register", "layouts/auth", mapData, http.StatusOK)
}

func generateToken() (string, error) {
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenBytes), nil
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req, ok := c.Locals("registerRequest").(requests.RegisterRequest)
	if !ok {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz kayıt isteği")
		return c.Redirect("/auth/register", fiber.StatusSeeOther)
	}

	if req.Password != req.ConfirmPassword {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifreler eşleşmiyor")
		return c.Redirect("/auth/register", fiber.StatusSeeOther)
	}

	user := &models.User{
		Name:     req.Name,
		Account:  req.Account,
		Password: req.Password,
		Status:   true,
		Type:     models.Panel,
	}

	resetToken, err := generateToken()
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Reset token oluşturulamadı")
		return c.Redirect("/auth/register", fiber.StatusSeeOther)
	}
	user.ResetToken = resetToken

	verificationToken, err := generateToken()
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Verification token oluşturulamadı")
		return c.Redirect("/auth/register", fiber.StatusSeeOther)
	}
	user.VerificationToken = verificationToken

	ctx := c.UserContext()
	if err := h.service.CreateUser(ctx, user); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kullanıcı oluşturulamadı. Lütfen tekrar deneyin.")
		return c.Redirect("/auth/register", fiber.StatusSeeOther)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kayıt işlemi başarıyla tamamlandı. Lütfen email adresinizi doğrulayın.")

	mailService := services.NewMailService()
	baseURL := os.Getenv("APP_BASE_URL")
	verificationLink := baseURL + "/auth/verify-email?token=" + verificationToken
	emailBody := "Lütfen email adresinizi doğrulamak için aşağıdaki bağlantıya tıklayın: " + verificationLink
	_ = mailService.SendMail(user.Account, "Email Doğrulama", emailBody)

	return renderer.Render(c, "auth/verify_email_notice", "layouts/auth", nil, http.StatusOK)
}

func (h *AuthHandler) ShowForgotPassword(c *fiber.Ctx) error {
	mapData := fiber.Map{
		"Title": "Şifremi Unuttum",
	}
	return renderer.Render(c, "auth/forgot_password", "layouts/auth", mapData, http.StatusOK)
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	req, ok := c.Locals("forgotPasswordRequest").(requests.ForgotPasswordRequest)
	if !ok {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz istek")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}

	if err := h.service.SendPasswordResetLink(req.Email); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifre sıfırlama bağlantısı gönderilemedi. Lütfen tekrar deneyin.")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Şifre sıfırlama bağlantısı başarıyla gönderildi. Lütfen emailinizi kontrol edin.")

	return c.Redirect("/auth/login", fiber.StatusSeeOther)
}

func (h *AuthHandler) ShowResetPassword(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz veya eksik token.")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}

	mapData := fiber.Map{
		"Title": "Şifre Sıfırla",
		"Token": token,
	}
	return renderer.Render(c, "auth/reset_password", "layouts/auth", mapData, http.StatusOK)
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	req, ok := c.Locals("resetPasswordRequest").(requests.ResetPasswordRequest)
	if !ok || req.Token == "" {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Geçersiz veya eksik token.")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}

	if req.NewPassword != req.ConfirmPassword {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifreler eşleşmiyor. Lütfen tekrar deneyin.")
		return c.Redirect("/auth/reset-password", fiber.StatusSeeOther)
	}

	if err := h.service.ResetPassword(req.Token, req.NewPassword); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Şifre sıfırlama işlemi başarısız oldu.")
		return c.Redirect("/auth/reset-password", fiber.StatusSeeOther)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Şifreniz başarıyla sıfırlandı. Lütfen giriş yapın.")
	return c.Redirect("/auth/login", fiber.StatusSeeOther)
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Doğrulama tokeni eksik veya geçersiz.")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}

	if err := h.service.VerifyEmail(token); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Email doğrulama başarısız.")
		return c.Redirect("/auth/forgot-password", fiber.StatusSeeOther)
	}

	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Email başarıyla doğrulandı.")
	return c.Redirect("/auth/login", fiber.StatusSeeOther)
}

func (h *AuthHandler) ShowResendVerification(c *fiber.Ctx) error {
    mapData := fiber.Map{"Title": "Email Doğrulama Linkini Yeniden Gönder"}
    return renderer.Render(c, "auth/resend_verification", "layouts/auth", mapData, http.StatusOK)
}

func (h *AuthHandler) ResendVerification(c *fiber.Ctx) error {
    email := c.FormValue("email")
    if email == "" {
        _ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Email adresi gerekli.")
        return c.Redirect("/auth/resend-verification", fiber.StatusSeeOther)
    }
    if err := h.service.ResendVerificationLink(email); err != nil {
        _ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Doğrulama linki gönderilemedi.")
        return c.Redirect("/auth/resend-verification", fiber.StatusSeeOther)
    }
    _ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Doğrulama linki e-posta adresinize gönderildi.")
    return c.Redirect("/auth/login", fiber.StatusSeeOther)
}
