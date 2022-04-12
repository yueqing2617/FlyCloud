package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

var store = base64Captcha.DefaultMemStore

// Captcha 生成验证码
func MakeCaptcha() (id, code string, err error) {
	var driver base64Captcha.Driver
	// 配置验证码驱动
	config := base64Captcha.DriverString{
		Height:          46,
		Width:           120,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890",
		BgColor:         &color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	// 按名称加载字体
	driver = config.ConvertFonts()
	// 生成验证码
	captcha := base64Captcha.NewCaptcha(driver, store)
	// 生成验证码ID
	id, code, err = captcha.Generate()
	// 返回验证码ID和验证码
	return id, code, err
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id, code string) bool {
	if store.Verify(id, code, true) {
		return true
	} else {
		return false
	}
}
