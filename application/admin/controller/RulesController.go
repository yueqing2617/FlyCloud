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
	"reflect"
)

type Tree struct {
	models.Rules
	Children []*Tree `json:"children"`
}

//==============================================================================================
// Structure: RulesController
//==============================================================================================
type RulesController interface {
	application.BaseController
}

//==============================================================================================
// Structure: RulesControllerImpl
//==============================================================================================
type RulesControllerImpl struct {
	Db    *gorm.DB
	Cache *bigcache.BigCache
}

func NewRulesController() *RulesControllerImpl {
	db := database.GetDB()
	models.InitRulesModel(db)
	return &RulesControllerImpl{Db: db, Cache: cache.GetCacheObj()}
}

// @Title Select
// @Description 获取所有规则列表，并递归获取子规则
// @Success 200 {data,total} data []*Tree,total int "获取成功"
// @Failure 403 :id is empty
// @router /admin/rules/list [get]
func (this *RulesControllerImpl) Select(ctx *gin.Context) {
	// 查询所有规则
	var rules []*Tree
	var total int
	this.Db.Model(models.Rules{}).Count(&total).Find(&rules)
	// 递归获取子规则
	data := treeData(rules, 0)
	// 返回数据
	response.Success(ctx, gin.H{
		"data":  data,
		"total": total,
	}, "获取成功")
}

// treeData 树形数据
func treeData(rules []*Tree, parentId int) []*Tree {
	var nodes []*Tree
	if reflect.ValueOf(rules).IsValid() {
		for _, v := range rules {
			if v.Pid == parentId {
				nodes = append(nodes, v)
				v.Children = append(v.Children, treeData(rules, v.ID)...)
			}
		}
	}
	return nodes
}
