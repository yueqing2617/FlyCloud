package response

import "github.com/gin-gonic/gin"

// Response is the response struct
func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{"code": code, "data": data, "message": msg})
}

// Success is the success response
func Success(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, 200, 200, data, msg)
}

// Error is the error response
func Error(ctx *gin.Context, msg string, code int) {
	Response(ctx, 200, code, nil, msg)
}
