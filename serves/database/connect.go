package database

import (
	"FlyCloud/serves/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// 声明一个数据库连接池
var DB *gorm.DB

// 从配置文件中读取数据库连接信息，并建立连接。初始化助手函数
func InitDB(config *config.DatabaseConfig) *gorm.DB {
	fmt.Println("------------init database----------")
	var err error
	// 读取配置文件
	// 判断连接方式
	if config.Type == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
			config.Suffix)
		// 建立连接
		DB, err = gorm.Open("mysql", dsn)
	} else if config.Type == "sqlite3" {
		dsn := "file:" + config.Database + "?cache=shared&mode=rwc"
		// 建立连接
		DB, err = gorm.Open("sqlite3", dsn)
	}

	if err != nil {
		fmt.Println(err)
	}
	//defer DB.Close()
	//DB.DB().SetConnMaxLifetime()
	fmt.Println("------------init database success----------")
	return DB
}

// 获取函数
func GetDB() *gorm.DB {
	return DB
}
