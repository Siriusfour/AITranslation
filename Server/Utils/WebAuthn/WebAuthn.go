package WebAuthn

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Global/Consts"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/redis/go-redis/v9"
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

type WebAuthn struct {
	Config    Config
	Challenge string
}

func CreateWebAuthnConfigFactory(cfg interf.ConfigInterface, UserID int64, UseName, Email string) *WebAuthn {
	return &WebAuthn{
		Config: Config{
			Name:             cfg.GetString("WebAuthn.rp.Name"),
			ID:               cfg.GetString("WebAuthn.rp.ID"),
			UserName:         UseName + "(" + Email + ")",
			PubKeyCredParams: Consts.PubKeyCredParams,
			Attestation:      cfg.GetString("WebAuthn.Attestation"),
			Attachment:       cfg.GetString("WebAuthn.AuthenticatorAttachment"),
			TimeOut:          time.Second * time.Duration(cfg.GetInt("WebAuthn.TimeOut")),
			UserID:           UserID,
			ChallengeTTL:     time.Now().Add(time.Hour * time.Duration(cfg.GetInt("WebAuthn.Challenge_TTL"))).Unix(),
		},
	}
}

func (w *WebAuthn) CreateChallenge(redisClient *redis.Client) error {
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

	// 使用事务 Pipeline
	Key := fmt.Sprintf("UserID:%d", w.Config.UserID)

	pipe := redisClient.TxPipeline()
	pipe.HSet(context.Background(), Key, map[string]interface{}{
		"challenge":         w.Challenge,
		"challenge_OutTime": w.Config.ChallengeTTL,
	})

	_, err = pipe.Exec(context.Background())
	if err != nil {
		return fmt.Errorf("存储会话失败: %w", err)
	}

	return nil

}
