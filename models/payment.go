package models

import "time"

type Payment struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	ProductID uint      `json:"product_id" gorm:"foreignkey:ProductID"`
	Price     uint      `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
