package apimodel

import (
	"fmt"
	"strings"
	"ticket-service/database/model"
	"ticket-service/httpserver/errcode"
)

type SeatInfo struct {
	Nums  int     `json:"nums"`
	Price float64 `json:"price"`
}

type TrainInfo struct {
	ID           int      `json:"ID" uri:"id"`
	Name         string   `json:"name"`
	PassCity     []string `json:"passCity"`
	Start        string   `json:"start"`
	End          string   `json:"end"`
	StartAt      string   `json:"start_at"`
	SeaTing      SeatInfo `json:"seating"`
	Sleeping     SeatInfo `json:"sleeping"`
	HighSleeping SeatInfo `json:"high_sleeping"`
	Business     SeatInfo `json:"business"`
}

type TrainInfoRequest struct {
	ID           int      `json:"id" uri:"id" form:"id"`
	Name         string   `json:"name" form:"name"`
	PassCity     []string `json:"pass_city"`
	Start        string   `json:"start" form:"start"`
	End          string   `json:"end" form:"end"`
	SeaTing      SeatInfo `json:"seating"`
	Sleeping     SeatInfo `json:"sleeping"`
	HighSleeping SeatInfo `json:"high_sleeping"`
	Business     SeatInfo `json:"business"`
	StartAt      string   `json:"start_at"`
	CreatedAt    string   `json:"created_time"`
	UpdatedAt    string   `json:"updated_time"`
	PaginationRequest
}

type TrainInfoResponse struct {
	List []TrainInfo `json:"list"`
	PaginationResponse
}

func (t *TrainInfo) Load(trainData model.Train) {
	t.ID = trainData.ID
	t.Name = trainData.Name
	t.PassCity = strings.Split(trainData.PassCity, "-")
	t.Start = trainData.Start
	t.End = trainData.End
	t.StartAt = trainData.StartAt.Format("2006-01-02 15:04:05")
	t.SeaTing = SeatInfo{Nums: trainData.SeaTingNums, Price: trainData.SeaTingPrice}
	t.Sleeping = SeatInfo{Nums: trainData.SleepingNums, Price: trainData.SleepingPrice}
	t.HighSleeping = SeatInfo{Nums: trainData.HighSleepingNums, Price: trainData.SleepingPrice}
	t.Business = SeatInfo{Nums: trainData.BusinessNums, Price: trainData.BusinessPrice}
}

func (resp *TrainInfoResponse) Load(total int64, list []model.Train) {
	resp.List = make([]TrainInfo, 0, len(list))
	for _, v := range list {
		info := TrainInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

// Handler valid

func (req TrainInfoRequest) Valid(opt string) error {
	if opt == ValidOptCreateOrUpdate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Name == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "train_name")
		}
		if req.Start == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "train_start")
		}
		if req.End == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "train_end")
		}
		if req.PassCity == nil {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "pass_city")
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
