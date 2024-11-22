package service

import (
	"fmt"
	"github.com/jinzhu/copier"
	log "github.com/wonderivan/logger"
	"ticket-service/api/apimodel"
	"ticket-service/database/model"
	"ticket-service/global"
	"ticket-service/httpserver/errcode"
	"ticket-service/pkg/utils"
	"time"
)

const layout = "2006-01-02 15:04:05"
const nilTime = "1970-01-01 00:00:00"

func (operator *ResourceOperator) LoadStation_CodeMap() error {
	var station_data []model.Station
	selector := make(map[string]interface{})
	err := operator.Database.ListEntityBySelectFilter(model.TableNameStation, selector, model.QueryParams{}, &station_data, []string{"id", "name"})
	if err != nil {
		log.Error("加载车站-id表失败. err:[%v]", err)
		return err
	}
	for _, v := range station_data {
		global.StationCodeMap[v.ID] = v.Name
	}
	return nil
}

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

	if opt.ID > 0 && opt.ID != req.ID {
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
	if req.TrainType != "" {
		selector[model.FieldTrainType] = req.TrainType
	}

	var count int64
	var trains []model.Train
	err := operator.Database.CountEntityByFilter(model.TableNameTrain, selector, model.OneQuery, &count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
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

		//车号模糊查询
		if req.Name != "" {
			var keyword []model.Keyword
			keyword = append(keyword, model.Keyword{Field: model.FieldName, Value: req.Name, Type: 0})
			subquery := &model.SubQuery{
				Keywords: keyword,
			}
			queryParams.SubQueries = append(queryParams.SubQueries, subquery)
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
	//刷新StationMap表
	operator.LoadStation_CodeMap()
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
			Field:     model.FieldUpdatedTime,
			Direction: apimodel.OrderDesc,
		}
		queryParams.Orders = append(queryParams.Orders, order)
		if req.PageSize > 0 {
			queryParams.Limit = &req.PageSize
			offset := (req.PageNo - 1) * req.PageSize
			queryParams.Offset = &offset
		}
		//车站名模糊查询
		if req.Name != "" {
			var keyword []model.Keyword
			keyword = append(keyword, model.Keyword{Field: model.FieldName, Value: req.Name, Type: 0})
			subquery := &model.SubQuery{
				Keywords: keyword,
			}
			queryParams.SubQueries = append(queryParams.SubQueries, subquery)
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
	if opt.ID > 0 && opt.ID != req.ID {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "车站名")
	}

	selector = make(map[string]interface{})
	selector[model.FieldStationCode] = req.Code
	err = operator.Database.ListEntityByFilter(model.TableNameStation, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("车站编码查找失败. err:[%v]", err)
		return err
	}
	if opt.ID > 0 && opt.ID != req.ID {
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
	var train model.Train
	selector := make(map[string]interface{})

	//校验列车是否存在
	selector[model.FieldID] = req.TrainID
	err := operator.Database.ListEntityByFilter(model.TableNameTrain, selector, model.OneQuery, &train)
	if err != nil {
		log.Error("列车查找失败. err:[%v]", err)
		return 0, err
	}
	if train.ID <= 0 {
		return 0, fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "列车")
	}

	selector = make(map[string]interface{})
	selector[model.FieldTrainID] = req.TrainID
	selector[model.FieldDepartureTime] = utils.ParseTime(layout, req.DepartureDate)
	//校验是否重复添加
	err = operator.Database.ListEntityByFilter(model.TableNameTrainSchedule, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("行驶计划查找失败. err:[%v]", err)
		return 0, err
	}

	if opt.ID > 0 {
		return 0, fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "同列车同时刻形式计划")
	}

	err = copier.Copy(&opt, req)
	if err != nil {
		return 0, err
	}
	opt.TrainName = train.TrainType + train.Name

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

func (operator *ResourceOperator) UpdateTrainSchedule(req *apimodel.TrainScheduleRequest) error {
	var opt model.TrainSchedule
	var train model.Train
	selector := make(map[string]interface{})

	//校验列车是否存在
	selector[model.FieldID] = req.TrainID
	err := operator.Database.ListEntityByFilter(model.TableNameTrain, selector, model.OneQuery, &train)
	if err != nil {
		log.Error("列车查找失败. err:[%v]", err)
		return err
	}
	if train.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "列车")
	}

	selector = make(map[string]interface{})
	selector[model.FieldID] = req.ID
	//校验是否存在行驶计划
	err = operator.Database.ListEntityByFilter(model.TableNameTrainSchedule, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("行驶计划查找失败. err:[%v]", err)
		return err
	}

	if opt.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "同列车同时刻形式计划")
	}
	end_data := opt.EndDate
	CreateTime := opt.CreatedAt
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}
	opt.CreatedAt = CreateTime
	opt.TrainName = train.TrainType + train.Name

	//时间格式化
	opt.DepartureDate = utils.ParseTime(layout, req.DepartureDate)

	//咱不变动结束时间
	////暂时设置为默认空值
	opt.EndDate = end_data

	//新增行驶计划
	err = operator.Database.SaveEntity(model.TableNameTrainSchedule, &opt)
	if err != nil {
		log.Error("创建行驶计划. err:[%v]", err)
		return err
	}

	return nil
}

func (operator *ResourceOperator) QueryTrainScheduleList(req *apimodel.TrainScheduleRequest) (*apimodel.TrainSchedulePageResponse, error) {
	var resp apimodel.TrainSchedulePageResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	//id查
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	if req.TrainID > 0 {
		selector[model.FieldTrainID] = req.TrainID
	}
	if req.DepartureDate != "" {
		selector[model.FieldDepartureTime] = utils.ParseTime(layout, req.DepartureDate)
	}

	var count int64
	var schedules []model.TrainSchedule
	err := operator.Database.CountEntityByFilter(model.TableNameTrainSchedule, selector, model.OneQuery, &count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
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
		//车号模糊查询
		if req.TrainName != "" {
			var keyword []model.Keyword
			keyword = append(keyword, model.Keyword{Field: model.FieldTrainName, Value: req.TrainName, Type: 0})
			subquery := &model.SubQuery{
				Keywords: keyword,
			}
			queryParams.SubQueries = append(queryParams.SubQueries, subquery)
		}
		err = operator.Database.PreloadEntityByFilter(model.TableNameTrainSchedule, selector, queryParams, &schedules, []string{"Stops", "Seats"})
		if err != nil {
			log.Error("行驶计划数据查询失败. err:[%v]", err)
			return nil, err
		}
	}
	resp.Load(count, schedules)
	return &resp, nil
}

func (operator *ResourceOperator) DeleteTrainSchedule(req *apimodel.TrainScheduleRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	// 开启事务
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("DeleteTrainSchedule TransactionBegin Error.err[%v]", err)
		return err
	}
	defer func() {
		_ = tx.TransactionRollback()
	}()
	selector[model.FieldID] = req.ID
	err = tx.Database.DeleteEntityByFilter(model.TableNameTrainSchedule, selector, queryParams, &model.TrainSchedule{})
	if err != nil {
		log.Error("行驶计划删除失败. err:[%v]", err)
		return err
	}

	//删除对应的停靠信息，座位信息
	selector = make(map[string]interface{})
	selector[model.FieldScheduleID] = req.ID
	err = tx.Database.DeleteEntityByFilter(model.TableNameTrainStop, selector, queryParams, &model.TrainStop{})
	if err != nil {
		log.Error("停靠计划删除失败. err:[%v]", err)
		return err
	}
	err = tx.Database.DeleteEntityByFilter(model.TableNameTrainSeat, selector, queryParams, &model.TrainSeat{})
	if err != nil {
		log.Error("座位信息删除失败. err:[%v]", err)
		return err
	}
	err = tx.TransactionCommit()
	if err != nil {
		log.Error("TransactionCommit Error.err[%v]", err)
		return err
	}
	return nil
}

// 停靠信息
func (operator *ResourceOperator) CreateTrainStopInfo(req *apimodel.TrainStopInfoRequest) error {
	var schedule model.TrainSchedule
	selector := make(map[string]interface{})
	//校验行驶计划是否存在
	// 开启事务
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("CreateTrainStopInfo TransactionBegin Error.err[%v]", err)
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

func (operator *ResourceOperator) UpdateTrainStopInfo(req *apimodel.TrainStopInfoRequest) error {
	var schedule model.TrainSchedule
	selector := make(map[string]interface{})
	//校验行驶计划是否存在
	// 开启事务
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("CreateTrainStopInfo TransactionBegin Error.err[%v]", err)
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
		baseMode := model.Model{ID: v.ID}
		stopList = append(stopList, model.TrainStop{
			Model:         baseMode,
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
	err = tx.Database.SaveEntity(model.TableNameTrainStop, stopList)
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

func (operator *ResourceOperator) DeleteTrainStopInfo(req *apimodel.TrainStopInfoRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameTrainStop, selector, queryParams, &model.TrainStop{})
	if err != nil {
		log.Error("列车停靠数据删除失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) QueryTrainStopInfoList(req *apimodel.TrainStopInfoRequest) (*apimodel.TrainStopInfoPageResponse, error) {
	var resp apimodel.TrainStopInfoPageResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	//id查
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	//行驶计划查
	if req.ScheduleID > 0 {
		selector[model.FieldScheduleID] = req.ScheduleID
	}

	var count int64
	var stopList []model.TrainStop
	err := operator.Database.CountEntityByFilter(model.TableNameTrainStop, selector, model.OneQuery, &count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
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
		err = operator.Database.ListEntityByFilter(model.TableNameTrainStop, selector, queryParams, &stopList)
		if err != nil {
			log.Error("行驶计划数据查询失败. err:[%v]", err)
			return nil, err
		}
	}
	resp.Load(count, stopList)
	return &resp, nil
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

func (operator *ResourceOperator) DeleteTrainSeatInfo(req *apimodel.TrainSeatInfoRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameTrainSeat, selector, queryParams, &model.TrainSeat{})
	if err != nil {
		log.Error("列车座位数据删除失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) QueryTrainSeatInfoList(req *apimodel.TrainSeatInfoRequest) (*apimodel.TrainSeatInfoPageResponse, error) {
	var resp apimodel.TrainSeatInfoPageResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	//id查
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	//行驶计划查
	if req.ScheduleID > 0 {
		selector[model.FieldScheduleID] = req.ScheduleID
	}

	var count int64
	var seatList []model.TrainSeat
	err := operator.Database.CountEntityByFilter(model.TableNameTrainSeat, selector, model.OneQuery, &count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
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
		err = operator.Database.ListEntityByFilter(model.TableNameTrainSeat, selector, queryParams, &seatList)
		if err != nil {
			log.Error("列车座位数据查询失败. err:[%v]", err)
			return nil, err
		}
	}
	resp.Load(count, seatList)
	return &resp, nil
}

func (operator *ResourceOperator) UpdateTrainSeatInfo(req *apimodel.TrainSeatInfoRequest) error {
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
		log.Error("行驶计划查找失败. err:[%v]", err)
		return err
	}
	if schedule.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "行驶计划")
	}

	//遍历req中的stopList
	var seatList []model.TrainSeat
	for _, v := range req.SeatInfoList {
		var baseMode model.Model
		baseMode.ID = v.ID
		seatList = append(seatList, model.TrainSeat{
			Model:       baseMode,
			ScheduleID:  req.ScheduleID,
			SeatNums:    v.SeatNums,
			SeatNowNums: v.SeatNums,
			SeatType:    v.SeatType,
			Price:       v.Price,
		})
	}

	//新增停靠信息
	err = tx.Database.SaveEntity(model.TableNameTrainSeat, seatList)
	if err != nil {
		log.Error("更新座位信息. err:[%v]", err)
		return err
	}
	err = tx.TransactionCommit()
	if err != nil {
		log.Error("更新座位信息 TransactionCommit Error.err[%v]", err)
		return err
	}
	return nil
}
