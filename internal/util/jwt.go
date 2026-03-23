package util

import (
	"gin_demo/internal/config"
	"gin_demo/internal/model"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-module/carbon"
)

type UserClaim struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func GenToken(user model.Member, userType string) (string, error) {
	cfg, err := config.GetJwtCfg(userType)
	if err != nil {
		return "", err
	}
	seconds := cfg.TTL
	jwtKey := cfg.Secret
	secretKey := []byte(jwtKey)
	claims := UserClaim{
		UserId:   strconv.FormatUint(user.ID, 10),
		UserName: user.MemberName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(carbon.Now().AddSeconds(seconds).ToStdTime()),
			IssuedAt:  jwt.NewNumericDate(carbon.Now().ToStdTime()),
			NotBefore: jwt.NewNumericDate(carbon.Now().ToStdTime()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func ParseToken(tokenString string, userType string) (*UserClaim, error) {
	cfg, err := config.GetJwtCfg(userType)
	if err != nil {
		return nil, err
	}
	jwtKey := cfg.Secret
	secretKey := []byte(jwtKey)
	token, err := jwt.ParseWithClaims(tokenString, &UserClaim{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
