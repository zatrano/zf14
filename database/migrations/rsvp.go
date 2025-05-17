package migrations

import (
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateRsvpTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Rsvp{})
}
