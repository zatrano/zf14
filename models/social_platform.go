package models

type SocialPlatform struct {
	BaseModel
	Name      string `gorm:"size:100;not null;unique;index"`
	IconClass string `gorm:"size:100"`
	IsActive  bool   `gorm:"default:true;index"`
	Cards     []Card `gorm:"many2many:card_social_platform;joinForeignKey:SocialPlatformID;joinReferences:CardID"`
}
