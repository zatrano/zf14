package migrations

import (
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateBankTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Bank{})
}
