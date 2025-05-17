package services

import (
	"context"

	"davet.link/models"
	"davet.link/repositories"
)

type IInvitationService interface {
	GetAllInvitations() ([]models.Invitation, error)
	GetInvitationByID(id uint) (*models.Invitation, error)
	CreateInvitation(ctx context.Context, invitation *models.Invitation) error
	UpdateInvitation(ctx context.Context, id uint, invitation *models.Invitation) error
	DeleteInvitation(ctx context.Context, id uint) error
}

type InvitationService struct {
	repo repositories.IInvitationRepository
}

func NewInvitationService() IInvitationService {
	return &InvitationService{repo: repositories.NewInvitationRepository()}
}

func (s *InvitationService) GetAllInvitations() ([]models.Invitation, error) {
	return s.repo.GetAllInvitations()
}

func (s *InvitationService) GetInvitationByID(id uint) (*models.Invitation, error) {
	return s.repo.GetInvitationByID(id)
}

func (s *InvitationService) CreateInvitation(ctx context.Context, invitation *models.Invitation) error {
	return s.repo.CreateInvitation(ctx, invitation)
}

func (s *InvitationService) UpdateInvitation(ctx context.Context, id uint, invitation *models.Invitation) error {
	return s.repo.UpdateInvitation(ctx, id, invitation)
}

func (s *InvitationService) DeleteInvitation(ctx context.Context, id uint) error {
	return s.repo.DeleteInvitation(ctx, id)
}
