package errcode

import (
	"strings"
)

const (
	SuccessCodeBusiness = 0
	SuccessMsgBusiness  = "success"

	// ErrorCodeInternal HTTP CODE
	ErrorCodeInternal = 500
	// ErrorMsgInternal 默认错误信息
	ErrorMsgInternal          = "系统错误"
	ErrorCodeInvalidParameter = 400
	ErrorCodeUnauthorized     = 401
	ErrorMsgUnauthorized      = "认证或授权失败"
	ErrorMsgUnknownToken      = "无法识别token，用户不存在或签发人错误"
	ErrorMsgUnknownAuthorized = "非法登录"
	ErrorMsgExpireToken       = "token过期"
	ErrorMsgNeedReLogin       = "登录失效，请重新登录"
	ErrorMsgInvalidToken      = "无效token"
	ErrorCodeNotfound         = 404
	ErrorMsgNotfound          = "无资源错误"
	ErrorMsgNotAuth           = "暂无权限"

	// ErrorMsgPrefixInvalidParameter 错误信息前缀
	ErrorMsgPrefixInvalidParameter = "参数验证错误%v"

	ErrorMsgSuffixParamExists    = "%v已经存在"
	ErrorMsgSuffixParamNotExists = "%v不存在"
	ErrorMsgNoTickets            = "%v已售罄"

	// ErrorCodeBusiness Business Code
	ErrorCodeBusiness = 9999

	ErrorMsgMethodNotFound = "请求方法不允许"
	ErrorMsgHandleNotFound = "请求URL不存在"

	ErrorMsgLoadParam     = "读取请求参数失败"
	ErrorMsgValidateParam = "参数验证错误"
	ErrorMsgAtoiParam     = "参数转换失败"

	ErrorMsgUserNameOrPassword = "用户名尚未注册"
	ErrorMsgUserLogin          = "用户登录失败"
	ErrorMsgUserPassword       = "密码错误"
	ErrorMsgUserChangePass     = "旧密码错误"
	ErrorMsgCaptcha            = "验证码生成错误"
	ErrorMsgValidateCaptcha    = "验证码错误"

	ErrorMsgGetUserInfo  = "用户信息获取失败"
	ErrorMsgUserLoginOut = "退出登录失败"

	ErrorMsgCreateData     = "创建数据失败"
	ErrorMsgListData       = "获取数据失败"
	ErrorMsgUpdateData     = "修改数据失败"
	ErrorMsgDeleteData     = "删除数据失败"
	ErrorMsgCancel         = "取消失败"
	ErrorMsgPay            = "支付失败"
	ErrorMsgCreateOrUpdate = "修改/创建数据失败"
	//ErrorMsgBatchCreate    = "批量创建数据失败"

	ErrorMsgTicketNoNum         = "车票售罄"
	ErrorMsgDataExists          = "记录已经存在"
	ErrorMsgDataNotExists       = "记录不存在"
	ErrorMsgTransactionOpen     = "事务开启失败"
	ErrorMsgTransactionCommit   = "事务提交失败"
	ErrorMsgTransactionRollback = "事务回滚失败"
	ErrorMsgRedisLock           = "获取分布式锁失败"
	ErrorMsgRedisUnLock         = "释放分布式锁失败"
	ErrorMsgSetNX               = "写入redis失败"
	ErrorMsgLogOut              = "退出登录失败"
	ErrorMsgWriteRedis          = "写入redis失败"
	ErrorMsgDelRedis            = "删除redis失败"
)

var (
	ErrCode = map[string]int{

		ErrorMsgMethodNotFound: 5000,
		ErrorMsgHandleNotFound: 5001,

		ErrorMsgLoadParam:     5002,
		ErrorMsgValidateParam: 5003,
		ErrorMsgAtoiParam:     5004,

		ErrorMsgCreateData: 5005,
		ErrorMsgListData:   5006,
		ErrorMsgUpdateData: 5007,
		ErrorMsgDeleteData: 5008,

		ErrorMsgUserNameOrPassword: 5010,
		ErrorMsgUserLogin:          5011,
		ErrorMsgGetUserInfo:        5012,
		ErrorMsgUserLoginOut:       5013,
		ErrorMsgUserPassword:       5014,
		ErrorMsgUserChangePass:     5015,

		ErrorMsgTransactionOpen:     6002,
		ErrorMsgTransactionCommit:   6003,
		ErrorMsgTransactionRollback: 6004,
	}

	// CommonErrorMsg 通用错误信息
	CommonErrorMsg = []string{
		ErrorMsgSuffixParamExists,
		ErrorMsgSuffixParamNotExists,
		ErrorMsgPrefixInvalidParameter,
		ErrorMsgNoTickets,
		ErrorMsgUnauthorized,
	}

	// PostProcessingMsg 通用的错误处理信息后置处理
	PostProcessingMsg = map[string]string{
		ErrorMsgSuffixParamExists:      ErrorMsgDataExists,
		ErrorMsgNoTickets:              ErrorMsgTicketNoNum,
		ErrorMsgSuffixParamNotExists:   ErrorMsgDataNotExists,
		ErrorMsgPrefixInvalidParameter: ErrorMsgValidateParam,
	}
)

// GetErrorCode 从已经记录在案的code中查询，如果有则返回，没有返回默认 ErrorCodeBusiness
func GetErrorCode(msg string) int {
	code, ok := ErrCode[msg]
	if !ok {
		for _, item := range CommonErrorMsg {
			t := strings.TrimPrefix(item, "%v")
			t = strings.TrimSuffix(t, "%v")
			if strings.Contains(msg, t) {
				return ErrCode[PostProcessingMsg[item]]
			}
		}
		return ErrorCodeBusiness
	}
	return code
}
