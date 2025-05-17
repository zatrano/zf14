package services

import (
	"context"

	"davet.link/models"
	"davet.link/repositories"
)

type ICardService interface {
	GetAllCards() ([]models.Card, error)
	GetCardByID(id uint) (*models.Card, error)
	CreateCard(ctx context.Context, card *models.Card) error
	UpdateCard(ctx context.Context, id uint, card *models.Card) error
	DeleteCard(ctx context.Context, id uint) error
}

type CardService struct {
	repo repositories.ICardRepository
}

func NewCardService() ICardService {
	return &CardService{repo: repositories.NewCardRepository()}
}

func (s *CardService) GetAllCards() ([]models.Card, error) {
	return s.repo.GetAllCards()
}

func (s *CardService) GetCardByID(id uint) (*models.Card, error) {
	return s.repo.GetCardByID(id)
}

func (s *CardService) CreateCard(ctx context.Context, card *models.Card) error {
	return s.repo.CreateCard(ctx, card)
}

func (s *CardService) UpdateCard(ctx context.Context, id uint, card *models.Card) error {
	return s.repo.UpdateCard(ctx, id, card)
}

func (s *CardService) DeleteCard(ctx context.Context, id uint) error {
	return s.repo.DeleteCard(ctx, id)
}
