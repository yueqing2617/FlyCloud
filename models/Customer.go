package models

import (
	"FlyCloud/pkg/Db"
	"github.com/jinzhu/gorm"
)

// Customer struct
type Customer struct {
	Db.Field
	Name    string `gorm:"type:varchar(100);not null" json:"name"`
	Email   string `gorm:"type:varchar(100);" json:"email"`
	Phone   string `gorm:"type:varchar(100);" json:"phone"`
	Address string `gorm:"type:varchar(100);" json:"address"`
	Company string `gorm:"type:varchar(100);" json:"company"`
	Notes   string `gorm:"type:varchar(100);" json:"notes"`
	Status  int    `gorm:"type:int(2);" json:"status"`
}

// TableName sets the insert table name for this struct type
func (g *Customer) TableName() string {
	return "customer"
}

// initGuestTable creates the table for the Customer model
func InitCustomerTable(db *gorm.DB) {
	// 判断表是否存在，不存在则创建
	if !db.HasTable(&Customer{}) {
		db.CreateTable(&Customer{})
	}
}
