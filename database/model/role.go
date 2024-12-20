package model

type Role struct {
	Model
	Name    string    `json:"name"`
	Pid     int       `json:"pid"`
	Comment string    `json:"comment"`
	Routers []Routers `json:"routers" gorm:"many2many:routers_roles"`
	Users   []User    `json:"users" gorm:"many2many:user_roles"`
}

type RoleRouters struct {
	RoleID    int `json:"role_id" gorm:"column:role_id"`
	RoutersID int `json:"routers_id" gorm:"column:routers_id"`
}

type UserRoles struct {
	UserID int `json:"user_id" gorm:"column:user_id"`
	RoleID int `json:"role_id" gorm:"column:role_id"`
}

func (m *Role) TableName() string {
	return TableNameRole
}

func (m *RoleRouters) TableName() string {
	return TableNameRoleRouters
}

func (m *UserRoles) TableName() string {
	return TableNameUserRoles
}
