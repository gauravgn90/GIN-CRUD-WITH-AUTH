package model

import (
	"time"
)

type UserRole struct {
	Id        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int        `gorm:"not null" json:"user_id"`
	RoleID    int        `gorm:"not null" json:"role_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	User User `gorm:"foreignKey:UserID;references:Id"`
	Role Role `gorm:"foreignKey:RoleID;references:Id"`
}
