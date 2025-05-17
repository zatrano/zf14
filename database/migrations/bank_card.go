package migrations

import (
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateBankCardTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.BankCard{})
}
