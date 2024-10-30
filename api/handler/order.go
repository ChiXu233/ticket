package handler

import (
	"github.com/gin-gonic/gin"
	"ticket-service/api/apimodel"
	"ticket-service/httpserver/app"
	"ticket-service/httpserver/errcode"
)

//@TODO 查找所有待支付订单、查找历史订单、删除订单、支付订单

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
	err = handler.Operator.CreateUserOrder(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateData, err)
		return
	}
	app.Success(c, nil)
}

//
//func (handler *RestHandler) DeleteUserOrder(c *gin.Context) {
//	var req apimodel.TrainStationRequest
//	err := c.ShouldBindUri(&req)
//	if err != nil {
//		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
//		return
//	}
//	err = req.Valid(apimodel.ValidOptDel)
//	if err != nil {
//		app.SendParameterErrorResponse(c, err.Error())
//		return
//	}
//	err = handler.Operator.DeleteUserOrder(&req)
//	if err != nil {
//		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
//		return
//	}
//	app.Success(c, nil)
//}
//
//func (handler *RestHandler) QueryUserOrderList(c *gin.Context) {
//	req := apimodel.TrainStationRequest{
//		PaginationRequest: apimodel.DefaultPaginationRequest,
//	}
//	err := c.ShouldBindQuery(&req)
//	if err != nil {
//		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
//		return
//	}
//	err = req.Valid(apimodel.ValidOptList)
//	if err != nil {
//		app.SendParameterErrorResponse(c, err.Error())
//		return
//	}
//	resp, err := handler.Operator.QueryUserOrderList(&req)
//	if err != nil {
//		app.SendServerErrorResponse(c, errcode.ErrorMsgListData, err)
//		return
//	}
//	app.Success(c, resp)
//}
//
//func (handler *RestHandler) CancelUserOrder(c *gin.Context) {
//	var req apimodel.TrainStationRequest
//	err := c.ShouldBindUri(&req)
//	if err != nil {
//		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
//		return
//	}
//	err = req.Valid(apimodel.ValidOptDel)
//	if err != nil {
//		app.SendParameterErrorResponse(c, err.Error())
//		return
//	}
//	err = handler.Operator.CancelUserOrder(&req)
//	if err != nil {
//		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
//		return
//	}
//	app.Success(c, nil)
//}
//
//func (handler *RestHandler) PayUserOrder(c *gin.Context) {
//	var req apimodel.TrainStationRequest
//	err := c.ShouldBindUri(&req)
//	if err != nil {
//		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
//		return
//	}
//	err = req.Valid(apimodel.ValidOptDel)
//	if err != nil {
//		app.SendParameterErrorResponse(c, err.Error())
//		return
//	}
//	err = handler.Operator.PayUserOrder(&req)
//	if err != nil {
//		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
//		return
//	}
//	app.Success(c, nil)
//}
