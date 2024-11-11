package apimodel

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"strconv"
	"ticket-service/database/model"
	"ticket-service/global"
	"ticket-service/httpserver/errcode"
)

// base struct

type UserOrderInfo struct {
	UUID             uuid.UUID     `json:"uuid"`
	ID               int           `json:"id"`
	UserID           int           `json:"user_id"`
	UserPhone        string        `json:"user_phone"`
	ScheduleID       int           `json:"schedule_id"`      //运行计划id
	StartStationID   int           `json:"start_station_id"` //起点站id
	StartInfo        TrainStopInfo `json:"start_info"`       //起点站
	StartStationName string        `json:"start_station_name"`
	EndStationID     int           `json:"end_station_id"` //终点站id
	EndInfo          TrainStopInfo `json:"end_info"`       //终点站
	EndStationName   string        `json:"end_station_name"`
	SeatType         string        `json:"seat_type"`      //座位类型
	SeatNum          string        `json:"seat_num"`       //座位编号
	Price            float64       `json:"price"`          //订单金额
	IsDepart         bool          `json:"is_depart"`      //出行状态
	IsPay            bool          `json:"is_pay"`         //支付状态
	IsCancel         bool          `json:"is_cancel"`      //是否取消 0:未取消 1:已取消
	DepartureTime    string        `json:"departure_time"` //出发时间
	ArrivalTime      string        `json:"arrival_time"`   //到达时间
	CreatedAt        string        `json:"created_time"`
	UpdatedAt        string        `json:"updated_time"`
}

//Request struct

type UserOrderRequest struct {
	UUID           uuid.UUID `json:"uuid" form:"uuid" uri:"uuid"`
	ID             int       `json:"id" form:"id"`
	UserID         int       `json:"user_id" form:"user_id"`
	UserPhone      string    `json:"user_phone"`
	ScheduleID     int       `json:"schedule_id" form:"schedule_id"` //运行计划id
	StartStationID int       `json:"start_station_id"`               //起点站id
	EndStationID   int       `json:"end_station_id"`                 //终点站id
	SeatType       string    `json:"seat_type"`                      //座位类型
	Price          float64   `json:"price"`                          //订单金额
	IsDepart       bool      `json:"is_depart"`                      //出行状态
	IsPay          bool      `json:"is_pay"`                         //支付状态
	IsCancel       bool      `json:"is_cancel"`                      //是否取消 0:未取消 1:已取消
	DepartureTime  string    `json:"departure_time"`                 //出发时间
	ArrivalTime    string    `json:"arrival_time"`                   //到达时间
	Tag            string    `json:"tag" form:"tag"`                 //分组条件
	CreatedAt      string    `json:"created_time"`
	UpdatedAt      string    `json:"updated_time"`
	PaginationRequest
}

// Response struct

type UserOrderPageResponse struct {
	WaitingPayList    []UserOrderInfo `json:"Waiting_pay"`    //待支付
	WaitingDepartList []UserOrderInfo `json:"Waiting_depart"` //待出行
	BeenDepartList    []UserOrderInfo `json:"Been_depart"`    //已出行
	BeenCancelList    []UserOrderInfo `json:"Been_cancel"`    //已取消
	List              []UserOrderInfo `json:"list"`           //未进行分组--id查 1个
	PaginationResponse
}

// DataUnmarshal

func (t *UserOrderInfo) Load(orderData model.UserOrder) {
	t.UserID = orderData.UserID
	t.UUID = orderData.UUID
	t.UserPhone = orderData.UserPhone
	t.ID = orderData.ID
	t.ScheduleID = orderData.ScheduleID
	t.StartStationID = orderData.StartStationID
	t.StartStationName = global.StationCodeMap[orderData.StartStationID]
	t.EndStationID = orderData.EndStationID
	t.EndStationName = global.StationCodeMap[orderData.EndStationID]
	t.Price = orderData.Price
	t.IsDepart = orderData.IsDepart
	t.IsPay = orderData.IsPay
	t.IsCancel = orderData.IsCancel
	t.DepartureTime = orderData.DepartureTime.String()
	t.ArrivalTime = orderData.ArrivalTime.String()
	t.CreatedAt = orderData.CreatedAt.String()
	t.UpdatedAt = orderData.UpdatedAt.String()
	t.SeatNum = strconv.Itoa(orderData.SeatNum)
	t.SeatType = orderData.SeatType
	t.StartInfo.Load(orderData.StartStation)
	t.EndInfo.Load(orderData.EndStation)
}

// DataLoading

// WaitingPayList     //待支付
// WaitingDepartList  //待出行
// BeenDepartList     //已出行
// BeenCancelList     //已取消

func (resp *UserOrderPageResponse) Load(total int64, list []model.UserOrder, tag string) {
	switch tag {
	case "WaitingPayList":
		resp.WaitingPayList = make([]UserOrderInfo, 0, len(list))
		for _, v := range list {
			info := UserOrderInfo{}
			info.Load(v)
			resp.WaitingPayList = append(resp.WaitingPayList, info)
		}
	case "WaitingDepartList":
		resp.WaitingDepartList = make([]UserOrderInfo, 0, len(list))
		for _, v := range list {
			info := UserOrderInfo{}
			info.Load(v)
			resp.WaitingDepartList = append(resp.WaitingDepartList, info)
		}
	case "BeenDepartList":
		resp.BeenDepartList = make([]UserOrderInfo, 0, len(list))
		for _, v := range list {
			info := UserOrderInfo{}
			info.Load(v)
			resp.BeenDepartList = append(resp.BeenDepartList, info)
		}
	case "BeenCancelList":
		resp.BeenCancelList = make([]UserOrderInfo, 0, len(list))
		for _, v := range list {
			info := UserOrderInfo{}
			info.Load(v)
			resp.BeenCancelList = append(resp.BeenCancelList, info)
		}
	default:
		resp.List = make([]UserOrderInfo, 0, len(list))
		for _, v := range list {
			info := UserOrderInfo{}
			info.Load(v)
			resp.List = append(resp.List, info)
		}
	}
	resp.TotalSize = int(total)
}

// Handler valid

func (req UserOrderRequest) Valid(opt string) error {
	if opt == ValidOptCreate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.ScheduleID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "schedule_id")
		}
		if req.StartStationID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "start_station_id")
		}
		if req.EndStationID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "end_station_id")
		}
		if req.UserID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "user_id")
		}
	} else if opt == ValidOptUpdate {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.ScheduleID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "schedule_id")
		}
		if req.StartStationID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "start_station_id")
		}
		if req.EndStationID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "end_station_id")
		}
		if req.UserID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "user_id")
		}
	} else if opt == ValidOptDel {
		if req.UUID == uuid.Nil {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "uuid")
		}
	} else if opt == ValidOptCancel {
		if req.UUID == uuid.Nil {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "uuid")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}
