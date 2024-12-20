package casbin

import (
	"ticket-service/database/model"
)

type RoleRel struct {
	PRole string
	Role  string
}

type RoleUserInfo struct {
	UserName string `gorm:"column:username"`
	RoleName string `gorm:"column:name"`
}

type RoleRouterInfo struct {
	RoleName string `gorm:"column:name"`
	Uri      string `gorm:"column:uri"`
	Method   string `gorm:"column:method"`
}

func (this *RoleRel) String() string {
	return this.PRole + ":" + this.Role
}

// 获取角色
func GetRoles(pid int, m *[]*RoleRel, pName string) {
	proles := make([]*model.Role, 0)
	DB.Where("pid=?", pid).Find(&proles)
	if len(proles) == 0 {
		return
	}
	for _, item := range proles {
		if pName != "" {
			*m = append(*m, &RoleRel{pName, item.Name})
		}
		GetRoles(item.ID, m, item.Name)
	}
}

// 获取用户和角色
func GetUserRoles() (userInfos []*RoleUserInfo) {
	DB.Table("user as a,user_roles b,role c").
		Select("a.username,c.name").
		Where("a.id=b.user_id and b.role_id=c.id").
		Order("a.id desc").
		Find(&userInfos)
	return
}

// 获取路由和角色
func GetRouterRoles() (routerInfos []*RoleRouterInfo) {
	DB.Table("routers as a,routers_roles b,role c").
		Select("a.uri,a.method,c.name").
		Where("a.id=b.routers_id and b.role_id=c.id").
		Order("a.id desc").
		Find(&routerInfos)
	return
}
