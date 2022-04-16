package Db

import "time"

// 构建一个基础表字段模型
type Field struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `gorm:"column:create_time;size:255" json:"create_time"`
	UpdatedAt time.Time  `gorm:"column:update_time;size:256" json:"update_time"`
	DeletedAt *time.Time `sql:"index" gorm:"column:delete_time;size:256" json:"delete_time"`
	PageNum   int        `gorm:"-" json:"pageNum"`
	PageSize  int        `gorm:"-" json:"pageSize"`
}
