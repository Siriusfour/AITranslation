package token

import (
	"AITranslatio/Global"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/SnowFlak"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type tokenConfig interface {
	GetInt(key string) int
}

// CreateTokenFactory 0-AccessToken  1-RefreshToken
func CreateTokenFactory(TokenType int, UserID int64) *Token {

	Token := &Token{
		"AccessToken",
		UserID,
		SnowFlak.CreateSnowflakeFactory().GetID(),
		Global.EncryptKey,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(Global.Config.GetInt("Token.AkOutTime")) * time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Global.Config,
	}

	if TokenType == 1 {
		Token.Type = "RefreshToken"
		Token.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(Global.Config.GetInt("Token.RkOutTime")) * time.Hour * 24))
	}

	return Token
}

type Token struct {
	Type       string //ak or rk
	UserID     int64  `json:"UserID"`  //用户ID
	TokenID    int64  `json:"TokenID"` //Tokne的唯一ID
	EncryptKey []byte
	jwt.RegisteredClaims
	config tokenConfig
}

// redis存储的结构
type UserAuth struct {
	TokenID      int64  `redis:"TokenID"`
	UserID       string `redis:"UserID"`
	RegisterTime int64  `redis:"RegisterTime"`
}

// GeneratedToken 生成token
func (t *Token) GeneratedToken() (string, error) {

	TokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, t) //创建token
	Token, err := TokenClaims.SignedString(t.EncryptKey)        //对其签名

	// 使用事务 Pipeline
	Key := fmt.Sprintf("UserID:%d", t.UserID)

	TTL := time.Duration(Global.Config.GetInt("Token.RkOutTime")) * time.Hour

	pipe := Global.RedisClient.TxPipeline()
	pipe.HSet(context.Background(), Key, map[string]interface{}{
		"TokenID":      t.TokenID,
		"UserID":       t.UserID,
		"RegisterTime": t.ExpiresAt.Unix(),
	})
	pipe.Expire(context.Background(), Key, TTL)

	_, err = pipe.Exec(context.Background())
	if err != nil {
		return "", fmt.Errorf("存储会话失败: %w", err)
	}

	return Token, nil
}

// ParseToken 判断token是否有效   0.签名是否相同（防篡改）  1.是否过期     3.是否异地登录
func ParseToken(VerifyToken string) error {

	//解析token信息到结构体
	claims := &Token{}
	token, err := jwt.ParseWithClaims(VerifyToken, claims, func(token *jwt.Token) (interface{}, error) {
		return Global.EncryptKey, nil
	})

	//如果是token过期，则返回自定义的token过期错误
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return MyErrors.ErrTokenExpired
		}
		return err
	}

	//0 / 1
	if !token.Valid {
		return MyErrors.ErrTokenInvalid
	}

	//2
	var UserInfo UserAuth //存储来自redis的数据
	UserID := claims.UserID
	TokenID := claims.TokenID

	//从redis获取数据到结构体
	cmd := Global.RedisClient.HGetAll(context.Background(), fmt.Sprintf("UserID:%d", UserID))
	if err := cmd.Err(); err != nil {

		return err
	}
	m := cmd.Val()
	if len(m) == 0 {
		// 说明这个 key 不存在，token已经到期删除
		return MyErrors.ErrTokenExpired
	}
	if err := cmd.Scan(&UserInfo); err != nil {
		return err
	}

	//3
	if UserInfo.TokenID != TokenID {
		return fmt.Errorf(
			"%w: 你的账号已于 %s 被其他地方登录",
			MyErrors.ErrAccountKicked,
			time.Unix(UserInfo.RegisterTime, 0).Format(time.DateTime),
		)
	}
	return nil
}

// GetDataFormToken 从token里面解析出某个值
func GetDataFormToken[T any](Token string, arg string) (error, T) {

	var zero T

	token, err := jwt.Parse(Token, func(token *jwt.Token) (interface{}, error) {
		return Global.EncryptKey, nil
	})
	if err != nil {
		return err, zero
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return MyErrors.ErrorAssert, zero
	}

	if arg == "UserID" {
		value, ok := claims[arg].(float64)
		if !ok {
			return MyErrors.ErrorAssert, zero
		}
		return nil, any(int64(value)).(T)
	}
	if arg == "TokenID" {
		value, ok := claims[arg].(float64)
		if !ok {
			return MyErrors.ErrorAssert, zero
		}
		return nil, any(int64(value)).(T)
	}

	value, ok := claims[arg].(T)
	if !ok {
		return MyErrors.ErrorAssert, zero
	}

	return nil, value
}

// Verify 验证token
func Verify(Token string) error {

	//解析并校验
	err := ParseToken(Token)
	if err != nil {
		return err
	}
	return nil

	//TODO ak=""且rk！=“”時刷新ak
}

// func Refresh(RefreshCTX *gin.Context, tokens DTO.AuthDTO) (error, string) {
//
//	if tokens.RefreshToken != "" {
//		//解析RefreshToken是否有效
//		Token, err := ParseToken(Global.EncryptKey, tokens.RefreshToken)
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
//		AccessToken, err := GeneratedToken(Global.EncryptKey, jwt.SigningMethodHS256, userID, time.Duration(10000))
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
