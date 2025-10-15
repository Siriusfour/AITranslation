package MyErrors

import "errors"

// token 错误
var (
	ErrTokenInvalid      = errors.New("token无效")
	ErrTokenExpired      = errors.New("token已过期")
	ErrTokenMalformed    = errors.New("token格式错误")
	ErrTokenNotActiveYet = errors.New("token尚未生效")
	ErrSessionExpired    = errors.New("会话已失效")
	ErrUnauthorized      = errors.New("未授权")
)

// 基础错误
var (
	ErrorAssert = errors.New("类型断言失败")
)

// webAthun错误
var (
	ErrorClientDataTypeIsFail      = errors.New("type不为webauthn.create！")
	ErrorClientDataChallengeIsFail = errors.New("challenge错误！")
	ErrorClientDataRPID_IsFail     = errors.New("域名错误，当前页面域名与配置文件的RPID不符！")
)
