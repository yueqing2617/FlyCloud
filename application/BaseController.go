package application

import (
	"github.com/gin-gonic/gin"
)

// 构建RESTful API模型
type BaseController interface {
	Select(ctx *gin.Context)
	Insert(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Find(ctx *gin.Context)
}
