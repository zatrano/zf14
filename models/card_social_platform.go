package models

type CardSocialPlatform struct {
	BaseModel
	CardID           uint   `gorm:"not null;index"`
	SocialPlatformID uint   `gorm:"not null;index"`
	LinkValue        string `gorm:"size:255"`

	Card           Card           `gorm:"foreignKey:CardID"`
	SocialPlatform SocialPlatform `gorm:"foreignKey:SocialPlatformID"`
}
