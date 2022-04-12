package controller

import (
	"FlyCloud/application"
	"FlyCloud/models"
	"FlyCloud/pkg/Db"
	"FlyCloud/pkg/response"
	"FlyCloud/pkg/system"
	"FlyCloud/serves/cache"
	acs "FlyCloud/serves/casbin"
	"FlyCloud/serves/database"
	"github.com/allegro/bigcache"
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// RoleController ...
type RoleController interface {
	application.BaseController
	GetAllRoles(ctx *gin.Context)
}

// RoleControllerImpl ...
type RoleControllerImpl struct {
	Db    *gorm.DB
	Acs   *casbin.Enforcer
	Cache *bigcache.BigCache
}

// @Title GetAllRoles
// @Description Get all roles
// @Success 200 {data,total} data  []models.Role, total int "success"
// @router /admin/roles/getAllRoles [get]
func (c *RoleControllerImpl) GetAllRoles(ctx *gin.Context) {
	var roles []models.Roles
	var total int
	var err error

	if err = c.Db.Model(models.Roles{}).Count(&total).Find(&roles).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Success(ctx, gin.H{"data": roles, "total": total}, "success")
}

// @Title Find
// @Description 根据alias查询角色
// @Param	alias	path	string	true	"角色别名"
// @Success 200 {object,ids} models.Roles,ids []uint "查询结果"
// @router /admin/roles/info/:id [get]
func (r RoleControllerImpl) Find(ctx *gin.Context) {
	// 获取参数
	var alias = ctx.Param("alias")
	// 查询
	var model models.Roles
	if err := r.Db.First(&model, "alias = ?", alias).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 从CasbinRule中获取角色别名的权限
	var casbinRule []gormadapter.CasbinRule
	if err := r.Db.Where("p_type = ?", "p").Find(&casbinRule, "v0 = ?", alias).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 以casbinRule中的v1和v2为条件，查询Rules表中的规则,并获取ids
	var ids []int
	for _, v := range casbinRule {
		var rule models.Rules
		if err := r.Db.First(&rule, "path = ? and method = ?", v.V1, v.V2).Error; err != nil {
			response.Error(ctx, err.Error(), http.StatusInternalServerError)
			return
		}
		ids = append(ids, rule.ID)
	}

	// 返回结果
	response.Success(ctx, gin.H{
		"roles": model,
		"ids":   ids,
	}, "查询成功")
}

// @Title Select
// @Description 查询角色列表
// @Param	model	query	models.Roles	false	"查询条件"
// @Success 200 {object,total} []models.Roles,total int "返回结果"
// @router /admin/roles/list [post]
func (r RoleControllerImpl) Select(ctx *gin.Context) {
	// 获取参数
	var model models.Roles
	if err := ctx.ShouldBind(&model); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 多模糊条件查询
	db := r.Db.Model(models.Roles{})
	if model.Alias != "" {
		db = db.Where("alias LIKE ?", "%"+model.Alias+"%")
	}
	if model.Name != "" {
		db = db.Where("name LIKE ?", "%"+model.Name+"%")
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
	var roles []models.Roles
	var total int
	if err := db.Count(&total).Offset((page - 1) * pageSize).Limit(pageSize).Find(&roles).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 返回结果
	response.Success(ctx, gin.H{
		"total": total,
		"data":  roles,
	}, "查询成功")
}

// @Title Insert
// @Description 新增角色
// @Param	model	body	models.Roles	true	"新增角色"
// @Param	ids	body	[]int	true	"权限id"
// @Success 200 {id,object} id uint,models.Roles "返回结果"
// @router /admin/roles/add [post]
func (r RoleControllerImpl) Insert(ctx *gin.Context) {
	// 获取参数
	var model = struct {
		models.Roles
		Ids []int `json:"ids"`
	}{}
	if err := ctx.ShouldBind(&model); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 如果角色Alias名称为空，则自动生成唯一的英文名称
	if model.Alias == "" {
		model.Alias = system.RandString(10)
	}
	// 查询是否存在相同的角色alias
	exist, _ := Db.IsExist(r.Db, "roles", map[string]interface{}{
		"alias": model.Alias,
	})
	if exist {
		response.Error(ctx, "角色Alias已存在", http.StatusBadRequest)
		return
	}

	// 新角色信息
	var newRole = models.Roles{
		Name:        model.Name,
		Alias:       model.Alias,
		Description: model.Description,
	}

	// 新增角色
	var id uint
	if err := r.Db.Model(models.Roles{}).Create(&newRole).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	// 新增角色权限
	if len(model.Ids) > 0 {
		// 根据Ids从权限菜单中获取权限Path Method,并新增角色权限
		var Permissions = make([]models.Rules, 0)
		if err := r.Db.Table("menu_rules").Where("id in (?)", model.Ids).Find(&Permissions).Error; err != nil {
			response.Error(ctx, err.Error(), http.StatusBadRequest)
			return
		}

		for _, v := range Permissions {
			if v.Pid != 0 {
				r.Acs.AddPolicy(model.Alias, v.Path, v.Method)
			}
		}

	}
	// 返回结果
	response.Success(ctx, gin.H{
		"id":   id,
		"data": model,
	}, "新增成功")
}

// @Title Update
// @Description 更新角色
// @Param	id	path	uint	true	"角色id"
// @param	ids json	[]int	true	"权限IDs"
// @Param	model	body	models.Roles	true	"更新角色"
// @Success 200 {id,object} id uint,models.Roles "返回结果"
// @router /admin/roles/edit/:id [put]
func (r RoleControllerImpl) Update(ctx *gin.Context) {
	// 获取参数
	var id = ctx.Param("id")
	var ids = []int{}
	type Model struct {
		models.Roles
		Ids []int `json:"ids"`
	}
	if err := ctx.ShouldBindUri(&id); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	var model Model
	if err := ctx.ShouldBind(&model); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	// 如果编辑的角色是超级管理员，则不允许修改
	if model.Alias == "super" {
		response.Error(ctx, "不允许修改超级管理员", http.StatusBadRequest)
		return
	}

	// 更新数据
	update := models.Roles{
		Name:        model.Name,
		Description: model.Description,
	}

	// 更新
	if err := r.Db.Model(&models.Roles{}).Where("id = ?", id).Update(&update).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	// 是否选择了权限
	if len(model.Ids) <= 0 {
		response.Error(ctx, "请选择权限", http.StatusBadRequest)
		return
	}
	// 删除原有权限
	var bak_rules []gormadapter.CasbinRule
	if err := r.Db.Model(&gormadapter.CasbinRule{}).Where("v0 = ?", model.Alias).Find(&bak_rules).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, v := range bak_rules {
		cacn := model.Alias + v.V1 + v.V2
		result := r.Acs.RemovePolicy(model.Alias, v.V1, v.V2)
		if result {
			_ = r.Cache.Delete(cacn)
		}
	}
	// 根据Ids从权限菜单中获取权限Path Method,并新增角色权限

	var rules []models.Rules
	if err := r.Db.Table("menu_rules").Where("id in (?)", model.Ids).Find(&rules).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	// 添加新的权限
	for _, v := range rules {
		if v.Pid != 0 {
			r.Acs.AddPolicy(model.Alias, v.Path, v.Method)
			ids = append(ids, v.ID)
		}
	}

	// 返回结果
	response.Success(ctx, gin.H{
		"id":   ids,
		"data": model,
	}, "更新成功")
}

// @Title Delete
// @Description 删除角色
// @Param	alias	path	uint	true	"角色alias"
// @Success 200  "返回结果"
// @router /admin/roles/delete/:id [delete]
func (r RoleControllerImpl) Delete(ctx *gin.Context) {
	// 获取参数
	var alias = ctx.Param("alias")

	// super admin 不能删除
	if alias == "super" {
		response.Error(ctx, "超级管理员不能删除", http.StatusBadRequest)
		return
	}

	// 删除
	if err := r.Db.Delete(&models.Roles{}, "alias = ?", alias).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	// 删除角色权限
	if err := Db.DeleteAll(r.Db, "casbin_rule", map[string]interface{}{
		"p_type": "p",
		"v0":     alias,
	}); err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	response.Success(ctx, gin.H{
		"alias": alias,
	}, "删除成功")
}

func NewRoleController() *RoleControllerImpl {
	db := database.GetDB()
	models.InitRolesTable(db)
	return &RoleControllerImpl{Db: db, Cache: cache.GetCacheObj(), Acs: acs.GetEnforcer()}
}
