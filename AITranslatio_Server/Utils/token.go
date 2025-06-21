package Utils

import (
	"AITranslatio/Src/DTO"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

func GeneratedToken(key []byte, method jwt.SigningMethod, LoginDTO *DTO.LoginDTO, day time.Duration) (string, error) {

	playLoad := jwt.MapClaims{
		"iss":  "Suis",
		"sub":  "AITranslatio.cn",
		"aud":  LoginDTO.UserID,
		"UUID": LoginDTO.UserID,
		"exp":  time.Now().Add(day * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(method, playLoad)
	return token.SignedString(key)
}

func ParseToken(key []byte, jwtStr string) (jwt.Claims, error) {
	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	fmt.Println("token:", token)
	if err != nil {
		return nil, err
	}
	// 校验 Claims 对象是否有效，0.签名是否正确， 1.是否过期，2.是否是原设备
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims, nil
}
