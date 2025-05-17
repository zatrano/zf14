package models

type InvitationDetail struct {
	BaseModel
	InvitationID       uint   `gorm:"not null;index"`
	Title              string `gorm:"size:255"`
	Person             string `gorm:"size:100"`
	IsMotherLive       bool   `gorm:"default:true"`
	MotherName         string `gorm:"size:100"`
	MotherSurname      string `gorm:"size:100"`
	IsFatherLive       bool   `gorm:"default:true"`
	FatherName         string `gorm:"size:100"`
	FatherSurname      string `gorm:"size:100"`
	BrideName          string `gorm:"size:100"`
	BrideSurname       string `gorm:"size:100"`
	IsBrideMotherLive  bool   `gorm:"default:true"`
	BrideMotherName    string `gorm:"size:100"`
	BrideMotherSurname string `gorm:"size:100"`
	IsBrideFatherLive  bool   `gorm:"default:true"`
	BrideFatherName    string `gorm:"size:100"`
	BrideFatherSurname string `gorm:"size:100"`
	GroomName          string `gorm:"size:100"`
	GroomSurname       string `gorm:"size:100"`
	IsGroomMotherLive  bool   `gorm:"default:true"`
	GroomMotherName    string `gorm:"size:100"`
	GroomMotherSurname string `gorm:"size:100"`
	IsGroomFatherLive  bool   `gorm:"default:true"`
	GroomFatherName    string `gorm:"size:100"`
	GroomFatherSurname string `gorm:"size:100"`

	Invitation Invitation `gorm:"foreignKey:InvitationID"`
}
