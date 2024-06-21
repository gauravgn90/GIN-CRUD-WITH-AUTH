package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	Id          int            `gorm:"type:bigint(20);primary_key" json:"id"`
	Name        string         `json:"name"  binding:"required,min=3,max=50"`
	Description string         `json:"description"  binding:"required,min=5,max=50"`
	Price       float64        `json:"price"  binding:"required"`
	Quantity    int            `json:"quantity"  binding:"required"`
	Category    string         `json:"category"  binding:"required"`
	CreatedAt   time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;type:timestamp"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP;type:timestamp"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"default:NULL;type:timestamp"`
}

func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
	// Set the updated_at field to the current time before saving the record
	p.UpdatedAt = time.Now()
	return nil
}
