package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"ticket-service/api/apimodel"
	"ticket-service/database"
)

var resourceOperator Operator

type ResourceOperator struct {
	database.Database
}

type Operator interface {
	//用户
	Login(c *gin.Context, req apimodel.UserInfoRequest) (*apimodel.LoginResponse, error)
	Register(req *apimodel.UserInfoRequest) error
	UpdateUserInfo(req *apimodel.UserInfoRequest) error
	DeleteUser(req *apimodel.UserInfoRequest) error
	QueryUserList(req *apimodel.UserInfoRequest) (*apimodel.UserPageResponse, error)
	ChangePassword(req *apimodel.UserChangePWRequest) error
	QueryUserByUUID(uuid uuid.UUID) error

	//车辆
	CreateTrain(req *apimodel.TrainInfoRequest) error
	QueryTrainList(req *apimodel.TrainInfoRequest) (*apimodel.TrainInfoPageResponse, error)
	DeleteTrain(req *apimodel.TrainInfoRequest) error
	UpdateTrain(req *apimodel.TrainInfoRequest) error

	//车站
	CreateStation(req *apimodel.TrainStationRequest) error
	QueryStationList(req *apimodel.TrainStationRequest) (*apimodel.StationInfoPageResponse, error)
	DeleteStation(req *apimodel.TrainStationRequest) error
	UpdateStation(req *apimodel.TrainStationRequest) error

	//行驶计划
	CreateTrainSchedule(req *apimodel.TrainScheduleRequest) (int, error)
	DeleteTrainSchedule(req *apimodel.TrainScheduleRequest) error
	UpdateTrainSchedule(req *apimodel.TrainScheduleRequest) error
	QueryTrainScheduleList(req *apimodel.TrainScheduleRequest) (*apimodel.TrainSchedulePageResponse, error)

	//停靠信息
	CreateTrainStopInfo(req *apimodel.TrainStopInfoRequest) error
	QueryTrainStopInfoList(req *apimodel.TrainStopInfoRequest) (*apimodel.TrainStopInfoPageResponse, error)
	DeleteTrainStopInfo(req *apimodel.TrainStopInfoRequest) error
	UpdateTrainStopInfo(req *apimodel.TrainStopInfoRequest) error

	//座位
	CreateTrainSeatInfo(req *apimodel.TrainSeatInfoRequest) error
	DeleteTrainSeatInfo(req *apimodel.TrainSeatInfoRequest) error
	QueryTrainSeatInfoList(req *apimodel.TrainSeatInfoRequest) (*apimodel.TrainSeatInfoPageResponse, error)
	UpdateTrainSeatInfo(req *apimodel.TrainSeatInfoRequest) error

	//订单
	CreateUserOrder(req *apimodel.UserOrderRequest) (uuid.UUID, error)
	QueryUserOrderList(req *apimodel.UserOrderRequest) (*apimodel.UserOrderPageResponse, error)
	CancelUserOrder(req *apimodel.UserOrderRequest) error
	DeleteUserOrder(req *apimodel.UserOrderRequest) error
	PayUserOrder(req *apimodel.UserOrderRequest) error

	//StationMap
	LoadStation_CodeMap() error
}

func GetOperator() Operator {
	if resourceOperator == nil {
		resourceOperator = &ResourceOperator{
			Database: database.GetDatabase(),
		}
	}
	return resourceOperator
}

func NewMockOperator() ResourceOperator {
	return ResourceOperator{
		Database: database.GetDatabase(),
	}
}

func (operator *ResourceOperator) TransactionBegin() (*ResourceOperator, error) {
	db, err := database.GetDatabase().Begin()
	if err != nil {
		return nil, err
	}
	return &ResourceOperator{
		Database: db,
	}, nil
}

func (operator *ResourceOperator) TransactionCommit() error {
	return operator.Database.Commit()
}

func (operator *ResourceOperator) TransactionRollback() error {
	return operator.Database.Rollback()
}
