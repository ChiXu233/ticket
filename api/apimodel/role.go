package apimodel

import (
	"fmt"
	"ticket-service/database/model"
	"ticket-service/httpserver/errcode"
	"ticket-service/pkg/utils"
)

type PreloadRole struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_time,omitempty"`
	UpdatedAt string `json:"updated_time,omitempty"`
	Name      string `json:"name"`
	Pid       int    `json:"pid"` //权限组 暂时没必要用
	Comment   string `json:"comment"`
}

type RoleInfo struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_time,omitempty"`
	UpdatedAt string `json:"updated_time,omitempty"`
	Name      string `json:"name"`
	Pid       int    `json:"pid"` //权限组 暂时没必要用
	Routers   []PreloadRouter
	Users     []PreloadUser
	Comment   string `json:"comment"`
}

type RoleInfoRequest struct {
	ID        int             `json:"id" form:"id" uri:"id"`
	CreatedAt string          `json:"created_time,omitempty"`
	UpdatedAt string          `json:"updated_time,omitempty"`
	Name      string          `json:"name" form:"name"`
	Pid       int             `json:"pid" form:"pid"`
	Routers   []PreloadRouter `json:"routers" form:"routers"`
	Users     []PreloadUser   `json:"users" form:"users"`
	Comment   string          `json:"comment" form:"comment"`
	PaginationRequest
}

type RoleInfoPageResponse struct {
	List []RoleInfo `json:"list"`
	PaginationResponse
}

func (t *RoleInfo) Load(roles model.Role) {
	t.ID = roles.ID
	t.CreatedAt = utils.TimeFormat(roles.CreatedAt)
	t.UpdatedAt = utils.TimeFormat(roles.UpdatedAt)
	t.Name = roles.Name
	t.Comment = roles.Comment
	for _, v := range roles.Routers {
		t.Routers = append(t.Routers, PreloadRouter{
			ID:        v.ID,
			CreatedAt: utils.TimeFormat(v.CreatedAt),
			UpdatedAt: utils.TimeFormat(v.UpdatedAt),
			Name:      v.Name,
			Method:    v.Method,
			Uri:       v.Uri,
		})
	}
	for _, v := range roles.Users {
		t.Users = append(t.Users, PreloadUser{
			ID:        v.ID,
			UUID:      v.UUID,
			Username:  v.Username,
			Password:  v.Password,
			NickName:  v.NickName,
			Phone:     v.Phone,
			Email:     v.Email,
			CreatedAt: utils.TimeFormat(v.CreatedAt),
			UpdatedAt: utils.TimeFormat(v.UpdatedAt),
		})
	}
}

func (resp *RoleInfoPageResponse) Load(total int64, list []model.Role) {
	resp.List = make([]RoleInfo, 0, len(list))
	for _, v := range list {
		info := RoleInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

func (req RoleInfoRequest) Valid(opt string) error {
	if opt == ValidOptCreate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Name == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
		}
	} else if opt == ValidOptUpdate {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Name == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "name")
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
