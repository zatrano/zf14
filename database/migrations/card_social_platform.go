package migrations

import (
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateCardSocialPlatformTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.CardSocialPlatform{})
}
