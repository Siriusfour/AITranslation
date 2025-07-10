package Utils

import (
	"AITranslatio/Global"
	"AITranslatio/Src/DTO"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

func GeneratedToken(key []byte, method jwt.SigningMethod, UserID int, day time.Duration) (string, error) {

	playLoad := jwt.MapClaims{
		"iss":    "Suis",
		"sub":    "AITranslatio.cn",
		"UserID": UserID,
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
	// 校验 Claims 对象是否有效，0.签名是否正确， 1.是否过期，2.是否被吊销
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token解析失败")
	}

	// EXP
	if exp, ok := claims["exp"].(float64); ok {
		if exp < float64(time.Now().Unix()) {
			return nil, errors.New("token已过期")
		}
	}

	//验证token是否存在，是否被吊销，是否有被异地登录
	userID := int(claims["UserID"].(float64))

	tokenInfo, ok := Global.TokenMap.TokenMap[userID]
	if !ok || tokenInfo == nil {
		return nil, errors.New("token不存在")
	}

	if Global.TokenMap.TokenMap[userID].Revoked == false {
		return nil, errors.New("token已被注销")
	}

	if Global.TokenMap.TokenMap[userID].AccessToken != jwtStr && Global.TokenMap.TokenMap[userID].RefreshToken != jwtStr {

		return nil, fmt.Errorf("你的账号已于%s被异地登录", Global.TokenMap.TokenMap[userID].RegisteredTime)
	}

	return token.Claims, nil
}

func Verify(AccessToken string) error {
	if AccessToken != "" {
		//解析并校验
		_, err := ParseToken(Global.PKEY, AccessToken)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("token不存在")
	}
}

func Refresh(tokens DTO.Auth) (error, string) {

	if tokens.RefreshToken != "" {
		TokenInfo, err := ParseToken(Global.PKEY, tokens.RefreshToken)
		if err != nil {
			return err, ""
		}
		//刷新map里面的AK
		claims, ok := TokenInfo.(jwt.MapClaims)
		if !ok {
			return errors.New("token is empty"), ""
		}
		userID := int(claims["UserID"].(float64))

		// 检查TokenMap是否初始化
		if Global.TokenMap == nil {
			return errors.New("token map not initialized"), ""
		}

		// 检查TokenMap.TokenMap是否初始化
		if Global.TokenMap.TokenMap == nil {
			return errors.New("token map not initialized"), ""
		}

		// 检查userID是否存在
		tokenInfo, exists := Global.TokenMap.TokenMap[userID]
		if !exists {
			return errors.New("token not found"), ""
		}

		// 检查tokenInfo是否为nil（如果是指针类型）
		if tokenInfo == nil {
			return errors.New("token info is nil"), ""
		}

		AccessToken, err := GeneratedToken(Global.PKEY, jwt.SigningMethodHS256, userID, time.Duration(10000))
		if err != nil {
			return err, ""
		}
		tokenInfo.AccessToken = AccessToken
		return nil, AccessToken
	} else {
		return errors.New("token is empty"), ""
	}

}
