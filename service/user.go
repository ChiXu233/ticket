package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	config "ticket-service/conf"
	"ticket-service/httpserver/errcode"
	util "ticket-service/utils"
	"time"
)

func (operator *ResourceOperator) Login(c *gin.Context) (map[string]interface{}, error) {
	j := &util.JWT{SigningKey: []byte(config.Conf.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(util.BaseClaims{
		ID:          1,
		Username:    "aa",
		AuthorityId: 1,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, fmt.Errorf(errcode.ErrorMsgUnauthorized)
	}
	util.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
	resp := make(map[string]interface{})
	resp["User"] = "aa"
	resp["Token"] = token
	resp["ExpiresAt"] = claims.RegisteredClaims.ExpiresAt.Unix() * 1000
	return resp, nil

}
