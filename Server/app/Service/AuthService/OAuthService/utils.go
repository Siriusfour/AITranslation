package OAuthService

import (
	"AITranslatio/app/types"
	"github.com/gin-gonic/gin"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type OAuthService interface {
	GetChallenge(*gin.Context) (string, error)     //生成随机挑战
	VerifyChallenge(*gin.Context) error            //验证随机挑战
	GetUserInfo(*gin.Context) (*types.Auth, error) //从OAuth提供方获取token
}

func CreateOAuthFactroy(server string) OAuthService {
	if server == "Github" {
		return &Github{}
	}
	if server == "WX" {
	}
	if server == "QQ" {
	}
	return nil
}
