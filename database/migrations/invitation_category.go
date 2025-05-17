package migrations

import (
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateInvitationCategoryTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.InvitationCategory{})
}
