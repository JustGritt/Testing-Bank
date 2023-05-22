package models

import "time"

type Product struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	Name      string `json:"name" gorm:"unique;not null"`
	Price     uint   `json:"price"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
