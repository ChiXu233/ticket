package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"ticket-service/database"
)

var E *casbin.Enforcer
var DB *gorm.DB

func InitCasbin() error {
	DB = database.GetDatabase().GetDB()
	adapter, err := gormadapter.NewAdapterByDB(DB)
	if err != nil {
		return err
	}
	e, err := casbin.NewEnforcer("./conf/casbin.conf", adapter)
	if err != nil {
		return err
	}
	// 必须执行
	err = e.LoadPolicy()
	if err != nil {
		return err
	}
	E = e
	initPolicy()
	//初始化匹配器
	InitEnforcer()
	return nil
}

func initPolicy() {
	//// 加载权限组
	//m := make([]*RoleRel, 0)
	//GetRoles(0, &m, "")
	//for _, r := range m {
	//	E.AddRoleForUser(r.PRole, r.Role)
	//}

	// 加载用户和权限
	userInfos := GetUserRoles()
	for _, user := range userInfos {
		E.AddRoleForUser(user.UserName, user.RoleName)
	}

	// 加载路由和权限
	routers := GetRouterRoles()
	for _, router := range routers {
		E.AddPolicy(router.RoleName, router.Uri, router.Method)
	}
}
