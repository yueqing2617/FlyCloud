package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Option func(*gin.Engine)

var options = []Option{}

// 注册app的路由
func Include(opts ...Option) {
	options = append(options, opts...)
}

// 初始化路由
func Init() *gin.Engine {
	fmt.Println("------------init router----------")
	r := gin.New()
	for _, opt := range options {
		opt(r)
	}
	fmt.Println("------------应用完成启动----------")
	// 返回路由对象
	return r
}
