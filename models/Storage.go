package models

import (
	"FlyCloud/pkg/Db"

	"github.com/jinzhu/gorm"
)

// define storage struct
type Storage struct {
	Db.Field
	Name     string `gorm:"column:name;type:varchar(255)" json:"name"`
	Location string `gorm:"column:location;type:varchar(255)" json:"location"`
	Type     string `gorm:"column:type;type:varchar(255)" json:"type"`
	UserId   uint   `gorm:"column:user_id;type:int(11)" json:"user_id"`
	Ext      string `gorm:"column:ext;type:varchar(255)" json:"ext"`
}

// TableName 设置表名
func (Storage) TableName() string {
	return "storage"
}

// InitStorageTable 初始化storage
func InitStorageTable(Db *gorm.DB) {
	// 判断表是否存在
	if !Db.HasTable(&Storage{}) {
		// 创建表
		Db.CreateTable(&Storage{})

	}
}

// 过滤空值，并生成查询条件
func (storage *Storage) Filter(Db *gorm.DB) *gorm.DB {
	if storage.Name != "" {
		Db = Db.Where("name = ?", storage.Name)
	}
	if storage.Location != "" {
		Db = Db.Where("location = ?", storage.Location)
	}
	if storage.Type != "" {
		Db = Db.Where("type = ?", storage.Type)
	}
	if storage.UserId != 0 {
		Db = Db.Where("user_id = ?", storage.UserId)
	}
	if storage.Ext != "" {
		Db = Db.Where("ext = ?", storage.Ext)
	}
	return Db
}
