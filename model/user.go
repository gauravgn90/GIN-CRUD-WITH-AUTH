package model

import "time"

type Status int

const (
	Active Status = 1 + iota
	Inactive
	Suspended
)

func (s Status) String() string {
	return [...]string{"Active", "Inactive", "Suspended"}[s-1]
}

type NewUser struct {
	Name     string `json:"name"  binding:"required,min=3,max=50"`
	Username string `json:"username"  binding:"required,min=5,max=50"`
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required,min=5,max=50"`
}

type User struct {
	Id int `gorm:"AUTO_INCREMENT;PRIMARY_KEY;not null" json:"id"`
	NewUser
	Status    Status     `gorm:"type:int;default:1" json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`

	Roles []Role `gorm:"many2many:user_roles;joinForeignKey:UserID;JoinReferences:RoleID" json:"roles"`
}
type UserUpdate struct {
	Name     string `json:"name"  binding:"required,min=3,max=50"`
	Username string `json:"username"  binding:"required,min=5,max=50"`
	Email    string `json:"email"  binding:"required,email"`
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

type SuccessResponse struct {
	Data       string `json:"data"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
}

type FailureResponse struct {
	Data       string `json:"data"`
	Status     string `json:"status" default:"error"`
	StatusCode int    `json:"status_code"`
}

type UserResponseSuccess struct {
	Data       Data   `json:"data"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
}

type UserResponseFailure struct {
	Data       string `json:"data"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
}

type Data struct {
	Token string `json:"token,omitempty"`
}

type UserType interface {
	isUserType()
}

func (u User) isUserType()        {}
func (uu UserUpdate) isUserType() {}

type UserResponseLogoutSuccess struct {
	Data       string `json:"data" default:"Logged out successfully"`
	Status     string `json:"status" default:"success"`
	StatusCode int    `json:"status_code" default:"200"`
}

type UserResponseLogoutUnAuthorizedFailure struct {
	Data       string `json:"data" default:"Authorization header is required"`
	Status     string `json:"status" default:"error"`
	StatusCode int    `json:"status_code" default:"400"`
}
type UserResponseLogoutInternalServerFailure struct {
	Data       string `json:"data" default:"Internal Server Error"`
	Status     string `json:"status" default:"error"`
	StatusCode int    `json:"status_code" default:"500"`
}
