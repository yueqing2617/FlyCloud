package middleware

import (
	"FlyCloud/pkg/jwt"
	"FlyCloud/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWTCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			response.Error(ctx, "token is empty", http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		// 声明jwt实例
		j := jwt.NewJwt()
		// 验证token
		_, claim, err := j.ParseToken(tokenString)
		if err != nil {
			response.Error(ctx, err.Error(), http.StatusUnauthorized)
			ctx.Abort()
			return
		}
		// 将验证通过的信息放入上下文
		ctx.Set("claim", claim)
		ctx.Next()
	}
}
