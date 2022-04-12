package app

import (
	"FlyCloud/application/admin"
	"FlyCloud/application/api"
	"FlyCloud/serves/cache"
	acs "FlyCloud/serves/casbin"
	"FlyCloud/serves/config"
	"FlyCloud/serves/database"
	"FlyCloud/serves/logging"
	"FlyCloud/serves/routers"
	"fmt"
)

// Start the application
func Start() {
	// 初始化配置
	config.InitConfig()
	// 初始化日志
	logging.InitLogger(config.Config.LoggerConfig)
	// 初始化数据库
	db := database.InitDB(config.Config.DatabaseConfig)
	// 初始化缓存
	cache.InitCache(config.Config.CacheConfig)
	// 加载Casbin
	acs.InitEnforcer(db)
	// 加载多个app的路由
	routers.Include(admin.Routes, api.Routes)
	// 初始化路由
	run := routers.Init()
	// 启动服务
	if err := run.Run(":8080"); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}
