package apimodel

import (
	"fmt"
	"ticket-service/database/model"
	"ticket-service/global"
	"ticket-service/httpserver/errcode"
)

// baseMode

// 车辆基本信息
type TrainInfo struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_time,omitempty"`
	UpdatedAt string `json:"updated_time,omitempty"`
	Name      string `json:"name"`        // 车次编号
	TrainType string `json:"train_type" ` // 型号 G、T、Z、K、D
	//Schedules    []TrainScheduleInfo `json:"schedules"`     //行驶计划
}

// 车站信息
type StationInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`     //车站名称
	Code      string `json:"code"`     //车站编码
	Province  string `json:"province"` //所属省份
	City      string `json:"city"`     //所属城市
	CreatedAt string `json:"created_time,omitempty"`
	UpdatedAt string `json:"updated_time,omitempty"`
}

// 行驶信息
type TrainScheduleInfo struct {
	ID            int             `json:"id"`
	TrainID       int             `json:"train_id"`       //列车id
	TrainName     string          `json:"train_name"`     //列车车号
	DepartureDate string          `json:"departure_date"` //出发日期
	Stops         []TrainStopInfo `json:"stops"`          //停靠信息
	Seats         []TrainSeatInfo `json:"seats"`          //座位信息
	CreatedAt     string          `json:"created_time,omitempty"`
	UpdatedAt     string          `json:"updated_time,omitempty"`
}

// 停靠信息
type TrainStopInfo struct {
	ID            int    `json:"id"`
	ScheduleID    int    `json:"schedule_id"` //运行计划id
	StationID     int    `json:"station_id"`  //车站id
	StationName   string `json:"station_name"`
	StopOrder     int    `json:"stop_order"`     //停靠顺序
	DepartureTime string `json:"departure_time"` //发车时间
	CreatedAt     string `json:"created_time,omitempty"`
	UpdatedAt     string `json:"updated_time,omitempty"`
}

// 座位信息
type TrainSeatInfo struct {
	ID          int     `json:"id"`
	ScheduleID  int     `json:"schedule_id"`   //运行计划id
	SeatNums    int     `json:"seat_nums"`     //座位数量
	SeatNowNums int     `json:"seat_now_nums"` //库存数量
	SeatType    string  `json:"seat_type"`     //座位类别
	Price       float64 `json:"price"`         //价格
	CreatedAt   string  `json:"created_time,omitempty"`
	UpdatedAt   string  `json:"updated_time,omitempty"`
}

// Request

type TrainInfoRequest struct {
	ID        int    `json:"id" uri:"id" form:"id"`
	Name      string `json:"name" form:"name"`             // 车次编号
	TrainType string `json:"train_type" form:"train_type"` // 型号 G、T、Z、K、D
	CreatedAt string `json:"created_time,omitempty"`
	UpdatedAt string `json:"updated_time,omitempty"`
	PaginationRequest
}

type TrainStationRequest struct {
	ID        int    `json:"id" uri:"id" form:"id"`
	Name      string `json:"name" form:"name"`         //车站名称
	Code      string `json:"code" form:"code"`         //车站编码
	Province  string `json:"province" form:"province"` //所属省份
	City      string `json:"city" form:"city"`         //所属城市
	CreatedAt string `json:"created_time,omitempty"`
	UpdatedAt string `json:"updated_time,omitempty"`
	PaginationRequest
}

//创建列车行驶计划 => 创建行驶计划,选择列车 => 填写停靠信息 => 填写座位信息

type TrainScheduleRequest struct {
	ID            int    `json:"id" uri:"id" form:"id"`
	TrainID       int    `json:"train_id" uri:"train_id" form:"train_id"` //列车id
	TrainName     string `json:"train_name"`
	DepartureDate string `json:"departure_date"` //出发时间
	CreatedAt     string `json:"created_time,omitempty"`
	UpdatedAt     string `json:"updated_time,omitempty"`
	PaginationRequest
}

type TrainStopInfoRequest struct {
	ID            int             `json:"id" uri:"id" form:"id"`
	ScheduleID    int             `json:"schedule_id" uri:"schedule_id" form:"schedule_id"` //运行计划id
	TrainStopList []TrainStopInfo `json:"train_stop_list"`
	PaginationRequest
}

type TrainSeatInfoRequest struct {
	ID           int             `json:"id" uri:"id" form:"id"`
	ScheduleID   int             `json:"schedule_id" uri:"schedule_id" form:"schedule_id"` //运行计划id
	SeatInfoList []TrainSeatInfo `json:"train_seat_list"`
	PaginationRequest
}

// Response struct

type TrainInfoPageResponse struct {
	List []TrainInfo `json:"list"`
	PaginationResponse
}

type StationInfoPageResponse struct {
	List []StationInfo `json:"list"`
	PaginationResponse
}

type TrainSchedulePageResponse struct {
	List []TrainScheduleInfo `json:"list"`
	PaginationResponse
}

type TrainStopInfoPageResponse struct {
	List []TrainStopInfo `json:"list"`
	PaginationResponse
}

type TrainSeatInfoPageResponse struct {
	List []TrainSeatInfo `json:"list"`
	PaginationResponse
}

//DataUnmarshal

func (t *TrainInfo) Load(TrainInfoData model.Train) {
	t.ID = TrainInfoData.ID
	t.CreatedAt = TrainInfoData.CreatedAt.String()
	t.UpdatedAt = TrainInfoData.UpdatedAt.String()
	t.Name = TrainInfoData.Name
	t.TrainType = TrainInfoData.TrainType
}

func (t *StationInfo) Load(TrainInfoData model.Station) {
	t.ID = TrainInfoData.ID
	t.CreatedAt = TrainInfoData.CreatedAt.String()
	t.UpdatedAt = TrainInfoData.UpdatedAt.String()
	t.Name = TrainInfoData.Name
	t.Province = TrainInfoData.Province
	t.City = TrainInfoData.City
	t.Code = TrainInfoData.Code
}

func (t *TrainScheduleInfo) Load(TrainInfoData model.TrainSchedule) {
	t.ID = TrainInfoData.ID
	t.CreatedAt = TrainInfoData.CreatedAt.String()
	t.UpdatedAt = TrainInfoData.UpdatedAt.String()
	t.TrainName = TrainInfoData.TrainName
	t.DepartureDate = TrainInfoData.DepartureDate.String()
	for _, v := range TrainInfoData.Stops {
		t.Stops = append(t.Stops, TrainStopInfo{
			ID:            v.ID,
			ScheduleID:    v.ScheduleID,
			StationID:     v.StationID,
			StationName:   global.StationCodeMap[v.ID],
			StopOrder:     v.StopOrder,
			DepartureTime: v.DepartureTime.String(),
			CreatedAt:     v.CreatedAt.String(),
			UpdatedAt:     v.UpdatedAt.String(),
		})
	}
	for _, v := range TrainInfoData.Seats {
		t.Seats = append(t.Seats, TrainSeatInfo{
			ID:          v.ID,
			ScheduleID:  v.ScheduleID,
			SeatNums:    v.SeatNums,
			SeatNowNums: v.SeatNowNums,
			SeatType:    v.SeatType,
			Price:       v.Price,
			CreatedAt:   v.CreatedAt.String(),
			UpdatedAt:   v.UpdatedAt.String(),
		})
	}
}

func (t *TrainStopInfo) Load(TrainStopData model.TrainStop) {
	t.ID = TrainStopData.ID
	t.CreatedAt = TrainStopData.CreatedAt.String()
	t.UpdatedAt = TrainStopData.UpdatedAt.String()
	t.ScheduleID = TrainStopData.ScheduleID
	t.StationName = global.StationCodeMap[TrainStopData.StationID]
	t.StationID = TrainStopData.StationID
	t.StopOrder = TrainStopData.StopOrder
	t.DepartureTime = TrainStopData.DepartureTime.String()
}

//ScheduleID  int     `json:"schedule_id"`
//SeatNums    int     `json:"seat_nums"`
//SeatNowNums int     `json:"seat_now_nums"`
//SeatType    string  `json:"seat_type"`
//Price       float64 `json:"price"`

func (t *TrainSeatInfo) Load(TrainSeatData model.TrainSeat) {
	t.ID = TrainSeatData.ID
	t.CreatedAt = TrainSeatData.CreatedAt.String()
	t.UpdatedAt = TrainSeatData.UpdatedAt.String()
	t.ScheduleID = TrainSeatData.ScheduleID
	t.SeatNums = TrainSeatData.SeatNums
	t.SeatNowNums = TrainSeatData.SeatNowNums
	t.SeatType = TrainSeatData.SeatType
	t.Price = TrainSeatData.Price
}

//DataLoading

func (resp *TrainInfoPageResponse) Load(total int64, list []model.Train) {
	resp.List = make([]TrainInfo, 0, len(list))
	for _, v := range list {
		info := TrainInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

func (resp *StationInfoPageResponse) Load(total int64, list []model.Station) {
	resp.List = make([]StationInfo, 0, len(list))
	for _, v := range list {
		info := StationInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

func (resp *TrainSchedulePageResponse) Load(total int64, list []model.TrainSchedule) {
	resp.List = make([]TrainScheduleInfo, 0, len(list))
	for _, v := range list {
		info := TrainScheduleInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

func (resp *TrainStopInfoPageResponse) Load(total int64, list []model.TrainStop) {
	resp.List = make([]TrainStopInfo, 0, len(list))
	for _, v := range list {
		info := TrainStopInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

func (resp *TrainSeatInfoPageResponse) Load(total int64, list []model.TrainSeat) {
	resp.List = make([]TrainSeatInfo, 0, len(list))
	for _, v := range list {
		info := TrainSeatInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

// Handler valid

func (req TrainInfoRequest) Valid(opt string) error {
	if opt == ValidOptCreate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Name == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
		}
		if req.TrainType == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "type")
		}
	} else if opt == ValidOptUpdate {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Name == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
		}
		if req.TrainType == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "type")
		}
	} else if opt == ValidOptDel {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}

func (req TrainStationRequest) Valid(opt string) error {
	if opt == ValidOptCreate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Name == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
		}
		if req.Code == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "code")
		}
		if req.Province == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "province")
		}
		if req.City == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "city")
		}
	} else if opt == ValidOptUpdate {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Name == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
		}
		if req.Code == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "code")
		}
		if req.Province == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "province")
		}
		if req.City == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "city")
		}
	} else if opt == ValidOptDel {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}

func (req TrainScheduleRequest) Valid(opt string) error {
	if opt == ValidOptCreate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.TrainID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "train_id")
		}
		if req.DepartureDate == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "departue_date")
		}
	} else if opt == ValidOptUpdate {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.TrainID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "train_id")
		}
		if req.DepartureDate == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "departue_date")
		}
	} else if opt == ValidOptDel {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}

func (req TrainStopInfoRequest) Valid(opt string) error {
	if opt == ValidOptCreate {
		if req.TrainStopList == nil {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "stopList")
		}
		for _, v := range req.TrainStopList {
			if v.ID < 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
			}
			if v.ScheduleID < 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "schedule_id")
			}
			if v.StationID <= 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "station_id")
			}
			if v.StopOrder < 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "stop_order")
			}
			if v.DepartureTime == "" {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "departure_time") //出发时间
			}
		}
		if req.ScheduleID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "schedule_id")
		}
	} else if opt == ValidOptUpdate {
		if req.TrainStopList == nil {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "stopList")
		}
		for _, v := range req.TrainStopList {
			if v.ID <= 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
			}
			if v.ScheduleID < 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "schedule_id")
			}
			if v.StationID <= 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "station_id")
			}
			if v.StopOrder < 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "stop_order")
			}
			if v.DepartureTime == "" {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "departure_time") //出发时间
			}
		}
		if req.ScheduleID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "schedule_id")
		}
	} else if opt == ValidOptDel {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}

func (req TrainSeatInfoRequest) Valid(opt string) error {
	if opt == ValidOptCreate {
		if req.SeatInfoList == nil {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "seatList")
		}
		for _, v := range req.SeatInfoList {
			if v.ID < 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
			}
			if v.SeatNums < 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "seat_nums")
			}
			if v.SeatType == "" {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "seat_type")
			}
			if v.Price <= 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "price")
			}
		}
		if req.ScheduleID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.ScheduleID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "schedule_id")
		}
	} else if opt == ValidOptUpdate {
		if req.SeatInfoList == nil {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "seatList")
		}
		for _, v := range req.SeatInfoList {
			if v.ID <= 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
			}
			if v.SeatNums < 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "seat_nums")
			}
			if v.SeatType == "" {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "seat_type")
			}
			if v.Price <= 0 {
				return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "price")
			}
		}
		if req.ScheduleID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.ScheduleID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "schedule_id")
		}
	} else if opt == ValidOptDel {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}
