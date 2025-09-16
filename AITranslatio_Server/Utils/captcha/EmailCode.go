package captcha

import (
	"crypto/rand"
	"math/big"
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

	//发送给目标邮箱

	return nil

}

func (Coder *Coder) SendPhoneCode() error {

	//存储在redis，并从setting获取duration过期时间

	//发送给目标邮箱

	return nil

}
