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

// CustomerController ...
type CustomerController interface {
	application.BaseController
	GetAll(ctx *gin.Context)
}

// CustomerControllerImpl ...
type CustomerControllerImpl struct {
	Db    *gorm.DB
	Cache *bigcache.BigCache
}

// @Title GetAll
// @Description get all Customer
// @Success 200 {data,total} data []models.Customer,total int
// @Failure 500
// @router /admin/Customer/getAll [get]
func (c *CustomerControllerImpl) GetAll(ctx *gin.Context) {
	var (
		customers []models.Customer
		total     int
	)
	if err := c.Db.Model(&models.Customer{}).Count(&total).Error; err != nil {
		response.Error(ctx, "服务器错误："+err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(ctx, gin.H{"data": customers, "total": total}, "获取成功")
}

// @Title Insert
// @Description 新增客户
// @Param	body		body 	models.Customer	true		"body for Customer content"
// @Success 200 {object} models.Customer
// @Failure 400 Invalid page supplied
// @Failure 404 data not found
// @router /admin/Customer/add [post]
func (c *CustomerControllerImpl) Insert(ctx *gin.Context) {
	// 获取参数
	var Customer models.Customer
	if err := ctx.ShouldBindJSON(&Customer); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 客户名称不能为空
	if Customer.Name == "" {
		response.Error(ctx, "客户名称不能为空", http.StatusBadRequest)
		return
	}
	// 将客户信息放入struct中
	Customer.Status = 1
	// 插入数据库
	newCustomer := models.Customer{
		Name:    Customer.Name,
		Status:  Customer.Status,
		Email:   Customer.Email,
		Phone:   Customer.Phone,
		Address: Customer.Address,
		Company: Customer.Company,
		Notes:   Customer.Notes,
	}
	// 插入数据库
	if err := c.Db.Create(&newCustomer).Error; err != nil {
		response.Error(ctx, "新增客户失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{
		"data": newCustomer,
	}, "新增客户成功")
}

// @Title Update
// @Description 更新客户
// @Param	body		body 	models.Customer	true		"body for Customer content"
// @Success 200 {data} models.Customer
// @Failure 400 Invalid page supplied
// @Failure 404 data not found
// @router /admin/Customer/edit [put]
func (c *CustomerControllerImpl) Update(ctx *gin.Context) {
	// 从参数中获取客户id
	CustomerId := ctx.Param("id")
	// 获取参数
	var Customer models.Customer
	if err := ctx.ShouldBindJSON(&Customer); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 客户名称不能为空
	if Customer.Name == "" {
		response.Error(ctx, "客户名称不能为空", http.StatusBadRequest)
		return
	}
	// 将客户信息放入struct中
	Customer.Status = 1
	// 更新数据库
	if err := c.Db.Model(&Customer).Where("id = ?", CustomerId).Updates(Customer).Error; err != nil {
		response.Error(ctx, "更新客户失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{
		"data": Customer,
	}, "更新客户成功")
}

// @Title Delete
// @Description 删除客户
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {string} delete success!
// @Failure 404 data not found
// @router /admin/Customer/delete/:id [delete]
func (c *CustomerControllerImpl) Delete(ctx *gin.Context) {
	// 从参数中获取客户id
	CustomerId := ctx.Param("id")
	// 删除数据库
	if err := c.Db.Where("id = ?", CustomerId).Delete(&models.Customer{}).Error; err != nil {
		response.Error(ctx, "删除客户失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{
		"data": "删除客户成功",
	}, "删除客户成功")
}

// @Title Find
// @Description 查询客户
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {data} models.Customer
// @Failure 404 data not found
// @router /admin/Customer/info/:id [get]
func (c *CustomerControllerImpl) Find(ctx *gin.Context) {
	// 从参数中获取客户id
	CustomerId := ctx.Param("id")
	// 查询数据库
	var Customer models.Customer
	if err := c.Db.Where("id = ?", CustomerId).First(&Customer).Error; err != nil {
		response.Error(ctx, "查询客户失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{
		"data": Customer,
	}, "查询客户成功")
}

func NewCustomerControllerImpl() *CustomerControllerImpl {
	db := database.GetDB()
	models.InitCustomerTable(db)
	return &CustomerControllerImpl{Db: db, Cache: cache.GetCacheObj()}
}

// @Title Select
// @Description 按条件查询
// @Param	body		body 	models.Customer	true		"body for Customer content"
// @Success 200 {data,total} data []models.Customer,total int  "返回数据"
// @Failure 400 查询失败
// @router general/Customer/list [post]
func (c *CustomerControllerImpl) Select(ctx *gin.Context) {
	var model models.Customer
	if err := ctx.ShouldBind(&model); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 声明查询对象
	var query = c.Db.Model(&models.Customer{})
	// 查询条件
	if model.Name != "" { // 姓名,模糊查询
		query = query.Where("name like ?", "%"+model.Name+"%")
	}
	if model.Phone != "" { // 手机号,精确查询
		query = query.Where("phone = ?", model.Phone)
	}
	if model.Email != "" { // 邮箱,精确查询
		query = query.Where("email = ?", model.Email)
	}
	if model.Company != "" { // 公司,模糊查询
		query = query.Where("company like ?", "%"+model.Company+"%")
	}
	if model.Address != "" { // 地址,模糊查询
		query = query.Where("address like ?", "%"+model.Address+"%")
	}
	if model.Status != 0 { // 状态,精确查询
		query = query.Where("status = ?", model.Status)
	}

	var data []models.Customer
	var total int
	// 查询数据
	if err := c.Db.Model(&models.Customer{}).Where(model).Count(&total).Offset((model.PageNum - 1) * model.PageSize).Limit(model.PageSize).Find(&data).Error; err != nil {
		response.Error(ctx, "查询失败:"+err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{"data": data, "total": total}, "查询成功")
}
