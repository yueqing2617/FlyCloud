package models

import (
	"FlyCloud/pkg/Db"
	"github.com/jinzhu/gorm"
)

// Color struct
type Color struct {
	Db.Field
	Name  string `gorm:"type:varchar(255);not null"`
	Value string `gorm:"type:text"`
}

// TableName set table name
func (Color) TableName() string {
	return "sample_color"
}

// InitColorTable
func InitColorTable(db *gorm.DB) {
	// 判断表是否存在
	if !db.HasTable(&Color{}) {
		db.CreateTable(&Color{})
	}
}
