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
// @Failure DATA_NOT_EXIST "数据不存在"
// @router /admin/rules/list [get]
func (this *RulesControllerImpl) Select(ctx *gin.Context) {
	// 查询所有规则
	var rules []*Tree
	var total int
	if err := this.Db.Model(&models.Rules{}).Count(&total).Find(&rules).Error; err != nil {
		response.Error(ctx, "获取规则列表失败"+err.Error(), http.StatusInternalServerError)
		return
	}
	// 树形结构
	data := treeData(rules, 0)
	//data := getParentNode(rules, 18) // 测试用, 获取节点的所有上级节点
	// 返回数据
	response.Success(ctx, gin.H{
		"data":  data,
		"total": total,
	}, "获取成功")
}

// 递归排序，父节点获取子节点
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

/*
 * 传入规则列表和规则ID
 * 根据节点的父节点ID，获取节点的父节点
 * 如果父节点的父节点ID为0，则说明该节点为根节点
 * 返回 []*Tree
 */
func getParentNode(rules []*Tree, id int) []*Tree {
	var nodes []*Tree
	if reflect.ValueOf(rules).IsValid() {
		for _, v := range rules {
			if v.ID == id {
				if v.Pid != 0 {
					v.Children = append(nodes, getParentNode(rules, v.Pid)...)
				}
				nodes = append(nodes, v)
			}
		}
	}
	return nodes
}
