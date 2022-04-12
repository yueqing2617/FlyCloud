package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CORS跨域请求头
func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE,PUT")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type,Content-Length,Accept-Encoding,X-Requested-with, Origin, Authorization")
		ctx.Header("Content-Type", "application/json;charset=UTF-8")
		// 放行OPTIONS
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		ctx.Next()
	}
}
