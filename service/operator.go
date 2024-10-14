package service

import (
	"github.com/gin-gonic/gin"
	"ticket-service/database"
)

var resourceOperator Operator

type ResourceOperator struct {
	database.Database
}

type Operator interface {
	Login(c *gin.Context) (map[string]interface{}, error)
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
