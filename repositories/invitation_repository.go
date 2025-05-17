package repositories

import (
	"context"

	"davet.link/configs/databaseconfig"
	"davet.link/models"
	"davet.link/pkg/queryparams"
)

type IInvitationRepository interface {
	GetAllInvitations() ([]models.Invitation, error)
	GetInvitationByID(id uint) (*models.Invitation, error)
	CreateInvitation(ctx context.Context, invitation *models.Invitation) error
	UpdateInvitation(ctx context.Context, id uint, invitation *models.Invitation) error
	DeleteInvitation(ctx context.Context, id uint) error
}

type InvitationRepository struct {
	*BaseRepository[models.Invitation]
}

func NewInvitationRepository() IInvitationRepository {
	repo := &InvitationRepository{BaseRepository: NewBaseRepository[models.Invitation](databaseconfig.DB)}
	repo.SetPreloads("User", "Category", "InvitationDetail", "Rsvp")
	return repo
}

func (r *InvitationRepository) GetAllInvitations() ([]models.Invitation, error) {
	invitations, _, err := r.BaseRepository.GetAll(queryparams.ListParams{})
	return invitations, err
}

func (r *InvitationRepository) GetInvitationByID(id uint) (*models.Invitation, error) {
	return r.BaseRepository.GetByID(id)
}

func (r *InvitationRepository) CreateInvitation(ctx context.Context, invitation *models.Invitation) error {
	return r.BaseRepository.Create(ctx, invitation)
}

func (r *InvitationRepository) UpdateInvitation(ctx context.Context, id uint, invitation *models.Invitation) error {
	data := map[string]interface{}{
		"invitation_key": invitation.InvitationKey,
		"is_confirmed":   invitation.IsConfirmed,
		"is_rsvp":        invitation.IsRSVP,
		"user_id":        invitation.UserID,
		"category_id":    invitation.CategoryID,
		"image":          invitation.Image,
		"description":    invitation.Description,
		"venue":          invitation.Venue,
		"address":        invitation.Address,
		"date":           invitation.Date,
		"time":           invitation.Time,
		"note":           invitation.Note,
		"telephone":      invitation.Telephone,
		"location":       invitation.Location,
		"link":           invitation.Link,
		"type":           invitation.Type,
		"template":       invitation.Template,
	}
	return r.BaseRepository.Update(ctx, id, data, 0)
}

func (r *InvitationRepository) DeleteInvitation(ctx context.Context, id uint) error {
	return r.BaseRepository.Delete(ctx, id)
}
