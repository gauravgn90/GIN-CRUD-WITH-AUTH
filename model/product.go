package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	Id          int       `gorm:"primary_key" json:"id"`
	Name        string    `json:"name"  binding:"required,min=3,max=50"`
	Description string    `json:"description"  binding:"required,min=5,max=50"`
	Price       float64   `json:"price"  binding:"required"`
	Quantity    int       `json:"quantity"  binding:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt   string    `json:"-" gorm:"default:NULL" `
}

func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
	// Set the updated_at field to the current time before saving the record
	p.UpdatedAt = time.Now()
	return nil
}
