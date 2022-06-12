package service

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("16849841325189456f487") //密钥
var jwtEffectTime = 2000 * time.Hour            //有效时间

type Claims struct {
	ID                 uint64
	jwt.StandardClaims //jwt-go提供的标准claim
}

func GenerateToken(id uint64) (string, error) {
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
	if err != nil {
		return nil, err
	}
	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		//claims中包含用户信息
		return claims, nil
	}

	return nil, err

}

// 根据token获得userID
func Token2ID(token string) (uint64, error) {
	claims, err := ParseToken(token)
	//超时 返回新token
	//校验错误 返回err
	if err != nil {
		// if ve, ok := err.(*jwt.ValidationError); ok {
		// 	//token 超出有效期
		// 	if ve.Errors&jwt.ValidationErrorExpired != 0 {
		// 		token, err1 := RefreshToken(token)
		// 		if err1 == nil {
		// 			return Token2ID(token)
		// 		} else {
		// 			return 0, err1
		// 		}
		// 	}
		// }
		return 0, err
	}
	return claims.ID, nil
}

func RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if c, ok := token.Claims.(*Claims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		c.StandardClaims.ExpiresAt = time.Now().Add(jwtEffectTime).Unix()
		return GenerateToken(c.ID)
	}

	return "", err

}
