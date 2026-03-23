package router

import (
	"gin_demo/internal/config"
	"gin_demo/internal/controller"
	middlewares "gin_demo/internal/middleware"
	"gin_demo/internal/service"
	"gin_demo/internal/util"
	"strconv"

	docs "gin_demo/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(engine *gin.Engine) *gin.Engine {
	engine.NoRoute(func(context *gin.Context) {
		util.HttpResponse(context, 404, "Object not found.", nil)
	})

	engine.GET("/", controller.HelloWorldHandler)
	engine.GET("/get", controller.GetHandler)
	engine.POST("/upload", controller.UploadHandler)
	engine.GET("/alioss", controller.AliOssHandler)

	userApi := engine.Group("/user")
	{
		userService := (&service.UserService{}).Create(config.DbClient)
		userHandler := &controller.UserHandler{IUserService: userService}
		userApi.POST("/reg", userHandler.UserReg)
		userApi.POST("/login", userHandler.UserLogin)
		userApi.Use(middlewares.AuthMiddleware()).GET("/index", userHandler.UserIndex)
	}

	userAddressApi := engine.Group("/useraddress").Use(middlewares.AuthMiddleware())
	{
		userAddressService := (&service.UserAddressService{}).Create(config.DbClient)
		userAddressHandler := &controller.UserAddressHandler{IUserAddressService: userAddressService}
		userAddressApi.GET("/index", userAddressHandler.AddressList)
		userAddressApi.GET("/info", userAddressHandler.AddressInfo)
		userAddressApi.POST("/add", userAddressHandler.AddAddress)
		userAddressApi.POST("/edit", userAddressHandler.UpdateAddress)
		userAddressApi.POST("/del", userAddressHandler.DeleteAddress)
		userAddressApi.POST("/setdefault", userAddressHandler.SetDefaultAddress)
	}
	registerSwagger(engine)

	return engine
}

func registerSwagger(r gin.IRouter) {
	// API文档访问地址: http://host/swagger/index.html
	// 注解定义可参考 https://github.com/swaggo/swag#declarative-comments-format
	// 样例 https://github.com/swaggo/swag/blob/master/example/basic/api/api.go
	port := strconv.Itoa(config.GetHttpPort())
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Title = "管理后台接口"
	docs.SwaggerInfo.Description = "实现一个管理系统的后端API服务"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "127.0.0.1:" + port
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//to gen docs use the command swag init -o ./docs -pdl 3
}
