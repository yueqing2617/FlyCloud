package models

import (
	"FlyCloud/pkg/Db"
	"github.com/jinzhu/gorm"
)

type Roles struct {
	Db.Field
	Name        string `gorm:"type:varchar(25);not null;" json:"name"`
	Alias       string `gorm:"type:varchar(55);not null;unique;unique_index" json:"alias"`
	Description string `gorm:"type:text" json:"description"`
}

// TableName 设置表名
func (Roles) TableName() string {
	return "roles"
}

// 初始化检查表是否存在，不存在则创建。并新增初始数据。
func InitRolesTable(db *gorm.DB) {
	// 判断表是否存在
	if db.HasTable("roles") != true {
		// 创建表
		db.CreateTable(&Roles{})
		// 新增初始数据
		db.Create(&Roles{
			Name:        "超级管理员",
			Alias:       "super",
			Description: "超级管理员",
		})
	}
}

// 获取所有角色
func GetAllRoles(db *gorm.DB) []Roles {
	var roles []Roles
	db.Find(&roles)
	return roles
}
