package apimodel

import (
	"fmt"
	"ticket-service/database/model"
	"ticket-service/global"
	"ticket-service/httpserver/errcode"
)

// base struct

type UserOrderInfo struct {
	ID               int           `json:"id"`
	UserID           int           `json:"user_id"`
	UserPhone        string        `json:"user_phone"`
	ScheduleID       int           `json:"schedule_id"`        //运行计划id
	StartStationID   int           `json:"start_station_id"`   //起点站id
	StartStationInfo TrainStopInfo `json:"start_station_info"` //起点站
	StartStationName string        `json:"start_station_name"`
	EndStationID     int           `json:"end_station_id"`   //终点站id
	EndStationInfo   TrainStopInfo `json:"end_station_info"` //终点站
	EndStationName   string        `json:"end_station_name"`
	SeatType         string        `json:"seat_type"`      //座位类型
	Price            float64       `json:"price"`          //订单金额
	IsDepart         bool          `json:"is_depart"`      //出行状态
	IsPay            bool          `json:"is_pay"`         //支付状态
	DepartureTime    string        `json:"departure_time"` //出发时间
	ArrivalTime      string        `json:"arrival_time"`   //到达时间
	CreatedAt        string        `json:"created_time"`
	UpdatedAt        string        `json:"updated_time"`
}

//Request struct

type UserOrderRequest struct {
	ID             int     `json:"id"`
	UserID         int     `json:"user_id"`
	UserPhone      string  `json:"user_phone"`
	ScheduleID     int     `json:"schedule_id"`      //运行计划id
	StartStationID int     `json:"start_station_id"` //起点站id
	EndStationID   int     `json:"end_station_id"`   //终点站id
	SeatType       string  `json:"seat_type"`        //座位类型
	Price          float64 `json:"price"`            //订单金额
	IsDepart       bool    `json:"is_depart"`        //出行状态
	IsPay          bool    `json:"is_pay"`           //支付状态
	DepartureTime  string  `json:"departure_time"`   //出发时间
	ArrivalTime    string  `json:"arrival_time"`     //到达时间
	CreatedAt      string  `json:"created_time"`
	UpdatedAt      string  `json:"updated_time"`
	PaginationRequest
}

// Response struct

type UserOrderPageResponse struct {
	List []UserOrderInfo `json:"list"`
	PaginationResponse
}

// DataUnmarshal

func (t *UserOrderInfo) Load(orderData model.UserOrder) {
	t.UserID = orderData.UserID
	t.UserPhone = orderData.UserPhone
	t.ID = orderData.ID
	t.ScheduleID = orderData.ScheduleID
	t.StartStationID = orderData.StartStationID
	//t.StartStationInfo = orderData.StartStationInfo
	t.StartStationName = global.StationCodeMap[orderData.StartStationID]
	t.EndStationID = orderData.EndStationID
	//t.EndStationInfo = orderData.EndStationInfo
	t.EndStationName = global.StationCodeMap[orderData.EndStationID]
	t.Price = orderData.Price
	t.IsDepart = orderData.IsDepart
	t.IsPay = orderData.IsPay
	t.DepartureTime = orderData.DepartureTime.String()
	t.ArrivalTime = orderData.ArrivalTime.String()
	t.CreatedAt = orderData.CreatedAt.String()
	t.UpdatedAt = orderData.UpdatedAt.String()
}

// DataLoading

func (resp *UserOrderPageResponse) Load(total int64, list []model.UserOrder) {
	resp.List = make([]UserOrderInfo, 0, len(list))
	for _, v := range list {
		info := UserOrderInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
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
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}
