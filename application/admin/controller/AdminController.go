package controller

import (
	"FlyCloud/application"
	"FlyCloud/models"
	"FlyCloud/pkg/jwt"
	"FlyCloud/pkg/response"
	"FlyCloud/pkg/system"
	"FlyCloud/serves/cache"
	"FlyCloud/serves/database"
	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// 管理员管理控制器
type AdminController interface {
	application.BaseController
}

// 管理员管理控制器实现
type adminController struct {
	Db    *gorm.DB
	Cache *bigcache.BigCache
}

// @Title Select
// @Description 查询管理员
// @param model	body	models.Admin	true	"查询条件"
// @Success 200 {object,total} []models.Admin,total int "查询结果"
// @router /admin/admins/list [post]
func (a adminController) Select(ctx *gin.Context) {
	// 获取参数
	var model models.Admin
	if err := ctx.ShouldBindJSON(&model); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 过滤条件 = 查询条件
	var where = map[string]interface{}{}
	if model.ID != 0 {
		where["id"] = model.ID
	}
	if model.Status != 0 {
		where["status"] = model.Status
	}
	if model.RolesName != "" {
		where["roles_name"] = model.RolesName
	}
	// 多条件查询
	db := a.Db.Model(models.Admin{})
	// 查询过滤 用户名 模糊查询
	if model.Username != "" {
		db = db.Where("username like ?", "%"+model.Username+"%")
	}
	// 查询过滤 用户昵称 模糊查询
	if model.Nickname != "" {
		db = db.Where("nickname like ?", "%"+model.Nickname+"%")
	}

	// 分页是否有值
	var page = 1
	if model.PageNum != 0 {
		page = model.PageNum
	}
	var pageSize = 10
	if model.PageSize != 0 {
		pageSize = model.PageSize
	}
	// 查询
	var admins []models.Admin
	var total int
	if err := db.Where(where).Count(&total).Limit(page).Offset((page - 1) * pageSize).Find(&admins).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{
		"total": total,
		"data":  admins,
	}, "查询成功")
}

// @Title Insert
// @Description 新增管理员
// @param model	body	models.Admin	true	"新增数据"
// @Success 200 {object} models.Admin "新增结果"
// @router /admin/admin/add [post]
func (a adminController) Insert(ctx *gin.Context) {
	// 获取参数
	var model models.Admin
	if err := ctx.ShouldBindJSON(&model); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 新增
	if err := a.Db.Create(&model).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{
		"data": model,
	}, "新增成功")
}

// @Title Update
// @Description 更新管理员
// @param model	body	models.Admin	true	"更新数据"
// @Success 200 {object} models.Admin "更新结果"
// @router /admin/admin/edit/:id [put]
func (a adminController) Update(ctx *gin.Context) {
	// 获取参数
	var id = ctx.Param("id")
	var model models.Admin
	if err := ctx.ShouldBindJSON(&model); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 更新
	if err := a.Db.Model(&model).Where("id = ?", id).Updates(model).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{
		"data": model,
	}, "更新成功")
}

// @Title Delete
// @Description 删除管理员
// @param id	path	int	true	"删除数据id"
// @Success 200 {object} models.Admin "删除结果"
// @router /admin/admin/delete/:id [delete]
func (a adminController) Delete(ctx *gin.Context) {
	// 获取参数
	var id = ctx.Param("id")
	// 从ctx中获取管理员信息
	claim := ctx.MustGet("claim").(*jwt.CustomClaims)
	// 判断当前当前登录的管理员是否是自己
	if claim.UserId == system.StrToUint(id) {
		response.Error(ctx, "不能删除自己", http.StatusBadRequest)
		return
	}
	// 删除
	if err := a.Db.Delete(&models.Admin{}, "id = ?", id).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{}, "删除成功")
}

// @Title Find
// @Description 根据ID查询管理员
// @param id	path	int	true	"查询数据id"
// @Success 200 {object} models.Admin "查询结果"
// @router /admin/admin/info/:id [get]
func (a adminController) Find(ctx *gin.Context) {
	// 获取参数
	var id = ctx.Param("id")
	// 查询
	var model models.Admin
	if err := a.Db.First(&model, "id = ?", id).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{
		"data": model,
	}, "查询成功")
}

// 构造函数
func NewAdminController() *adminController {
	db := database.GetDB()
	return &adminController{Db: db, Cache: cache.GetCacheObj()}
}
