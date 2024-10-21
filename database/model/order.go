package model

type PayStatus int

const (
	Waiting PayStatus = 0 //未支付
	Paid    PayStatus = 1 //已支付
)

type UserOrder struct {
	Model
	TrainID int       `json:"train_id"`
	Status  PayStatus `json:"status" `
}
