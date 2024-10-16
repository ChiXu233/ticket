package service

import (
	"fmt"
	"github.com/jinzhu/copier"
	log "github.com/wonderivan/logger"
	"strings"
	"ticket-service/api/apimodel"
	"ticket-service/database/model"
	"ticket-service/httpserver/errcode"
	"time"
)

func (operator *ResourceOperator) CreateTrain(req *apimodel.TrainInfoRequest) error {
	var opt model.Train
	selector := make(map[string]interface{})
	//校验车次名
	selector[model.FieldName] = req.Name
	err := operator.Database.ListEntityByFilter(model.TableNameTrain, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("车次名查找失败. err:[%v]", err)
		return err
	}

	if opt.ID > 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "车次")
	}

	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}
	pc := strings.Join(req.PassCity, "-")
	opt.PassCity = pc
	opt.StartAt, _ = time.Parse("2006-01-02 15:04:05", req.StartAt)

	//load
	Load(&opt, req)
	err = operator.Database.CreateEntity(model.TableNameTrain, &opt)
	if err != nil {
		log.Error("创建列车. err:[%v]", err)
		return err
	}

	return nil
}

func (operator *ResourceOperator) QueryTrainList(req *apimodel.TrainInfoRequest) (*apimodel.TrainInfoResponse, error) {
	var resp apimodel.TrainInfoResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	//id查
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	//车号查
	if req.Name != "" {
		selector[model.FieldName] = req.Name
	}
	if req.Start != "" && req.End != "" {
		selector[model.FieldStartPosition] = req.Start
		selector[model.FieldEndPosition] = req.End
	}
	fmt.Println(selector)
	var count int64
	var trains []model.Train
	err := operator.Database.CountEntityByFilter(model.TableNameTrain, selector, model.OneQuery, &count)
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
		err = operator.Database.ListEntityByFilter(model.TableNameTrain, selector, queryParams, &trains)
		if err != nil {
			log.Error("地图数据查询失败. err:[%v]", err)
			return nil, err
		}
	}
	resp.Load(count, trains)
	return &resp, nil
}

func (operator *ResourceOperator) DeleteTrain(req *apimodel.TrainInfoRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameTrain, selector, queryParams, &model.Train{})
	if err != nil {
		log.Error("列车数据删除失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) UpdateTrain(req *apimodel.TrainInfoRequest) error {
	var opt model.Train
	//修改用户信息
	selector := make(map[string]interface{})
	selector[model.FieldID] = req.ID
	err := operator.Database.ListEntityByFilter(model.TableNameTrain, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("查找列车失败. err:[%v]", err)
		return err
	}
	if opt.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "id")
	}

	//保持ID、账号密码不变,暂存createTime（save方法全字段更新）
	req.ID = opt.ID
	CreateTime := opt.CreatedAt
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}
	opt.CreatedAt = CreateTime
	opt.PassCity = strings.Join(req.PassCity, "-")
	//load
	Load(&opt, req)

	err = operator.Database.SaveEntity(model.TableNameTrain, &opt)
	if err != nil {
		log.Error("列车数据更新失败. err:[%v]", err)
		return err
	}
	return nil
}

func Load(opt *model.Train, req *apimodel.TrainInfoRequest) {
	opt.SeaTingNums = req.SeaTing.Nums
	opt.SeaTingPrice = req.SeaTing.Price
	opt.SleepingNums = req.Sleeping.Nums
	opt.SleepingPrice = req.Sleeping.Price
	opt.HighSleepingNums = req.HighSleeping.Nums
	opt.HighSleepingPrice = req.HighSleeping.Price
	opt.BusinessNums = req.Business.Nums
	opt.BusinessPrice = req.Business.Price
}
