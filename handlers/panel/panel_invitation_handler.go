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

type PanelInvitationHandler struct {
	invitationService services.IInvitationService
}

func NewPanelInvitationHandler() *PanelInvitationHandler {
	return &PanelInvitationHandler{invitationService: services.NewInvitationService()}
}

func (h *PanelInvitationHandler) ListPanelInvitations(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	invitations, err := h.invitationService.GetAllInvitations()
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiyeler alınamadı.")
	}
	// Sadece ilgili kullanıcıya ait davetiyeleri filtrele
	var userInvitations []models.Invitation
	for _, inv := range invitations {
		if inv.UserID == userID {
			userInvitations = append(userInvitations, inv)
		}
	}
	return renderer.Render(c, "panel/invitations/list", "layouts/panel", fiber.Map{
		"Title":       "Davetiyeler",
		"Invitations": userInvitations,
	}, http.StatusOK)
}

func (h *PanelInvitationHandler) ShowCreatePanelInvitation(c *fiber.Ctx) error {
	return renderer.Render(c, "panel/invitations/create", "layouts/panel", fiber.Map{
		"Title": "Yeni Davetiye Ekle",
	})
}

func (h *PanelInvitationHandler) CreatePanelInvitation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var req models.Invitation
	if err := c.BodyParser(&req); err != nil {
		return renderer.Render(c, "panel/invitations/create", "layouts/panel", fiber.Map{
			"Title":                     "Yeni Davetiye Ekle",
			flashmessages.FlashErrorKey: "Form verileri okunamadı.",
		}, http.StatusBadRequest)
	}
	req.UserID = userID
	if err := h.invitationService.CreateInvitation(c.UserContext(), &req); err != nil {
		return renderer.Render(c, "panel/invitations/create", "layouts/panel", fiber.Map{
			"Title":                     "Yeni Davetiye Ekle",
			flashmessages.FlashErrorKey: "Davetiye kaydedilemedi.",
		}, http.StatusInternalServerError)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla eklendi.")
	return c.Redirect("/panel/invitations", http.StatusFound)
}

func (h *PanelInvitationHandler) ShowUpdatePanelInvitation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	invitation, err := h.invitationService.GetInvitationByID(uint(id))
	if err != nil || invitation.UserID != userID {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye bulunamadı.")
		return c.Redirect("/panel/invitations", http.StatusSeeOther)
	}
	return renderer.Render(c, "panel/invitations/update", "layouts/panel", fiber.Map{
		"Title":      "Davetiye Düzenle",
		"Invitation": invitation,
	})
}

func (h *PanelInvitationHandler) UpdatePanelInvitation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	var req models.Invitation
	if err := c.BodyParser(&req); err != nil {
		return renderer.Render(c, "panel/invitations/update", "layouts/panel", fiber.Map{
			"Title":                     "Davetiye Düzenle",
			flashmessages.FlashErrorKey: "Form verileri okunamadı.",
		}, http.StatusBadRequest)
	}
	invitation, err := h.invitationService.GetInvitationByID(uint(id))
	if err != nil || invitation.UserID != userID {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye bulunamadı.")
		return c.Redirect("/panel/invitations", http.StatusSeeOther)
	}
	if err := h.invitationService.UpdateInvitation(c.UserContext(), uint(id), &req); err != nil {
		return renderer.Render(c, "panel/invitations/update", "layouts/panel", fiber.Map{
			"Title":                     "Davetiye Düzenle",
			flashmessages.FlashErrorKey: "Davetiye güncellenemedi.",
		}, http.StatusInternalServerError)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla güncellendi.")
	return c.Redirect("/panel/invitations", http.StatusFound)
}

func (h *PanelInvitationHandler) DeletePanelInvitation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	invitation, err := h.invitationService.GetInvitationByID(uint(id))
	if err != nil || invitation.UserID != userID {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye bulunamadı.")
		return c.Redirect("/panel/invitations", http.StatusSeeOther)
	}
	if err := h.invitationService.DeleteInvitation(c.UserContext(), uint(id)); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Davetiye silinemedi.")
		return c.Redirect("/panel/invitations", http.StatusSeeOther)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Davetiye başarıyla silindi.")
	return c.Redirect("/panel/invitations", http.StatusFound)
}
