package service

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("16849841325189456f487") //密钥
var jwtEffectTime = 2 * time.Hour               //有效时间

type Claims struct {
	ID                 uint
	jwt.StandardClaims //jwt-go提供的标准claim
}

func GenerateToken(id uint) (string, error) {
	//过期时间
	nowTime := time.Now()

	//定义payload
	claims := Claims{
		id,
		jwt.StandardClaims{
			ExpiresAt: nowTime.Add(jwtEffectTime).Unix(),
			Issuer:    "zjcy",
		},
	}
	//生成签名字符串
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//生成token
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

//解析和校验
func ParseToken(token string) (*Claims, error) {
	//进行格式的校验，并检查token是否有效
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			//claims中包含用户信息
			return claims, nil
		}
	}

	return nil, err
}
