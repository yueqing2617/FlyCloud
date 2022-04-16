package controller

import (
	"FlyCloud/application"
	"FlyCloud/models"
	"FlyCloud/pkg/response"
	"FlyCloud/serves/cache"
	"FlyCloud/serves/database"
	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// ColorController ...
type ColorController interface {
	application.BaseController
	GetAll()
}

// ColorControllerImpl ...
type ColorControllerImpl struct {
	Db    *gorm.DB
	Cache *bigcache.BigCache
}

func NewColorControllerImpl() *ColorControllerImpl {
	db := database.GetDB()
	models.InitColorTable(db)
	return &ColorControllerImpl{Db: db, Cache: cache.GetCacheObj()}
}

// @Title GetAll
// @Description 获取所有颜色
// @Success 200 {data,total} data []models.Color,total int "返回的数据"
// @router /admin/clothes/color/getAll [get]
func (c *ColorControllerImpl) GetAll(ctx *gin.Context) {
	// 获取所有颜色
	var colors []models.Color
	var total int
	if err := c.Db.Model(&models.Color{}).Count(&total).Find(&colors).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(ctx, gin.H{"data": colors, "total": total}, "获取所有颜色成功")
}

// @Title Select
// @Description 按条件查询颜色
// @Param   body     body   models.Color  false    "颜色"
// @Success 200 {data,total} data []models.Color,total int "返回的数据"
// @router /admin/clothes/color/list [post]
func (c *ColorControllerImpl) Select(ctx *gin.Context) {
	// 按条件查询颜色
	var color models.Color
	if err := ctx.ShouldBind(&color); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 声明查询条件
	query := c.Db.Model(&models.Color{})
	if color.Name != "" {
		query = query.Where("name like ?", "%"+color.Name+"%")
	}

	var colors []models.Color
	var total int
	if err := query.Where(color).Count(&total).Find(&colors).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(ctx, gin.H{"data": colors, "total": total}, "按条件查询颜色成功")
}

// @Title Insert
// @Description 新增颜色
// @Param   body     body   models.Color  true    "颜色"
// @Success 200 {data} models.Color "返回的数据"
// @router /admin/clothes/color/add [post]
func (c *ColorControllerImpl) Insert(ctx *gin.Context) {
	// 新增颜色
	var color models.Color
	if err := ctx.ShouldBind(&color); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 定义颜色数据
	newColor := models.Color{
		Name:  color.Name,
		Value: color.Value,
	}
	if err := c.Db.Create(&newColor).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(ctx, gin.H{"data": color}, "新增颜色成功")
}

// @Title Update
// @Description 更新颜色
// @Param   body     body   models.Color  true    "颜色"
// @Success 200 {data} models.Color "返回的数据"
// @router /admin/clothes/color/edit/:id [put]
func (c *ColorControllerImpl) Update(ctx *gin.Context) {
	// 更新颜色
	var color models.Color
	if err := ctx.ShouldBind(&color); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 定义颜色数据
	newColor := models.Color{
		Name:  color.Name,
		Value: color.Value,
	}
	if err := c.Db.Model(&models.Color{}).Where("id = ?", ctx.Param("id")).Updates(newColor).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(ctx, gin.H{"data": color}, "更新颜色成功")
}

// @Title Delete
// @Description 删除颜色
// @Param   body     body   models.Color  true    "颜色"
// @Success 200 {data} models.Color "返回的数据"
// @router /admin/clothes/color/delete/:id [delete]
func (c *ColorControllerImpl) Delete(ctx *gin.Context) {
	// 删除颜色
	var color models.Color
	if err := ctx.ShouldBind(&color); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.Db.Delete(&color).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(ctx, gin.H{"data": color}, "删除颜色成功")
}
