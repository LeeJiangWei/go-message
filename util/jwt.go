package util

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type Claims struct {
	UserID uint `json:"userid"`
	jwt.StandardClaims
}

func GetJWTSecret() []byte {
	return []byte(JWTConfig.Secret)
}

func GenerateJWT(userID uint) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(JWTConfig.Expire)

	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    JWTConfig.Issuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(GetJWTSecret()) // 用 secret 生成签名
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
