package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"

	"ticket-service/api/apimodel"
	"ticket-service/httpserver/app"
	"ticket-service/httpserver/errcode"
)

func (handler *RestHandler) CreateUserOrder(c *gin.Context) {
	var req apimodel.UserOrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptCreate)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	orderUID, err := handler.Operator.CreateUserOrder(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateData, err)
		return
	}
	app.Success(c, map[string]interface{}{"order_uuid": orderUID})
}

func (handler *RestHandler) DeleteUserOrder(c *gin.Context) {
	var req apimodel.UserOrderRequest
	var err error
	uuidStr := c.Param("uuid")
	req.UUID, err = uuid.FromString(uuidStr)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptDel)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.DeleteUserOrder(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) QueryUserOrderList(c *gin.Context) {
	req := apimodel.UserOrderRequest{
		PaginationRequest: apimodel.DefaultPaginationRequest,
	}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptList)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	resp, err := handler.Operator.QueryUserOrderList(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgListData, err)
		return
	}
	app.Success(c, resp)
}

func (handler *RestHandler) CancelUserOrder(c *gin.Context) {
	var req apimodel.UserOrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptCancel)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.CancelUserOrder(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCancel, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) PayUserOrder(c *gin.Context) {
	var req apimodel.UserOrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptCancel)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.PayUserOrder(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgPay, err)
		return
	}
	app.Success(c, nil)
}
