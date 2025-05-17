package models

type Rsvp struct {
	BaseModel
	InvitationID     uint   `gorm:"not null;index"`
	Name             string `gorm:"size:100;not null"`
	Telephone        string `gorm:"size:30"`
	ParticipantCount int    `gorm:"default:1"`

	Invitation Invitation `gorm:"foreignKey:InvitationID"`
}
