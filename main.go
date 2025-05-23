package main

import (
	"fmt"
	log "github.com/wonderivan/logger"
	"ticket-service/api/handler"
	config "ticket-service/conf"
	"ticket-service/database"
	"ticket-service/database/casbin"
	"ticket-service/httpserver"
	"ticket-service/pkg/utils/captcha"
	"ticket-service/pkg/utils/redis"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		panic("load config with error:" + err.Error())
	}

	err = config.InitLog()
	if err != nil {
		panic("init log with error:" + err.Error())
	}
	//
	err = database.InitDB()
	if err != nil {
		panic("init database with error:" + err.Error())
	}

	err = redis.InitRedis()
	if err != nil {
		panic("init redis with error:" + err.Error())
	}

	err = casbin.InitCasbin()
	if err != nil {
		panic("init casbin with error:" + err.Error())
	}

	//读取登录验证码配置
	captcha.InitCaptchaDriver()

	err = handler.NewHandler().Operator.LoadStation_CodeMap() //车站id-name静态资源对应表，减少sql压力
	if err != nil {
		panic("init LoadStationMap with error:" + err.Error())
	}
	server := httpserver.CreateHttpServer()
	listenAddress := fmt.Sprintf("0.0.0.0:%s", config.Conf.APP.Port)
	if err = server.Run(listenAddress); err != nil {
		log.Error("ticket_service exit with error: %v", err)
	}

}
