package token

import (
	"AITranslatio/Global"
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/SnowFlak"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"time"
)

type TokenProvider interface {
	GeneratedToken(int64, int) (string, error)
	ParseToken(string) error
}

type CreateToken struct {
	Key             []byte
	AkExp           time.Duration
	RkExp           time.Duration
	SnowFlakManager SnowFlak.SnowFlakManager
	RedisClient     redis.Cmdable
}

// CreateTokenFactory 0-AccessToken  1-RefreshToken
func CreateTokenFactory(c *CreateToken) *JWTGenerator {
	return &JWTGenerator{
		c.Key,
		c.AkExp,
		c.RkExp,
		c.SnowFlakManager,
		c.RedisClient,
	}
}

type JWTGenerator struct {
	encryptKey      []byte        // 密钥
	accessExpire    time.Duration // AK 过期时间
	refreshExpire   time.Duration // RK 过期时间
	SnowFlakManager SnowFlak.SnowFlakManager
	redisClient     redis.Cmdable
}

// redis存储的结构
type jwtInfo struct {
	UserID       int64 `redis:"UserID"`
	TokenID      int64 `redis:"TokenID"`
	RegisterTime int64 `redis:"RegisterTime"`
	jwt.RegisteredClaims
}

// GeneratedToken 生成token
func (t *JWTGenerator) GeneratedToken(userID int64, tokenType int) (string, error) {

	now := time.Now()
	expire := t.accessExpire
	subject := "AccessToken"

	if tokenType == Consts.RefreshToken {
		expire = t.refreshExpire
		subject = "RefreshToken"
	}

	jwtInfo := &jwtInfo{
		userID,
		t.SnowFlakManager.GetID(),
		t.SnowFlakManager.GetID(),
		jwt.RegisteredClaims{
			Subject:   subject,
			Issuer:    "MyProject",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expire)),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtInfo) //创建token
	token, err := tokenClaims.SignedString(t.encryptKey)              //对其签名

	// 使用事务 Pipeline
	Key := fmt.Sprintf("UserID:%d", userID)

	pipe := t.redisClient.TxPipeline()
	pipe.HSet(context.Background(), Key, map[string]interface{}{
		"TokenID":      jwtInfo.TokenID,
		"UserID":       jwtInfo.UserID,
		"RegisterTime": jwtInfo.ExpiresAt.Unix(),
	})
	pipe.Expire(context.Background(), Key, expire*time.Hour)

	_, err = pipe.Exec(context.Background())
	if err != nil {
		return "", fmt.Errorf("存储会话失败: %w", err)
	}

	return token, nil
}

// ParseToken 判断token是否有效   0.签名是否相同（防篡改）  1.是否过期     3.是否异地登录
func (t *JWTGenerator) ParseToken(VerifyToken string) error {

	//解析token信息到结构体
	jwtInfo := &jwtInfo{}
	token, err := jwt.ParseWithClaims(VerifyToken, jwtInfo, func(token *jwt.Token) (interface{}, error) {
		return t.encryptKey, nil
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

	UserID := jwtInfo.UserID
	TokenID := jwtInfo.TokenID

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
	if err := cmd.Scan(&jwtInfo); err != nil {
		return err
	}

	//3
	if jwtInfo.TokenID != TokenID {
		return fmt.Errorf(
			"%w: 你的账号已于 %s 被其他地方登录",
			MyErrors.ErrAccountKicked,
			time.Unix(jwtInfo.RegisterTime, 0).Format(time.DateTime),
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
func (t *JWTGenerator) Verify(Token string) error {

	//解析并校验
	err := t.ParseToken(Token)
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
