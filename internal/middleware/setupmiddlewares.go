package middlewares

import (
	"github.com/gin-gonic/gin"
)

func SetupMiddlewares(engine *gin.Engine) *gin.Engine {
	engine.Use(SecurityHeaders())
	engine.Use(Cors())
	engine.Use(ProcessTraceID())
	return engine
}
