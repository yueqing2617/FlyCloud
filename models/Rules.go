package models

import (
	"FlyCloud/pkg/Db"

	"github.com/jinzhu/gorm"
)

// Rules is a struct
type Rules struct {
	ID     int    `json:"id" gorm:"primary_key"`
	Name   string `json:"label" gorm:"type:varchar(255);not null"`
	Path   string `json:"value" gorm:"type:varchar(255)"`
	Method string `json:"method" gorm:"type:varchar(255)"`
	Pid    int    `json:"pid" gorm:"type:int;not null"`
}

// TableName is a function
func (Rules) TableName() string {
	return "menu_rules"
}

// InitRulesModel is a function
func InitRulesModel(DB *gorm.DB) {
	if !DB.HasTable(&Rules{}) {
		DB.CreateTable(&Rules{})
		// 判断菜单表是否创建成功
		is, _ := Db.IsTableEmpty(DB, "menu_rules")
		if is {
			// 初始化管理员规则
			DB.Create(&Rules{ID: 1, Name: "管理员管理", Path: "/admin/admin", Method: "", Pid: 0})
			DB.Create(&Rules{ID: 2, Name: "管理员列表", Path: "/admin/admin/list", Method: "POST", Pid: 1})
			DB.Create(&Rules{ID: 3, Name: "管理员添加", Path: "/admin/admin/add", Method: "POST", Pid: 1})
			DB.Create(&Rules{ID: 4, Name: "管理员编辑", Path: "/admin/admin/edit", Method: "PUT", Pid: 1})
			DB.Create(&Rules{ID: 5, Name: "管理员删除", Path: "/admin/admin/delete", Method: "DELETE", Pid: 1})
			DB.Create(&Rules{ID: 6, Name: "管理员信息", Path: "/admin/admin/info", Method: "GET", Pid: 1})
			// 初始化角色规则
			DB.Create(&Rules{ID: 7, Name: "角色管理", Path: "/admin/roles", Method: "", Pid: 0})
			DB.Create(&Rules{ID: 8, Name: "角色列表", Path: "/admin/roles/list", Method: "POST", Pid: 7})
			DB.Create(&Rules{ID: 9, Name: "角色添加", Path: "/admin/roles/add", Method: "POST", Pid: 7})
			DB.Create(&Rules{ID: 10, Name: "角色编辑", Path: "/admin/roles/edit", Method: "PUT", Pid: 7})
			DB.Create(&Rules{ID: 11, Name: "角色删除", Path: "/admin/roles/delete", Method: "DELETE", Pid: 7})
			DB.Create(&Rules{ID: 12, Name: "角色信息", Path: "/admin/roles/info", Method: "GET", Pid: 7})
			DB.Create(&Rules{ID: 13, Name: "获取所有角色", Path: "/admin/roles/getAll", Method: "GET", Pid: 7})
			// 初始化菜单规则
			DB.Create(&Rules{ID: 14, Name: "菜单管理", Path: "/admin/menu", Method: "", Pid: 0})
			DB.Create(&Rules{ID: 15, Name: "菜单列表", Path: "/admin/menu/list", Method: "POST", Pid: 14})
			DB.Create(&Rules{ID: 16, Name: "菜单添加", Path: "/admin/menu/add", Method: "POST", Pid: 14})
			DB.Create(&Rules{ID: 17, Name: "菜单编辑", Path: "/admin/menu/edit", Method: "PUT", Pid: 14})
			DB.Create(&Rules{ID: 18, Name: "菜单删除", Path: "/admin/menu/delete", Method: "DELETE", Pid: 14})
			DB.Create(&Rules{ID: 19, Name: "菜单信息", Path: "/admin/menu/info", Method: "GET", Pid: 14})
			// 初始化存储管理规则
			DB.Create(&Rules{ID: 20, Name: "存储管理", Path: "/admin/storage", Method: "", Pid: 0})
			DB.Create(&Rules{ID: 21, Name: "存储列表", Path: "/admin/storage", Method: "POST", Pid: 20})
			DB.Create(&Rules{ID: 22, Name: "存储删除", Path: "/admin/storage/delete", Method: "Delete", Pid: 20})

			// 初始化客户管理规则
			DB.Create(&Rules{ID: 23, Name: "客户管理", Path: "/admin/customer", Method: "", Pid: 0})
			DB.Create(&Rules{ID: 24, Name: "客户列表", Path: "/admin/customer/list", Method: "POST", Pid: 23})
			DB.Create(&Rules{ID: 25, Name: "客户添加", Path: "/admin/customer/add", Method: "POST", Pid: 23})
			DB.Create(&Rules{ID: 26, Name: "客户编辑", Path: "/admin/customer/edit", Method: "PUT", Pid: 23})
			DB.Create(&Rules{ID: 27, Name: "客户删除", Path: "/admin/customer/delete", Method: "DELETE", Pid: 23})
			DB.Create(&Rules{ID: 28, Name: "客户信息", Path: "/admin/customer/info", Method: "GET", Pid: 23})
			DB.Create(&Rules{ID: 29, Name: "获取所有客户", Path: "/admin/customer/getAll", Method: "GET", Pid: 23})
			// 初始化服装管理规则
			DB.Create(&Rules{ID: 30, Name: "服装管理", Path: "/admin/clothes", Method: "", Pid: 0})
			// 初始化款式管理规则
			DB.Create(&Rules{ID: 31, Name: "款式管理", Path: "/admin/clothes/sample", Method: "", Pid: 30})
			DB.Create(&Rules{ID: 32, Name: "款式列表", Path: "/admin/clothes/sample/list", Method: "POST", Pid: 31})
			DB.Create(&Rules{ID: 33, Name: "款式添加", Path: "/admin/clothes/sample/add", Method: "POST", Pid: 31})
			DB.Create(&Rules{ID: 34, Name: "款式编辑", Path: "/admin/clothes/sample/edit", Method: "PUT", Pid: 31})
			DB.Create(&Rules{ID: 35, Name: "款式删除", Path: "/admin/clothes/sample/delete", Method: "DELETE", Pid: 31})
			DB.Create(&Rules{ID: 36, Name: "款式信息", Path: "/admin/clothes/sample/info", Method: "GET", Pid: 31})
			DB.Create(&Rules{ID: 37, Name: "获取所有款式", Path: "/admin/clothes/sample/getAll", Method: "GET", Pid: 31})
			// 初始化服装颜色规则
			DB.Create(&Rules{ID: 31, Name: "颜色管理", Path: "/admin/clothes/color", Method: "", Pid: 30})
			DB.Create(&Rules{ID: 32, Name: "颜色列表", Path: "/admin/clothes/color/list", Method: "POST", Pid: 31})
			DB.Create(&Rules{ID: 33, Name: "颜色添加", Path: "/admin/clothes/color/add", Method: "POST", Pid: 31})
			DB.Create(&Rules{ID: 34, Name: "颜色编辑", Path: "/admin/clothes/color/edit", Method: "PUT", Pid: 31})
			DB.Create(&Rules{ID: 35, Name: "颜色删除", Path: "/admin/clothes/color/delete", Method: "DELETE", Pid: 31})
			DB.Create(&Rules{ID: 37, Name: "获取所有颜色", Path: "/admin/clothes/color/getAll", Method: "GET", Pid: 31})
		}

	}
}
