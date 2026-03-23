package middlewares

import "github.com/gin-gonic/gin"

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
		//context.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		context.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Add("Access-Control-Max-Age", "86400")
		context.Writer.Header().Add("Access-Control-Expose-Headers", "Authorization")
		context.Writer.Header().Add("Access-Control-Allow-Headers", "Origin, Access-Control-Request-Headers, Access-Control-Allow-Headers, cache-control, Authorization, Access-Token, Content-Type, Accept, Connection, User-Agent, Cookie,X-Requested-With,X-Request-From,X-Trace-Id,x-request-form,request-platform")
		context.Next()
	}
}
