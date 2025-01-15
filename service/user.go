package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/jinzhu/copier"
	log "github.com/wonderivan/logger"
	"strconv"
	"ticket-service/api/apimodel"
	config "ticket-service/conf"
	"ticket-service/database/model"
	"ticket-service/httpserver/errcode"
	"ticket-service/pkg/utils"
	"ticket-service/pkg/utils/redis"
	"time"
)

const ResetPassword = "123456"
const TokenRedisKey = "Authorization-user:"
const RedisLockKey = "redisLock"

func (operator *ResourceOperator) Login(c *gin.Context, req apimodel.UserInfoRequest) (*apimodel.LoginResponse, error) {
	var opt model.User
	var resp apimodel.LoginResponse
	selector := make(map[string]interface{})
	selector[model.FieldUserName] = req.Username
	err := operator.Database.PreloadEntityByFilter(model.TableNameUser, selector, model.OneQuery, &opt, []string{model.PreloadRoles})
	if err != nil {
		log.Error("用户名查找失败. err:[%v]", err)
		return nil, err
	}
	var roleName string
	if opt.Roles != nil {
		roleName = opt.Roles[0].Name
	}
	if opt.ID <= 0 {
		return nil, fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "用户名")
	}
	if ok := utils.BcryptCheck(req.Password, opt.Password); !ok {
		return nil, fmt.Errorf(errcode.ErrorMsgUserPassword)
	}
	j := &utils.JWT{SigningKey: []byte(config.Conf.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(utils.BaseClaims{
		UUID:     opt.UUID,
		ID:       uint64(opt.ID),
		UserID:   opt.ID,
		Username: opt.Username,
		RoleName: roleName,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, fmt.Errorf(errcode.ErrorMsgUnauthorized)
	}
	_, err = redis.RedisClient.Set(TokenRedisKey+strconv.Itoa(opt.ID), token, time.Hour*24*7).Result()
	if err != nil {
		log.Error("Login.token写入redis失败 err:[%v]", err)
		return nil, err
	}
	utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
	resp.Load(opt)
	resp.Token = token
	resp.ExpireAt = claims.RegisteredClaims.ExpiresAt.Unix() * 1000

	return &resp, nil
}

func (operator *ResourceOperator) Register(req *apimodel.UserInfoRequest) error {
	var opt model.User
	selector := make(map[string]interface{})
	//校验用户名
	selector[model.FieldUserName] = req.Username
	err := operator.Database.ListEntityByFilter(model.TableNameUser, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("用户名查找失败. err:[%v]", err)
		return err
	}

	if opt.UUID != uuid.Nil {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamExists, "用户名")
	}

	//保持id自增
	req.ID = 0
	//密码哈希加密，设置uuid
	req.Password = utils.BcryptHash(req.Password)
	err = copier.Copy(&opt, req)
	opt.UUID = uuid.Must(uuid.NewV4())
	if err != nil {
		return err
	}

	err = operator.Database.CreateEntity(model.TableNameUser, &opt)
	if err != nil {
		log.Error("用户注册失败. err:[%v]", err)
		return err
	}
	return nil
}
func (operator *ResourceOperator) UpdateUserInfo(req *apimodel.UserInfoRequest) error {
	var opt model.User
	//修改用户信息
	selector := make(map[string]interface{})
	selector[model.FieldID] = req.ID
	err := operator.Database.ListEntityByFilter(model.TableNameUser, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("查找用户失败. err:[%v]", err)
		return err
	}
	if opt.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "id")
	}
	if opt.UUID == uuid.Nil {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "uuid")
	}

	//保持ID、账号密码不变,暂存createTime（save方法全字段更新）
	req.ID = opt.ID
	req.UUID = opt.UUID.String()
	req.Username = opt.Username
	req.Password = opt.Password
	CreateTime := opt.CreatedAt
	err = copier.Copy(&opt, req)
	if err != nil {
		return err
	}
	opt.CreatedAt = CreateTime

	err = operator.Database.SaveEntity(model.TableNameUser, &opt)
	if err != nil {
		log.Error("用户数据更新失败. err:[%v]", err)
		return err
	}
	return nil
}
func (operator *ResourceOperator) DeleteUser(req *apimodel.UserInfoRequest) error {
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	selector[model.FieldID] = req.ID
	err := operator.Database.DeleteEntityByFilter(model.TableNameUser, selector, queryParams, &model.User{})
	if err != nil {
		log.Error("用户数据删除失败. err:[%v]", err)
		return err
	}
	selector = make(map[string]interface{})
	selector[model.FieldUserID] = req.ID
	err = operator.Database.DeleteEntityByFilter(model.TableNameUserRoles, selector, queryParams, &model.User{})
	if err != nil {
		log.Error("用户数据删除失败. err:[%v]", err)
		return err
	}
	return nil
}
func (operator *ResourceOperator) QueryUserList(req *apimodel.UserInfoRequest) (*apimodel.UserPageResponse, error) {
	var resp apimodel.UserPageResponse
	selector := make(map[string]interface{})
	queryParams := model.QueryParams{}
	//uuid查
	if req.UUID != "" {
		selector[model.FieldUUID] = req.UUID
	}
	if req.ID > 0 {
		selector[model.FieldID] = req.ID
	}
	//账号查
	if req.Username != "" {
		selector[model.FieldUserName] = req.Username
	}
	order := model.Order{
		Field:     model.FieldUpdatedTime,
		Direction: apimodel.OrderDesc,
	}
	queryParams.Orders = append(queryParams.Orders, order)
	if req.PageSize > 0 {
		queryParams.Limit = &req.PageSize
		offset := (req.PageNo - 1) * req.PageSize
		queryParams.Offset = &offset
	}
	//模糊查询
	//用户名
	if req.NickName != "" {
		var keyword []model.Keyword
		keyword = append(keyword, model.Keyword{Field: model.FieldNickName, Value: req.NickName, Type: 0})
		subquery := &model.SubQuery{
			Keywords: keyword,
		}
		queryParams.SubQueries = append(queryParams.SubQueries, subquery)
	}
	//电话号
	if req.Phone != "" {
		var keyword []model.Keyword
		keyword = append(keyword, model.Keyword{Field: model.FieldUserPhone, Value: req.Phone, Type: 0})
		subquery := &model.SubQuery{
			Keywords: keyword,
		}
		queryParams.SubQueries = append(queryParams.SubQueries, subquery)
	}
	//邮箱
	if req.Email != "" {
		var keyword []model.Keyword
		keyword = append(keyword, model.Keyword{Field: model.FieldUserEmail, Value: req.Email, Type: 0})
		subquery := &model.SubQuery{
			Keywords: keyword,
		}
		queryParams.SubQueries = append(queryParams.SubQueries, subquery)
	}
	if req.StartTime != "" && req.EndTime != "" {
		rangeQuery := &model.RangeQuery{
			Field: model.FieldCreatedTime,
			Start: req.StartTime,
			End:   req.EndTime,
		}
		queryParams.RangeQueries = append(queryParams.RangeQueries, rangeQuery)
	}

	var count int64
	var users []model.User
	err := operator.Database.CountEntityByFilter(model.TableNameUser, selector, queryParams, &count)
	if err != nil {
		return nil, err
	}
	err = operator.Database.PreloadEntityByFilter(model.TableNameUser, selector, queryParams, &users, []string{model.PreloadRoles})
	if err != nil {
		log.Error("用户数据查询失败. err:[%v]", err)
		return nil, err
	}

	resp.Load(count, users)
	return &resp, nil
}
func (operator *ResourceOperator) ChangePassword(req *apimodel.UserChangePWRequest) error {
	//验证码确认机制待添入
	var opt model.User
	selector := make(map[string]interface{})
	//selector[model.FieldID] = req.ID
	selector[model.FieldUUID] = req.UUID
	err := operator.Database.ListEntityByFilter(model.TableNameUser, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("数据查询失败. err:[%v]", err)
		return err
	}
	if opt.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "uuid")
	}
	if ok := utils.BcryptCheck(req.OldPass, opt.Password); !ok {
		return fmt.Errorf(errcode.ErrorMsgUserChangePass)
	}
	opt.Password = utils.BcryptHash(req.NewPass)
	err = operator.Database.SaveEntity(model.TableNameUser, &opt)
	if err != nil {
		log.Error("用户数据更新失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) QueryUserByUUID(uuid uuid.UUID) error {
	var user model.User
	selector := make(map[string]interface{})
	selector[model.FieldUUID] = uuid
	err := operator.Database.ListEntityByFilter(model.TableNameUser, selector, model.OneQuery, &user)
	if err != nil {
		log.Error("数据查询失败. err:[%v]", err)
		return err
	}
	if user.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "uuid")
	}
	return nil
}

func (operator *ResourceOperator) ResetPassword(req *apimodel.UserChangePWRequest) error {
	//验证码确认机制待添入
	var opt model.User
	selector := make(map[string]interface{})
	//selector[model.FieldID] = req.ID
	selector[model.FieldUUID] = req.UUID
	err := operator.Database.ListEntityByFilter(model.TableNameUser, selector, model.OneQuery, &opt)
	if err != nil {
		log.Error("数据查询失败. err:[%v]", err)
		return err
	}
	if opt.ID <= 0 {
		return fmt.Errorf(errcode.ErrorMsgSuffixParamNotExists, "uuid")
	}

	opt.Password = utils.BcryptHash(ResetPassword)
	err = operator.Database.SaveEntity(model.TableNameUser, &opt)
	if err != nil {
		log.Error("用户数据更新失败. err:[%v]", err)
		return err
	}
	return nil
}

func (operator *ResourceOperator) FreshToken(c *gin.Context) (string, error) {
	token := utils.GetToken(c)
	j := utils.NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		if errors.Is(err, utils.TokenExpired) {
			claims, err = j.ParseExpireToken(token)
			if err != nil {
				log.Error("解析token失败,err:[%v]", err)
				return "", fmt.Errorf(errcode.ErrorMsgInvalidToken)
			}
		} else {
			log.Error("解析token失败,err:[%v]", err)
			return "", fmt.Errorf(errcode.ErrorMsgInvalidToken)
		}
	}

	// 获取uid
	// 拼装redis key,验证有效性在30s内 return
	// set nx 获取分布式锁
	// 生成新的token
	// 写入redis
	// 释放分布式锁

	tokenDB, _ := redis.RedisClient.Get(TokenRedisKey + strconv.Itoa(claims.UserID)).Result()
	if tokenDB == "" {
		//redis中没有token => 代表用户退出登录清除token 或 redis中token过期 => 告诉前端需要重新登录
		return "", fmt.Errorf(errcode.ErrorMsgNeedReLogin)
	}
	j = utils.NewJWT()
	claimsDB, err := j.ParseToken(tokenDB)
	if err != nil {
		if errors.Is(err, utils.TokenExpired) {
			claimsDB, err = j.ParseExpireToken(tokenDB)
			if err != nil {
				log.Error("解析token失败,err:[%v]", err)
				return "", fmt.Errorf(errcode.ErrorMsgInvalidToken)
			}
		} else {
			log.Error("解析token失败,err:[%v]", err)
			return "", fmt.Errorf(errcode.ErrorMsgInvalidToken)
		}
	}
	if time.Now().Unix()-claimsDB.IssuedAt.Unix() < 30 {
		//若30秒内则判断为重复请求,返回原token
		return tokenDB, nil
	}
	//生成新token返回
	j = &utils.JWT{SigningKey: []byte(config.Conf.JWT.SigningKey)} // 唯一签名
	newClaims := j.CreateClaims(utils.BaseClaims{
		UUID:     claims.UUID,
		ID:       uint64(claims.UserID),
		UserID:   claims.UserID,
		Username: claims.Username,
		RoleName: claims.RoleName,
	})
	tokenRes, err := j.CreateToken(newClaims)
	if err != nil {
		return "", fmt.Errorf(errcode.ErrorMsgUnauthorized)
	}
	//直接使用setNX，新token不会覆盖旧token，没办法做到废弃之前的token，从而达到刷新的效果
	maxRetries := 3
	var uid string
	for i := 0; i < maxRetries; i++ {
		uid, err = redis.LockWithTimeout(RedisLockKey, time.Second*3, time.Second*1)
		if errors.Is(err, redis.LockError) {
			log.Error("redis分布式锁获取失败，继续重试 err:[%v] ", err)
			time.Sleep(time.Second * 1)
		}
	}
	defer func() {
		//释放分布式锁
		if err = redis.UnLock(RedisLockKey, uid); err != nil {
			if errors.Is(err, redis.UnLockError) {
				log.Error("释放分布式锁失败 err:[%v] uid:[%v]", err, uid)
				return
			}
		}
	}()
	if uid == "" {
		log.Error("获取分布式锁失败,err:[%v]", err)
		return "", fmt.Errorf(errcode.ErrorMsgRedisLock)
	}
	_, err = redis.RedisClient.Set(TokenRedisKey+strconv.Itoa(claims.UserID), tokenRes, time.Hour*24*3).Result()
	if err != nil {
		log.Error("fresh token写入redis失败 err:[%v]", err)
		return "", fmt.Errorf(errcode.ErrorMsgWriteRedis)
	}
	utils.SetToken(c, tokenRes, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
	return tokenRes, nil
}

func (operator *ResourceOperator) LogOut(c *gin.Context) error {
	token := utils.GetToken(c)
	j := utils.NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		if errors.Is(err, utils.TokenExpired) {
			claims, err = j.ParseExpireToken(token)
			if err != nil {
				log.Error("解析token失败,err:[%v]", err)
				return fmt.Errorf(errcode.ErrorMsgInvalidToken)
			}
		} else {
			log.Error("解析token失败,err:[%v]", err)
			return fmt.Errorf(errcode.ErrorMsgInvalidToken)
		}
	}
	_, err = redis.RedisClient.Del(TokenRedisKey + strconv.Itoa(claims.UserID)).Result()
	if err != nil {
		log.Error("删除redis缓存失败,err:[%v]", err)
		return fmt.Errorf(errcode.ErrorMsgDelRedis)
	}
	utils.ClearToken(c)
	return nil
}
