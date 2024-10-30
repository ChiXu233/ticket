package service

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid/v5"
	log "github.com/wonderivan/logger"
	"strconv"
	"sync"
	"ticket-service/api/apimodel"
	"ticket-service/database/model"
	"ticket-service/httpserver/errcode"
	"ticket-service/pkg/utils/redis"
	"time"
)

func (operator *ResourceOperator) CreateUserOrder(req *apimodel.UserOrderRequest) error {
	var opt model.UserOrder
	var userDB model.User
	var mutex sync.Mutex
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("DeleteTrainSchedule TransactionBegin Error.err[%v]", err)
		return err
	}
	defer func() {
		_ = tx.TransactionRollback()
	}()
	//校验用户ID 以及手机号

	selector := make(map[string]interface{})
	selector[model.FieldID] = req.UserID
	err = tx.Database.ListEntityByFilter(model.TableNameUser, selector, model.OneQuery, &userDB)
	if err != nil {
		log.Error("用户数据查询失败. err:[%v]", err)
		return err
	}

	//校验行驶计划
	var schedule model.TrainSchedule

	selector = make(map[string]interface{})
	selector[model.FieldID] = req.ScheduleID
	err = tx.Database.PreloadEntityByFilter(model.TableNameTrainSchedule, selector, model.OneQuery, &schedule, []string{"Stops", "Seats"})
	if err != nil {
		log.Error("行驶计划数据查询失败. err:[%v]", err)
		return err
	}

	var startInfo, endInfo model.TrainStop
	//校验所属站点
	for _, v := range schedule.Stops {
		if v.ID == req.StartStationID {
			opt.StartStationID = v.ID
			startInfo = v
		} else if v.ID == req.EndStationID {
			opt.EndStationID = v.ID
			endInfo = v
		}
	}
	if userDB.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "用户")
	}
	if schedule.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "行驶计划")
	}
	var seatPrice float64
	for _, v := range schedule.Seats {
		if v.SeatType == req.SeatType {
			opt.SeatType = v.SeatType
			seatPrice = v.Price
		}
	}
	if opt.SeatType == "" {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "座位类型")
	}
	if opt.StartStationID <= 0 || opt.EndStationID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "起点或终点")
	}
	opt.UUID = uuid.Must(uuid.NewV4())
	opt.UserID = userDB.ID
	opt.UserPhone = userDB.Phone
	opt.ScheduleID = schedule.ID
	opt.Price = float64(endInfo.StopOrder-startInfo.StopOrder) * seatPrice
	opt.DepartureTime = startInfo.DepartureTime
	opt.ArrivalTime = endInfo.DepartureTime
	selector = make(map[string]interface{})
	selector[model.FieldScheduleID] = req.ScheduleID
	selector[model.FieldSeatType] = req.SeatType

	//加锁
	mutex.Lock()
	err = tx.Database.ReduceEntityRowsByFilter(model.TableNameTrainSeat, selector, model.OneQuery, "seat_now_nums", "1")
	if err != nil {
		if err.Error() == "NoNums" {
			mutex.Unlock()
			return fmt.Errorf(errcode.ErrorMsgNoTickets, req.SeatType)
		}
		log.Error("座位数量-1失败. err:[%v]", err)
		mutex.Unlock()
		return err
	}
	mutex.Unlock()
	opt.CreatedAt = time.Now()
	//新创建的订单写入redis中
	jsonData, err := json.Marshal(&opt)
	if err != nil {
		log.Error("新增订单 序列化失败 Error.err[%v]", err)
		return err
	}
	_, err = redis.RedisClient.HSet("order_"+strconv.Itoa(opt.UserID), opt.UUID.String(), string(jsonData)).Result()
	if err != nil {
		log.Error("创建订单.写入redis失败 err:[%v]", err)
		return err
	}

	//注册定时任务,写入redis中
	operator.TimerFreeOrder(time.Minute*10, opt)

	err = tx.TransactionCommit()
	if err != nil {
		log.Error("新增订单 TransactionCommit Error.err[%v]", err)
		return err
	}
	return nil
}

//func (operator *ResourceOperator) QueryStationList(req *apimodel.TrainStationRequest) (*apimodel.StationInfoPageResponse, error) {
//	var resp apimodel.StationInfoPageResponse
//	selector := make(map[string]interface{})
//	queryParams := model.QueryParams{}
//	//id查
//	if req.ID > 0 {
//		selector[model.FieldID] = req.ID
//	}
//	if req.Code != "" {
//		selector[model.FieldStationCode] = req.Code
//	}
//	if req.City != "" {
//		selector[model.FieldStationCity] = req.City
//	}
//	if req.Province != "" {
//		selector[model.FieldStationProvince] = req.Province
//	}
//	var count int64
//	var stations []model.Station
//	err := operator.Database.CountEntityByFilter(model.TableNameStation, selector, model.OneQuery, &count)
//	if err != nil {
//		return nil, err
//	}
//	if count > 0 {
//		order := model.Order{
//			Field:     model.FieldID,
//			Direction: apimodel.OrderAsc,
//		}
//		queryParams.Orders = append(queryParams.Orders, order)
//		if req.PageSize > 0 {
//			queryParams.Limit = &req.PageSize
//			offset := (req.PageNo - 1) * req.PageSize
//			queryParams.Offset = &offset
//		}
//		//车站名模糊查询
//		if req.Name != "" {
//			var keyword []model.Keyword
//			keyword = append(keyword, model.Keyword{Field: model.FieldName, Value: req.Name, Type: 0})
//			subquery := &model.SubQuery{
//				Keywords: keyword,
//			}
//			queryParams.SubQueries = append(queryParams.SubQueries, subquery)
//		}
//		err = operator.Database.ListEntityByFilter(model.TableNameStation, selector, queryParams, &stations)
//		if err != nil {
//			log.Error("车站数据查询失败. err:[%v]", err)
//			return nil, err
//		}
//	}
//	resp.Load(count, stations)
//	return &resp, nil
//}
//
//func (operator *ResourceOperator) DeleteStation(req *apimodel.TrainStationRequest) error {
//	selector := make(map[string]interface{})
//	queryParams := model.QueryParams{}
//	selector[model.FieldID] = req.ID
//	err := operator.Database.DeleteEntityByFilter(model.TableNameStation, selector, queryParams, &model.Station{})
//	if err != nil {
//		log.Error("车站数据删除失败. err:[%v]", err)
//		return err
//	}
//	return nil
//}
//
//func (operator *ResourceOperator) UpdateStation(req *apimodel.TrainStationRequest) error {
//	var opt model.Station
//	//修改车站信息
//	selector := make(map[string]interface{})
//	//校验车站名
//	selector[model.FieldName] = req.Name
//	err := operator.Database.ListEntityByFilter(model.TableNameStation, selector, model.OneQuery, &opt)
//	if err != nil {
//		log.Error("车站名查找失败. err:[%v]", err)
//		return err
//	}
//	if opt.ID > 0 && opt.ID != req.ID {
//		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "车站名")
//	}
//
//	selector = make(map[string]interface{})
//	selector[model.FieldStationCode] = req.Code
//	err = operator.Database.ListEntityByFilter(model.TableNameStation, selector, model.OneQuery, &opt)
//	if err != nil {
//		log.Error("车站编码查找失败. err:[%v]", err)
//		return err
//	}
//	if opt.ID > 0 && opt.ID != req.ID {
//		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "车站编码")
//	}
//
//	selector = make(map[string]interface{})
//	selector[model.FieldID] = req.ID
//	err = operator.Database.ListEntityByFilter(model.TableNameStation, selector, model.OneQuery, &opt)
//	if err != nil {
//		log.Error("查找车站失败. err:[%v]", err)
//		return err
//	}
//	if opt.ID <= 0 {
//		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "id")
//	}
//
//	//保持ID不变,暂存createTime（save方法全字段更新）
//	req.ID = opt.ID
//	CreateTime := opt.CreatedAt
//	err = copier.Copy(&opt, req)
//	if err != nil {
//		return err
//	}
//	opt.CreatedAt = CreateTime
//
//	err = operator.Database.SaveEntity(model.TableNameStation, &opt)
//	if err != nil {
//		log.Error("列车数据更新失败. err:[%v]", err)
//		return err
//	}
//	return nil
//}
