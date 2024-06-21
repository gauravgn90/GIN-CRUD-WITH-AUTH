package model

type Role struct {
	Id          int    `gorm:"AUTO_INCREMENT;PRIMARY_KEY;not null" json:"id"`
	Name        string `gorm:"type:varchar(255);not null;unique" json:"name" binding:"required,min=3,max=50"`
	RoleCode    string `gorm:"type:varchar(255);not null;unique" json:"role_code" binding:"required,min=3,max=50"`
	Description string `json:"description"  binding:"required,min=5,max=50"`

	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Users       []User       `gorm:"many2many:user_roles;joinForeignKey:RoleID;JoinReferences:UserID" json:"-"`
}
