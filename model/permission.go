package model

type Permission struct {
	Id             int    `gorm:"AUTO_INCREMENT;PRIMARY_KEY;not null" json:"id"`
	Name           string `gorm:"type:varchar(255);not null;unique" json:"name"  binding:"required,min=3,max=50"`
	PermissionCode string `gorm:"type:varchar(255);not null;unique" json:"permission_code"  binding:"required,min=3,max=50"`
	Description    string `json:"description"  binding:"required,min=5,max=50"`
	Role           []Role `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
