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
		middleware.Cors(),
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
	v1.POST("/register", restHandler.Register)
	//v1.Use(middleware.JWTAuth())
	v1.Group("")
	{

		v1.GET("/ping", restHandler.Ping)
		v1.POST("/update_user", restHandler.UpdateUserInfo)
		v1.GET("/query_user", restHandler.QueryUserList)
		v1.DELETE("/delete_user/:id", restHandler.DeleteUser)
		v1.POST("/change_password", restHandler.ChangePassword)
		v1.POST("/reset_password", restHandler.ResetPassword)

		//列车基本信息
		v1.POST("/create_train_info", restHandler.CreateTrain)
		v1.GET("/query_train_info", restHandler.QueryTrainList)
		v1.POST("/update_train_info", restHandler.UpdateTrain)
		v1.DELETE("/delete_train_info/:id", restHandler.DeleteTrain)

		//车站信息
		v1.POST("/create_train_station", restHandler.CreateStation)
		v1.GET("/query_train_station", restHandler.QueryStationList)
		v1.POST("/update_train_station", restHandler.UpdateStation)
		v1.DELETE("/delete_train_station/:id", restHandler.DeleteStation)

		//创建列车行驶计划
		v1.POST("/create_train_schedule", restHandler.CreateTrainSchedule)
		v1.POST("/create_train_stop", restHandler.CreateTrainStopInfo)
		v1.POST("/create_train_seat", restHandler.CreateTrainSeatInfo)

		//行驶计划
		v1.DELETE("/delete_train_schedule/:id", restHandler.DeleteTrainSchedule)
		v1.GET("/query_train_schedule", restHandler.QueryTrainScheduleList)
		v1.POST("/update_train_schedule", restHandler.UpdateTrainSchedule)

		//停靠信息
		v1.DELETE("/delete_train_stop/:id", restHandler.DeleteTrainStopInfo)
		v1.GET("/query_train_stop", restHandler.QueryTrainStopInfoList)
		v1.POST("/update_train_stop/", restHandler.UpdateTrainStopInfo)
		//
		//座位信息
		v1.DELETE("/delete_train_seat/:id", restHandler.DeleteTrainSeatInfo)
		v1.GET("/query_train_seat", restHandler.QueryTrainSeatInfoList)
		v1.POST("/update_train_seat/", restHandler.UpdateTrainSeatInfo)

		//@TODO 待测试
		//订单
		v1.POST("/create_user_order", restHandler.CreateUserOrder)
		v1.GET("/query_user_order", restHandler.QueryUserOrderList)
		v1.POST("/cancel_user_order/:uuid", restHandler.CancelUserOrder)
		v1.POST("/pay_user_order/:uuid", restHandler.PayUserOrder)
		v1.DELETE("/delete_user_order/:uuid", restHandler.DeleteUserOrder)
	}

	if config.Conf.APP.Mode == gin.DebugMode {
		debug := contextPath.Group(ApiDebug)
		debug.Group("")
	}
}
