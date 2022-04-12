package jwt

import (
	"FlyCloud/models"
	"FlyCloud/serves/config"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// 定义一个jwt的结构体
type Jwt struct {
	// 获取配置文件
	config *config.JwtConfig
}

// 初始化JWT对象
func NewJwt() *Jwt {
	return &Jwt{
		config: config.Config.JwtConfig,
	}
}

// 自定义有效载荷
type CustomClaims struct {
	// 自定义字段
	UserId uint `json:"userId"`
	// 这里如果不设置jwt的过期时间，那么签名就会失败
	UserRole string `json:"userRole"`
	// StandardClaims包含了jwt的一些标准信息，如生成时间，签名，过期时间等
	jwt.StandardClaims
}

// 创建一个jwt的方法
func (j *Jwt) CreateToken(obj *models.Admin) (string, error) {

	// 创建一个token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		UserId:   obj.ID,
		UserRole: obj.RolesName,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(j.config.ExpiresAt)).Unix(),
			Audience:  j.config.Audience,
			Issuer:    j.config.Issuer,
			Subject:   j.config.Subject,
		},
	})
	// 生成一个token
	return token.SignedString([]byte(j.config.PrivateKey))
}

// []byte(j.config.PrivateKey
// 解析jwt
//func (j *Jwt) ParseToken(tokenString string) (*CustomClaims, error) {
//	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return []byte(j.config.PrivateKey), nil
//	})
//
//	if err != nil {
//		// https://gowalker.org/github.com/dgrijalva/jwt-go#ValidationError
//		// jwt.ValidationError 是一个无效token的错误结构
//		if ve, ok := err.(*jwt.ValidationError); ok {
//			// ValidationErrorMalformed是一个uint常量，表示token不可用
//			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
//				return nil, fmt.Errorf("token不可用")
//				// ValidationErrorExpired表示Token过期
//			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
//				return nil, fmt.Errorf("token过期")
//				// ValidationErrorNotValidYet表示无效token
//			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
//				return nil, fmt.Errorf("无效的token")
//			} else {
//				return nil, fmt.Errorf("token不可用")
//			}
//
//		}
//	}
//
//	// 将token中的claims信息解析出来并断言成用户自定义的有效载荷结构
//	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
//		return claims, nil
//	}
//
//	return nil, fmt.Errorf("token无效")
//
//}

func (j *Jwt) ParseToken(tokenString string) (*jwt.Token, *CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.PrivateKey), nil
	})

	return token, claims, err
}
