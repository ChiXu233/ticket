package handler

import "ticket-service/service"

type RestHandler struct {
	Operator service.Operator
}

func NewHandler() *RestHandler {
	return &RestHandler{
		Operator: service.GetOperator(),
	}
}
