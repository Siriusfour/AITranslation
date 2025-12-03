package OAuthService

import (
	"AITranslatio/Config/interf"

	"AITranslatio/Global/Consts"
	"AITranslatio/Utils/SnowFlak"
	"AITranslatio/Utils/token"
	"AITranslatio/app/DAO/AuthDAO"
	"AITranslatio/app/types"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
)

type OAuthService interface {
	GetChallenge(*gin.Context) (*types.Challenge, error) //生成随机挑战
	VerifyChallenge(*gin.Context) error                  //验证随机挑战
	GetUserInfo(*gin.Context) (*types.LoginInfo, error)  //从OAuth提供方获取token
}

func CreateOAuthServiceFactroy(cfg interf.ConfigInterface, logger *zap.Logger, JWTGenerator *token.JWTGenerator, snowFlakManager *SnowFlak.SnowFlakeGenerator, redis *redis.Client, DAO AuthDAO.Inerf, server string) OAuthService {

	switch server {
	case "QQ":

	case "WX":

	case "Github":
		return &GithubService{
			cfg,
			logger,
			JWTGenerator,
			snowFlakManager,
			redis,
			DAO,
		}

	default:
		return nil
	}
	return nil
}

func GetToken(cfg interf.ConfigInterface, URL string, ctx *gin.Context) (*GitHubAppTokenResponse, error) {

	data := url.Values{}
	data.Set("client_id", cfg.GetString("OAuth.Github.client_id"))
	data.Set("client_secret", cfg.GetString("OAuth.Github.client_secret"))
	data.Set("code", ctx.GetString(Consts.ValidatorPrefix+"code"))

	req, err := http.NewRequest("POST", URL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
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

	var Resp GitHubAppTokenResponse
	if err := json.Unmarshal(body, &Resp); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}
	fmt.Println(string(body))

	if Resp.AccessToken == "" {
		return nil, fmt.Errorf("empty AccessToken: %s", string(body))
	}
	if Resp.RefreshToken == "" {
		return nil, fmt.Errorf("empty RefreshToken: %s", string(body))
	}

	return &Resp, nil

}

func GetUserInfo[T any](URL, Accept string, accessToken string) (*T, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", Accept)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(": %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read github user body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github user api status=%d, body=%s", resp.StatusCode, string(body))
	}

	var out T
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("unmarshal github user failed: %w", err)
	}

	return &out, nil
}
