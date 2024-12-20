package service

import (
	"fmt"
	"github.com/jinzhu/copier"
	log "github.com/wonderivan/logger"
	"ticket-service/api/apimodel"
	"ticket-service/database/model"
	"ticket-service/httpserver/errcode"
)

func (operator *ResourceOperator) CreateRole(req *apimodel.RoleInfoRequest) error {
	var opt model.Role
	selector := make(map[string]interface{})
	//校验
	selector[model.FieldName] = req.Name
	err := operator.Database.ListEntityByFilter(model.TableNameRole, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("route查找失败. err:[%v]", err)
		return err
	}
	var userIds []int
	if req.Users != nil {
		for _, v := range req.Users {
			userIds = append(userIds, v.ID)
		}
		var userDB []model.User
		queryParams := model.QueryParams{}
		inQueries := &model.InQuery{
			Field:  model.FieldID,
			Values: userIds,
		}
		queryParams.InQueries = append(queryParams.InQueries, inQueries)
		err = operator.Database.ListEntityByFilter(model.TableNameUser, map[string]interface{}{}, queryParams, &userDB)
		if err != nil {
			log.Error("查找 role 对应 user_id. err:[%v]", err)
			return err
		}
		if len(userDB) != len(userIds) {
			return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "users")
		}

	}
	var routerIds []int
	if req.Routers != nil {
		for _, v := range req.Routers {
			routerIds = append(routerIds, v.ID)
		}
		var routerDB []model.Routers
		queryParams := model.QueryParams{}
		inQueries := &model.InQuery{
			Field:  model.FieldID,
			Values: routerIds,
		}
		queryParams.InQueries = append(queryParams.InQueries, inQueries)
		err = operator.Database.ListEntityByFilter(model.TableNameRouters, map[string]interface{}{}, queryParams, &routerDB)
		if err != nil {
			log.Error("查找 role 对应 router_id. err:[%v]", err)
			return err
		}
		if len(routerDB) != len(routerIds) {
			return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "routers")
		}
	}

	if opt.ID > 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "role_name")
	}
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}

	err = operator.Database.CreateEntity(model.TableNameRole, &opt)
	if err != nil {
		log.Error("创建route. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) QueryRoleList(req *apimodel.RoleInfoRequest) (*apimodel.RoleInfoPageResponse, error) {
	var resp apimodel.RoleInfoPageResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	//id查
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	if req.Name != "" {
		selector[model.FieldName] = req.Name
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
	var role []model.Role
	err := operator.Database.CountEntityByFilter(model.TableNameRole, selector, queryParams, &count)
	if err != nil {
		return nil, err
	}
	err = operator.Database.PreloadEntityByFilter(model.TableNameRole, selector, queryParams, &role, []string{model.PreloadUser, model.PreloadRouters})
	if err != nil {
		log.Error("数据查询失败. err:[%v]", err)
		return nil, err
	}

	resp.Load(count, role)
	return &resp, nil
}

func (operator *ResourceOperator) DeleteRole(req *apimodel.RoleInfoRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameRole, selector, queryParams, &model.Routers{})
	if err != nil {
		log.Error("数据删除失败. err:[%v]", err)
		return err
	}
	selector = make(map[string]interface{})
	selector[model.FieldRoleID] = req.ID
	err = operator.Database.DeleteEntityByFilter(model.TableNameRoleRouters, selector, queryParams, &model.RoleRouters{})
	if err != nil {
		log.Error("role_users硬删除失败. err:[%v]", err)
		return err
	}
	selector = make(map[string]interface{})
	selector[model.FieldRoleID] = req.ID
	err = operator.Database.DeleteEntityByFilter(model.TableNameUserRoles, selector, queryParams, &model.RoleRouters{})
	if err != nil {
		log.Error("role_users硬删除失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) UpdateRole(req *apimodel.RoleInfoRequest) error {
	var opt model.Role
	//修改用户信息
	selector := make(map[string]interface{})
	selector[model.FieldID] = req.ID
	err := operator.Database.ListEntityByFilter(model.TableNameRole, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("查找角色失败. err:[%v]", err)
		return err
	}
	if opt.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "id")
	}
	var userIds []int
	if req.Users != nil {
		for _, v := range req.Users {
			userIds = append(userIds, v.ID)
		}
		var userDB []model.User
		queryParams := model.QueryParams{}
		inQueries := &model.InQuery{
			Field:  model.FieldID,
			Values: userIds,
		}
		queryParams.InQueries = append(queryParams.InQueries, inQueries)
		err = operator.Database.ListEntityByFilter(model.TableNameUser, map[string]interface{}{}, queryParams, &userDB)
		if err != nil {
			log.Error("查找 role 对应 user_id. err:[%v]", err)
			return err
		}
		if len(userDB) != len(userIds) {
			return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "users")
		}

	}
	var routerIds []int
	if req.Routers != nil {
		for _, v := range req.Routers {
			routerIds = append(routerIds, v.ID)
		}
		var routerDB []model.Routers
		queryParams := model.QueryParams{}
		inQueries := &model.InQuery{
			Field:  model.FieldID,
			Values: routerIds,
		}
		queryParams.InQueries = append(queryParams.InQueries, inQueries)
		err = operator.Database.ListEntityByFilter(model.TableNameRouters, map[string]interface{}{}, queryParams, &routerDB)
		if err != nil {
			log.Error("查找 role 对应 router_id. err:[%v]", err)
			return err
		}
		if len(routerDB) != len(routerIds) {
			return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "routers")
		}
	}
	//保持ID不变,暂存createTime、role字段（save方法全字段更新）
	req.ID = opt.ID
	CreateTime := opt.CreatedAt
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}
	opt.CreatedAt = CreateTime
	err = operator.Database.SaveEntity(model.TableNameRole, &opt)
	if err != nil {
		log.Error("数据更新失败. err:[%v]", err)
		return err
	}
	return nil
}
