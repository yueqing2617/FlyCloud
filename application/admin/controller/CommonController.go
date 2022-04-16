package controller

import (
	"FlyCloud/models"
	"FlyCloud/pkg/Db"
	"FlyCloud/pkg/captcha"
	"FlyCloud/pkg/jwt"
	"FlyCloud/pkg/md5"
	"FlyCloud/pkg/response"
	"FlyCloud/serves/cache"
	"FlyCloud/serves/database"
	"FlyCloud/serves/logging"
	"net/http"

	"github.com/allegro/bigcache"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/jinzhu/gorm"
)

// @Title CommonController
// @Description 公共操作控制器

// 定义公共操作控制器
type CommonController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	GetCaptcha(ctx *gin.Context)
	GetUserInfo(ctx *gin.Context)
}

// 定义公共操作控制器
type commonController struct {
	Db    *gorm.DB
	Cache *bigcache.BigCache
}

// @Title Login
// @Description 登录
// @Param	username,telephone	json	string	true	"用户名或手机号"
// @Param	password		json 	string	true		"密码"
// @Param	captcha			json 	string	true		"验证码"
// @Param	appid			json 	string	true		"appid"
// @Success 200 {token,userInfo} token string,userInfo gin.H "登录成功"
// @Failure 0 "登录失败"
// @router /common/login [post]
func (c commonController) Login(ctx *gin.Context) {
	// 获取参数
	type param struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Captcha  string `json:"captcha"`
		Appid    string `json:"appid"`
	}
	var p param
	if err := ctx.ShouldBindJSON(&p); err != nil {
		response.Error(ctx, "参数错误："+err.Error(), http.StatusBadRequest)
		return
	}
	// 验证用户名或密码是否为空
	if err := validation.Validate(p.Username, validation.Required); err != nil {
		response.Error(ctx, "用户名或密码不能为空", http.StatusBadRequest)
		return
	}
	if err := validation.Validate(p.Password, validation.Required); err != nil {
		response.Error(ctx, "用户名或密码不能为空", http.StatusBadRequest)
		return
	}
	// 验证验证码是否为空
	if err := validation.Validate(p.Captcha, validation.Required); err != nil {
		response.Error(ctx, "验证码不能为空", http.StatusBadRequest)
		return
	}
	// 验证验证码
	if captcha.VerifyCaptcha(p.Appid, p.Captcha) != true {
		response.Error(ctx, "验证码错误!", http.StatusBadRequest)
		return
	}
	// 声明管理员模型
	var admin models.Admin
	// 验证用户名或密码是否正确
	if err := c.Db.Table("admin").Where("username = ? or telephone = ?", p.Username, p.Username).First(&admin).Error; err != nil {
		response.Error(ctx, "用户名或密码错误!", http.StatusBadRequest)
		return
	}
	// 验证密码是否正确
	if md5.Compare(p.Password, admin.Password) != true {
		response.Error(ctx, "用户名或密码错误!", http.StatusBadRequest)
		return
	}
	// 声明jwt
	j := jwt.NewJwt()
	token, err := j.CreateToken(&admin)
	if err != nil {
		response.Error(ctx, "生成token失败!", http.StatusInternalServerError)
		logging.Error("err:", err)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{
		"token": token,
		"userInfo": gin.H{
			"id":        admin.ID,
			"username":  admin.Username,
			"telephone": admin.Telephone,
			"nickname":  admin.Nickname,
			"img_src":   admin.ImgSrc,
			"role":      admin.RolesName,
		},
	}, "登录成功")
}

// @Title GetUserInfo
// @Description 获取用户信息
// @Param	token		json 	string	true		"token"
// @Success 200 {data} data models.Admin "获取成功"
// @Failure 0 "获取失败"
// @router /common/userinfo [get]
func (c commonController) GetUserInfo(ctx *gin.Context) {
	// 从ctx中获取claim
	claim := ctx.MustGet("claim").(*jwt.CustomClaims)
	// 获取用户信息
	var admin models.Admin
	if err := c.Db.Table("admin").Where("id = ?", claim.UserId).First(&admin).Error; err != nil {
		response.Error(ctx, "获取用户信息失败!", http.StatusInternalServerError)
		return
	}
	// 构建返回数据
	var data = map[string]interface{}{}
	data["id"] = admin.ID
	data["username"] = admin.Username
	data["nickname"] = admin.Nickname
	data["img_src"] = admin.ImgSrc
	data["role_name"] = admin.RolesName
	// 返回数据
	response.Success(ctx, gin.H{
		"data": data,
	}, "获取成功！")
}

// @Title Register
// @Description 注册
// @Param	model		json	models.Admin	true	"appid"
// @Param	appid		json	string	true	"appid"
// @Param	captcha		json	string	true	"验证码"
// @Success 200 {token} "注册成功"
// @Failure 0 "注册失败"
// @router /common/register [post]
func (c commonController) Register(ctx *gin.Context) {
	// 获取参数
	type param struct {
		models.Admin
		Appid   string `json:"appid"`
		Captcha string `json:"captcha"`
	}
	// 绑定参数
	var p param
	if err := ctx.ShouldBindJSON(&p); err != nil {
		response.Error(ctx, "参数错误："+err.Error(), http.StatusBadRequest)
		return
	}
	// 验证用户名是否为空
	if err := validation.Validate(p.Username, validation.Required); err != nil {
		response.Error(ctx, "用户名不能为空！", http.StatusBadRequest)
		return
	}
	// 验证手机号是否为空
	if err := validation.Validate(p.Telephone, validation.Required); err != nil && len(p.Telephone) != 11 {
		response.Error(ctx, "手机号验证不正确", http.StatusBadRequest)
		return
	}
	// 验证密码是否为空，两次输入是否相等
	if err := validation.Validate(p.Password, validation.Required); err != nil || p.Password != p.ConfirmPassword {
		response.Error(ctx, "两次输入的密码不一致", http.StatusBadRequest)
		return
	}
	// 验证用户昵称是否为空
	if err := validation.Validate(p.Nickname, validation.Required); err != nil {
		response.Error(ctx, "用户昵称不能为空！", http.StatusBadRequest)
		return
	}
	// 验证验证码是否为空
	if err := validation.Validate(p.Captcha, validation.Required); err != nil {
		response.Error(ctx, "验证码不能为空", http.StatusBadRequest)
		return
	}
	// 验证验证码
	if captcha.VerifyCaptcha(p.Appid, p.Captcha) != true {
		response.Error(ctx, "验证码错误!", http.StatusBadRequest)
		return
	}
	// 验证用户名是否已注册
	if models.IsExistAdminByUsername(c.Db, p.Username) {
		response.Error(ctx, "该用户名已被注册！", http.StatusBadRequest)
		return
	}
	// 验证手机号是否已注册
	if models.IsExistAdminByTelephone(c.Db, p.Telephone) {
		response.Error(ctx, "该手机号已被注册！", http.StatusBadRequest)
		return
	}
	data := models.Admin{
		Username:    p.Username,
		Password:    md5.Encry(p.ConfirmPassword),
		Nickname:    p.Nickname,
		Sex:         p.Sex,
		Description: p.Description,
		Department:  p.Department,
		ImgSrc:      p.ImgSrc,
		Telephone:   p.Telephone,
		Status:      1,
		RolesName:   "admin",
	}
	add, err := Db.InsertGetId(c.Db, "admin", &data)
	if err != nil {
		response.Error(ctx, "创建用户失败："+err.Error(), http.StatusInternalServerError)
		return
	}
	// 获取token
	// 声明jwt
	j := jwt.NewJwt()
	data.ID = add
	token, err := j.CreateToken(&data)
	if err != nil {
		response.Error(ctx, "生成token失败!", http.StatusInternalServerError)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{
		"token": token,
	}, "注册成功!")
}

// @Title GetCaptcha
// @Description 获取验证码
// @Success 200 {code,rand} "获取验证码成功"
// @Failure 0 "生成验证码失败!"
// @router /common/captcha [get]
func (c commonController) GetCaptcha(ctx *gin.Context) {
	// 生成验证码
	id, code, err := captcha.MakeCaptcha()
	if err != nil {
		response.Error(ctx, "生成验证码失败!", http.StatusInternalServerError)
		return
	}
	// 返回数据
	response.Success(ctx, gin.H{
		"code": code,
		"rand": id,
	}, "获取验证码成功!")
}

// 实例化公共操作控制器
func NewCommonController() *commonController {
	db := database.GetDB()
	// 初始化Admin表
	models.InitAdminTable(db)

	return &commonController{
		Db:    db,
		Cache: cache.GetCacheObj(),
	}
}
