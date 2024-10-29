package model

import "time"

type UserOrder struct {
	Model
	UserID           int       `json:"user_id" gorm:"index;not null;comment:用户ID"`
	UserPhone        string    `json:"user_phone" gorm:"index;not null;comment:联系方式"`
	ScheduleID       int       `json:"schedule_id" gorm:"index;not null;comment:运行计划id"`
	StartStationID   int       `json:"start_station_id" gorm:"index;not null;comment:起点站id"`
	StartStationInfo TrainStop `json:"start_station_info"`
	EndStationID     int       `json:"end_station_id" gorm:"index;not null;comment:终点站id"`
	EndStationInfo   TrainStop `json:"end_station_info"`
	SeatType         string    `json:"seat_type" gorm:"not null;comment:座位类别"`
	Price            float64   `json:"price" gorm:"comment:订单金额"`
	IsDepart         bool      `json:"is_depart" gorm:"default:false"` //出行状态
	IsPay            bool      `json:"is_pay" gorm:"default:false"`    //支付状态
	DepartureTime    time.Time `json:"departure_time" gorm:"comment:出发时间"`
	ArrivalTime      time.Time `json:"arrival_time" gorm:"comment:到达时间"`
}

func (m *UserOrder) TableName() string {
	return TableNameUserOrder
}
