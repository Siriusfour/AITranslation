package token

import (
	"AITranslatio/Global"
	"AITranslatio/Global/CustomErrors"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

// 0-AccessToken  1-RefreshToken
func CreateTokenFactory(OutTime int) *TokenInfo {

	return &TokenInfo{
		"Suis",
		time.Now().Add(time.Duration(OutTime) * time.Hour * 24).Unix(),
		time.Now().Unix(),
		Global.PKEY,
		0,
	}
}

type TokenInfo struct {
	iss       string
	exp       int64
	iat       int64
	EncrypKey []byte
	UserID    int
}

type Token struct {
	RefreshToken   string `redis:"RefreshToken"`
	AccessToken    string `redis:"AccessToken"`
	RegisteredTime string `redis:"RegisteredTime"`
}

// 生成token  秘钥 userID
func (t *TokenInfo) GeneratedToken(UserID int64) (string, error) {

	playLoad := jwt.MapClaims{
		"iss":    t.iss,
		"exp":    t.exp,
		"iat":    t.iat,
		"UserID": UserID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, playLoad) //创建token
	return token.SignedString(t.EncrypKey)                       //生成签名防止伪造
}

// 判断token是否有效   0.签名是否相同（防篡改） 1.是否过期  2.是否被吊销（是否存在与redis）   3.是否异地登录
func ParseToken(PKey []byte, VerifyToken string) error {

	token, err := jwt.Parse(VerifyToken, func(token *jwt.Token) (interface{}, error) {
		return PKey, nil
	})

	if err != nil {
		return err
	}

	//0  1
	if !token.Valid {
		return errors.New(CustomErrors.ErrorsTokenInvalid)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New(CustomErrors.ErrorAssert)
	}

	//2
	var TokenFromRedis Token //存储来自redis的数据
	userID := int(claims["UserID"].(float64))
	err = Global.RedisClient.HGetAll(context.Background(), "userID_"+strconv.Itoa(userID)).Scan(&TokenFromRedis)
	if err != nil {
		return errors.New(CustomErrors.ErrorRedisGetDATA + ":" + err.Error())
	}

	if TokenFromRedis.RefreshToken == "" && TokenFromRedis.AccessToken == "" {
		return errors.New(CustomErrors.ErrorsTokenNotActiveYet)
	}

	//3
	if TokenFromRedis.RefreshToken != VerifyToken && TokenFromRedis.AccessToken != VerifyToken {
		return fmt.Errorf("你的账号已于%s被异地登录", TokenFromRedis.RegisteredTime)
	}

	return nil
}

// 验证token
func Verify(Token string) error {

	//解析并校验
	err := ParseToken(Global.PKEY, Token)
	if err != nil {
		return err
	}
	return nil

	//TODO ak=“且rk！=“”時刷新ak
}

// func Refresh(RefreshCTX *gin.Context, tokens DTO.Auth) (error, string) {
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
// }
