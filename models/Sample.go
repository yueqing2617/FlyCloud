package models

import (
	"FlyCloud/pkg/Db"
	"github.com/jinzhu/gorm"
)

// 服装款式 struct
type Sample struct {
	Db.Field
	Name       string   `gorm:"type:varchar(100);not null" json:"name"`
	Year       int      `gorm:"type:int(4);" json:"year"`
	CustomerId int      `gorm:"type:int(4);not null" json:"customer_id"`
	Customer   Customer `gorm:"foreignkey:CustomerId" json:"customer"`
	Season     string   `gorm:"type:varchar(100);" json:"season"`
	Style      string   `gorm:"type:varchar(100);" json:"style"`
	Color      string   `gorm:"type:text;" json:"color"`
	Size       string   `gorm:"type:text;" json:"size"`
	Price      float64  `gorm:"type:decimal(10,2);" json:"price"`
	ImgSrc     string   `gorm:"type:text;" json:"img_src"`
	Status     int      `gorm:"type:int(2);" json:"status"`
	IsStorage  int      `gorm:"type:int(2);" json:"is_storage"`
}

// TableName 设置表名
func (Sample) TableName() string {
	return "sample"
}

// InitSampleTable 初始化表
func InitSampleTable(db *gorm.DB) {
	// 判断表是否存在
	if !db.HasTable(&Sample{}) {
		// 创建表
		db.CreateTable(&Sample{})
		// 创建索引
		db.Model(&Sample{}).AddIndex("idx_sample_customer_id", "customer_id")
	}
}
