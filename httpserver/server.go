package httpserver

import (
	"github.com/gin-gonic/gin"
	"ticket-service/api/handler"
	config "ticket-service/conf"
	"ticket-service/httpserver/app"
	"ticket-service/httpserver/middleware"
)

const (
	ApiDebug   = "debug"
	ApiVersion = ""
)

func CreateHttpServer() *gin.Engine {
	gin.SetMode(config.Conf.APP.Mode)
	engine := gin.New()
	middlewareList := []gin.HandlerFunc{
		gin.Logger(),
		// 日志组件增强，用来打印gin的入参
		middleware.RequestInfo(),
		gin.Recovery(),
	}
	// 路由注册，中间件引入
	RegisterRoutes(engine, middlewareList)
	return engine
}

func RegisterRoutes(router *gin.Engine, middlewares []gin.HandlerFunc) {
	// 为全局路由注册中间件
	router.Use(middlewares...)
	// 捕捉不允许的方法
	router.NoMethod(app.MethodNotFound)
	router.NoRoute(app.HandleNotFound)
	// 静态路由
	router.Static("/files", "./files")

	// 设置系统路径上下文
	contextPath := router.Group(config.Conf.APP.ContextPath)

	v1 := contextPath.Group(ApiVersion)
	// api接口注册鉴权中间件
	restHandler := handler.NewHandler()
	v1.POST("/login", restHandler.Login)
	v1.Use(middleware.JWTAuth())
	v1.Group("")
	{
		v1.GET("/ping", restHandler.Ping)

		v1.POST("/register", restHandler.Register)
		v1.POST("/update_user", restHandler.UpdateUserInfo)
		v1.GET("/query_user", restHandler.QueryUserList)
		v1.DELETE("/delete_user/:id", restHandler.DeleteUser)
		v1.POST("/change_password", restHandler.ChangePassword)
		v1.POST("/create_train", restHandler.CreateTrain)
		v1.GET("/query_train", restHandler.QueryTrainList)
		v1.POST("/update_train", restHandler.UpdateTrain)
		v1.DELETE("/delete_train/:id", restHandler.DeleteTrain)
	}
	//
	if config.Conf.APP.Mode == gin.DebugMode {
		debug := contextPath.Group(ApiDebug)
		debug.Group("")
	}
}
