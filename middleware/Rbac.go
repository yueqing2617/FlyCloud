package middleware

import (
	"FlyCloud/pkg/jwt"
	"FlyCloud/pkg/response"
	"FlyCloud/serves/cache"
	acs "FlyCloud/serves/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Rbac 鉴权中间件, 前置条件为JWTCheck 中间件鉴权通过
func AdminPrivilege() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从ctx中获取用户信息
		claim := ctx.MustGet("claim").(*jwt.CustomClaims)
		// 判断claim是否为空
		if claim.UserId <= 0 {
			response.Error(ctx, "用户信息获取失败", http.StatusBadRequest)
			ctx.Abort()
			return
		} else {
			// 如果用户是超级管理员，则直接放行
			if claim.UserRole == "super" {
				ctx.Next()
				return
			} else { // 如果不是超级管理员，则需要判断用户是否有权限访问该资源
				// 实例化缓存
				ce := cache.GetCacheObj()
				// 实例化enforce实例
				enforcer := acs.GetEnforcer()
				// 从ctx中获取路由信息
				path := ctx.Request.URL.Path
				method := ctx.Request.Method
				// 定义缓存key
				cache_key := claim.UserRole + path + method
				// 判断缓存中是否存在该key
				entry, err := ce.Get(cache_key)
				if err == nil && entry != nil {
					if string(entry) == "true" {
						ctx.Next()
					} else {
						response.Error(ctx, "没有权限访问该资源", http.StatusForbidden)
						ctx.Abort()
						return
					}
				} else {
					// 加载策略
					err := enforcer.LoadPolicy()
					if err != nil {
						response.Error(ctx, "加载策略失败", http.StatusInternalServerError)
					}
					// 判断用户是否有权限访问该资源
					result, err := enforcer.EnforceSafe(claim.UserRole, path, method)
					if err != nil {
						response.Error(ctx, "权限表找不到该资源", http.StatusForbidden)
						ctx.Abort()
						return
					}
					if !result {
						ce.Set(cache_key, []byte("false"))
						response.Error(ctx, "没有权限访问该资源", http.StatusForbidden)
						ctx.Abort()
						return
					} else {
						ce.Set(cache_key, []byte("true"))
					}
					ctx.Next()
				}
			}
		}

	}
}
