package models

type BankCard struct {
	BaseModel
	BankID     uint   `gorm:"not null;index"`
	CardID     uint   `gorm:"not null;index"`
	IbanNumber string `gorm:"size:34;not null"`

	Bank Bank `gorm:"foreignKey:BankID"`
	Card Card `gorm:"foreignKey:CardID"`
}
