package OAuthService

import (
	"AITranslatio/Global"
	"AITranslatio/Global/Consts"
	"AITranslatio/Utils/token"
	"AITranslatio/app/types"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Github struct{}

type GitHubAppTokenResponse struct {
	GithubID              int    `json:"github_id"`
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	Scope                 string `json:"scope"`
	ExpiresIn             int    `json:"expires_in"`               // 8 小时（秒）
	RefreshToken          string `json:"refresh_token"`            // ghr_...
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"` // 6 个月（秒）
}

// 生成随机数，存入redis
func (Github *Github) GetChallenge(ctx *gin.Context) (string, error) {
	// 生成随机部分
	randomPart := make([]byte, 24)
	_, err := rand.Read(randomPart)
	if err != nil {
		return "", err
	}

	// 生成时间戳
	timestamp := time.Now().Unix()

	// 组合挑战
	randomBytes := make([]byte, 32)
	copy(randomBytes[:24], randomPart)
	binary.BigEndian.PutUint64(randomBytes[24:], uint64(timestamp))

	challenge := base64.StdEncoding.EncodeToString(randomBytes)

	// 使用事务 Pipeline
	Key := fmt.Sprintf("UserID:%d", ctx.GetInt64("UserID"))

	pipe := Global.RedisClient.TxPipeline()
	pipe.HSet(context.Background(), Key, map[string]interface{}{
		"OAuth_Github_challenge":         "" + challenge,
		"OAuth_Github_challenge_OutTime": time.Now().Add(time.Hour * time.Duration(Global.Config.GetInt("OAuth.Challenge_TTL"))).Unix(),
	})

	_, err = pipe.Exec(context.Background())
	if err != nil {
		return "", fmt.Errorf("存储会话失败: %w", err)
	}

	return challenge, nil
}

func (Github *Github) VerifyChallenge(ctx *gin.Context) error {
	UserID := strconv.FormatInt(ctx.GetInt64(Consts.ValidatorPrefix+"UserID"), 10)
	exists, err := Global.RedisClient.HExists(ctx, UserID, "OAuth_Github_challenge").Result()
	if err != nil || !exists {
		return err
	}
	challenge, err := Global.RedisClient.HGet(ctx, UserID, "OAuth_Github_challenge").Result()
	if err != nil {
		return fmt.Errorf("redis获取challenge失败：%w", err)
	}
	OutTime, err := Global.RedisClient.HGet(ctx, UserID, "OAuth_Github_challenge_OutTime").Result()
	if err != nil {
		return fmt.Errorf("redis获取OutTime：%w", err)
	}

	if ctx.GetString(Consts.ValidatorPrefix+"challenge") != challenge {
		{
			return errors.New("challenge不存在")
		}
	}
	ts, _ := strconv.ParseInt(OutTime, 10, 64)
	out := time.Unix(ts, 0)

	if time.Now().After(out) {
		return errors.New("challeng过期")
	}

	return nil
}

func (Github *Github) GetUserInfo(ctx *gin.Context) (*types.Auth, error) {

	githubTokenURL := "https://api.github.com/oauth/token"

	data := url.Values{}
	data.Set("client_id", Global.Config.GetString("Github.client_id"))
	data.Set("client_secret", Global.Config.GetString("Github.client_secret"))
	data.Set("code", ctx.GetString(Consts.ValidatorPrefix+"code"))

	req, err := http.NewRequest("POST", githubTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp GitHubAppTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("empty access_token: %s", string(body))
	}

	AccessToken, errAk := token.CreateTokenFactory(Consts.AccessToken, ctx.GetInt64("UserID")).GeneratedToken()
	if errAk != nil {
		return nil, fmt.Errorf("生成token失败: %w", errAk)
	}
	RefreshToken, errRk := token.CreateTokenFactory(Consts.RefreshToken, ctx.GetInt64("UserID")).GeneratedToken()
	if errRk != nil {
		return nil, fmt.Errorf("生成token失败: %w", errAk)
	}

	return &types.Auth{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
	}, nil

}
