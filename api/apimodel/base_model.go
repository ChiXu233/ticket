package apimodel

import (
	"fmt"
	"gorm.io/gorm/utils"
	"ticket-service/httpserver/errcode"
	util "ticket-service/pkg/utils"
)

const (
	DefaultPageNo   = 1
	DefaultPageSize = 0
	DefaultOrderBy  = "created_at"
	OrderDesc       = "desc"
	OrderAsc        = "asc"

	ValidOptCreate   = "save"
	ValidOptUpdate   = "update"
	ValidOptList     = "query"
	ValidOptDel      = "del"
	ValidOptLogin    = "login"
	ValidOptCancel   = "cancel"
	ValidOptResetPwd = "reset"
)

var (
	DefaultPaginationRequest = PaginationRequest{
		PageNo:   DefaultPageNo,
		PageSize: DefaultPageSize,
		OrderBy:  DefaultOrderBy,
		Order:    OrderDesc,
	}
)

type PaginationRequest struct {
	PageNo    int    `json:"page_no" form:"page_no"`
	PageSize  int    `json:"page_size" form:"page_size"`
	OrderBy   string `json:"order_by" form:"order_by"`
	Order     string `json:"order" form:"order"`
	StartTime string `json:"start_time" form:"start_time"`
	EndTime   string `json:"end_time" form:"end_time"`
}

func (req PaginationRequest) Valid(orderByList []string) error {
	// pageSize为0代表不分页
	if req.PageSize < 0 {
		return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "[page_size]")
	}
	if req.PageNo <= 0 {
		return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "[page_no]")
	}
	if !utils.Contains(orderByList, req.OrderBy) {
		return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "[order_by]")
	}
	if req.Order != OrderDesc && req.Order != OrderAsc {
		return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "[order]")
	}
	if req.StartTime != "" && util.TimeFormat(util.ParseTime("2006-01-02 15:04:05", req.StartTime)) != req.StartTime {
		return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "[start_time]")
	}
	if req.EndTime != "" && util.TimeFormat(util.ParseTime("2006-01-02 15:04:05", req.EndTime)) != req.EndTime {
		return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "[end_time]")
	}

	return nil
}

type PaginationResponse struct {
	TotalSize int `json:"total_size"`
}
