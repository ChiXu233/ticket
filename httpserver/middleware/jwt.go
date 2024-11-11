package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"ticket-service/api/handler"
	config "ticket-service/conf"
	"ticket-service/httpserver/app"
	"ticket-service/httpserver/errcode"
	utils "ticket-service/pkg/utils"
	"time"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := utils.GetToken(c)
		if token == "" {
			app.SendAuthorizedErrorResponse(c, errcode.ErrorMsgUnknownAuthorized)
			c.Abort()
			return
		}

		//黑名单
		//if jwtService.IsBlacklist(token) {
		//	response.FailWithDetailed(gin.H{"reload": true}, "您的帐户异地登陆或令牌失效", c)
		//	utils.ClearToken(c)
		//	c.Abort()
		//	return
		//}

		j := utils.NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				//授权过期
				app.SendAuthorizedErrorResponse(c, errcode.ErrorMsgExpireToken)
				utils.ClearToken(c)
				c.Abort()
				return
			}
			//验证失败
			app.SendAuthorizedErrorResponse(c, errcode.ErrorMsgUnauthorized)
			//response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
			utils.ClearToken(c)
			c.Abort()
			return
		}

		//判断该token所携带用户信息是否正确

		if claims.Issuer != config.Conf.JWT.Issuer {
			app.SendAuthorizedErrorResponse(c, errcode.ErrorMsgUnknownToken)
			utils.ClearToken(c)
			c.Abort()
			return
		}

		err = handler.NewHandler().Operator.QueryUserByUUID(claims.UUID)
		if claims.UUID == uuid.Nil || err != nil {
			app.SendAuthorizedErrorResponse(c, errcode.ErrorMsgUnknownToken)
			utils.ClearToken(c)
			c.Abort()
			return
		}

		// 已登录用户被管理员禁用 需要使该用户的jwt失效 此处比较消耗性能 如果需要 请自行打开
		// 用户被删除的逻辑 需要优化 此处比较消耗性能 如果需要 请自行打开
		//if user, err := userService.FindUserByUuid(claims.UUID.String()); err != nil || user.Enable == 2 {
		//	_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
		//	response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
		//	c.Abort()
		//}

		c.Set("claims", claims)
		c.Next()
		//若处于缓冲时间则颁发new_token
		if claims.ExpiresAt.Unix()-time.Now().Unix() < claims.BufferTime {
			dr, _ := utils.ParseDuration(config.Conf.JWT.ExpiresTime)
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(dr))
			newToken, _ := j.CreateTokenByOldToken(token, *claims)
			newClaims, _ := j.ParseToken(newToken)
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt.Unix(), 10))
			utils.SetToken(c, newToken, int(dr.Seconds()))

			//单点登录
			//if global.CONFIG.System.UseMultipoint {
			//	RedisJwtToken, err := jwtService.GetRedisJWT(newClaims.Username)
			//	if err != nil {
			//		global.LOG.Error("get redis jwt failed", zap.Error(err))
			//	} else { // 当之前的取成功时才进行拉黑操作
			//		_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: RedisJwtToken})
			//	}
			//	// 无论如何都要记录当前的活跃状态
			//	_ = jwtService.SetRedisJWT(newToken, newClaims.Username)
			//}

		}
	}
}
