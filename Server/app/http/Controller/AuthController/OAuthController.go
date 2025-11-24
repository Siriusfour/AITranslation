package AuthController

import (
	"AITranslatio/app/Service/AuthService/OAuthService"
	"AITranslatio/app/http/reposen"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type OAuthController interface {
	GetChallenge(*gin.Context)    //生成随机挑战
	VerifyChallenge(*gin.Context) //验证随机挑战

	Login(*gin.Context) //登录

}

func CreateOAuthFactroy(server string) OAuthController {
	if server == "Github" {
		return &Github{}
	}
	return nil
}

type Github struct{}

func (g *Github) GetChallenge(ctx *gin.Context) {
	challenge, err := OAuthService.CreateOAuthFactroy("Github").GetChallenge(ctx)
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("生成随机数失败：%w", err))
	}

	reposen.OK(ctx, challenge)

}

func (g *Github) VerifyChallenge(ctx *gin.Context) {
	err := OAuthService.CreateOAuthFactroy("Github").VerifyChallenge(ctx)
	if err != nil {
		reposen.ErrorParam(ctx, err)
	}

	reposen.OK(ctx, struct{}{})
}

func (g *Github) Login(ctx *gin.Context) {

	//验证challenge
	err := OAuthService.CreateOAuthFactroy("Github").VerifyChallenge(ctx)
	if err != nil {
		reposen.Fail(ctx, errors.New("请求过期"))
	}

	//用code换取Github的AccessToken和RefreshToken

	//生成自己的token

	//返回

}
