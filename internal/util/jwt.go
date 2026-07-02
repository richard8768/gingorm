package util

import (
	"errors"
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

func GenToken(user model.Member, userType string, tokenType string) (string, error) {
	if tokenType != "login" && tokenType != "logout" {
		return "", errors.New("invalid token type")
	}
	cfg, err := config.GetJwtCfg(userType)
	if err != nil {
		return "", err
	}
	seconds := cfg.TTL
	jwtKey := cfg.Secret
	secretKey := []byte(jwtKey)

	var expiresAt *jwt.NumericDate
	if tokenType == "login" {
		expiresAt = jwt.NewNumericDate(carbon.Now().AddSeconds(seconds).ToStdTime())
	} else {

		expiresAt = jwt.NewNumericDate(carbon.Now().SubSeconds(seconds).ToStdTime())

	}
	claims := UserClaim{
		UserId:   strconv.FormatUint(user.ID, 10),
		UserName: user.MemberName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
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
