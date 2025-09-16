package NotAuth

import (
	"AITranslatio/Utils/WebAuthn"
	"AITranslatio/app/http/DTO/NotAuthDTO"
	"AITranslatio/app/http/reposen"
	"github.com/gin-gonic/gin"
)

func (Controller *NotAuthController) ApplicationWebAuthn(RegisterCtx *gin.Context) {

	var WebAuthnDTO *NotAuthDTO.WebAuthnRegisterDTO

	WebAuthnDTO.UserID = RegisterCtx.GetInt64("UserID")
	WebAuthnDTO.Email = RegisterCtx.GetString("Email")

	err := WebAuthn.ApplicationWebAuthn(WebAuthnDTO)
	if err != nil {
		reposen.ErrorSystem(RegisterCtx, "", err.Error())
	}

	//生成随机挑战
	challenge, err := WebAuthn.CreateChallenge()
	if err != nil {
		reposen.ErrorSystem(RegisterCtx, "", err.Error())
	}
	reposen.OK(RegisterCtx, challenge, "")

}

func (Controller *NotAuthController) GetChallenge(RegisterCtx *gin.Context) {

	challenge, err := WebAuthn.CreateChallenge()
	if err != nil {
		return
	}

	reposen.OK(RegisterCtx, challenge, "")
}

func (Controller *NotAuthController) VerifyWebAuthn(RegisterCtx *gin.Context) {

}
