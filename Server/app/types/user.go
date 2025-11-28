package types

type LoginInfo struct {

	//鉴权信息
	Auth Auth

	//用户信息
	Nickname string
	UserID   int64
	Avatar   interface{}
}

type Challenge struct {
	Challenge string
}
