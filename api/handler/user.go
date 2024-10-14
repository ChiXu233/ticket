package handler

import (
	"github.com/gin-gonic/gin"
	"ticket-service/httpserver/app"
	"ticket-service/httpserver/errcode"
)

func (handler *RestHandler) Login(c *gin.Context) {
	resp, err := handler.Operator.Login(c)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgUnauthorized, err)
		return
	}
	app.Success(c, resp)
}
