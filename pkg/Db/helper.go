package Db

import (
	"github.com/jinzhu/gorm"
)

// 封装了一些常用的数据库操作

// 插入一条数据，返回插入的id
func InsertGetId(db *gorm.DB, table string, data interface{}) (uint, error) {
	var id uint
	if err := db.Table(table).Create(data).Scan(&id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

// 批量插入数据，返回插入的id
func InsertGetIds(db *gorm.DB, table string, data []interface{}) ([]uint, error) {
	var ids []uint
	if err := db.Table(table).Create(data).Scan(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// 查询数据并分页，返回未分页前的数据总数和分页后的数据
func Paginate(db *gorm.DB, table string, page int, pageSize int, where map[string]interface{}, order string) (int64, []interface{}, error) {
	var count int64
	var data []interface{}
	if err := db.Table(table).Where(where).Count(&count).Error; err != nil {
		return 0, nil, err
	}
	if err := db.Table(table).Where(where).Order(order).Offset((page - 1) * pageSize).Limit(pageSize).Find(&data).Error; err != nil {
		return 0, nil, err
	}
	return count, data, nil
}

// 判断是否为空表
func IsTableEmpty(db *gorm.DB, table string) (bool, error) {
	var count int64
	if err := db.Table(table).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

// 批量删除符合条件的数据
func DeleteAll(db *gorm.DB, table string, where map[string]interface{}) error {
	return db.Table(table).Where(where).Delete(nil).Error
}

// Sum 求和
func Sum(db *gorm.DB, table string, field string, where map[string]interface{}) (int64, error) {
	var result []int64
	var count int64
	if err := db.Table(table).Where(where).Pluck(field, &result).Error; err != nil {
		return 0, err
	}
	for _, v := range result {
		count += v
	}
	return count, nil
}

// 是否存在符合条件的数据
func IsExist(db *gorm.DB, table string, where map[string]interface{}) (bool, error) {
	var count int64
	if err := db.Table(table).Where(where).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// 根据Ids查询数据
func FindByIds(db *gorm.DB, table string, ids []int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	if err := db.Table(table).Where("id in (?)", ids).Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
