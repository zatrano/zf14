package handlers

import (
	"net/http"
	"strconv"

	"davet.link/models"
	"davet.link/pkg/flashmessages"
	"davet.link/pkg/renderer"
	"davet.link/services"

	"github.com/gofiber/fiber/v2"
)

type InvitationHandler struct {
	invitationService services.IInvitationService
}

func NewInvitationHandler() *InvitationHandler {
	return &InvitationHandler{invitationService: services.NewInvitationService()}
}

func (h *InvitationHandler) ListInvitations(c *fiber.Ctx) error {
	invitations, err := h.invitationService.GetAllInvitations()
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiyeler alınamadı.")
	}
	return renderer.Render(c, "dashboard/invitations/list", "layouts/dashboard", fiber.Map{
		"Title":       "Davetiyeler",
		"Invitations": invitations,
	}, http.StatusOK)
}

func (h *InvitationHandler) ShowCreateInvitation(c *fiber.Ctx) error {
	return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
		"Title": "Yeni Davetiye Ekle",
	})
}

func (h *InvitationHandler) CreateInvitation(c *fiber.Ctx) error {
	var req models.Invitation
	if err := c.BodyParser(&req); err != nil {
		return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
			"Title":                     "Yeni Davetiye Ekle",
			flashmessages.FlashErrorKey: "Form verileri okunamadı.",
		}, http.StatusBadRequest)
	}
	if err := h.invitationService.CreateInvitation(c.UserContext(), &req); err != nil {
		return renderer.Render(c, "dashboard/invitations/create", "layouts/dashboard", fiber.Map{
			"Title":                     "Yeni Davetiye Ekle",
			flashmessages.FlashErrorKey: "Davetiye kaydedilemedi.",
		}, http.StatusInternalServerError)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla eklendi.")
	return c.Redirect("/dashboard/invitations", http.StatusFound)
}

func (h *InvitationHandler) ShowUpdateInvitation(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	invitation, err := h.invitationService.GetInvitationByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye bulunamadı.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}
	return renderer.Render(c, "dashboard/invitations/update", "layouts/dashboard", fiber.Map{
		"Title":      "Davetiye Düzenle",
		"Invitation": invitation,
	})
}

func (h *InvitationHandler) UpdateInvitation(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req models.Invitation
	if err := c.BodyParser(&req); err != nil {
		return renderer.Render(c, "dashboard/invitations/update", "layouts/dashboard", fiber.Map{
			"Title":                     "Davetiye Düzenle",
			flashmessages.FlashErrorKey: "Form verileri okunamadı.",
		}, http.StatusBadRequest)
	}
	if err := h.invitationService.UpdateInvitation(c.UserContext(), uint(id), &req); err != nil {
		return renderer.Render(c, "dashboard/invitations/update", "layouts/dashboard", fiber.Map{
			"Title":                     "Davetiye Düzenle",
			flashmessages.FlashErrorKey: "Davetiye güncellenemedi.",
		}, http.StatusInternalServerError)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla güncellendi.")
	return c.Redirect("/dashboard/invitations", http.StatusFound)
}

func (h *InvitationHandler) DeleteInvitation(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.invitationService.DeleteInvitation(c.UserContext(), uint(id)); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye silinemedi.")
		return c.Redirect("/dashboard/invitations", http.StatusSeeOther)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla silindi.")
	return c.Redirect("/dashboard/invitations", http.StatusFound)
}
