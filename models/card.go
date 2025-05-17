package models

type Card struct {
	BaseModel
	UserID          uint             `gorm:"not null;index"`
	Slug            string           `gorm:"size:100;unique;not null;index"`
	Name            string           `gorm:"size:100;not null"`
	Title           string           `gorm:"size:100"`
	ProfilePicture  string           `gorm:"size:255"`
	Phone           string           `gorm:"size:30"`
	Email           string           `gorm:"size:100"`
	LocationURL     string           `gorm:"size:255"`
	WebsiteURL      string           `gorm:"size:255"`
	IsActive        bool             `gorm:"default:true;index"`
	Banks           []Bank           `gorm:"many2many:bank_card;joinForeignKey:CardID;joinReferences:BankID"`
	SocialPlatforms []SocialPlatform `gorm:"many2many:card_social_platform;joinForeignKey:CardID;joinReferences:SocialPlatformID"`
	User            *User            `gorm:"foreignKey:UserID"`
}
