package CustomErrors

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
