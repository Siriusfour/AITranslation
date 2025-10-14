package WebAuthn

import (
	"AITranslatio/Global"
	"AITranslatio/Global/Consts"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"strconv"
	"time"
)

// webAuthn服务器配置信息
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
}

type WebAuthn struct {
	Config    Config
	Challenge string
}

func CreateWebAuthnConfigFactory(UseName, Email string) *WebAuthn {
	return &WebAuthn{
		Config: Config{
			Name:             Global.Config.GetString("WebAuthn.rp.Name"),
			ID:               Global.Config.GetString("WebAuthn.rp.ID"),
			UserName:         UseName + "(" + Email + ")",
			PubKeyCredParams: Consts.PubKeyCredParams,
			Attestation:      Global.Config.GetString("WebAuthn.Attestation"),
			Attachment:       Global.Config.GetString("WebAuthn.AuthenticatorAttachment"),
			TimeOut:          time.Second * time.Duration(Global.Config.GetInt("WebAuthn.TimeOut")),
		},
		//Challenge:"",
	}
}

func (w *WebAuthn) CreateChallenge(UserID int64) error {
	// 生成随机部分
	randomPart := make([]byte, 24)
	_, err := rand.Read(randomPart)
	if err != nil {
		return err
	}

	// 生成时间戳
	timestamp := time.Now().Unix()

	// 组合挑战
	randomBytes := make([]byte, 32)
	copy(randomBytes[:24], randomPart)
	binary.BigEndian.PutUint64(randomBytes[24:], uint64(timestamp))

	challenge := base64.StdEncoding.EncodeToString(randomBytes)

	w.Challenge = challenge

	//把随机数存放在redis,有效期为5分钟
	Global.RedisClient.Set(context.Background(), strconv.FormatInt(UserID, 10), challenge, 5*time.Minute)

	return nil

}

//
//func (w *WebAuthn) VerifyChallenge(Challenge string) (bool, error) {
//
//}
