package models

import (
	"FlyCloud/pkg/Db"
	"FlyCloud/pkg/md5"

	"github.com/jinzhu/gorm"
)

// Admin struct
type Admin struct {
	Db.Field
	Username        string `gorm:"type:varchar(100);unique_index" json:"username"`
	Password        string `gorm:"type:varchar(255)" json:"password"`
	Sex             string `gorm:"type:varchar(4);not null;DEFAULT:'未知'" json:"sex"`
	Nickname        string `gorm:"type:varchar(15)" json:"nickname"`
	Telephone       string `gorm:"type:varchar(15);not null;unique" json:"telephone"`
	Department      string `gorm:"type:varchar(45)" json:"department"`
	ImgSrc          string `gorm:"type:text" json:"img_src"`
	Description     string `gorm:"type:text" json:"description"`
	Status          int    `gorm:"type:int(1);default(1)" json:"status"`
	RolesName       string `gorm:"type:varchar(255)" json:"roles_name"`
	Roles           Roles  `gorm:"foreignKey:RolesName;association_foreignkey:Alias" json:"roles"`
	ConfirmPassword string `gorm:"-" json:"confirm_password"`
	Captcha         string `gorm:"-" json:"captcha"`
	Appid           string `gorm:"-" json:"appid"`
}

// TableName 设置表名
func (Admin) TableName() string {
	return "admin"
}

// 检查表是否存在，不存在则创建。并新增初始数据
func InitAdminTable(db *gorm.DB) {
	// 判断表是否存在
	if db.HasTable("admin") != true {
		// 创建表
		db.CreateTable(&Admin{})
		// 新增初始数据
		db.Create(&Admin{
			Username:    "admin",
			Password:    md5.Encry("123456"),
			Nickname:    "管理员",
			Telephone:   "12345678901",
			Sex:         "男",
			ImgSrc:      "https://q.qlogo.cn/g?b=qq&nk=804966813&s=640",
			Status:      1,
			RolesName:   "super",
			Department:  "管理员",
			Description: "超级管理员",
		})
	}
}

// 根据用户名判断用户是否存在
func IsExistAdminByUsername(db *gorm.DB, username string) bool {
	var admin Admin
	db.Where("username = ?", username).First(&admin)
	if admin.ID > 0 {
		return true
	}
	return false
}

// 根据手机号判断用户是否存在
func IsExistAdminByTelephone(db *gorm.DB, telephone string) bool {
	var admin Admin
	db.Where("telephone = ?", telephone).First(&admin)
	if admin.ID > 0 {
		return true
	}
	return false
}
