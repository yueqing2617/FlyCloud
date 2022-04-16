package controller

import (
	"FlyCloud/models"
	"FlyCloud/pkg/jwt"
	"FlyCloud/pkg/response"
	"FlyCloud/serves/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// 系统设置控制器
type SettingsController interface {
	GetSettings(ctx *gin.Context)
	UpdateSettings(ctx *gin.Context)
}

// 系统设置控制器
type settingsController struct {
	Db *gorm.DB
}

// 实例化系统设置控制器
func NewSettingsController() *settingsController {
	db := database.GetDB()
	// 初始化系统设置表
	models.InitSettingsTable(db)
	return &settingsController{
		Db: db,
	}
}

// @Title GetSettings
// @Description 获取所有系统设置
// @Success 200 {object} []models.Settings "获取成功"
// @Failure 0 "获取失败"
// @router /settings [get]
func (c *settingsController) GetSettings(ctx *gin.Context) {
	settings := []models.Settings{}
	if err := c.Db.Find(&settings).Error; err != nil {
		response.Error(ctx, "获取系统设置失败："+err.Error(), http.StatusBadRequest)
		return
	}

	// 构造返回数据，key的值作为键值对的键，value的值作为键值对的值
	data := make(map[string]interface{})
	for _, setting := range settings {
		data[setting.Key] = setting.Val
	}

	response.Success(ctx, gin.H{"data": data}, "获取系统设置成功")
}

// @Title 更新系统设置
// @Description 更新系统设置
// @Success 200 {string} url string "更新成功"
// @Failure 0 "更新失败"
// @router /admin/settings/update [put]
func (c *settingsController) UpdateSettings(ctx *gin.Context) {
	// 从ctx中获取claims
	claim := ctx.MustGet("claim").(*jwt.CustomClaims)
	// 判断是否为超级管理员
	if claim.UserRole != "super" {
		response.Error(ctx, "您没有权限更新系统设置", http.StatusBadRequest)
		return
	}
	settings := models.Settings{}
	if err := ctx.ShouldBindJSON(&settings); err != nil {
		response.Error(ctx, "获取系统设置失败："+err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.Db.Save(&settings).Error; err != nil {
		response.Error(ctx, "更新系统设置失败："+err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(ctx, gin.H{"data": settings}, "更新系统设置成功")
}
