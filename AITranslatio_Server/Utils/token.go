package Utils

import (
	"AITranslatio/Global"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

type TokenInfo struct {
	RefreshToken   string `redis:"RefreshToken"`
	AccessToken    string `redis:"AccessToken"`
	RegisteredTime string `redis:"RegisteredTime"`
}

// 生成token
func GeneratedToken(key []byte, method jwt.SigningMethod, UserID int, day time.Duration) (string, error) {

	playLoad := jwt.MapClaims{
		"iss":    "Suis",
		"sub":    "AITranslatio.cn",
		"UserID": UserID,
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
	// 校验 Claims ，0.签名是否正确， 1.是否存在于redis 2.是否被异地登录(传入token与redis存在的token的值不同)
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token解析失败")
	}

	//验证token是否存在于redis
	userID := int(claims["UserID"].(float64))

	var tokenInfo TokenInfo

	err = Global.RedisClient.HGetAll(context.Background(), "userID:"+strconv.Itoa(userID)).Scan(&tokenInfo)
	if err != nil {
		return nil, errors.New("token不存在")
	}

	if tokenInfo.AccessToken != jwtStr && tokenInfo.RefreshToken != jwtStr {
		return nil, fmt.Errorf("你的账号已于%s被异地登录", tokenInfo.RegisteredTime)
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
		return errors.New("token值为空")
	}

	//TODO ak=“且rk！=“”時刷新ak
}

//func Refresh(RefreshCTX *gin.Context, tokens DTO.Auth) (error, string) {
//
//	if tokens.RefreshToken != "" {
//		//解析RefreshToken是否有效
//		TokenInfo, err := ParseToken(Global.PKEY, tokens.RefreshToken)
//		if err != nil {
//			return err, ""
//		}
//
//		//从token获取信息
//		claims, ok := TokenInfo.(jwt.MapClaims)
//		if !ok {
//			return errors.New("token is empty"), ""
//		}
//		userID := int(claims["UserID"].(float64))
//
//		// 检查TokenMap是否初始化
//		if Global.TokenMap == nil {
//			return errors.New("token map not initialized"), ""
//		}
//		if Global.TokenMap.TokenMap == nil {
//			return errors.New("token map not initialized"), ""
//		}
//
//		// 检查userID是否存在
//		tokenInfo, exists := Global.TokenMap.TokenMap[userID]
//		if !exists {
//			return errors.New("token not found"), ""
//		}
//
//		AccessToken, err := GeneratedToken(Global.PKEY, jwt.SigningMethodHS256, userID, time.Duration(10000))
//		if err != nil {
//			return err, ""
//		}
//		tokenInfo.AccessToken = AccessToken
//		return nil, AccessToken
//	} else {
//		return errors.New("token is empty"), ""
//	}
//
//}
