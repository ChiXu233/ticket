package model

import "time"

// Train 代表列车模型。
type Train struct {
	Model
	Name      string `json:"name" gorm:"index;not null;comment:车次编号"`
	TrainType string `json:"train_type" gorm:"not null;comment:型号"` // G、T、Z、K、D
}

// Station 代表车站信息。
type Station struct {
	Model
	Name     string `json:"name" gorm:"index;not null;comment:车站名称"`
	Code     string `json:"code" gorm:"not null;comment:车站编码"`
	Province string `json:"province" gorm:"comment:所属省份"`
	City     string `json:"city" gorm:"comment:所属城市"`
}

// TrainSchedule 代表列车驾驶信息。
type TrainSchedule struct {
	Model
	TrainID       int         `json:"train_id" gorm:"index;not null;comment:列车id"`
	DepartureDate time.Time   `json:"departure_date" gorm:"index;not null;comment:出发日期"`
	EndDate       time.Time   `json:"end_date" gorm:"index";comment:"结束日期"`
	Stops         []TrainStop `json:"stops" gorm:"foreignKey:ScheduleID"`
	Seats         []TrainSeat `json:"seats" gorm:"foreignKey:ScheduleID"`
}

// TrainStop 代表停靠信息。
type TrainStop struct {
	Model
	ScheduleID    int       `json:"schedule_id" gorm:"index;not null;comment:运行计划id"`
	StationID     int       `json:"station_id" gorm:"index;not null;comment:车站id"`
	StopOrder     int       `json:"stop_order" gorm:"not null;comment:停靠顺序"`
	DepartureTime time.Time `json:"departure_time" gorm:"comment:发车时间"`
}

// TrainSeat 代表座位信息。
type TrainSeat struct {
	Model
	ScheduleID  int     `json:"schedule_id" gorm:"index;not null;comment:运行计划id"`
	SeatNums    int     `json:"seat_nums" gorm:"not null;comment:座位数量"`
	SeatNowNums int     `json:"seat_now_nums" gorm:"comment:库存数量"`
	SeatType    string  `json:"seat_type" gorm:"not null;comment:座位类别"`
	Price       float64 `json:"price" gorm:"not null;comment:价格"`
}

func (m *Train) TableName() string {
	return TableNameTrain
}
func (m *Station) TableName() string {
	return TableNameStation
}
func (m *TrainSchedule) TableName() string {
	return TableNameTrainSchedule
}
func (m *TrainStop) TableName() string {
	return TableNameTrainStop
}
func (m *TrainSeat) TableName() string {
	return TableNameTrainSeat
}
