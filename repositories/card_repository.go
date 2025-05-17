package repositories

import (
	"context"

	"davet.link/configs/databaseconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"
)

type ICardRepository interface {
	GetAllCards() ([]models.Card, error)
	GetCardByID(id uint) (*models.Card, error)
	CreateCard(ctx context.Context, card *models.Card) error
	UpdateCard(ctx context.Context, id uint, card *models.Card) error
	DeleteCard(ctx context.Context, id uint) error
}

type CardRepository struct {
	*BaseRepository[models.Card]
}

func NewCardRepository() ICardRepository {
	repo := &CardRepository{BaseRepository: NewBaseRepository[models.Card](databaseconfig.DB)}
	repo.SetPreloads("User", "Banks", "SocialPlatforms")
	return repo
}

func (r *CardRepository) GetAllCards() ([]models.Card, error) {
	cards, _, err := r.BaseRepository.GetAll(queryparams.ListParams{})
	return cards, err
}

func (r *CardRepository) GetCardByID(id uint) (*models.Card, error) {
	return r.BaseRepository.GetByID(id)
}

func (r *CardRepository) CreateCard(ctx context.Context, card *models.Card) error {
	return r.BaseRepository.Create(ctx, card)
}

func (r *CardRepository) UpdateCard(ctx context.Context, id uint, card *models.Card) error {
	data := map[string]interface{}{
		"name":            card.Name,
		"slug":            card.Slug,
		"title":           card.Title,
		"profile_picture": card.ProfilePicture,
		"phone":           card.Phone,
		"email":           card.Email,
		"location_url":    card.LocationURL,
		"website_url":     card.WebsiteURL,
		"is_active":       card.IsActive,
	}
	return r.BaseRepository.Update(ctx, id, data, 0)
}

func (r *CardRepository) DeleteCard(ctx context.Context, id uint) error {
	return r.BaseRepository.Delete(ctx, id)
}
