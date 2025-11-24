package types

type UserInfo struct {

	//鉴权信息
	Auth Auth

	//用户信息
	UserName string
	UserID   int64
	Avatar   interface{}
}
