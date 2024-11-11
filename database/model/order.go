package model

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type UserOrder struct {
	//@TODO 座位编号可用 总票数-当前票数 表示;需后端维护一个 map[int]string;
	Model
	UUID           uuid.UUID `json:"uuid" gorm:"comment:UUID"` // 唯一标识
	UserID         int       `json:"user_id" gorm:"index;not null;comment:用户ID"`
	UserPhone      string    `json:"user_phone" gorm:"index;not null;comment:联系方式"`
	ScheduleID     int       `json:"schedule_id" gorm:"index;not null;comment:运行计划id"`
	StartStationID int       `json:"start_station_id" gorm:"index;not null;comment:起点站id"`
	StartStation   TrainStop `json:"start_station"`
	EndStationID   int       `json:"end_station_id" gorm:"index;not null;comment:终点站id"`
	EndStation     TrainStop `json:"end_station"`
	SeatType       string    `json:"seat_type" gorm:"not null;comment:座位类别"`
	SeatNum        int       `json:"seat_num" gorm:"comment:座位编号"`
	SeatCarriage   string    `json:"seat_carriage" gorm:"所在车厢"`
	Price          float64   `json:"price" gorm:"comment:订单金额"`
	IsDepart       bool      `json:"is_depart" gorm:"default:false"` //出行状态 0:未出发 1:已出发
	IsPay          bool      `json:"is_pay" gorm:"default:false"`    //支付状态 0:未支付 1:已支付
	IsCancel       bool      `json:"is_cancel" gorm:"default:false"` //是否取消 0:未取消 1:已取消
	DepartureTime  time.Time `json:"departure_time" gorm:"comment:出发时间"`
	ArrivalTime    time.Time `json:"arrival_time" gorm:"comment:到达时间"`
}

func (m *UserOrder) TableName() string {
	return TableNameUserOrder
}
