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
	E, err = casbin.NewEnforcer("./conf/casbin.conf", adapter)
	if err != nil {
		return err
	}
	// 必须执行
	err = E.LoadPolicy()
	if err != nil {
		return err
	}
	initPolicy()
	//初始化匹配器
	InitEnforcer()
	return nil
}

func initPolicy() {

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

func DeleteRoleForUser(uName, rName string) error {
	_, err := E.DeleteRoleForUser(uName, rName)
	if err != nil {
		return err
	}
	return nil
}

func DeletePolicy(roleName, Uri, Method string) error {
	_, err := E.RemovePolicy(roleName, Uri, Method)
	if err != nil {
		return err
	}
	return nil
}

func AddRoleForUser(uName, rName string) error {
	_, err := E.AddRoleForUser(uName, rName)
	if err != nil {
		return err
	}
	return nil
}

func AddPolicy(roleName, Uri, Method string) error {
	_, err := E.AddPolicy(roleName, Uri, Method)
	if err != nil {
		return err
	}
	return nil
}
