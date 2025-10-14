package Captcha

import (
	"AITranslatio/Global"
	"context"
	"crypto/rand"
	"math/big"
	"time"
)

type Coder struct {
	Code  string
	Email string
	Phone string
}

func CreateCoderFactory(Email string) (*Coder, error) {

	Max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, Max)
	if err != nil {
		return nil, nil
	}

	return &Coder{
		Code:  n.String(),
		Email: Email,
	}, nil

}

func (Coder *Coder) SendEmailCode() error {

	//存储在redis，并从setting获取duration过期时间
	Global.RedisClient.Set(context.Background(), Coder.Code, 1, 5*60*60*time.Second)

	//发送给目标邮箱
	return nil

}

func (Coder *Coder) SendPhoneCode() error {

	//存储在redis，并从setting获取duration过期时间

	//发送给目标邮箱

	return nil

}
