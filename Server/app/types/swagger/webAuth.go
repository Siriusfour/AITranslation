package swagger

import (
	"AITranslatio/Global/Consts"
	"time"
)

// WebAuthn 服务器配置信息
type WebAuthnInfo struct {
	//网站信息
	RPName string
	RPID   string

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

type LoginWebAuthnInfo struct {
	Challenge string
	RPName    string
	RPID      string
	TimeOut   time.Duration
}

type RegisterWebAuthnInfo struct {
	Challenge string
	Info      WebAuthnInfo
}
