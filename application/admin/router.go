package admin

import (
	"FlyCloud/application/admin/controller"
	"FlyCloud/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	// 注册图片文件访问路由
	r.StaticFS("./storage", http.Dir("storage"))

	r.Use(middleware.CorsMiddleware())
	// 注册公共控制器路由分组
	common := r.Group("/admin/common")
	{
		common_controller := controller.NewCommonController()
		common.POST("/login", common_controller.Login)
		common.POST("/register", common_controller.Register)
		common.GET("/captcha", common_controller.GetCaptcha)
	}
	// 注册上传控制器路由分组
	upload := r.Group("/upload")
	{
		upload.Use(middleware.JWTCheck())
		upload_controller := controller.NewStorageController()
		upload.POST("/image", upload_controller.UploadImage)
		upload.POST("/file", upload_controller.UploadFile)
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
			roles.GET("/getAll", roles_controller.GetAllRoles)
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

		// 注册存储控制器路由分组
		storage := admin.Group("/storage")
		{
			storage_controller := controller.NewStorageController()
			storage.DELETE("/delete/:id", storage_controller.Delete)
			storage.POST("/list", storage_controller.Select)
		}
		// 注册系统设置控制器路由分组
		system := admin.Group("/setting")
		{
			system_controller := controller.NewSettingsController()
			system.PUT("/update", system_controller.UpdateSettings)
		}
		// 注册客户控制器路由分组
		customer := admin.Group("/customer")
		{
			customer_controller := controller.NewCustomerControllerImpl()
			customer.POST("/add", customer_controller.Insert)
			customer.PUT("/edit/:id", customer_controller.Update)
			customer.DELETE("/delete/:id", customer_controller.Delete)
			customer.POST("/list", customer_controller.Select)
			customer.GET("/info/:id", customer_controller.Find)
			customer.GET("/getAll", customer_controller.GetAll)
		}
		// 注册服装控制器路由分组
		clothes := admin.Group("/clothes")
		{
			// 注册服装款式控制器路由分组
			sample := clothes.Group("/sample")
			{
				sample_controller := controller.NewSampleController()
				sample.POST("/add", sample_controller.Insert)
				sample.PUT("/edit/:id", sample_controller.Update)
				sample.DELETE("/delete/:id", sample_controller.Delete)
				sample.POST("/list", sample_controller.Select)
				sample.GET("/info/:id", sample_controller.Find)
				sample.GET("/getAll", sample_controller.GetAll)
			}

			// 注册服装颜色控制器路由分组
			color := clothes.Group("/color")
			{
				color_controller := controller.NewColorControllerImpl()
				color.POST("/add", color_controller.Insert)
				color.PUT("/edit/:id", color_controller.Update)
				color.DELETE("/delete/:id", color_controller.Delete)
				color.POST("/list", color_controller.Select)
				color.GET("/getAll", color_controller.GetAll)
			}
		}
	}
	// 注册后台APP心跳接口
	app := r.Group("/admin/common")
	{
		app.Use(middleware.JWTCheck())
		common_controller := controller.NewCommonController()
		app.GET("/getUserInfo", common_controller.GetUserInfo)
	}
	// 注册系统设置控制器路由分组
	system := r.Group("/settings")
	{
		system_controller := controller.NewSettingsController()
		system.GET("/get", system_controller.GetSettings)
	}
}
