package controller

import (
	"FlyCloud/application"
	"FlyCloud/models"
	"FlyCloud/pkg/jwt"
	"FlyCloud/pkg/response"
	"FlyCloud/pkg/system"
	"FlyCloud/serves/cache"
	"FlyCloud/serves/database"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// 定义存储控制器
type StorageController interface {
	application.BaseController
	UploadImage(ctx *gin.Context)
	UploadFile(ctx *gin.Context)
	GetImage(ctx *gin.Context)
	GetFile(ctx *gin.Context)
}

// 定义存储控制器
type storageController struct {
	Db    *gorm.DB
	Cache *bigcache.BigCache
}

// 实例化存储控制器
func NewStorageController() *storageController {
	db := database.GetDB()
	// 初始化储存表
	models.InitStorageTable(db)
	return &storageController{
		Db:    db,
		Cache: cache.GetCacheObj(),
	}
}

// @Title UploadImage
// @Description 上传图片
// @Success 200 {string} url string "上传成功"
// @Failure 0 "上传失败"
// @router /upload/image [post]
func (c *storageController) UploadImage(ctx *gin.Context) {
	// 从ctx中获取claims
	claim := ctx.MustGet("claim").(*jwt.CustomClaims)
	// 获取图片文件
	file, err := ctx.FormFile("file")
	if err != nil {
		response.Error(ctx, "获取图片文件失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 从系统设置中获取文件上传大小
	settings, err := models.GetSettingsByKeys(c.Db, []string{"site_upload_image_size", "site_upload_image_ext"})
	if err != nil {
		response.Error(ctx, "获取系统设置失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 判断文件大小
	if file.Size > system.StrToInt64(settings["site_upload_image_size"]) {
		response.Error(ctx, "图片文件大小超过限制："+err.Error(), http.StatusBadRequest)
		return
	}
	// 获取图片文件名
	filename := file.Filename

	// 获取文件后缀
	ext := filepath.Ext(filename)
	// 以,分割从settings中获取允许上传的文件类型
	exts := strings.Split(settings["site_upload_image_ext"], ",")
	// 去除ext前面的.,并转换成小写
	ext = strings.ToLower(ext[1:])
	// 判断文件类型
	if !system.InArray(exts, ext) {
		response.Error(ctx, "图片文件类型不允许："+err.Error(), http.StatusBadRequest)
		return
	}

	/**
	 * 判断 storage/upload/image 文件夹是否存在
	 * 如果不存在则创建
	**/
	// 获取项目根目录
	rootPath := "./storage/upload/image"
	// 判断 storage/upload/image 文件夹是否存在
	isE, _ := system.IsExist(rootPath)
	if !isE {
		// 创建文件夹
		err := system.MkDir(rootPath)
		if err != nil {
			response.Error(ctx, "创建文件夹失败："+err.Error(), http.StatusBadRequest)
			return
		}
	}
	/**
	 * 判断 storage/upload/image/{date} 文件夹是否存在
	 * 如果不存在则创建
	**/
	// 获取当前日期
	date := system.GetDate()
	// 判断 storage/upload/image/{date} 文件夹是否存在
	isE, _ = system.IsExist(rootPath + "/" + date)
	if !isE {
		// 创建文件夹
		err := system.MkDir(rootPath + "/" + date)
		if err != nil {
			response.Error(ctx, "创建文件夹失败："+err.Error(), http.StatusBadRequest)
			return
		}
	}
	/**
	 * 保存图片
	**/
	// 获取随机字符串
	randStr := system.RandString(32)
	// savePath 保存路径
	savePath := rootPath + "/" + date + "/" + randStr + "." + ext
	// 保存图片
	err = ctx.SaveUploadedFile(file, savePath)
	if err != nil {
		response.Error(ctx, "保存图片失败："+err.Error(), http.StatusBadRequest)
		return
	}
	/**
	 * 将图片保存路径保存到数据库
	**/
	// 去除savePath前面的.
	savePath = strings.Replace(savePath, ".", "", 1)

	// 创建存储对象
	storage := models.Storage{
		Name:     filename,
		Location: savePath,
		UserId:   claim.UserId,
		Ext:      ext,
		Type:     "image",
	}
	// 保存到数据库,并获取id
	err = c.Db.Create(&storage).Error
	if err != nil {
		response.Error(ctx, "保存图片路径到数据库失败："+err.Error(), http.StatusBadRequest)
		return
	}
	/**
	 * 将域名拼接到图片路径
	**/
	// 获取域名
	domain := ctx.Request.Host
	// 将域名拼接到图片路径
	url := "//" + domain + savePath
	// 返回图片路径
	response.Success(ctx, gin.H{"url": url}, url)
}

// @Title UploadFile
// @Description 上传文件
// @Success 200 {string} url string "上传成功"
// @Failure 0 "上传失败"
// @router /upload/file [post]
func (c *storageController) UploadFile(ctx *gin.Context) {
	// 从ctx中获取claims
	claim := ctx.MustGet("claim").(*jwt.CustomClaims)
	// 获取文件
	file, err := ctx.FormFile("file")
	if err != nil {
		response.Error(ctx, "获取文件失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 从系统设置中获取文件上传大小
	settings, err := models.GetSettingsByKeys(c.Db, []string{"site_upload_ext", "site_upload_file_size"})
	if err != nil {
		response.Error(ctx, "获取系统设置失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 判断文件大小
	if file.Size > system.StrToInt64(settings["site_upload_file_size"]) {
		response.Error(ctx, "图片文件大小超过限制："+err.Error(), http.StatusBadRequest)
		return
	}
	// 获取图片文件名
	filename := file.Filename

	// 获取文件后缀
	ext := filepath.Ext(filename)
	// 以,分割从settings中获取允许上传的文件类型
	exts := strings.Split(settings["site_upload_ext"], ",")
	// 去除ext前面的.,并转换成小写
	ext = strings.ToLower(ext[1:])
	// 判断文件类型
	if !system.InArray(exts, ext) {
		response.Error(ctx, "文件类型不允许："+err.Error(), http.StatusBadRequest)
		return
	}

	/**
	 * 判断 storage/upload/file 文件夹是否存在
	 * 如果不存在则创建
	**/
	// 获取项目根目录
	rootPath := "./storage/upload/file"
	// 判断 storage/upload/file 文件夹是否存在
	isE, _ := system.IsExist(rootPath)
	if !isE {
		// 创建文件夹
		err := system.MkDir(rootPath)
		if err != nil {
			response.Error(ctx, "创建文件夹失败："+err.Error(), http.StatusBadRequest)
			return
		}
	}
	/**
	 * 判断 storage/upload/file/{date} 文件夹是否存在
	 * 如果不存在则创建
	**/
	// 获取当前日期
	date := system.GetDate()
	// 判断 storage/upload/file/{date} 文件夹是否存在
	isE, _ = system.IsExist(rootPath + "/" + date)
	if !isE {
		// 创建文件夹
		err := system.MkDir(rootPath + "/" + date)
		if err != nil {
			response.Error(ctx, "创建文件夹失败："+err.Error(), http.StatusBadRequest)
			return
		}
	}
	/**
	 * 保存文件
	**/
	// 获取随机字符串
	randStr := system.RandString(32)
	// savePath 保存路径
	savePath := rootPath + "/" + date + "/" + randStr + "." + ext
	// 保存图片
	err = ctx.SaveUploadedFile(file, savePath)
	if err != nil {
		response.Error(ctx, "保存文件失败："+err.Error(), http.StatusBadRequest)
		return
	}
	/**
	 * 将图片保存路径保存到数据库
	**/
	// 去除savePath前面的.
	savePath = strings.Replace(savePath, ".", "", 1)

	// 创建存储对象
	storage := models.Storage{
		Name:     filename,
		Location: savePath,
		UserId:   claim.UserId,
		Ext:      ext,
		Type:     "file",
	}
	// 保存到数据库,并获取id
	err = c.Db.Create(&storage).Error
	if err != nil {
		response.Error(ctx, "保存文件路径到数据库失败："+err.Error(), http.StatusBadRequest)
		return
	}
	/**
	 * 将域名拼接到文件路径
	**/
	// 获取域名
	domain := ctx.Request.Host
	// 将域名拼接到文件路径
	url := "//" + domain + savePath
	// 返回图片路径
	response.Success(ctx, gin.H{"url": url}, url)
}

// @Title GetImage
// @Description 获取图片
// @Success 200 {string} url string "获取成功"
// @Failure 0 "获取失败"
// @router /image [get]
func (c *storageController) GetImage(ctx *gin.Context) {
	// 获取图片路径
	location := ctx.Query("location")
	// 获取数据库中的图片路径
	storage := models.Storage{}
	// 查询数据库
	err := c.Db.Where("location = ?", location).First(&storage).Error
	if err != nil {
		response.Error(ctx, "获取图片失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 获取图片路径
	url := storage.Location
	// 返回图片对象
	response.Success(ctx, gin.H{"url": url}, "获取图片成功")
}

// @Title GetFile
// @Description 获取文件
// @Success 200 {string} url string "获取成功"
// @Failure 0 "获取失败"
// @router /file [get]
func (c *storageController) GetFile(ctx *gin.Context) {
	// 获取文件路径
	location := ctx.Query("location")
	// 获取数据库中的文件路径
	storage := models.Storage{}
	// 查询数据库
	err := c.Db.Where("location = ?", location).First(&storage).Error
	if err != nil {
		response.Error(ctx, "获取文件失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 获取文件路径
	url := storage.Location
	// 返回文件对象
	response.Success(ctx, gin.H{"url": url}, "获取文件成功")
}

// @Title Select
// @Description 获取文件列表
// @Param model json models.Storage true "查询条件"
// @Success 200 {data,count} data []models.Storage,count int "获取成功"
// @Failure 0 "获取失败"
// @router /storage/select [get]
func (c *storageController) Select(ctx *gin.Context) {
	// 获取查询条件
	var model models.Storage
	// 获取查询条件
	if err := ctx.ShouldBindQuery(&model); err != nil {
		response.Error(ctx, "获取查询条件失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 调用 storage 模型的 Filter 方法
	db := model.Filter(c.Db)
	// 获取总数
	var count int
	var data []models.Storage
	// 查询数据，并分页
	if err := db.Count(&count).Limit(model.PageSize).Offset((model.PageNum - 1) * model.PageSize).Find(&data).Error; err != nil {
		response.Error(ctx, "获取数据失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{"data": data, "count": count}, "获取数据成功")
}

// @Title Delete
// @Description 删除文件
// @Param id path int true "文件id"
// @Success 200 {string} string "删除成功"
// @Failure 0 "删除失败"
// @router /storage/delete/:id [delete]
func (c *storageController) Delete(ctx *gin.Context) {
	// 获取文件id
	id := ctx.Param("id")
	// 创建存储对象
	storage := models.Storage{}
	// 查询数据库
	err := c.Db.Where("id = ?", id).First(&storage).Error
	if err != nil {
		response.Error(ctx, "删除失败："+err.Error(), http.StatusInternalServerError)
		return
	}
	// 删除文件
	err = os.Remove(storage.Location)
	if err != nil {
		response.Error(ctx, "删除失败："+err.Error(), http.StatusInternalServerError)
		return
	}
	// 删除数据库中的文件路径
	err = c.Db.Delete(&storage).Error
	if err != nil {
		response.Error(ctx, "删除失败："+err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回成功
	response.Success(ctx, nil, "删除成功")
}
