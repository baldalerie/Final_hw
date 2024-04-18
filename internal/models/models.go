package models

import (
	"time"
)

type Users struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `json:"name"`
	Email        string         `gorm:"unique" json:"email"`
	Password     string         `json:"password"`
	Transactions []Transactions `gorm:"foreignKey:UserID"`
}

type Transactions struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          string    `gorm:"not null" json:"user_id"`
	Amount          float64   `gorm:"not null;default:0" json:"amount"`
	Currency        string    `gorm:"not null" json:"currency"`
	TransactionType string    `json:"type"`
	Category        string    `json:"category"`
	TransactionDate time.Time `json:"date"`
	Description     string    `json:"description"`
	Commission      float64   `gorm:"not null;default:0" json:"commission"`
}
type ConvertedTransaction struct {
	Amount       float64      `json:"amountConverted"`
	Currency     string       `json:"currencyConverted"`
	Transactions Transactions `json:"transactions"`
}
