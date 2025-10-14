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

type TokenMothd interface {
	GeneratedToken() (string, error)
	ParseToken(string) error
}

// 0-AccessToken  1-RefreshToken
func CreateTokenFactory(TokenType int, UserID int64) TokenMothd {

	Token := &Token{
		"AccessToken",
		"Suis",
		time.Now().Add(time.Duration(Global.Config.GetInt("Token.AkOutTime")) * time.Hour * 24).Unix(),
		time.Now().Unix(),
		Global.PKEY,
		UserID,
	}

	if TokenType == 1 {
		Token.Type = "RefreshToken"
		Token.exp = time.Now().Add(time.Duration(Global.Config.GetInt("Token.RkOutTime")) * time.Hour * 24).Unix()
	}

	return Token
}

type Token struct {
	Type      string //ak or rk
	iss       string //颁发人
	exp       int64  //过期时间
	iat       int64  //颁发时间
	EncrypKey []byte //签名秘钥
	UserID    int64  //用户ID
}

type TokenDTO struct {
	RefreshToken string `redis:"RefreshToken"`
	AccessToken  string `redis:"AccessToken"`
}

// 生成token  秘钥 userID
func (t *Token) GeneratedToken() (string, error) {

	PlayLoad := jwt.MapClaims{
		"iss":    t.iss,
		"exp":    t.exp,
		"iat":    t.iat,
		"UserID": t.UserID,
	}

	TokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, PlayLoad) //创建token
	Token, err := TokenClaims.SignedString(t.EncrypKey)

	//存入redis
	err = Global.RedisClient.HSet(context.Background(), strconv.FormatInt(t.UserID, 10), map[string]interface{}{
		t.Type: Token,
	}).Err()

	if err != nil {
		return "", err
	}

	return Token, err //生成签名防止伪造

}

// 判断token是否有效   0.签名是否相同（防篡改）  1.是否过期   2.是否被吊销（是否存在与redis）   3.是否异地登录
func (t *Token) ParseToken(VerifyToken string) error {

	token, err := jwt.Parse(VerifyToken, func(token *jwt.Token) (interface{}, error) {
		return t.EncrypKey, nil
	})

	if err != nil {
		return err
	}

	//0  、  1
	if !token.Valid {
		return errors.New(CustomErrors.ErrorsTokenInvalid)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New(CustomErrors.ErrorAssert)
	}

	//2
	var TokenFromRedis TokenDTO //存储来自redis的数据
	userID := int(claims["UserID"].(float64))

	//从redis获取数据到TokenDTO
	err = Global.RedisClient.HGetAll(context.Background(), strconv.Itoa(userID)).Scan(&TokenFromRedis)
	if err != nil {
		return errors.New(CustomErrors.ErrorRedisGetDATA + err.Error())
	}

	if TokenFromRedis.RefreshToken == "" && TokenFromRedis.AccessToken == "" {
		return errors.New(CustomErrors.ErrorsTokenNotActiveYet)
	}

	//3
	if TokenFromRedis.RefreshToken != VerifyToken && TokenFromRedis.AccessToken != VerifyToken {
		err, RegisterTime := GetDataFormToken[time.Time](TokenFromRedis.AccessToken, "iat")
		if err != nil {
			return err
		}

		return fmt.Errorf("你的账号已于%s被异地登录", time.Unix(RegisterTime.Unix(), 0).Format(time.DateTime))
	}

	return nil
}

// 从token里面解析出某个值
func GetDataFormToken[T any](Token string, arg string) (error, T) {

	var zero T

	token, err := jwt.Parse(Token, func(token *jwt.Token) (interface{}, error) {
		return Global.PKEY, nil
	})
	if err != nil {

		return err, zero
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New(CustomErrors.ErrorAssert), zero
	}

	if arg == "UserID" {
		value, ok := claims[arg].(float64)
		if !ok {
			return errors.New(CustomErrors.ErrorAssert), zero
		}

		return nil, any(int64(value)).(T)

	}

	value, ok := claims[arg].(T)
	if !ok {
		return errors.New(CustomErrors.ErrorAssert), zero
	}

	return nil, value

}

// 验证token
func (t *Token) Verify(Token string) error {

	//解析并校验
	err := t.ParseToken(Token)
	if err != nil {
		return err
	}
	return nil

	//TODO ak=“且rk！=“”時刷新ak
}

// func Refresh(RefreshCTX *gin.Context, tokens DTO.AuthDTO) (error, string) {
//
//	if tokens.RefreshToken != "" {
//		//解析RefreshToken是否有效
//		Token, err := ParseToken(Global.PKEY, tokens.RefreshToken)
//		if err != nil {
//			return err, ""
//		}
//
//		//从token获取信息
//		claims, ok := Token.(jwt.MapClaims)
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
