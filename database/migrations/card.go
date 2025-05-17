package migrations

import (
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateCardTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Card{})
}
