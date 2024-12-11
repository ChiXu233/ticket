package apimodel

import (
	"fmt"
	"github.com/gofrs/uuid/v5"
	"ticket-service/database/model"
	"ticket-service/httpserver/errcode"
	"ticket-service/pkg/utils"
)

//base struct

type UserInfo struct {
	ID        int       `json:"id"`
	UUID      uuid.UUID `json:"uuid"`     // 用户UUID
	Username  string    `json:"username"` // 用户登录名
	Password  string    `json:"password"` // 用户登录密码
	NickName  string    `json:"nickName"` // 用户昵称
	Phone     string    `json:"phone" `   // 用户手机号
	Email     string    `json:"email"`    // 用户邮箱
	CreatedAt string    `json:"created_time"`
	UpdatedAt string    `json:"updated_time"`
}

type TokenInfo struct {
	Token    string
	ExpireAt interface{}
}

//Request struct

type UserInfoRequest struct {
	ID        int    `json:"id" uri:"id" form:"id"`
	UUID      string `json:"uuid" uri:"uuid" form:"uuid"` // 用户UUID
	Username  string `json:"username" form:"username"`    // 用户登录名
	Password  string `json:"password"`                    // 用户登录密码
	NickName  string `json:"nickName" form:"nickname"`    // 用户昵称
	Phone     string `json:"phone" form:"phone"`          // 用户手机号
	Email     string `json:"email" form:"email"`          // 用户邮箱
	CreatedAt string `json:"created_time"`
	UpdatedAt string `json:"updated_time"`
	PicId     string `json:"pic_id" form:"pic_id"`
	Answer    string `json:"answer" form:"answer"`
	PaginationRequest
}

type UserChangePWRequest struct {
	ID      int       `json:"id"`
	UUID    uuid.UUID `json:"uuid"`
	OldPass string    `json:"oldPass"`
	NewPass string    `json:"newPass"`
}

// Response struct

type UserPageResponse struct {
	List []UserInfo `json:"list"`
	PaginationResponse
}

type LoginResponse struct {
	ID        int
	UUID      uuid.UUID
	Username  string
	NickName  string
	Phone     string
	Email     string
	CreatedAt string
	UpdatedAt string
	TokenInfo
}

//DataUnmarshal

func (u *UserInfo) Load(userData model.User) {
	u.UUID = userData.UUID
	u.ID = userData.ID
	u.Username = userData.Username
	u.Password = userData.Password
	u.NickName = userData.NickName
	u.Phone = userData.Phone
	u.Email = userData.Email
	u.CreatedAt = utils.TimeFormat(userData.CreatedAt)
	u.UpdatedAt = utils.TimeFormat(userData.UpdatedAt)
}

func (t *LoginResponse) Load(userData model.User) {
	t.UUID = userData.UUID
	t.ID = userData.ID
	t.Username = userData.Username
	t.NickName = userData.NickName
	t.Phone = userData.Phone
	t.Email = userData.Email
	t.CreatedAt = utils.TimeFormat(userData.CreatedAt)
	t.UpdatedAt = utils.TimeFormat(userData.UpdatedAt)
}

//DataLoading

func (resp *UserPageResponse) Load(total int64, list []model.User) {
	resp.List = make([]UserInfo, 0, len(list))
	for _, v := range list {
		info := UserInfo{}
		info.Load(v)
		resp.List = append(resp.List, info)
	}
	resp.TotalSize = int(total)
}

// Handler valid

func (req UserInfoRequest) Valid(opt string) error {
	if opt == ValidOptCreate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Username == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "username")
		}
		if req.Password == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "password")
		}
		if req.Phone == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "phone")
		}
		if req.Email == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "email")
		}
	} else if opt == ValidOptUpdate {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.Username == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "username")
		}
		if req.Password == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "password")
		}
		if req.Phone == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "phone")
		}
		if req.Email == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "email")
		}
	} else if opt == ValidOptDel {
		if req.ID <= 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
	} else if opt == ValidOptLogin {
		if req.Username == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "username")
		}
		if req.Password == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "password")
		}
		if req.PicId == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "验证码图片id")
		}
		if req.Answer == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "验证码")
		}
	} else {
		orderByFields := []string{model.FieldID, model.FieldName, model.FieldCreatedTime, model.FieldUpdatedTime}
		return req.PaginationRequest.Valid(orderByFields)
	}
	return nil
}

func (req UserChangePWRequest) Valid(opt string) error {
	if opt == ValidOptCreate {
		if req.ID < 0 {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "id")
		}
		if req.UUID == uuid.Nil {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "uuid")
		}
		if req.OldPass == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "oldPass")
		}
		if req.NewPass == "" {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "newPass")
		}
	} else if opt == ValidOptResetPwd {
		if req.UUID == uuid.Nil {
			return fmt.Errorf(errcode.ErrorMsgPrefixInvalidParameter, "uuid")
		}
	}
	return nil
}
