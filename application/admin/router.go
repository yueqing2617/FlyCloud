package admin

import (
	"FlyCloud/application/admin/controller"
	"FlyCloud/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.Use(middleware.CorsMiddleware())
	// 注册公共控制器路由分组
	common := r.Group("/admin/common")
	{
		common_controller := controller.NewCommonController()
		common.POST("/login", common_controller.Login)
		common.POST("/register", common_controller.Register)
		common.GET("/captcha", common_controller.GetCaptcha)
	}
	// 后台API分组
	admin := r.Group("/admin")
	{
		// JWT验证中间件
		admin.Use(middleware.JWTCheck())
		// Casbin验证中间件
		admin.Use(middleware.AdminPrivilege())
		// 注册管理员控制器路由分组
		admins := admin.Group("/admin")
		{

			admins_controller := controller.NewAdminController()
			admins.POST("/add", admins_controller.Insert)
			admins.PUT("/edit/:id", admins_controller.Update)
			admins.DELETE("/delete/:id", admins_controller.Delete)
			admins.POST("/list", admins_controller.Select)
			admins.GET("/info/:id", admins_controller.Find)
		}
		// 注册角色控制器路由分组
		roles := admin.Group("/roles")
		{
			roles_controller := controller.NewRoleController()
			roles.POST("/add", roles_controller.Insert)
			roles.PUT("/edit/:id", roles_controller.Update)
			roles.DELETE("/delete/:alias", roles_controller.Delete)
			roles.POST("/list", roles_controller.Select)
			roles.GET("/info/:alias", roles_controller.Find)
			roles.GET("/getAllRoles", roles_controller.GetAllRoles)
		}
		// 注册规则控制器路由分组
		rules := admin.Group("/rules")
		{
			rules_controller := controller.NewRulesController()
			//rules.POST("/add", rules_controller.Insert)
			//rules.PUT("/edit/:id", rules_controller.Update)
			//rules.DELETE("/delete/:id", rules_controller.Delete)
			rules.GET("/list", rules_controller.Select)
			//rules.GET("/info/:id", rules_controller.Find)
		}
	}
	// 注册后台APP心跳接口
	app := r.Group("/admin/common")
	{
		app.Use(middleware.JWTCheck())
		common_controller := controller.NewCommonController()
		app.GET("/getUserInfo", common_controller.GetUserInfo)
	}
}
