package model

type Routers struct {
	Model
	Name   string `json:"name" gorm:"column:name"`
	Uri    string `json:"uri" gorm:"column:uri"`
	Method string `json:"method" gorm:"column:method"`
	Roles  []Role `json:"roles" gorm:"many2many:routers_roles"`
}

func (m *Routers) TableName() string {
	return TableNameRouters
}
