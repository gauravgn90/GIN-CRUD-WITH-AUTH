package model

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey;autoIncrement:false;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"role_id"`
	PermissionID uint `gorm:"primaryKey;autoIncrement:false;not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"permission_id"`

	Role       Role       `gorm:"foreignKey:RoleID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Permission Permission `gorm:"foreignKey:PermissionID;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

type CreateRolePermissionRequest struct {
	Role        Role         `json:"role" binding:"required"`
	Permissions []Permission `json:"permissions" binding:"required"`
}

type AssignPermissionsToRoleRequest struct {
	RoleID      uint  `json:"role_id" binding:"required"`
	Permissions []int `json:"permissions" binding:"required"`
}

type AssignRolesToUserRequest struct {
	UserID int   `json:"user_id" binding:"required"`
	Roles  []int `json:"roles" binding:"required"`
}
