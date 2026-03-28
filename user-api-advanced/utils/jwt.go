package utils

import (
	"errors"
	"time"
	"user-api-advanced/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func jwtSecret() []byte {
	if config.AppConfig != nil && config.AppConfig.JWT.Secret != "" {
		return []byte(config.AppConfig.JWT.Secret)
	}
	return []byte("change-me-in-production")
}

// jwtExpire 从配置读取过期时长，解析失败时降级为 24h
func jwtExpire() time.Duration {
	if config.AppConfig != nil && config.AppConfig.JWT.Expire != "" {
		if d, err := time.ParseDuration(config.AppConfig.JWT.Expire); err == nil {
			return d
		}
	}
	return 24 * time.Hour
}

func GenerateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpire())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-api-advanced",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
