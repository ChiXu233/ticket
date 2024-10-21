package service

import (
	"fmt"
	"github.com/jinzhu/copier"
	log "github.com/wonderivan/logger"
	"ticket-service/api/apimodel"
	"ticket-service/database/model"
	"ticket-service/httpserver/errcode"
	"ticket-service/pkg/utils"
	"time"
)

const layout = "2006-01-02 15:04:05"
const nilTime = "1970-01-01 00:00:00"

// 车辆信息

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

	err = operator.Database.CreateEntity(model.TableNameTrain, &opt)
	if err != nil {
		log.Error("创建列车. err:[%v]", err)
		return err
	}

	return nil
}

func (operator *ResourceOperator) QueryTrainList(req *apimodel.TrainInfoRequest) (*apimodel.TrainInfoPageResponse, error) {
	var resp apimodel.TrainInfoPageResponse
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
			log.Error("列车数据查询失败. err:[%v]", err)
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

	err = operator.Database.SaveEntity(model.TableNameTrain, &opt)
	if err != nil {
		log.Error("列车数据更新失败. err:[%v]", err)
		return err
	}
	return nil
}

// 车站信息
func (operator *ResourceOperator) CreateStation(req *apimodel.TrainStationRequest) error {
	var opt model.Station
	selector := make(map[string]interface{})
	//校验车站名
	selector[model.FieldName] = req.Name
	err := operator.Database.ListEntityByFilter(model.TableNameStation, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("车站名查找失败. err:[%v]", err)
		return err
	}
	if opt.ID > 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "车站名")
	}

	selector = make(map[string]interface{})
	selector[model.FieldStationCode] = req.Code
	err = operator.Database.ListEntityByFilter(model.TableNameStation, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("车站编码查找失败. err:[%v]", err)
		return err
	}
	if opt.ID > 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "车站编码")
	}

	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}

	err = operator.Database.CreateEntity(model.TableNameStation, &opt)
	if err != nil {
		log.Error("创建车站. err:[%v]", err)
		return err
	}

	return nil
}

func (operator *ResourceOperator) QueryStationList(req *apimodel.TrainStationRequest) (*apimodel.StationInfoPageResponse, error) {
	var resp apimodel.StationInfoPageResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	//id查
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	//站名查
	if req.Name != "" {
		selector[model.FieldName] = req.Name
	}
	if req.Code != "" {
		selector[model.FieldStationCode] = req.Code
	}
	if req.City != "" {
		selector[model.FieldStationCity] = req.City
	}
	if req.Province != "" {
		selector[model.FieldStationProvince] = req.Province
	}
	var count int64
	var stations []model.Station
	err := operator.Database.CountEntityByFilter(model.TableNameStation, selector, model.OneQuery, &count)
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
		err = operator.Database.ListEntityByFilter(model.TableNameStation, selector, queryParams, &stations)
		if err != nil {
			log.Error("车站数据查询失败. err:[%v]", err)
			return nil, err
		}
	}
	resp.Load(count, stations)
	return &resp, nil
}

func (operator *ResourceOperator) DeleteStation(req *apimodel.TrainStationRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameStation, selector, queryParams, &model.Station{})
	if err != nil {
		log.Error("车站数据删除失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) UpdateStation(req *apimodel.TrainStationRequest) error {
	var opt model.Station
	//修改车站信息
	selector := make(map[string]interface{})
	//校验车站名
	selector[model.FieldName] = req.Name
	err := operator.Database.ListEntityByFilter(model.TableNameStation, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("车站名查找失败. err:[%v]", err)
		return err
	}
	if opt.ID > 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "车站名")
	}

	selector = make(map[string]interface{})
	selector[model.FieldStationCode] = req.Code
	err = operator.Database.ListEntityByFilter(model.TableNameStation, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("车站编码查找失败. err:[%v]", err)
		return err
	}
	if opt.ID > 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "车站编码")
	}

	selector = make(map[string]interface{})
	selector[model.FieldID] = req.ID
	err = operator.Database.ListEntityByFilter(model.TableNameStation, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("查找车站失败. err:[%v]", err)
		return err
	}
	if opt.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "id")
	}

	//保持ID不变,暂存createTime（save方法全字段更新）
	req.ID = opt.ID
	CreateTime := opt.CreatedAt
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}
	opt.CreatedAt = CreateTime

	err = operator.Database.SaveEntity(model.TableNameStation, &opt)
	if err != nil {
		log.Error("列车数据更新失败. err:[%v]", err)
		return err
	}
	return nil
}

// 运行计划
//创建列车行驶计划 => [创建行驶计划,选择列车(保存) => 填写停靠信息(保存)]=>相同发车时间 => 填写座位信息(提交)

func (operator *ResourceOperator) CreateTrainSchedule(req *apimodel.TrainScheduleRequest) (int, error) {
	var opt model.TrainSchedule
	selector := make(map[string]interface{})
	//校验列车是否存在
	selector[model.FieldID] = req.TrainID
	err := operator.Database.ListEntityByFilter(model.TableNameTrain, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("列车查找失败. err:[%v]", err)
		return 0, err
	}
	if opt.ID <= 0 {
		return 0, fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "列车")
	}

	err = copier.Copy(&opt, req)
	if err != nil {
		return 0, err
	}
	//时间格式化
	opt.DepartureDate = utils.ParseTime(layout, req.DepartureDate)

	//暂时设置为默认空值
	opt.EndDate = utils.ParseTime(layout, nilTime)

	//新增行驶计划
	err = operator.Database.CreateEntity(model.TableNameTrainSchedule, &opt)
	if err != nil {
		log.Error("创建行驶计划. err:[%v]", err)
		return 0, err
	}

	return opt.ID, nil
}

// 停靠信息
func (operator *ResourceOperator) CreateTrainStopInfo(req *apimodel.TrainStopInfoRequest) error {
	var schedule model.TrainSchedule
	selector := make(map[string]interface{})
	//校验行驶计划是否存在
	// 开启事务
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("CreateOrUpdateTrainType TransactionBegin Error.err[%v]", err)
		return err
	}
	defer func() {
		_ = tx.TransactionRollback()
	}()

	selector[model.FieldID] = req.ScheduleID
	err = tx.Database.ListEntityByFilter(model.TableNameTrainSchedule, selector, model.OneQuery, &schedule)
	if err != nil {
		log.Error("行驶计划查找失败,无相同发车时间或无此列车id. err:[%v]", err)
		return err
	}
	if schedule.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "行驶计划")
	}

	//遍历req中的stopList
	var stopList []model.TrainStop
	var end_date time.Time
	for i, v := range req.TrainStopList {
		stopList = append(stopList, model.TrainStop{
			ScheduleID:    req.ScheduleID,
			StationID:     v.StationID,
			StopOrder:     v.StopOrder,
			DepartureTime: utils.ParseTime(layout, v.DepartureTime),
		})
		if i == len(req.TrainStopList)-1 {
			//同步更新列车终点站时间
			end_date = utils.ParseTime(layout, v.DepartureTime)
		}
	}

	//新增停靠信息
	err = tx.Database.BatchCreateEntity(model.TableNameTrainStop, stopList)
	if err != nil {
		log.Error("创建停靠信息. err:[%v]", err)
		return err
	}
	//更新行驶计划终点站到达时间
	selector = make(map[string]interface{})
	selector[model.FieldID] = req.ScheduleID
	err = tx.Database.UpdateEntityByFilter(model.TableNameTrainSchedule, selector, model.OneQuery, &map[string]interface{}{"end_date": end_date})
	if err != nil {
		log.Error("更新终点站到达时间失败 Error.err[%v]", err)
		return err
	}
	err = tx.TransactionCommit()
	if err != nil {
		log.Error("新增停靠信息 TransactionCommit Error.err[%v]", err)
		return err
	}
	return nil
}

// 座位
func (operator *ResourceOperator) CreateTrainSeatInfo(req *apimodel.TrainSeatInfoRequest) error {
	var schedule model.TrainSchedule
	selector := make(map[string]interface{})
	//校验行驶计划是否存在
	// 开启事务
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("CreateOrUpdateTrainType TransactionBegin Error.err[%v]", err)
		return err
	}
	defer func() {
		_ = tx.TransactionRollback()
	}()

	selector[model.FieldID] = req.ScheduleID
	err = tx.Database.ListEntityByFilter(model.TableNameTrainSchedule, selector, model.OneQuery, &schedule)
	if err != nil {
		log.Error("行驶计划查找失败,无相同发车时间或无此列车id. err:[%v]", err)
		return err
	}
	if schedule.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "行驶计划")
	}

	//遍历req中的stopList
	var seatList []model.TrainSeat
	for _, v := range req.SeatInfoList {
		seatList = append(seatList, model.TrainSeat{
			ScheduleID:  req.ScheduleID,
			SeatNums:    v.SeatNums,
			SeatNowNums: v.SeatNums,
			SeatType:    v.SeatType,
			Price:       v.Price,
		})
	}

	fmt.Println(seatList)
	//新增停靠信息
	err = tx.Database.BatchCreateEntity(model.TableNameTrainSeat, seatList)
	if err != nil {
		log.Error("创建座位信息. err:[%v]", err)
		return err
	}
	err = tx.TransactionCommit()
	if err != nil {
		log.Error("新增座位信息 TransactionCommit Error.err[%v]", err)
		return err
	}
	return nil
}
