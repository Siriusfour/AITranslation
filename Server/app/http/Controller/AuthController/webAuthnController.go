package AuthController

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/http/reposen"

	"fmt"
	"github.com/gin-gonic/gin"
)

// 申请凭证，获取webAunth服务器的信息（RPID,TimeOut,Challenge...）
// ApplicationWebAuthn
// @Summary      ApplicationWebAuthn-获取server配置项
// @Description  在注册一个凭证之前，先向服务器获取其WebAuthn的配置如：RPID,Alg ,详见：https://zhuanlan.zhihu.com/p/1966472631474717161
// @Tags         NotAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}    swagger.WebAuthnInfo             "获取成功"
// @Failure      400  {string}    string                           "获取失败"
// @Router       /Auth/ApplicationWebAuthn [GET]
func (Controller *AuthController) ApplicationWebAuthn(WebAuthnCtx *gin.Context) {

	//从token解析出UserID
	UserID := WebAuthnCtx.GetInt64("UserID")

	//从config获取webAuthn配置项
	WebAuthn, err := Controller.Service.ApplicationWebAuthn(UserID)
	if err != nil {
		reposen.ErrorSystem(WebAuthnCtx, fmt.Errorf("AuthService创建失败:", err))
		return
	}

	reposen.OK(WebAuthnCtx, WebAuthn)
	return
}

func (Controller *AuthController) GetUserAllCredential(WebAuthnCtx *gin.Context) {

	//生成随机挑战,置于redis分钟
	UserID := WebAuthnCtx.GetInt64(Consts.ValidatorPrefix + "UserID")
	w, err := Controller.Service.ApplicationWebAuthn(UserID)
	if err != nil {
		reposen.ErrorSystem(WebAuthnCtx, err)
		return
	}
	WebAuthnCtx.Set(Consts.ValidatorPrefix+"challenge", w.Challenge)

	//获取该用户所有的凭证
	data, err := Controller.Service.GetUserAllCredentialDTO(WebAuthnCtx)
	if err != nil {
		return
	}
	reposen.OK(WebAuthnCtx, data)
	return
}

// WebAuthnByLogin
// @Summary      WebAuth登录
// @Description  使用WebAuth安全密钥登录，详见：https://zhuanlan.zhihu.com/p/1966472631474717161
// @Tags         NotAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}    types.LoginInfo             "获取成功"
// @Failure      400  {string}    string                           "获取失败"
// @Router       /Auth/WebAuthnByLogin [POST]
func (Controller *AuthController) WebAuthnByLogin(WebAuthnCtx *gin.Context) {

	err := Controller.Service.WebAuthnToLogin(WebAuthnCtx)
	if err != nil {
		reposen.ErrorSystem(WebAuthnCtx, err)
		return
	}

	//查询登录成功的业务数据

	reposen.OK(WebAuthnCtx, nil)

}
