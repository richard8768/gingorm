package middlewares

import "github.com/gin-gonic/gin"

func SecurityHeaders() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Add("X-Content-Type-Options", "nosniff")
		context.Writer.Header().Add("X-XSS-Protection", "1; mode=block")
		context.Writer.Header().Add("X-Frame-Options", "SAMEORIGIN")
		context.Writer.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubdomains;preload")
		context.Next()
	}
}
