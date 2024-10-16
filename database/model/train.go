package model

import "time"

type Train struct {
	Model
	Name              string    `json:"name" gorm:"index;comment:车次名"`
	PassCity          string    `json:"pass_city" gorm:"index;comment:途径城市"`
	StartAt           time.Time `json:"start_time,omitempty" gorm:"column:start_at"`
	Start             string    `json:"start" gorm:"comment:始发站"`
	End               string    `json:"end" gorm:"comment:终点站"`
	SeaTingNums       int       `json:"sea_ting_nums"`
	SeaTingPrice      float64   `json:"sea_ting_price"`
	SleepingNums      int       `json:"sleeping_nums"`
	SleepingPrice     float64   `json:"sleeping_price"`
	HighSleepingNums  int       `json:"high_sleeping_nums"`
	HighSleepingPrice float64   `json:"high_sleeping_price"`
	BusinessNums      int       `json:"business_nums"`
	BusinessPrice     float64   `json:"business_price"`
}

func (m *Train) TableName() string {
	return TableNameTrain
}
