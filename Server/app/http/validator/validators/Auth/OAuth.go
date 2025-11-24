package Auth

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/http/Controller/AuthController"
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/http/validator/comon/data_transfer"
	"errors"
	"github.com/gin-gonic/gin"
)

type GetChallengeDTO struct{}

func (DTO GetChallengeDTO) CheckParams(WebAuthnContext *gin.Context) {
	//2.使用调用控制器
	extraAddBindDataContext := data_transfer.DataAddContext(DTO, Consts.ValidatorPrefix, WebAuthnContext)
	if extraAddBindDataContext == nil {
		reposen.ErrorSystem(WebAuthnContext, errors.New("DataAddContext无法绑定到*gin.contex"))
	} else {
		provider := extraAddBindDataContext.GetString(Consts.ValidatorPrefix + "OAuth_provider")
		AuthController.CreateOAuthFactroy(provider).GetChallenge(extraAddBindDataContext)
	}
}

type LoginByOAuth struct {
	OAuthProvider string
}

func (DTO LoginByOAuth) CheckParams(WebAuthnContext *gin.Context) {
	//2.使用调用控制器
	extraAddBindDataContext := data_transfer.DataAddContext(DTO, Consts.ValidatorPrefix, WebAuthnContext)
	if extraAddBindDataContext == nil {
		reposen.ErrorSystem(WebAuthnContext, errors.New("DataAddContext无法绑定到*gin.contex"))
	} else {
		provider := extraAddBindDataContext.GetString(Consts.ValidatorPrefix + "OAuth_provider")
		AuthController.CreateOAuthFactroy(provider).Login(extraAddBindDataContext)
	}
}
