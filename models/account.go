package models

import "time"

type Account struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	IBAN      string    `json:"iban"`
	Balance   uint      `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
