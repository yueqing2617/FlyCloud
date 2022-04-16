package controller

import (
	"FlyCloud/application"
	"FlyCloud/models"
	"FlyCloud/pkg/jwt"
	"FlyCloud/pkg/md5"
	"FlyCloud/pkg/response"
	"FlyCloud/pkg/system"
	"FlyCloud/serves/cache"
	"FlyCloud/serves/database"
	"net/http"

	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	db := a.Db.Model(models.Admin{}).Where("id > 0")
	// 查询条件
	if model.Username != "" {
		db = db.Where("username like ?", "%"+model.Username+"%")
	}
	if model.Nickname != "" { // 模糊查询
		db = db.Where("nickname LIKE ?", "%"+model.Nickname+"%")
	}
	if model.Telephone != "" {
		db = db.Where("telephone = ?", model.Telephone)
	}
	if model.Department != "" { // 模糊查询
		db = db.Where("department LIKE ?", "%"+model.Department+"%")
	}
	if model.Sex != "" {
		db = db.Where("sex = ?", model.Sex)
	}
	if model.Status != 0 {
		db = db.Where("status = ?", model.Status)
	}
	if model.RolesName != "" {
		db = db.Where("roles_name = ?", model.RolesName)
	}
	// 查询
	var total int
	var admins []models.Admin
	if err := db.Count(&total).Find(&admins).Error; err != nil {
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
	// 更改状态
	model.Status = 1
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
	// 如果密码不为空，则更新密码
	if model.Password != "" {
		model.Password = md5.Encry(model.Password)
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
