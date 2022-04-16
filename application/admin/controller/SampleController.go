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

// 服装款式管理
type SampleController interface {
	application.BaseController
	GetAll(ctx *gin.Context)
}

// 实现接口
type sampleController struct {
	Db    *gorm.DB
	Cache *bigcache.BigCache
}

// @Title GetAll
// @Description 获取所有服装款式
// @Success 200 {data,total} data []models.Sample,total int "获取成功"
// @router /admin/clothes/sample/getALL [get]
func (s *sampleController) GetAll(ctx *gin.Context) {
	var samples []models.Sample
	var total int

	// 获取所有服装款式
	if err := s.Db.Model(&models.Sample{}).Count(&total).Find(&samples).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Success(ctx, gin.H{"data": samples, "total": total}, "获取成功")
}

// @Title Select
// @Description 按条件获取服装款式列表
// @Param body body models.Sample true "body"
// @Success 200 {data,total} data []models.Sample,total int "获取成功"
// @Failure 400 {data} string "获取失败"
// @router /admin/clothes/sample/list [post]
func (s *sampleController) Select(ctx *gin.Context) {
	var sample models.Sample
	if err := ctx.ShouldBind(&sample); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 声明查询条件
	query := s.Db.Model(&models.Sample{})
	// 按条件查询
	if sample.Name != "" { // 按名称查询,模糊查询
		query = query.Where("name like ?", "%"+sample.Name+"%")
	}
	if sample.Status != 0 { // 按状态查询
		query = query.Where("status = ?", sample.Status)
	}
	if sample.Season != "" { // 按季节查询
		query = query.Where("season = ?", sample.Season)
	}
	if sample.Year != 0 { // 按年份查询
		query = query.Where("year = ?", sample.Year)
	}
	if sample.CustomerId != 0 { // 按客户id查询
		query = query.Where("customer_id = ?", sample.CustomerId)
	}
	if sample.IsStorage != 0 { // 按是否入库查询
		query = query.Where("is_storage = ?", sample.IsStorage)
	}
	// 查询
	var samples []models.Sample
	var total int
	if err := query.Count(&total).Offset((sample.PageNum - 1) * sample.PageSize).Limit(sample.PageSize).Find(&samples).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(ctx, gin.H{"data": samples, "total": total}, "获取成功")
}

// @Title Insert
// @Description 新增服装款式
// @Param body body models.Sample true "body"
// @Success 200 {data} string "新增成功"
// @Failure 400 {data} string "新增失败"
// @router /admin/clothes/sample/add [post]
func (s *sampleController) Insert(ctx *gin.Context) {
	var sample models.Sample
	if err := ctx.ShouldBind(&sample); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 判断款式名称是否为空
	if sample.Name == "" {
		response.Error(ctx, "款式名称不能为空", http.StatusBadRequest)
		return
	}
	// 判断是否选择了客户
	if sample.CustomerId == 0 {
		response.Error(ctx, "请选择客户", http.StatusBadRequest)
		return
	}
	// 判断是否选择了季节
	if sample.Season == "" {
		response.Error(ctx, "请选择季节", http.StatusBadRequest)
		return
	}
	// 判断是否选择了年份
	if sample.Year == 0 {
		response.Error(ctx, "请选择年份", http.StatusBadRequest)
		return
	}
	// 定义新增款式的数据
	newSample := models.Sample{
		Name:       sample.Name,
		Season:     sample.Season,
		Year:       sample.Year,
		CustomerId: sample.CustomerId,
		Status:     1,
		IsStorage:  1,
		Style:      sample.Style,
		Color:      sample.Color,
		Size:       sample.Size,
		Price:      sample.Price,
		ImgSrc:     sample.ImgSrc,
	}
	// 新增
	if err := s.Db.Create(&newSample).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(ctx, gin.H{"data": newSample}, "新增成功")
}

// @Title Update
// @Description 更新服装款式
// @Param body body models.Sample true "body"
// @Success 200 {data} string "更新成功"
// @Failure 400 {data} string "更新失败"
// @router /admin/clothes/sample/edit/:id [put]
func (s *sampleController) Update(ctx *gin.Context) {
	// 获取id
	id := ctx.Param("id")
	// 获取参数
	var sample models.Sample
	if err := ctx.ShouldBind(&sample); err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 更新
	if err := s.Db.Model(&models.Sample{}).Where("id = ?", id).Updates(sample).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(ctx, nil, "更新成功")
}

// @Title Delete
// @Description 删除服装款式
// @Param id path int true "id"
// @Success 200 {data} string "删除成功"
// @Failure 400 {data} string "删除失败"
// @router /admin/clothes/sample/delete/:id [delete]
func (s *sampleController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	// 获取待删除的服装款式
	var sample models.Sample
	if err := s.Db.First(&sample, id).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	// 未入库的服装款式不能删除
	if sample.IsStorage == 1 {
		response.Error(ctx, "服装款式未入库，不能删除", http.StatusBadRequest)
		return
	}
	// 正在使用的服装款式不能删除
	if sample.Status == 1 {
		response.Error(ctx, "服装款式正在使用，不能删除", http.StatusBadRequest)
		return
	}

	if err := s.Db.Delete(&models.Sample{}, "id = ?", id).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(ctx, nil, "删除成功")
}

// @Title Find
// @Description 获取服装款式详情
// @Param id path int true "id"
// @Success 200 {data} models.Sample "获取成功"
// @Failure 400 {data} string "获取失败"
// @router /admin/clothes/sample/info/:id [get]
func (s *sampleController) Find(ctx *gin.Context) {
	id := ctx.Param("id")
	// 获取待删除的服装款式
	var sample models.Sample
	if err := s.Db.First(&sample, id).Error; err != nil {
		response.Error(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	response.Success(ctx, gin.H{"data": sample}, "获取成功")
}

func NewSampleController() *sampleController {
	db := database.GetDB()
	models.InitSampleTable(db)
	return &sampleController{Db: db, Cache: cache.GetCacheObj()}
}
