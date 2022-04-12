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
			DB.Create(&Rules{ID: 1, Name: "管理员管理", Path: "", Method: "", Pid: 0})
			DB.Create(&Rules{ID: 2, Name: "管理员列表", Path: "/admin/admin/list", Method: "POST", Pid: 1})
			DB.Create(&Rules{ID: 3, Name: "管理员添加", Path: "/admin/admin/add", Method: "POST", Pid: 1})
			DB.Create(&Rules{ID: 4, Name: "管理员编辑", Path: "/admin/admin/edit", Method: "PUT", Pid: 1})
			DB.Create(&Rules{ID: 5, Name: "管理员删除", Path: "/admin/admin/delete", Method: "DELETE", Pid: 1})
			DB.Create(&Rules{ID: 6, Name: "管理员信息", Path: "/admin/admin/info", Method: "GET", Pid: 1})
			// 初始化角色规则
			DB.Create(&Rules{ID: 7, Name: "角色管理", Path: "", Method: "", Pid: 0})
			DB.Create(&Rules{ID: 8, Name: "角色列表", Path: "/admin/roles/list", Method: "POST", Pid: 7})
			DB.Create(&Rules{ID: 9, Name: "角色添加", Path: "/admin/roles/add", Method: "POST", Pid: 7})
			DB.Create(&Rules{ID: 10, Name: "角色编辑", Path: "/admin/roles/edit", Method: "PUT", Pid: 7})
			DB.Create(&Rules{ID: 11, Name: "角色删除", Path: "/admin/roles/delete", Method: "DELETE", Pid: 7})
			DB.Create(&Rules{ID: 12, Name: "角色信息", Path: "/admin/roles/info", Method: "GET", Pid: 7})
			DB.Create(&Rules{ID: 13, Name: "获取所有角色", Path: "/admin/roles/getAllRoles", Method: "GET", Pid: 7})
			// 初始化菜单规则
			DB.Create(&Rules{ID: 14, Name: "菜单管理", Path: "", Method: "", Pid: 0})
			DB.Create(&Rules{ID: 15, Name: "菜单列表", Path: "/admin/menu/list", Method: "POST", Pid: 14})
			DB.Create(&Rules{ID: 16, Name: "菜单添加", Path: "/admin/menu/add", Method: "POST", Pid: 14})
			DB.Create(&Rules{ID: 17, Name: "菜单编辑", Path: "/admin/menu/edit", Method: "PUT", Pid: 14})
			DB.Create(&Rules{ID: 18, Name: "菜单删除", Path: "/admin/menu/delete", Method: "DELETE", Pid: 14})
			DB.Create(&Rules{ID: 19, Name: "菜单信息", Path: "/admin/menu/info", Method: "GET", Pid: 14})
		}

	}
}
