package middleware

import (
	"github.com/gin-gonic/gin"
	"ticket-service/database/casbin"
	"ticket-service/httpserver/app"
	"ticket-service/httpserver/errcode"
)

// 结合jwt实现token校验、鉴权
func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader("user_name")
		if has, err := casbin.E.Enforce(user, c.Request.RequestURI, c.Request.Method); err != nil || !has {
			app.SendAuthorizedErrorResponse(c, errcode.ErrorMsgNotAuth)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
