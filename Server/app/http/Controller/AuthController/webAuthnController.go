package AuthController

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/Service/AuthService"
	"AITranslatio/app/http/reposen"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 获取webAunth服务器的信息（RPI,TimeOut,Challenge...）
func (Controller *AuthController) ApplicationWebAuthn(WebAuthnCtx *gin.Context) {

	//从token解析出UserID
	UserID := WebAuthnCtx.GetInt64(Consts.ValidatorPrefix + "UserID")

	//从config获取webAuthn配置项
	WebAuthn, err := AuthService.CreateAuthService().ApplicationWebAuthn(UserID)
	if err != nil {
		reposen.ErrorSystem(WebAuthnCtx, fmt.Errorf("AuthService创建失败:", err))
		return
	}

	reposen.OK(WebAuthnCtx, WebAuthn)
	return
}

// WebAuthnToRegister 验证WebAuthn休息是否合规
func (Controller *AuthController) WebAuthnToRegister(WebAuthnCtx *gin.Context) {

	//校验WebAuthn信息
	err := AuthService.CreateAuthService().VerifyWebAuthnToRegister(WebAuthnCtx)
	if err != nil {
		return
	}

	reposen.OK(WebAuthnCtx, nil)
	return

}

func (Controller *AuthController) GetUserAllCredential(WebAuthnCtx *gin.Context) {

	//生成随机挑战,置于redis分钟
	UserID := WebAuthnCtx.GetInt64(Consts.ValidatorPrefix + "UserID")
	w, err := AuthService.CreateAuthService().ApplicationWebAuthn(UserID)
	if err != nil {
		reposen.ErrorSystem(WebAuthnCtx, err)
		return
	}
	WebAuthnCtx.Set(Consts.ValidatorPrefix+"challenge", w.Challenge)

	//获取该用户所有的凭证
	data, err := AuthService.CreateAuthService().GetUserAllCredentialDTO(WebAuthnCtx)
	if err != nil {
		return
	}
	reposen.OK(WebAuthnCtx, data)
	return
}

func (Controller *AuthController) WebAuthnByLogin(WebAuthnCtx *gin.Context) {

	err := AuthService.CreateAuthService().WebAuthnToLogin(WebAuthnCtx)
	if err != nil {
		reposen.ErrorSystem(WebAuthnCtx, err)
		return
	}

	//查询登录成功的业务数据

	reposen.OK(WebAuthnCtx, nil)

}
