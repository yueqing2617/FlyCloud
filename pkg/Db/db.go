package Db

// 构建一个通用字段模型
type Field struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	CreatedAt XTime  `gorm:"column:create_time;type:int(36)" json:"create_time"`
	UpdatedAt XTime  `gorm:"column:update_time;type:int(36)" json:"update_time"`
	DeletedAt *XTime `sql:"index" gorm:"column:delete_time;type:int(36)" json:"delete_time"`
	PageNum   int    `gorm:"-" json:"pageNum"`
	PageSize  int    `gorm:"-" json:"pageSize"`
}
