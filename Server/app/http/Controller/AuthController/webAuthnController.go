package AuthController

import (
	"AITranslatio/Utils"
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/types/DTO"
	reponseTypes "AITranslatio/app/types/reponse"
	"fmt"
	"github.com/gin-gonic/gin"
)

// RegisterGetWebAuthnInfo 申请凭证，获取webAunth服务器的信息（RPID,TimeOut,Challenge...）
// @Summary      ApplicationWebAuthn-获取server配置项
// @Description  在注册一个凭证之前，先向服务器获取其WebAuthn的配置如：RPID,Alg ,详见：https://zhuanlan.zhihu.com/p/1966472631474717161
// @Tags         NotAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}    swagger.RegisterWebAuthnInfo             "获取成功"
// @Failure      400  {string}    string                           "获取失败"
// @Router       /Auth/ApplicationWebAuthn [GET]
func (Controller *AuthController) RegisterGetWebAuthnInfo(WebAuthnCtx *gin.Context) {

	UserID := WebAuthnCtx.GetInt64("UserID")
	SessionID := WebAuthnCtx.GetInt64("SessionID")

	//从config获取webAuthn配置项
	RegisterWebAuthn, err := Controller.Service.RegisterGetWebAuthnInfo(SessionID, UserID)
	if err != nil {
		reposen.ErrorSystem(WebAuthnCtx, fmt.Errorf("AuthService创建失败:", err))
		return
	}

	reposen.OK(WebAuthnCtx, RegisterWebAuthn)
	return
}
func (c *AuthController) RegisterWebAuthn(ctx *gin.Context) {

	var RegisterInfo *DTO.RegisterWebAuthn
	err := ctx.ShouldBind(&RegisterInfo)
	if err != nil {
		reposen.ErrorParam(ctx, err)
		return
	}

	//校验clientDataJSON
	err = c.Service.RegisterWebAuthn(ctx, RegisterInfo)
	if err != nil {
		return
	}
	//校验attestationObject
	reposen.OK(ctx, "")
}

// LoginGetWebAuthnInfo
// @Summary      WebAuth登录
// @Description  使用WebAuth安全秘钥登录前需要预先获知服务端信息，如challenge，RPID
// @Tags         NotAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}    swagger.LoginWebAuthnInfo           "获取成功"
// @Failure      400  {string}    string                      "获取失败"
// @Router       /Auth/LoginGetWebAuthnInfo [GET]
func (c *AuthController) LoginGetWebAuthnInfo(ctx *gin.Context) {

	SessionID := ctx.GetInt64("SessionID")

	//获取webAuthn配置项,生成challenge
	LoginWebAuthn, err := c.Service.LoginGetWebAuthnInfo(SessionID)
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("获取webAuthn配置项!%w", err))
		return
	}

	reposen.OK(ctx, LoginWebAuthn)

}

// LoginByWebAuthn
// @Summary      WebAuth登录
// @Description  使用WebAuth安全密钥登录，
// @Tags         NotAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}    types.LoginInfo             "success"
// @Failure      400  {string}    string                      "获取失败"
// @Router       /Auth/WebAuthnByLogin [POST]
func (Controller *AuthController) LoginByWebAuthn(ctx *gin.Context) {

	//绑定数据
	var LoginInfo *DTO.LoginWebAuthn
	err := ctx.ShouldBind(&LoginInfo)
	if err != nil {
		reposen.ErrorParam(ctx, fmt.Errorf("参数错误！", err))
		return
	}

	//从base64URL提取userID
	userID, err := Utils.Base64urlToInt64(LoginInfo.Response.UserHandle)
	if err != nil {
		reposen.ErrorParam(ctx, err)
	}

	sessionID := ctx.GetInt64("SessionID")
	userInfo, AccessToken, RefreshToken, err := Controller.Service.LoginByWebAuthn(LoginInfo.Response.ClientDataJSON, LoginInfo.Response.AuthenticatorData, LoginInfo.Response.Signature, LoginInfo.RawID, sessionID, userID)
	if err != nil {
		reposen.ErrorSystem(ctx, err)
		return
	}

	loginInfo := &reponseTypes.LoginInfo{
		reponseTypes.Auth{
			AccessToken,
			RefreshToken,
		},
		userInfo.Nickname,
		userInfo.UserID,
		userInfo.Avatar,
	}

	//查询登录成功的业务数据
	reposen.OK(ctx, loginInfo)

}
