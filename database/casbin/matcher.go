package casbin

import (
	"errors"
)

func InitEnforcer() {
	//自定义匹配器匹配策略
	//E.AddFunction("MethodMatch", func(arguments ...interface{}) (interface{}, error) {
	//	fmt.Println(arguments, "arguments1")
	//	if len(arguments) == 2 {
	//		key1, key2 := arguments[0].(string), arguments[1].(string)
	//		return MethodMatch(key1, key2), nil
	//	}
	//	return nil, errors.New("MethodMatch error")
	//})
	E.AddFunction("SuperMatch", func(arguments ...interface{}) (interface{}, error) {
		if len(arguments) == 1 {
			user := arguments[0].(string)
			return SuperMatch(user), nil
		}
		return nil, errors.New("SuperMatch error")
	})
}

//func MethodMatch(key1, key2 string) bool {
//	if key1 == key2 {
//		return true
//	}
//	return false
//}

// 这里我写死了，实际可以从数据库中读取
// 超级权限不用走鉴权
var roots = []string{"root", "admin", "super_admin"}

func SuperMatch(userName string) bool {
	for _, user := range roots {
		if user == userName {
			return true
		}
	}
	return false
}
