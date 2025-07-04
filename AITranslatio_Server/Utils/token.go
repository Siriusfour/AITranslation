package Utils

import (
	"AITranslatio/Global"
	"AITranslatio/Src/DTO"
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

func GeneratedToken(key []byte, method jwt.SigningMethod, LoginDTO *DTO.LoginDTO, day time.Duration) (string, error) {

	playLoad := jwt.MapClaims{
		"iss":    "Suis",
		"sub":    "AITranslatio.cn",
		"UserID": LoginDTO.UserID,
		"UUID":   LoginDTO.UUID,
		"exp":    time.Now().Add(day * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(method, playLoad)
	return token.SignedString(key)
}

func ParseToken(key []byte, jwtStr string) (jwt.Claims, error) {
	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}
	// 校验 Claims 对象是否有效，0.签名是否正确， 1.是否过期，2.是否是原设备
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// EXP
	if exp, ok := claims["exp"].(float64); ok {
		if exp < float64(time.Now().Unix()) {
			return nil, errors.New("token expired")
		}
	}

	userID := claims["UserID"].(string)
	//原设备验证（解析出来的token与map的ak、rk都不相等则不为原设备）
	Global.TokenMap.MU.RLock()
	defer Global.TokenMap.MU.RUnlock()
	if Global.TokenMap.TokenMap[userID].AccessToken != jwtStr && Global.TokenMap.TokenMap[userID].RefreshToken != jwtStr {
		return nil, errors.New("invalid token")
	}

	return token.Claims, nil
}
