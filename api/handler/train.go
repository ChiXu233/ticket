package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ticket-service/api/apimodel"
	"ticket-service/httpserver/app"
	"ticket-service/httpserver/errcode"
)

//车辆基本信息

func (handler *RestHandler) CreateTrain(c *gin.Context) {
	var req apimodel.TrainInfoRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptCreateOrUpdate)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.CreateTrain(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateData, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) UpdateTrain(c *gin.Context) {
	var req apimodel.TrainInfoRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptCreateOrUpdate)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.UpdateTrain(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgUpdateData, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) DeleteTrain(c *gin.Context) {
	var req apimodel.TrainInfoRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptDel)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.DeleteTrain(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
		return
	}
	app.Success(c, nil)
}
func (handler *RestHandler) QueryTrainList(c *gin.Context) {
	req := apimodel.TrainInfoRequest{
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
	resp, err := handler.Operator.QueryTrainList(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgListData, err)
		return
	}
	app.Success(c, resp)
}

//车站信息

func (handler *RestHandler) CreateStation(c *gin.Context) {
	var req apimodel.TrainStationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptCreateOrUpdate)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.CreateStation(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateData, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) UpdateStation(c *gin.Context) {
	var req apimodel.TrainStationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptCreateOrUpdate)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.UpdateStation(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgUpdateData, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) DeleteStation(c *gin.Context) {
	var req apimodel.TrainStationRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptDel)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.DeleteStation(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
		return
	}
	app.Success(c, nil)
}
func (handler *RestHandler) QueryStationList(c *gin.Context) {
	req := apimodel.TrainStationRequest{
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
	resp, err := handler.Operator.QueryStationList(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgListData, err)
		return
	}
	app.Success(c, resp)
}

//运行计划

func (handler *RestHandler) CreateTrainSchedule(c *gin.Context) {
	var req apimodel.TrainScheduleRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptCreateOrUpdate)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	scheduleID, err := handler.Operator.CreateTrainSchedule(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateData, err)
		return
	}
	resp := make(map[string]interface{})
	resp["schedule_id"] = scheduleID
	app.Success(c, resp)
}

//func (handler *RestHandler) UpdateStation(c *gin.Context) {
//	var req apimodel.TrainStationRequest
//	err := c.ShouldBindJSON(&req)
//	if err != nil {
//		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
//		return
//	}
//	err = req.Valid(apimodel.ValidOptCreateOrUpdate)
//	if err != nil {
//		app.SendParameterErrorResponse(c, err.Error())
//		return
//	}
//	err = handler.Operator.UpdateStation(&req)
//	if err != nil {
//		app.SendServerErrorResponse(c, errcode.ErrorMsgUpdateData, err)
//		return
//	}
//	app.Success(c, nil)
//}
//
//func (handler *RestHandler) DeleteStation(c *gin.Context) {
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
//	err = handler.Operator.DeleteStation(&req)
//	if err != nil {
//		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
//		return
//	}
//	app.Success(c, nil)
//}
//func (handler *RestHandler) QueryStationList(c *gin.Context) {
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
//	resp, err := handler.Operator.QueryStationList(&req)
//	if err != nil {
//		app.SendServerErrorResponse(c, errcode.ErrorMsgListData, err)
//		return
//	}
//	app.Success(c, resp)
//}

// 停靠信息
func (handler *RestHandler) CreateTrainStopInfo(c *gin.Context) {
	var req apimodel.TrainStopInfoRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptCreateOrUpdate)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.CreateTrainStopInfo(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateData, err)
		return
	}
	app.Success(c, nil)
}

// 座位
func (handler *RestHandler) CreateTrainSeatInfo(c *gin.Context) {
	var req apimodel.TrainSeatInfoRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	fmt.Println(req)
	err = req.Valid(apimodel.ValidOptCreateOrUpdate)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.CreateTrainSeatInfo(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateData, err)
		return
	}
	app.Success(c, nil)
}
