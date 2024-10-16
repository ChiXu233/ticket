package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/jinzhu/copier"
	log "github.com/wonderivan/logger"
	"ticket-service/api/apimodel"
	config "ticket-service/conf"
	"ticket-service/database/model"
	"ticket-service/httpserver/errcode"
	"ticket-service/pkg/utils"
	"time"
)

func (operator *ResourceOperator) Login(c *gin.Context, req apimodel.UserInfoRequest) (*apimodel.LoginResponse, error) {
	var opt model.User
	var resp apimodel.LoginResponse
	selector := make(map[string]interface{})
	selector[model.FieldUserName] = req.Username
	err := operator.Database.ListEntityByFilter(model.TableNameUser, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("用户名查找失败. err:[%v]", err)
		return nil, err
	}
	if opt.ID <= 0 {
		return nil, fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "用户名")
	}
	if ok := utils.BcryptCheck(req.Password, opt.Password); !ok {
		return nil, fmt.Errorf(errcode.ErrorMsgUserPassword)
	}
	j := &utils.JWT{SigningKey: []byte(config.Conf.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(utils.BaseClaims{
		UUID:     opt.UUID,
		ID:       uint64(opt.ID),
		Username: opt.Username,
		//AuthorityId: 1,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, fmt.Errorf(errcode.ErrorMsgUnauthorized)
	}
	utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
	resp.Load(opt)
	resp.Token = token
	resp.ExpireAt = claims.RegisteredClaims.ExpiresAt.Unix() * 1000
	return &resp, nil

}

func (operator *ResourceOperator) Register(req *apimodel.UserInfoRequest) error {
	var opt model.User
	selector := make(map[string]interface{})
	//校验用户名
	selector[model.FieldUserName] = req.Username
	err := operator.Database.ListEntityByFilter(model.TableNameUser, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("用户名查找失败. err:[%v]", err)
		return err
	}

	if opt.UUID != uuid.Nil {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "用户名")
	}

	//保持id自增
	req.ID = 0
	//密码哈希加密，设置uuid
	req.Password = utils.BcryptHash(req.Password)
	req.UUID = uuid.Must(uuid.NewV4())

	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}

	err = operator.Database.CreateEntity(model.TableNameUser, &opt)
	if err != nil {
		log.Error("用户注册失败. err:[%v]", err)
		return err
	}
	return nil
}
func (operator *ResourceOperator) UpdateUserInfo(req *apimodel.UserInfoRequest) error {
	var opt model.User
	//修改用户信息
	selector := make(map[string]interface{})
	selector[model.FieldID] = req.ID
	err := operator.Database.ListEntityByFilter(model.TableNameUser, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("查找用户失败. err:[%v]", err)
		return err
	}
	if opt.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "id")
	}
	if opt.UUID == uuid.Nil {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "uuid")
	}

	//保持ID、账号密码不变,暂存createTime（save方法全字段更新）
	req.ID = opt.ID
	req.UUID = opt.UUID
	req.Username = opt.Username
	req.Password = opt.Password
	CreateTime := opt.CreatedAt
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}
	opt.CreatedAt = CreateTime

	err = operator.Database.SaveEntity(model.TableNameUser, &opt)
	if err != nil {
		log.Error("用户数据更新失败. err:[%v]", err)
		return err
	}
	return nil
}
func (operator *ResourceOperator) DeleteUser(req *apimodel.UserInfoRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameUser, selector, queryParams, &model.User{})
	if err != nil {
		log.Error("用户数据删除失败. err:[%v]", err)
		return err
	}
	return nil
}
func (operator *ResourceOperator) QueryUserList(req *apimodel.UserInfoRequest) (*apimodel.UserPageResponse, error) {
	var resp apimodel.UserPageResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	//uuid查
	if req.UUID != uuid.Nil {
		selector[model.FieldUUID] = req.UUID
	}
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	//账号查
	if req.Username != "" {
		selector[model.FieldUserName] = req.Username
	}
	var count int64
	var users []model.User
	err := operator.Database.CountEntityByFilter(model.TableNameUser, selector, model.OneQuery, &count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		order := model.Order{
			Field:     model.FieldID,
			Direction: apimodel.OrderAsc,
		}
		queryParams.Orders = append(queryParams.Orders, order)
		if req.PageSize > 0 {
			queryParams.Limit = &req.PageSize
			offset := (req.PageNo - 1) * req.PageSize
			queryParams.Offset = &offset
		}
		err = operator.Database.ListEntityByFilter(model.TableNameUser, selector, queryParams, &users)
		if err != nil {
			log.Error("地图数据查询失败. err:[%v]", err)
			return nil, err
		}
	}
	resp.Load(count, users)
	return &resp, nil
}
func (operator *ResourceOperator) ChangePassword(req *apimodel.UserChangePWRequest) error {
	//验证码确认机制待添入
	var opt model.User
	selector := make(map[string]interface{})
	//selector[model.FieldID] = req.ID
	selector[model.FieldUUID] = req.UUID
	err := operator.Database.ListEntityByFilter(model.TableNameUser, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("数据查询失败. err:[%v]", err)
		return err
	}
	if opt.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "uuid")
	}
	if ok := utils.BcryptCheck(req.OldPass, opt.Password); !ok {
		return fmt.Errorf(errcode.ErrorMsgUserChangePass)
	}
	opt.Password = utils.BcryptHash(req.NewPass)
	err = operator.Database.SaveEntity(model.TableNameUser, &opt)
	if err != nil {
		log.Error("用户数据更新失败. err:[%v]", err)
		return err
	}
	return nil
}
