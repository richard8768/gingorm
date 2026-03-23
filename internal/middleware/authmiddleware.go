package middlewares

import (
	"gin_demo/internal/config"
	"strings"
	"time"

	"gin_demo/internal/util"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		authorization := context.Request.Header.Get("Authorization")
		if strings.Trim(authorization, "") == "" {
			util.HttpResponse(context, 500, "请登陆后再操作", nil)
			context.Abort()
			return
		}
		parts := strings.Split(authorization, " ")

		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			util.HttpResponse(context, 500, "意外的Token", nil)
			context.Abort()
			return
		}

		claims, err := util.ParseToken(parts[1], "user")
		if err != nil {
			util.HttpResponse(context, 500, "意外的Token.", nil)
			context.Abort()
			return
		}

		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			util.HttpResponse(context, 500, "登陆已过期,请重新登陆", nil)
			context.Abort()
			return
		}

		_, err = config.RedisClient.Get("UserId").Result()
		if err != nil {
			config.RedisClient.Set("UserId", claims.UserId, 1440)
			config.RedisClient.Set("UserName", claims.UserName, 1440)
		}
		context.Set("UserId", claims.UserId)
		context.Set("UserName", claims.UserName)

		context.Next()
	}
}
