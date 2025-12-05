package AuthController

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/token"
	"AITranslatio/app/Service/AuthService"
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/types"
	"AITranslatio/app/types/DTO"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	cfg          interf.ConfigInterface
	logger       *zap.Logger
	jwtGenerator *token.JWTGenerator
	Service      *AuthService.AuthService
	tracer       types.TracerInterf
	OAuthMap     map[string]OAuthController
}

func NewController(cfg interf.ConfigInterface, logger *zap.Logger, Service *AuthService.AuthService, tracer *zipkin.Tracer, oauthMap map[string]OAuthController, jwtGenerator *token.JWTGenerator) *AuthController {
	return &AuthController{
		cfg:          cfg,
		logger:       logger,
		jwtGenerator: jwtGenerator,
		Service:      Service,
		tracer:       tracer,
		OAuthMap:     oauthMap,
	}
}

// Login
// @Summary      用户登录
// @Description  1.Email，password登录  2.Token登录
// @Tags         NotAuth
// @Accept       json
// @Produce      json
// @Param        email            formData  string  false  "邮箱"
// @Param        password         formData  string  false  "密码"
// @Param        AccessToken      header    string  false  "jwt"
// @Param        RefreshToken     header    string  false  "jwt"
// @Success      200  {object}    types.LoginInfo      "登录成功"
// @Failure      400  {object}    swagger.PasswordIsFail    "密码错误"
// @Failure      4001  {object}   swagger.JwtTokenInvalid    "token无效"
// @Failure      4002  {object}   swagger.JwtTokenExpired    "token过期"
// @Router       /Auth/Login [POST]
func (Controller *AuthController) Login(ctx *gin.Context) {

	var LoginDTO DTO.Login
	//校验，绑定参数
	err := ctx.ShouldBind(&LoginDTO)
	if err != nil {
		Controller.logger.Error("非法参数", zap.String("Email:", LoginDTO.Email), zap.Error(err))
		reposen.ErrorParam(ctx, err)
		return
	}
	span := zipkin.SpanFromContext(ctx.Request.Context())
	if span != nil {
		span.Tag("Login", "")
		defer span.Finish()
	}

	//有账号密码优先账号密码，没有就获取token登录
	if (LoginDTO.Email != "" || LoginDTO.Password != "") && ctx.GetHeader("Authorization") == "" {
		Controller.LoginByPassword(ctx, LoginDTO)
	} else {
		Controller.LoginByToken(ctx)
	}

}

func (Controller *AuthController) LoginByPassword(ctx *gin.Context, DTO DTO.Login) {
	//zipkin打点
	childSpan, newCtx := Controller.tracer.StartSpanFromContext(ctx, "GetUserInfo")
	defer childSpan.Finish()
	ctx.Request = ctx.Request.WithContext(newCtx)

	childSpan.Tag("db.system", "MySQL")
	childSpan.Tag("sql:", "FindUser")

	LoginInfo, err := Controller.Service.LoginByPassWord(DTO.Email, DTO.Password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			reposen.Fail(ctx, fmt.Errorf("登录失败！密码错误%w", err))
			return
		}
		reposen.Fail(ctx, fmt.Errorf("登录失败！%w", err))
		return
	}

	reposen.OK(ctx, LoginInfo)
	return

}

func (Controller *AuthController) LoginByToken(ctx *gin.Context) {

	token := ctx.GetHeader("Authorization")
	if token == "" {
		reposen.ErrorTokenAuthFail(ctx, fmt.Errorf("登录消息无效"), Consts.JwtTokenInvalid)
		return
	}

	jwtInfo, err := Controller.jwtGenerator.ParseToken(token)
	if err != nil {
		if errors.Is(err, MyErrors.ErrTokenExpired) {
			reposen.ErrorTokenAuthFail(ctx, fmt.Errorf("登录失败:登录信息已过期%w", err), Consts.JwtTokenExpired)
			return
		}
	}

	UserID := jwtInfo.UserID

	childSpan, newCtx := Controller.tracer.StartSpanFromContext(ctx, "GetUserInfo")
	defer childSpan.Finish()
	ctx.Request = ctx.Request.WithContext(newCtx)

	childSpan.Tag("db.system", "MySQL")
	childSpan.Tag("sql:", "FindUser")
	loginInfo, err := Controller.Service.FindUser(UserID, "UserID")
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("获取用户信息失败：%w", err))
		return
	}

	reposen.OK(ctx, loginInfo)
	return

}
