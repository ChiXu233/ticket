package model

import "github.com/gofrs/uuid/v5"

type User struct {
	Model
	UUID     uuid.UUID `json:"uuid" gorm:"index;comment:用户UUID"`          // 用户UUID
	Username string    `json:"username" gorm:"index;comment:用户登录名"`       // 用户登录名
	Password string    `json:"password"  gorm:"comment:用户登录密码"`           // 用户登录密码
	NickName string    `json:"nickName" gorm:"default:系统用户;comment:用户昵称"` // 用户昵称
	Phone    string    `json:"phone"  gorm:"comment:用户手机号"`               // 用户手机号
	Email    string    `json:"email"  gorm:"comment:用户邮箱"`                // 用户邮箱
	//Orders   []UserOrder `json:"orders" gorm:"foreignKey:UserID"` //@TODO 用户与订单是两个模块，暂时没有主外键关联的必要。
	//HeaderImg string `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.to p/header.jpg;comment:用户头像"` // 用户头像
}

func (m *User) TableName() string {
	return TableNameUser
}
