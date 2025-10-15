package AuthController

import (
	"AITranslatio/app/Service/AuthService"
	"fmt"

	"AITranslatio/Utils/token"
	"AITranslatio/app/http/reposen"
	"github.com/gin-gonic/gin"
)

func (Controller *AuthController) ApplicationWebAuthn(ctx *gin.Context) {

	//从token解析出UserID
	err, UserID := token.GetDataFormToken[int64](ctx.GetHeader("Authorization"), "UserID")
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("token解析UserID失败:", err))
		return
	}

	//从config获取webAuthn配置项

	WebAuthn, err := AuthService.CreateAuthService().ApplicationWebAuthn(UserID)
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("AuthService创建失败:", err))
		return
	}

	reposen.OK(ctx, WebAuthn)

}

// VerifyWebAuthn 验证WebAuthn休息是否合规
func (Controller *AuthController) VerifyWebAuthn(WebAuthnCtx *gin.Context) {

	//RawID := WebAuthnCtx.GetString(Consts.ValidatorPrefix + "RawID")
	//ID := WebAuthnCtx.GetString(Consts.ValidatorPrefix + "ID")
	//ClientDataJSON := WebAuthnCtx.GetString(Consts.ValidatorPrefix + "clientDataJSON")
	//AttestationObject := WebAuthnCtx.GetString(Consts.ValidatorPrefix + "attestationObject")
	//UserID := WebAuthnCtx.GetString(Consts.ValidatorPrefix + "UserID")

	//校验WebAuthn信息
	err := AuthService.CreateAuthService().VerifyWebAuthn(WebAuthnCtx)
	if err != nil {
		return
	}

}
