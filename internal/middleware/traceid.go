package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
)

const XTRACEID string = "X-Trace-Id"

func ProcessTraceID() gin.HandlerFunc {
	return func(context *gin.Context) {
		traceId := context.Request.Header.Get(XTRACEID)
		if traceId == "" {
			u4id, _ := uuid.GenerateUUID()
			traceId = u4id
		}
		context.Set(XTRACEID, traceId)
		context.Writer.Header().Set(XTRACEID, traceId)
		context.Next()
	}
}
