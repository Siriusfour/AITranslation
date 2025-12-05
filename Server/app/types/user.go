package types

type LoginInfo struct {

	//鉴权信息
	Auth Auth `json:"auth" `

	//用户信息
	Nickname string `json:"nickname" example:"Suis"`
	UserID   int64  `json:"userid" example:"1162667863010451456"`
	Avatar   interface{}
}

type Challenge struct {
	Challenge string
}
