package apimodel

import (
	"fmt"
	"ticket-service/database/model"
	"ticket-service/httpserver/errcode"
	"ticket-service/pkg/utils"
)

type RouterInfos struct {
	ID        int           `json:"id"`
	CreatedAt string        `json:"created_time,omitempty"`
	UpdatedAt string        `json:"updated_time,omitempty"`
	Name      string        `json:"name"`
	Uri       string        `json:"uri"`
	Method    string        `json:"method"`
	Roles     []PreloadRole `json:"roles"`
}

type PreloadRouter struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_time,omitempty"`
	UpdatedAt string `json:"updated_time,omitempty"`
	Name      string `json:"name"`
	Uri       string `json:"uri"`
	Method    string `json:"method"`
}

type RoutersInfoRequest struct {
	ID        int           `json:"id" uri:"id" form:"id"`
	Name      string        `json:"name" form:"name"`     //
	Uri       string        `json:"uri" form:"uri"`       // 路由
	Method    string        `json:"method" form:"method"` //请求方法
	Roles     []PreloadRole `json:"roles"`
	CreatedAt string        `json:"created_time,omitempty"`
	UpdatedAt string        `json:"updated_time,omitempty"`
	PaginationRequest
}

type RoutersInfoPageResponse struct {
	List []RouterInfos `json:"list"`
	PaginationResponse
}

func (t *RouterInfos) Load(routers model.Routers) {
	t.ID = routers.ID
	t.CreatedAt = utils.TimeFormat(routers.CreatedAt)
	t.UpdatedAt = utils.TimeFormat(routers.UpdatedAt)
	t.Name = routers.Name
	t.Uri = routers.Uri
	t.Method = routers.Method
	for _, v := range routers.Roles {
		t.Roles = append(t.Roles, PreloadRole{
			ID:        v.ID,
			CreatedAt: utils.TimeFormat(v.CreatedAt),
			UpdatedAt: utils.TimeFormat(v.UpdatedAt),
			Name:      v.Name,
			Pid:       v.Pid,
			Comment:   v.Comment,
		})
	}
}

func (resp *RoutersInfoPageResponse) Load(total int64, list []model.Routers) {
	resp.List = make([]RouterInfos, 0, len(list))
	for _, v := range list {
		info := RouterInfos{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

func (req RoutersInfoRequest) Valid(opt string) error {
	if opt == ValidOptCreate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Name == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
		}
		if req.Uri == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "uri")
		}
		if req.Method == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "method")
		}
	} else if opt == ValidOptUpdate {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Name == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
		}
		if req.Uri == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "uri")
		}
		if req.Method == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "method")
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
