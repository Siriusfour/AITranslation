package OAuthService

import (
	"AITranslatio/Global"
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/SnowFlak"
	"AITranslatio/Utils/token"
	"AITranslatio/app/DAO/AuthDAO"
	"AITranslatio/app/Model/User"
	"AITranslatio/app/types"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type Github struct {
	DAO *AuthDAO.authDAO
}

type GitHubAppTokenResponse struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	Scope                 string `json:"scope"`
	ExpiresIn             int    `json:"expires_in"`               // 8 小时（秒）
	RefreshToken          string `json:"refresh_token"`            // ghr_...
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"` // 6 个月（秒）
}

type GitHubUser struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	Blog      string `json:"blog"`
	Location  string `json:"location"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
}

// GetChallenge 生成随机数，存入redis
func (Github *Github) GetChallenge(ctx *gin.Context) (*types.Challenge, error) {
	// 生成随机部分
	randomPart := make([]byte, 24)
	_, err := rand.Read(randomPart)
	if err != nil {
		return nil, err
	}

	// 生成时间戳
	timestamp := time.Now().Unix()

	// 组合挑战
	randomBytes := make([]byte, 32)
	copy(randomBytes[:24], randomPart)
	binary.BigEndian.PutUint64(randomBytes[24:], uint64(timestamp))

	challenge := base64.RawURLEncoding.EncodeToString(randomBytes)
	//以事务存入redis
	Key := fmt.Sprintf("SessionID:%s", ctx.GetString(Consts.ValidatorPrefix+"SessionID"))

	pipe := Global.RedisClient.TxPipeline()
	pipe.HSet(context.Background(), Key, map[string]interface{}{
		"OAuth_Github_challenge": "" + challenge,
	})

	pipe.Expire(context.Background(), Key, time.Duration(Global.Config.GetInt64("OAuth.Challenge_TTL"))*time.Second)

	_, err = pipe.Exec(context.Background())
	if err != nil {
		return nil, fmt.Errorf("存储会话失败: %w", err)
	}

	repsData := &types.Challenge{
		challenge,
	}

	return repsData, nil
}

func (Github *Github) VerifyChallenge(ctx *gin.Context) error {

	Key := fmt.Sprintf("SessionID:%s", ctx.GetString(Consts.ValidatorPrefix+"SessionID"))

	challenge, err := Global.RedisClient.HGet(ctx, Key, "OAuth_Github_challenge").Result()
	if err != nil {
		return fmt.Errorf("challenge不存在或过期：%w", err)
	}

	//验证传入的challenge和redis的challenge是否一致
	if ctx.GetString(Consts.ValidatorPrefix+"state") != challenge {
		return fmt.Errorf("两次challenge不一致，redis_challenge:%s,参数challenge为：%s", challenge, ctx.GetString(Consts.ValidatorPrefix+"state"))
	}

	return nil
}

func (Github *Github) GetUserInfo(ctx *gin.Context) (*types.LoginInfo, error) {

	//获取githubToken（存储到redis）
	Resp, err := GetToken("https://github.com/login/oauth/access_token", ctx)
	if err != nil {
		return nil, err
	}

	//以GithubToken换取UserInfo
	GithubUserInfo, err := GetUserInfo[GitHubUser]("https://api.github.com/user", "application/vnd.github+json", Resp.AccessToken)
	if err != nil {
		return nil, err
	}

	//由OAuthID判断是否存在该账号
	err, UserInfo := Github.DAO.CheckUserID(GithubUserInfo.ID, "GithubId")
	if err != nil {
		if errors.Is(err, MyErrors.ErrorOAuthIDrNotFound) { //如果是不存在该OAuthID，则直接创建账号
			userID := SnowFlak.CreateSnowflakeFactory().GetID()
			auth, err := Github.CreateToken(userID)
			LoginInfo, err := Github.CreateUser(GithubUserInfo, userID)
			LoginInfo.Auth.AccessToken = auth.AccessToken
			LoginInfo.Auth.RefreshToken = auth.RefreshToken
			return LoginInfo, err
		}
		return nil, err
	}

	auth, err := Github.CreateToken(UserInfo.UserID)
	if err != nil {
		return nil, err
	}

	//存在，直接返回在数据库找到用户
	return &types.LoginInfo{
		Auth: types.Auth{
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
		},
		Nickname: UserInfo.Nickname,
		UserID:   UserInfo.UserID,
		Avatar:   UserInfo.Avatar,
	}, nil
}

func (Github *Github) CreateUser(GithubUserInfo *GitHubUser, UserID int64) (*types.LoginInfo, error) {

	UserInfo := &User.User{
		GithubID: GithubUserInfo.ID,
		UserID:   UserID,
		Nickname: GithubUserInfo.Name,
		Email:    GithubUserInfo.Email,
		Avatar:   GithubUserInfo.AvatarURL,
	}

	err := Github.DAO.CreateUser(UserInfo)
	if err != nil {
		return nil, err
	}

	return &types.LoginInfo{
		Auth:     types.Auth{},
		Nickname: UserInfo.Nickname,
		UserID:   UserInfo.UserID,
		Avatar:   UserInfo.Avatar,
	}, nil

}

func (Github *Github) CreateToken(ID int64) (*types.Auth, error) {

	//创建APP的token
	AccessToken, errAk := token.CreateTokenFactory(Consts.AccessToken, ID).GeneratedToken()
	if errAk != nil {
		return nil, fmt.Errorf("生成token失败: %w", errAk)
	}
	RefreshToken, errRk := token.CreateTokenFactory(Consts.RefreshToken, ID).GeneratedToken()
	if errRk != nil {
		return nil, fmt.Errorf("生成token失败: %w", errAk)
	}

	return &types.Auth{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
	}, nil
}
