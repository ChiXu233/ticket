package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"ticket-service/api/apimodel"
	"ticket-service/httpserver/app"
	"ticket-service/httpserver/errcode"
	"ticket-service/pkg/utils/captcha"
)

func (handler *RestHandler) Login(c *gin.Context) {
	var req apimodel.UserInfoRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptLogin)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	param := captcha.Data{
		CaptchaId: req.PicId,
		Answer:    req.Answer,
	}
	if !captcha.Verify(param) {
		app.SendServerErrorResponse(c, errcode.ErrorMsgValidateCaptcha, nil)
		return
	}
	resp, err := handler.Operator.Login(c, req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgUnauthorized, err)
		return
	}
	app.Success(c, resp)
}

func (handler *RestHandler) Register(c *gin.Context) {
	var req apimodel.UserInfoRequest
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
	err = handler.Operator.Register(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCreateData, err)
		return
	}
	app.Success(c, nil)
}

// UpdateUserInfo 修改用户基础信息，不包括密码
func (handler *RestHandler) UpdateUserInfo(c *gin.Context) {
	var req apimodel.UserInfoRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	if req.UUID == "" {
		err = errors.New("参数验证错误[uuid]")
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = req.Valid(apimodel.ValidOptUpdate)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.UpdateUserInfo(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgUpdateData, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) DeleteUser(c *gin.Context) {
	var req apimodel.UserInfoRequest
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
	err = handler.Operator.DeleteUser(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgDeleteData, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) QueryUserList(c *gin.Context) {
	req := apimodel.UserInfoRequest{
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
	resp, err := handler.Operator.QueryUserList(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgListData, err)
		return
	}
	app.Success(c, resp)
}

func (handler *RestHandler) ChangePassword(c *gin.Context) {
	var req apimodel.UserChangePWRequest
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
	err = handler.Operator.ChangePassword(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgUpdateData, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) ResetPassword(c *gin.Context) {
	var req apimodel.UserChangePWRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.SendParameterErrorResponse(c, errcode.ErrorMsgLoadParam)
		return
	}
	err = req.Valid(apimodel.ValidOptResetPwd)
	if err != nil {
		app.SendParameterErrorResponse(c, err.Error())
		return
	}
	err = handler.Operator.ResetPassword(&req)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgUpdateData, err)
		return
	}
	app.Success(c, nil)
}

func (handler *RestHandler) GetCaptcha(c *gin.Context) {
	captchaData, err := captcha.Generate()
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgCaptcha, err)
		return
	}
	res := map[string]interface{}{
		"pic_id":   captchaData.CaptchaId,
		"pic_data": captchaData.Data,
	}
	app.Success(c, res)
}

func (handler *RestHandler) CheckCaptcha(c *gin.Context) {
	var param captcha.Data
	if err := c.ShouldBind(&param); err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgValidateParam, err)
		return
	}
	if !captcha.Verify(param) {
		app.SendServerErrorResponse(c, errcode.ErrorMsgValidateCaptcha, nil)
		return
	}
	app.Success(c, nil)
}

// FreshToken 刷新token接口
func (handler *RestHandler) FreshToken(c *gin.Context) {
	resp, err := handler.Operator.FreshToken(c)
	if err != nil {
		app.SendAuthorizedErrorResponse(c, err.Error())
		return
	}
	app.Success(c, resp)
}

func (handler *RestHandler) LogOut(c *gin.Context) {
	err := handler.Operator.LogOut(c)
	if err != nil {
		app.SendServerErrorResponse(c, errcode.ErrorMsgUserLogin, err)
		return
	}
	app.Success(c, nil)
}
