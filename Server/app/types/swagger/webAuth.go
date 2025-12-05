package swagger

import (
	"AITranslatio/Global/Consts"
	"time"
)

// WebAuthn 服务器配置信息
type Config struct {
	//网站信息
	Name string
	ID   string

	//注册服务的用户信息
	UserID   int64
	UserName string

	//webAuthn服务器配置信息
	PubKeyCredParams []Consts.Alg
	Attestation      string //认证器类型
	Attachment       string //认证器证明方式
	TimeOut          time.Duration
	ChallengeTTL     int64
}

type WebAuthnInfo struct {
	Config    Config
	Challenge string
}
