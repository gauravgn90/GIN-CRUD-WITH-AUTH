package model

import "time"

type User struct {
	Id        int        `gorm:"primary_key" json:"id"`
	Name      string     `json:"name"  binding:"required,min=3,max=50"`
	Username  string     `json:"username"  binding:"required,min=5,max=50"`
	Email     string     `json:"email"  binding:"required,email"`
	Password  string     `json:"password"  binding:"required,min=5,max=50"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}
type UserUpdate struct {
	Name     string  `json:"name"  binding:"required,min=3,max=50"`
	Username *string `json:"username"  binding:"required,min=5,max=50"`
	Email    *string `json:"email"  binding:"required,email"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserLogin struct {
	Username string `json:"username"  binding:"required,min=5,max=50"`
	Password string `json:"password"  binding:"required,min=5,max=50"`
}

type UserType interface {
	isUserType()
}

func (u User) isUserType()        {}
func (uu UserUpdate) isUserType() {}
