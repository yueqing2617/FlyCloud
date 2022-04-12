package acs

import (
	"fmt"
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/jinzhu/gorm"
)

// 声明Enforcer全局变量
var Enforcer *casbin.Enforcer

// 初始化Enforcer
func InitEnforcer(db *gorm.DB) *casbin.Enforcer {
	fmt.Println("------------------InitEnforcer------------------")
	// 创建Adapters
	adapter := gormadapter.NewAdapterByDB(db)
	// 创建Enforcer
	Enforcer = casbin.NewEnforcer("./config/rbac_model.conf", adapter)
	// 加载策略
	Enforcer.EnableLog(true)
	fmt.Println("------------------InitEnforcer-Success-----------------")
	return Enforcer
}

// 判断用户是否有权限
func CheckPermission(role_name string, path string, method string) bool {
	fmt.Println("------------------CheckPermission------------------")
	// 判断用户是否有权限
	return Enforcer.Enforce(fmt.Sprintf("%s", role_name), path, method)
}

// 获取Enforcer
func GetEnforcer() *casbin.Enforcer {
	return Enforcer
}
