package service

import (
	"fmt"
	"github.com/jinzhu/copier"
	log "github.com/wonderivan/logger"
	"ticket-service/api/apimodel"
	"ticket-service/database/model"
	"ticket-service/httpserver/errcode"
)

func (operator *ResourceOperator) CreateRouter(req *apimodel.RoutersInfoRequest) error {
	var opt model.Routers
	selector := make(map[string]interface{})
	//校验车次名
	selector[model.FieldMethod] = req.Method
	selector[model.FieldUri] = req.Uri
	err := operator.Database.ListEntityByFilter(model.TableNameRouters, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("route查找失败. err:[%v]", err)
		return err
	}
	if opt.ID > 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "route")
	}
	if req.Roles == nil {
		//super_admin默认赋予最高权限
		req.Roles = append(req.Roles, apimodel.PreloadRole{ID: 1, Name: "super_admin"})
	}
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}

	err = operator.Database.CreateEntity(model.TableNameRouters, &opt)
	if err != nil {
		log.Error("创建route. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) QueryRouterList(req *apimodel.RoutersInfoRequest) (*apimodel.RoutersInfoPageResponse, error) {
	var resp apimodel.RoutersInfoPageResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	//id查
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	if req.Method != "" {
		selector[model.FieldMethod] = req.Method
	}
	if req.Uri != "" {
		selector[model.FieldUri] = req.Uri
	}
	order := model.Order{
		Field:     model.FieldUpdatedTime,
		Direction: apimodel.OrderDesc,
	}
	queryParams.Orders = append(queryParams.Orders, order)
	if req.PageSize > 0 {
		queryParams.Limit = &req.PageSize
		offset := (req.PageNo - 1) * req.PageSize
		queryParams.Offset = &offset
	}

	//name模糊查询
	if req.Name != "" {
		var keyword []model.Keyword
		keyword = append(keyword, model.Keyword{Field: model.FieldName, Value: req.Name, Type: 0})
		subquery := &model.SubQuery{
			Keywords: keyword,
		}
		queryParams.SubQueries = append(queryParams.SubQueries, subquery)
	}
	if req.StartTime != "" && req.EndTime != "" {
		rangeQuery := &model.RangeQuery{
			Field: model.FieldCreatedTime,
			Start: req.StartTime,
			End:   req.EndTime,
		}
		queryParams.RangeQueries = append(queryParams.RangeQueries, rangeQuery)
	}
	var count int64
	var routers []model.Routers
	err := operator.Database.CountEntityByFilter(model.TableNameRouters, selector, queryParams, &count)
	if err != nil {
		return nil, err
	}
	err = operator.Database.PreloadEntityByFilter(model.TableNameRouters, selector, queryParams, &routers, []string{model.PreloadRoles})
	if err != nil {
		log.Error("数据查询失败. err:[%v]", err)
		return nil, err
	}

	resp.Load(count, routers)
	return &resp, nil
}

func (operator *ResourceOperator) DeleteRouter(req *apimodel.RoutersInfoRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameRouters, selector, queryParams, &model.Routers{})
	if err != nil {
		log.Error("数据删除失败. err:[%v]", err)
		return err
	}
	selector = make(map[string]interface{})
	selector[model.FieldRoutersID] = req.ID
	err = operator.Database.DeleteEntityByFilter(model.TableNameRoleRouters, selector, queryParams, &model.RoleRouters{})
	if err != nil {
		log.Error("硬删除失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) UpdateRouter(req *apimodel.RoutersInfoRequest) error {
	//todo 编辑route是不能编辑route所属角色的;想要修改权限只能去修改角色所属route
	var opt model.Routers
	//修改用户信息
	selector := make(map[string]interface{})
	selector[model.FieldID] = req.ID
	err := operator.Database.ListEntityByFilter(model.TableNameRouters, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("查找列车失败. err:[%v]", err)
		return err
	}
	if opt.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "id")
	}

	//保持ID不变,暂存createTime、role字段（save方法全字段更新）
	req.ID = opt.ID
	CreateTime := opt.CreatedAt
	roles := opt.Roles
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}
	opt.CreatedAt = CreateTime
	opt.Roles = roles
	err = operator.Database.SaveEntity(model.TableNameRouters, &opt)
	if err != nil {
		log.Error("数据更新失败. err:[%v]", err)
		return err
	}
	return nil
}
