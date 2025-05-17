package models

type InvitationCategory struct {
	BaseModel
	Name        string       `gorm:"size:100;not null;index"`
	Template    string       `gorm:"size:255"`
	IsActive    bool         `gorm:"default:true;index"`
	Invitations []Invitation `gorm:"foreignKey:CategoryID"`
}
