package validators

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/http/Controller/Auth"
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/http/validator/comon/data_transfer"
	"github.com/gin-gonic/gin"
)

type WebAuthnDTO struct {
	Challenge string `form:"Challenge" json:"Challenge"`
}

func (DTO WebAuthnDTO) CheckParams(WebAuthnContext *gin.Context) {

	//1.基础绑定验证
	if err := WebAuthnContext.ShouldBindJSON(&DTO); err != nil {
		reposen.ErrorTokenAuthFail(WebAuthnContext)
	}

	//2.使用调用控制器
	extraAddBindDataContext := data_transfer.DataAddContext(DTO, Consts.ValidatorPrefix, WebAuthnContext)
	if extraAddBindDataContext == nil {
		reposen.ErrorTokenAuthFail(WebAuthnContext)
	} else {
		(&Auth.NotAuthController{}).ApplicationWebAuthn(extraAddBindDataContext)
	}

}
