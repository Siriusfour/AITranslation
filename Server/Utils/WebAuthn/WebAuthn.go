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
type WebAuthn struct {
	//网站信息
	RPName string
	RPID   string

	PubKeyCredParams []Consts.Alg
	Attestation      string //认证器类型
	Attachment       string //认证器证明方式
	TimeOut          time.Duration
	ChallengeTTL     int64
}

func CreateWebAuthnConfigFactory(cfg interf.ConfigInterface) *WebAuthn {
	return &WebAuthn{
		RPName: cfg.GetString("WebAuthn.rp.RPName"),
		RPID:   cfg.GetString("WebAuthn.rp.RPID"),

		PubKeyCredParams: Consts.PubKeyCredParams,
		Attestation:      cfg.GetString("WebAuthn.Attestation"),
		Attachment:       cfg.GetString("WebAuthn.AuthenticatorAttachment"),
		TimeOut:          time.Minute * time.Duration(cfg.GetInt("WebAuthn.TimeOut")),

		ChallengeTTL: time.Now().Add(time.Hour * time.Duration(cfg.GetInt("WebAuthn.Challenge_TTL"))).Unix(),
	}
}

func (w *WebAuthn) CreateChallenge(redisClient *redis.Client, sessionID int64) (string, error) {
	// 生成随机部分
	randomPart := make([]byte, 24)
	_, err := rand.Read(randomPart)
	if err != nil {
		return "", err
	}

	// 生成时间戳
	timestamp := time.Now().Unix()

	// 组合挑战
	randomBytes := make([]byte, 32)
	copy(randomBytes[:24], randomPart)
	binary.BigEndian.PutUint64(randomBytes[24:], uint64(timestamp))

	challenge := base64.RawURLEncoding.EncodeToString(randomBytes)

	// 使用事务 Pipeline
	Key := fmt.Sprintf("SessionID:%d", sessionID)

	pipe := redisClient.TxPipeline()
	pipe.HSet(context.Background(), Key, map[string]interface{}{
		"webAuthn_challenge":         challenge,
		"webAuthn_challenge_OutTime": w.ChallengeTTL,
	})

	_, err = pipe.Exec(context.Background())
	if err != nil {
		return "", fmt.Errorf("存储会话失败: %w", err)
	}

	return challenge, nil
}
