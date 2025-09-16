package WebAuthn

import (
	"AITranslatio/Global"
	"AITranslatio/Utils/captcha"
	"AITranslatio/app/http/DTO/NotAuthDTO"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"strconv"
	"time"
)

type Challenge struct {
	Verify    string `json:"challenge"`
	Timestamp int64  `json:"timestamp"`
	Domain    string `json:"Domain"`
}

func CreateChallenge() (*Challenge, error) {
	// 生成随机部分
	randomPart := make([]byte, 24)
	_, err := rand.Read(randomPart)
	if err != nil {
		return nil, err
	}

	// 生成时间戳
	timestamp := time.Now().Unix()

	// 组合挑战
	randomBytes := make([]byte, 32)
	copy(randomBytes[:24], randomPart)
	binary.BigEndian.PutUint64(randomBytes[24:], uint64(timestamp))

	challenge := base64.StdEncoding.EncodeToString(randomBytes)

	// 创建挑战数据
	challengeData := &Challenge{
		Verify:    challenge,
		Timestamp: timestamp,
		Domain:    Global.Config.GetString("DoMain"),
	}

	//把随机数存放在redis,有效期为5分钟
	Global.RedisClient.Set(context.Background(), challenge, 1, 5*time.Minute)

	return challengeData, nil
}

func ApplicationWebAuthn(DTO *NotAuthDTO.WebAuthnRegisterDTO) error {

	//1.发送邮箱验证码，存入redis
	coder, err := captcha.CreateCoderFactory(DTO.Email)
	if err != nil {
		return err
	}

	err = coder.SendEmailCode()
	if err != nil {
		return err
	}

	if err := Global.RedisClient.HSet(context.Background(), strconv.FormatInt(DTO.UserID, 10), "captcha", coder.Code); err != nil {
		return err.Err()
	}

	return nil
}

func VerifyChallenge(Challenge *Challenge) (bool, error) {

	//

	//1.校验redis里是否有该随机数，value值是否为EmailCode
	result := Global.RedisClient.GetDel(context.Background(), Challenge.Verify)
	if result.Err() != nil {
		return false, result.Err()
	} else {
		return true, nil
	}

}
