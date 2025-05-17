package migrations

import (
	"davet.link/models"
	"gorm.io/gorm"
)

func MigrateSocialPlatformTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.SocialPlatform{})
}
