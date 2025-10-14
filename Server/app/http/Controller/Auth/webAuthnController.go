package Auth

import (
	"AITranslatio/Global/CustomErrors"
	"AITranslatio/app/Service/AuthService"

	"AITranslatio/Utils/token"
	"AITranslatio/app/http/reposen"
	"github.com/gin-gonic/gin"
)

func (Controller *NotAuthController) ApplicationWebAuthn(WebAuthnCtx *gin.Context) {

	//从token解析出UserID
	err, UserID := token.GetDataFormToken[int64](WebAuthnCtx.GetHeader("Authorization"), "UserID")
	if err != nil {
		reposen.ErrorSystem(WebAuthnCtx, CustomErrors.ErrorGetUserIDIsFail+err.Error())
		return
	}

	//从config获取webAuthn配置项

	WebAuthn, err := AuthService.CreateNotAuthService().ApplicationWebAuthn(UserID)
	if err != nil {
		reposen.ErrorSystem(WebAuthnCtx, err.Error())
		return
	}

	reposen.OK(WebAuthnCtx, WebAuthn)

}

// AddWebAuthn  参数：用户密码，公钥，(公钥签名之后的)挑战验证数据
func (Controller *NotAuthController) AddWebAuthn(WebAuthnCtx *gin.Context) {

	err, _ := token.GetDataFormToken[int64](WebAuthnCtx.GetHeader("Authorization"), "UserID")
	if err != nil {
		reposen.ErrorSystem(WebAuthnCtx, CustomErrors.ErrorGetUserIDIsFail+err.Error())
		return
	}

	//1.验证随机挑战(公钥解密随机挑战之后，验证查看redis有没有这段字符)

	//2.验证密码

	//2.把公钥加入该用户的条目

}
