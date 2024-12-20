package service

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid/v5"
	log "github.com/wonderivan/logger"
	"sync"
	"ticket-service/api/apimodel"
	"ticket-service/database/model"
	"ticket-service/httpserver/errcode"
	"ticket-service/pkg/utils"
	"ticket-service/pkg/utils/redis"
	"time"
)

func (operator *ResourceOperator) CreateUserOrder(req *apimodel.UserOrderRequest) (uuid.UUID, error) {
	var opt model.UserOrder
	var userDB model.User
	var mutex sync.Mutex
	tx, err := operator.TransactionBegin()
	if err != nil {
		log.Error("DeleteTrainSchedule TransactionBegin Error.err[%v]", err)
		return uuid.Nil, err
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
		return uuid.Nil, err
	}

	//校验行驶计划
	var schedule model.TrainSchedule
	selector = make(map[string]interface{})
	selector[model.FieldID] = req.ScheduleID
	err = tx.Database.PreloadEntityByFilter(model.TableNameTrainSchedule, selector, model.OneQuery, &schedule, []string{model.PreloadStops, model.PreloadSeats})
	if err != nil {
		log.Error("行驶计划数据查询失败. err:[%v]", err)
		return uuid.Nil, err
	}

	var startInfo, endInfo model.TrainStop
	//校验所属站点
	if schedule.Stops == nil || schedule.Seats == nil {
		log.Error("行驶计划关联座位或停靠信息数据查询失败. err:[%v]", err)
		return uuid.Nil, err
	}
	for _, v := range schedule.Stops {
		if v.StationID == req.StartStationID {
			opt.StartStationID = v.StationID
			startInfo = v
		} else if v.StationID == req.EndStationID {
			opt.EndStationID = v.StationID
			endInfo = v
		}
	}
	if userDB.ID <= 0 {
		return uuid.Nil, fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "用户")
	}
	if schedule.ID <= 0 {
		return uuid.Nil, fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "行驶计划")
	}
	var seatPrice float64
	for _, v := range schedule.Seats {
		if v.SeatType == req.SeatType {
			opt.SeatType = v.SeatType
			seatPrice = v.Price
		}
	}
	if opt.SeatType == "" {
		return uuid.Nil, fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "座位类型")
	}
	if opt.StartStationID <= 0 || opt.EndStationID <= 0 {
		return uuid.Nil, fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "起点或终点")
	}
	opt.UUID = uuid.Must(uuid.NewV4())
	opt.UserID = userDB.ID
	opt.UserPhone = userDB.Phone
	opt.ScheduleID = schedule.ID
	opt.Price = float64(endInfo.StopOrder-startInfo.StopOrder) * seatPrice
	opt.DepartureTime = startInfo.DepartureTime
	opt.ArrivalTime = endInfo.DepartureTime
	opt.StartStation = startInfo
	opt.EndStation = endInfo
	selector = make(map[string]interface{})
	selector[model.FieldScheduleID] = req.ScheduleID
	selector[model.FieldSeatType] = req.SeatType

	var seat_num_data model.TrainSeat

	//加锁
	mutex.Lock()
	err = tx.Database.ReduceEntityRowsByFilter(model.TableNameTrainSeat, selector, model.OneQuery, model.FieldSeatNowNums, "1")
	if err != nil {
		if err.Error() == "NoNums" {
			mutex.Unlock()
			return uuid.Nil, fmt.Errorf(errcode.ErrorMsgNoTickets, req.SeatType)
		}
		log.Error("座位数量-1失败. err:[%v]", err)
		mutex.Unlock()
		return uuid.Nil, err
	}
	err = tx.Database.ListEntityBySelectFilter(model.TableNameTrainSeat, selector, model.OneQuery, &seat_num_data, []string{model.FieldID, model.FieldSeatNums, model.FieldSeatNowNums})
	if err != nil {
		return uuid.Nil, fmt.Errorf(errcode.ErrorMsgListData)
	}
	mutex.Unlock()
	if seat_num_data.ID <= 0 {
		return uuid.Nil, fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "座位类型")
	}
	//座位编号用总票数-当前票数替代。
	opt.SeatNum = seat_num_data.SeatNums - seat_num_data.SeatNowNums
	opt.CreatedAt = time.Now()
	//新创建的订单写入redis中
	jsonData, err := json.Marshal(&opt)
	if err != nil {
		log.Error("新增订单 序列化失败 Error.err[%v]", err)
		return uuid.Nil, err
	}

	_, err = redis.RedisClient.HSet(redisKey, opt.UUID.String(), string(jsonData)).Result()
	if err != nil {
		log.Error("创建订单.写入redis失败 err:[%v]", err)
		return uuid.Nil, err
	}

	//设置过期时间
	_, err = redis.RedisClient.Expire(opt.UUID.String(), expireDuration).Result()
	if err != nil {
		log.Error("设置过期时间失败 err:[%v]", err)
		return uuid.Nil, err
	}

	//注册定时任务,写入redis中
	operator.TimerFreeOrder(time.Minute*3, opt)

	err = tx.TransactionCommit()
	if err != nil {
		log.Error("新增订单 TransactionCommit Error.err[%v]", err)
		return uuid.Nil, err
	}
	return opt.UUID, nil
}

func (operator *ResourceOperator) QueryUserOrderList(req *apimodel.UserOrderRequest) (*apimodel.UserOrderPageResponse, error) {
	//@TODO 查找订单：分四类 1待支付 2待出行 3已出行 4已取消 //对应需要参数tag:   WaitingPayList   WaitingDepartList   BeenDepartList   BeenCancelList
	//@TODO 按照order_id查或未输入tag: 返回默认list
	//@TODO 待支付订单 id、user_id、schedule_id三选一，若传入多个数据可能会重复。
	var resp apimodel.UserOrderPageResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	//订单id查
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	//user_id查
	if req.UserID > 0 {
		selector[model.FieldUserID] = req.UserID
	}
	//行驶计划查
	if req.ScheduleID > 0 {
		selector[model.FieldScheduleID] = req.ScheduleID
	}
	//创建分组查询条件： 待出行
	if req.Tag == "WaitingDepartList" {
		selector[model.FieldOrderIsPay] = true
		selector[model.FieldOrderIsDepart] = false
	}
	//创建分组查询条件： 已出行
	if req.Tag == "BeenDepartList" {
		selector[model.FieldOrderIsPay] = true
		selector[model.FieldOrderIsDepart] = true
	}
	//创建分组查询条件： 已取消
	if req.Tag == "BeenCancelList" {
		selector[model.FieldOrderIsCancel] = true
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
	if req.StartTime != "" && req.EndTime != "" {
		rangeQuery := &model.RangeQuery{
			Field: model.FieldCreatedTime,
			Start: req.StartTime,
			End:   req.EndTime,
		}
		queryParams.RangeQueries = append(queryParams.RangeQueries, rangeQuery)
	}

	if req.Tag == "WaitingPayList" {
		//待支付
		var orders []model.UserOrder
		hashData, err := redis.RedisClient.HGetAll(redisKey).Result()
		if err != nil {
			log.Error("查询订单 待支付订单 Error.err[%v]", err)
			return nil, err
		}
		for _, v := range hashData {
			var order model.UserOrder
			_ = json.Unmarshal([]byte(v), &order)
			orders = append(orders, order)
		}
		//根据传入id、user_id、schedule_id筛选
		var targetOrder []model.UserOrder
		for _, v := range orders {
			if req.ID > 0 {
				if v.ID == req.ID {
					targetOrder = append(targetOrder, v)
				}
			}
			if req.UserID > 0 {
				if v.UserID == req.UserID {
					targetOrder = append(targetOrder, v)
				}
			}
			if req.ScheduleID > 0 {
				if v.ScheduleID == req.ScheduleID {
					targetOrder = append(targetOrder, v)
				}
			}
		}
		//若未传入任何参数
		if req.ID == 0 && req.UserID == 0 && req.ScheduleID == 0 {
			targetOrder = orders
		}

		//分页处理
		if req.PageSize != 0 && req.PageNo != 0 {
			skip, end, err := utils.GetPage(len(targetOrder), req.PageNo, req.PageSize)
			if err != nil {
				return nil, fmt.Errorf("获取页面失败err:[%s]", err)
			}
			resp.Load(int64(len(targetOrder)), targetOrder[skip:end], req.Tag)
			return &resp, nil
		} else {
			resp.Load(int64(len(targetOrder)), targetOrder, req.Tag)
			return &resp, nil
		}

	} else if req.Tag == "WaitingDepartList" || req.Tag == "BeenDepartList" || req.Tag == "BeenCancelList" {
		var count int64
		var userOrders []model.UserOrder
		err := operator.Database.CountEntityByFilter(model.TableNameUserOrder, selector, queryParams, &count)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			err = operator.Database.ListEntityByFilter(model.TableNameUserOrder, selector, queryParams, &userOrders)
			if err != nil {
				log.Error("订单数据查询失败. err:[%v]", err)
				return nil, err
			}
		}
		resp.Load(count, userOrders, req.Tag)
		return &resp, nil
	} else if req.Tag == "" {
		//拼接所有订单。
		var orders []model.UserOrder
		hashData, err := redis.RedisClient.HGetAll(redisKey).Result()
		if err != nil {
			log.Error("查询订单 待支付订单 Error.err[%v]", err)
			return nil, err
		}
		for _, v := range hashData {
			var order model.UserOrder
			_ = json.Unmarshal([]byte(v), &order)
			orders = append(orders, order)
		}
		//根据传入id、user_id、schedule_id筛选
		var targetOrder []model.UserOrder
		for _, v := range orders {
			if req.ID > 0 {
				if v.ID == req.ID {
					targetOrder = append(targetOrder, v)
				}
			}
			if req.UserID > 0 {
				if v.UserID == req.UserID {
					targetOrder = append(targetOrder, v)
				}
			}
			if req.ScheduleID > 0 {
				if v.ScheduleID == req.ScheduleID {
					targetOrder = append(targetOrder, v)
				}
			}
		}
		//若未传入任何参数
		if req.ID == 0 && req.UserID == 0 && req.ScheduleID == 0 {
			targetOrder = orders
		}

		var count int64
		var userOrders []model.UserOrder
		err = operator.Database.CountEntityByFilter(model.TableNameUserOrder, selector, queryParams, &count)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			err = operator.Database.ListEntityByFilter(model.TableNameUserOrder, selector, queryParams, &userOrders)
			if err != nil {
				log.Error("订单数据查询失败. err:[%v]", err)
				return nil, err
			}
		}
		userOrders = append(userOrders, targetOrder...)

		//分页处理
		if req.PageSize != 0 && req.PageNo != 0 {
			skip, end, err := utils.GetPage(len(userOrders), req.PageNo, req.PageSize)
			if err != nil {
				return nil, fmt.Errorf("获取页面失败err:[%s]", err)
			}
			resp.Load(int64(len(userOrders)), userOrders[skip:end], req.Tag)
			return &resp, nil
		} else {
			resp.Load(int64(len(userOrders)), userOrders, req.Tag)
			return &resp, nil
		}

	} else {
		log.Error("订单 输入tag无效:[%#v]", req.Tag)
		return nil, fmt.Errorf("订单数据 tag输入无效:%s", req.Tag)
	}
}

func (operator *ResourceOperator) CancelUserOrder(req *apimodel.UserOrderRequest) error {
	//去redis中拿取订单，修改订单状态持久化到mysql中
	var orderDB model.UserOrder
	exists, err := redis.RedisClient.HExists(redisKey, req.UUID.String()).Result()
	if err != nil {
		log.Error("定时任务.读取redis待支付订单失败 err:[%v]", err)
		return err
	}
	selector := make(map[string]interface{})
	selector[model.FieldUUID] = req.UUID
	err = operator.Database.ListEntityByFilter(model.TableNameUserOrder, selector, model.OneQuery, &orderDB)
	if err != nil {
		log.Error("定时任务.读取redis待支付订单失败 err:[%v]", err)
		return err
	}

	if orderDB.ID > 0 {
		//已支付订单取消
		tx, err := operator.TransactionBegin()
		if err != nil {
			log.Error("DeleteTrainSchedule TransactionBegin Error.err[%v]", err)
			return err
		}
		defer func() {
			_ = tx.TransactionRollback()
		}()
		//针对已支付订单先做校验
		selector = make(map[string]interface{})
		selector[model.FieldScheduleID] = orderDB.ScheduleID
		selector[model.FieldSeatType] = orderDB.SeatType
		err = tx.Database.AddEntityRowsByFilter(model.TableNameTrainSeat, selector, model.OneQuery, "seat_now_nums", "1")
		if err != nil {
			log.Error("定时任务.恢复库存失败 err:[%v]", err)
			return err
		}
		//已取消
		orderDB.IsCancel = true
		err = tx.Database.SaveEntity(model.TableNameUserOrder, &orderDB)
		if err != nil {
			log.Error("取消订单失败 err:[%v]", err)
			return err
		}

		if orderDB.IsPay == true {
			//@TODO 退款功能；待定
		}

		err = tx.TransactionCommit()
		if err != nil {
			log.Error("定时任务.恢复库存 Commit 失败 err:[%v]", err)
			return err
		}

		return nil

	} else if exists {
		var order model.UserOrder
		//待支付订单取消
		orderData, err := redis.RedisClient.HGet(redisKey, req.UUID.String()).Result()
		if err != nil {
			log.Error("获取redis待支付订单失败 err:[%v]", err)
			return err
		}
		err = json.Unmarshal([]byte(orderData), &order)
		if err != nil {
			log.Error("获取redis获取待支付订单失败 err:[%v]", err)
			return err
		}
		tx, err := operator.TransactionBegin()
		if err != nil {
			log.Error("DeleteTrainSchedule TransactionBegin Error.err[%v]", err)
			return err
		}
		defer func() {
			_ = tx.TransactionRollback()
		}()

		selector = make(map[string]interface{})
		selector[model.FieldScheduleID] = order.ScheduleID
		selector[model.FieldSeatType] = order.SeatType
		err = tx.Database.AddEntityRowsByFilter(model.TableNameTrainSeat, selector, model.OneQuery, model.FieldSeatNowNums, "1")
		if err != nil {
			log.Error("定时任务.恢复库存失败 err:[%v]", err)
			return err
		}
		//已取消
		order.IsCancel = true
		err = tx.Database.SaveEntity(model.TableNameUserOrder, &order)
		if err != nil {
			log.Error("取消订单失败 err:[%v]", err)
			return err
		}

		_, err = redis.RedisClient.HDel(redisKey, req.UUID.String()).Result()
		if err != nil {
			log.Error("删除redis待支付订单失败 err:[%v]", err)
			return err
		}
		err = tx.TransactionCommit()
		if err != nil {
			log.Error("订单取消 持久化已取消 Commit 失败 err:[%v]", err)
			return err
		}
		return nil
	} else {
		return fmt.Errorf("取消订单 输入无效:%s", req.UUID)
	}
}

func (operator *ResourceOperator) PayUserOrder(req *apimodel.UserOrderRequest) error {
	//@TODO 假支付。
	//去redis中拿取订单，修改订单状态持久化到mysql中
	exists, err := redis.RedisClient.HExists(redisKey, req.UUID.String()).Result()
	if err != nil {
		log.Error("读取redis待支付订单失败 err:[%v]", err)
		return err
	}
	if exists {
		//订单过期后，释放库存，删除redis中键值。
		var order model.UserOrder
		orderData, err := redis.RedisClient.HGet(redisKey, req.UUID.String()).Result()
		if err != nil {
			log.Error("获取redis待支付订单失败 err:[%v]", err)
			return err
		}
		err = json.Unmarshal([]byte(orderData), &order)
		if err != nil {
			log.Error("获取redis获取待支付订单失败 err:[%v]", err)
			return err
		}

		// 支付完成
		order.IsPay = true
		err = operator.Database.SaveEntity(model.TableNameUserOrder, &order)
		if err != nil {
			log.Error("定时任务.恢复库存失败 err:[%v]", err)
			return err
		}
		_, err = redis.RedisClient.HDel(redisKey, req.UUID.String()).Result()
		if err != nil {
			log.Error("删除redis待支付订单失败 err:[%v]", err)
			return err
		}
		return nil
	} else {
		return fmt.Errorf("支付订单-uuid输入无效:%s", req.UUID)
	}
}

func (operator *ResourceOperator) DeleteUserOrder(req *apimodel.UserOrderRequest) error {
	var order model.UserOrder
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldUUID] = req.UUID
	err := operator.Database.DeleteEntityByFilter(model.TableNameUserOrder, selector, queryParams, &model.UserOrder{})
	if err != nil {
		log.Error("订单数据删除失败. err:[%v]", err)
		return err
	}
	//若删除的是待支付订单
	exists, err := redis.RedisClient.HExists(redisKey, req.UUID.String()).Result()
	if err != nil {
		log.Error("定时任务.读取redis待支付订单失败 err:[%v]", err)
		return err
	}
	if exists {
		_, err = redis.RedisClient.HDel(redisKey, req.UUID.String()).Result()
		if err != nil {
			log.Error("删除redis待支付订单失败 err:[%v]", err)
			return err
		}
		tx, _ := operator.TransactionBegin()
		defer func() {
			_ = tx.TransactionRollback()
		}()

		var mutex sync.Mutex
		mutex.Lock()
		selector = make(map[string]interface{})
		selector[model.FieldScheduleID] = order.ScheduleID
		selector[model.FieldSeatType] = order.SeatType
		err = tx.Database.AddEntityRowsByFilter(model.TableNameTrainSeat, selector, model.OneQuery, model.FieldSeatNowNums, "1")
		if err != nil {
			log.Error("定时任务.恢复库存失败 err:[%v]", err)
			return err
		}
		mutex.Unlock()
		err = tx.TransactionCommit()
		if err != nil {
			log.Error("定时任务.恢复库存 Commit 失败 err:[%v]", err)
			return err
		}
		return nil
	} else {
		return nil
	}
}
