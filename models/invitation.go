package models

type Invitation struct {
	BaseModel
	InvitationKey    string             `gorm:"size:100;unique;not null;index"`
	IsConfirmed      bool               `gorm:"default:false;index"`
	IsRSVP           bool               `gorm:"default:false;index"`
	UserID           uint               `gorm:"not null;index"`
	CategoryID       uint               `gorm:"not null;index"`
	Image            string             `gorm:"size:255"`
	Description      string             `gorm:"type:text"`
	Venue            string             `gorm:"size:255"`
	Address          string             `gorm:"size:255"`
	Date             string             `gorm:"size:20"`
	Time             string             `gorm:"size:20"`
	Note             string             `gorm:"type:text"`
	Telephone        string             `gorm:"size:30"`
	Location         string             `gorm:"size:255"`
	Link             string             `gorm:"size:255"`
	Type             string             `gorm:"size:50;index"`
	Template         string             `gorm:"size:255"`
	User             *User              `gorm:"foreignKey:UserID"`
	Category         InvitationCategory `gorm:"foreignKey:CategoryID"`
	InvitationDetail *InvitationDetail  `gorm:"foreignKey:InvitationID"`
	Rsvp             []Rsvp             `gorm:"foreignKey:InvitationID"`
}
