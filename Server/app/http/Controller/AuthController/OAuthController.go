package AuthController

import (
	"AITranslatio/app/Service/AuthService/OAuthService"
	"AITranslatio/app/http/reposen"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OAuthController interface {
	GetChallenge(*gin.Context)    //生成随机挑战
	VerifyChallenge(*gin.Context) //验证随机挑战

	Login(*gin.Context) //登录
}

func NewOAuthControllerFactroy(server OAuthService.OAuthService, logger *zap.Logger) OAuthController {
	return &GithubController{
		server,
		logger,
	}
}

type GithubController struct {
	GithubService OAuthService.OAuthService
	Logger        *zap.Logger
}

func (c *GithubController) GetChallenge(ctx *gin.Context) {
	challenge, err := c.GithubService.GetChallenge(ctx)
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("生成随机数失败：%w", err))
	}

	reposen.OK(ctx, challenge)
}

func (c *GithubController) VerifyChallenge(ctx *gin.Context) {
	err := c.GithubService.VerifyChallenge(ctx)
	if err != nil {
		reposen.ErrorParam(ctx, err)
	}

	reposen.OK(ctx, struct{}{})
}

func (c *GithubController) Login(ctx *gin.Context) {

	//验证challenge
	err := c.GithubService.VerifyChallenge(ctx)
	if err != nil {
		reposen.Fail(ctx, errors.New("请求过期"))
	}

	//用code换取Github的用户信息
	loginInfo, err := c.GithubService.GetUserInfo(ctx)
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("code换取Github的用户信息错误：%w", err))
	}

	reposen.OK(ctx, loginInfo)

	//返回

}
