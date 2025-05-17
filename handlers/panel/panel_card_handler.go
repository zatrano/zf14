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

type PanelCardHandler struct {
	cardService services.ICardService
}

func NewPanelCardHandler() *PanelCardHandler {
	return &PanelCardHandler{cardService: services.NewCardService()}
}

func (h *PanelCardHandler) ListPanelCards(c *fiber.Ctx) error {
	cards, err := h.cardService.GetAllCards()
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kartlar alınamadı.")
	}
	return renderer.Render(c, "panel/cards/list", "layouts/panel", fiber.Map{
		"Title": "Kartlar",
		"Cards": cards,
	}, http.StatusOK)
}

func (h *PanelCardHandler) ShowCreatePanelCard(c *fiber.Ctx) error {
	return renderer.Render(c, "panel/cards/create", "layouts/panel", fiber.Map{
		"Title": "Yeni Kart Ekle",
	})
}

func (h *PanelCardHandler) CreatePanelCard(c *fiber.Ctx) error {
	var req models.Card
	if err := c.BodyParser(&req); err != nil {
		return renderer.Render(c, "panel/cards/create", "layouts/panel", fiber.Map{
			"Title":                     "Yeni Kart Ekle",
			flashmessages.FlashErrorKey: "Form verileri okunamadı.",
		}, http.StatusBadRequest)
	}
	if err := h.cardService.CreateCard(c.UserContext(), &req); err != nil {
		return renderer.Render(c, "panel/cards/create", "layouts/panel", fiber.Map{
			"Title":                     "Yeni Kart Ekle",
			flashmessages.FlashErrorKey: "Kart kaydedilemedi.",
		}, http.StatusInternalServerError)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kart başarıyla eklendi.")
	return c.Redirect("/panel/cards", http.StatusFound)
}

func (h *PanelCardHandler) ShowUpdatePanelCard(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	card, err := h.cardService.GetCardByID(uint(id))
	if err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kart bulunamadı.")
		return c.Redirect("/panel/cards", http.StatusSeeOther)
	}
	return renderer.Render(c, "panel/cards/update", "layouts/panel", fiber.Map{
		"Title": "Kart Düzenle",
		"Card":  card,
	})
}

func (h *PanelCardHandler) UpdatePanelCard(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req models.Card
	if err := c.BodyParser(&req); err != nil {
		return renderer.Render(c, "panel/cards/update", "layouts/panel", fiber.Map{
			"Title":                     "Kart Düzenle",
			flashmessages.FlashErrorKey: "Form verileri okunamadı.",
		}, http.StatusBadRequest)
	}
	if err := h.cardService.UpdateCard(c.UserContext(), uint(id), &req); err != nil {
		return renderer.Render(c, "panel/cards/update", "layouts/panel", fiber.Map{
			"Title":                     "Kart Düzenle",
			flashmessages.FlashErrorKey: "Kart güncellenemedi.",
		}, http.StatusInternalServerError)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kart başarıyla güncellendi.")
	return c.Redirect("/panel/cards", http.StatusFound)
}

func (h *PanelCardHandler) DeletePanelCard(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.cardService.DeleteCard(c.UserContext(), uint(id)); err != nil {
		_ = flashmessages.SetFlashMessage(c, flashmessages.FlashErrorKey, "Kart silinemedi.")
		return c.Redirect("/panel/cards", http.StatusSeeOther)
	}
	_ = flashmessages.SetFlashMessage(c, flashmessages.FlashSuccessKey, "Kart başarıyla silindi.")
	return c.Redirect("/panel/cards", http.StatusFound)
}
