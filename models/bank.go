package models

type Bank struct {
	BaseModel
	Name     string `gorm:"size:100;not null;unique;index"`
	IsActive bool   `gorm:"default:true;index"`
	Cards    []Card `gorm:"many2many:bank_card;joinForeignKey:BankID;joinReferences:CardID"`
}
