package migrations

import (
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateInvitationTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Invitation{})
}
