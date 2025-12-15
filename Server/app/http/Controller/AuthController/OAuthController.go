package AuthController

import (
	"AITranslatio/app/Service/AuthService/OAuthService"
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/types/DTO"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OAuthController interface {
	GetChallenge(*gin.Context) //生成随机挑战
	//VerifyChallenge(*gin.Context) //验证随机挑战

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

// GetChallenge
// @Summary      OAuth2.0-在服务端获取一个随机数
// @Description  在 OAuth2.0规范中，首先client端要在server获取一个随机数，server会缓存下来，在第二次正式登录client会附带这个随机数在请求url（state），由server验证这个随机数是否与缓存的一致，用于防重放
// @Tags         NotAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}    types.Challenge              "获取成功"
// @Failure      400  {object}    swagger.GetChallengeIsFail    "获取失败"
// @Router       /Auth/GetChallenge [GET]
func (c *GithubController) GetChallenge(ctx *gin.Context) {
	challenge, err := c.GithubService.GetChallenge(ctx)
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("生成随机数失败：%w", err))
	}

	reposen.OK(ctx, challenge)
}

//func (c *GithubController) VerifyChallenge(ctx *gin.Context) {
//	err := c.GithubService.VerifyChallenge(ctx)
//	if err != nil {
//		reposen.ErrorParam(ctx, err)
//	}
//
//	reposen.OK(ctx, struct{}{})
//}

// LoginByOAuth
// @Summary      OAuth2.0-验证并登录
// @Description  客户端在url传递state，code， 服务端验证state，并用code向OAuth提供方换取AcessToken和RefreshToken，并返回用户信息
// @Tags         NotAuth
// @Accept       json
// @Produce      json
// @Param        state      query    string  true  "作用：防重放"
// @Param        code       query    string  true  "作用：向OAuth提供商换取token"
// @Success      200  {object}    types.LoginInfo              "获取成功"
// @Failure      400  {string}    string                       "获取失败"
// @Router       /Auth/LoginByWebAuthn [GET]
func (c *GithubController) Login(ctx *gin.Context) {

	//获取url的code和state
	OAuth := &DTO.OAuth{}
	err := ctx.ShouldBind(&OAuth)
	if err != nil {
		//c.logger.Error("非法参数", zap.String("Email:", OAuth.Email), zap.Error(err))
		reposen.ErrorParam(ctx, errors.New("非法参数"))
		return
	}

	//验证challenge
	err = c.GithubService.VerifyChallenge(ctx, OAuth)
	if err != nil {
		reposen.Fail(ctx, errors.New("请求过期"))
	}

	//用code换取Github的用户信息
	loginInfo, err := c.GithubService.GetUserInfo(ctx, OAuth)
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("code换取Github的用户信息错误：%w", err))
	}

	reposen.OK(ctx, loginInfo)

	//返回

}
