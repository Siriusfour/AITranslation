package validators

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/http/Controller/AuthController"
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/http/validator/comon/data_transfer"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type WebAuthnDTO struct{}

type VerifyWebAuthnDTO struct {
	RawID             string `json:"rawId"`
	ID                string `json:"id"`
	Type              string `json:"type"`
	ClientDataJSON    string `json:"clientDataJSON"`    // Base64URL
	AttestationObject string `json:"attestationObject"` // Base64URL

}

func (DTO WebAuthnDTO) CheckParams(WebAuthnContext *gin.Context) {

	//1.基础绑定验证
	//if err := WebAuthnContext.ShouldBindJSON(&DTO); err != nil {
	//	reposen.ErrorParam(WebAuthnContext, fmt.Errorf("参数基础校验失败%w", err))
	//	return
	//}

	//2.使用调用控制器
	extraAddBindDataContext := data_transfer.DataAddContext(DTO, Consts.ValidatorPrefix, WebAuthnContext)
	if extraAddBindDataContext == nil {
		reposen.ErrorSystem(WebAuthnContext, errors.New("DataAddContext无法绑定到*gin.contex"))
	} else {
		(&AuthController.AuthController{}).ApplicationWebAuthn(extraAddBindDataContext)
	}
}

func (DTO VerifyWebAuthnDTO) CheckParams(WebAuthnContext *gin.Context) {
	//1.基础绑定验证
	if err := WebAuthnContext.ShouldBindJSON(&DTO); err != nil {
		reposen.ErrorParam(WebAuthnContext, fmt.Errorf("参数基础校验失败%w", err))
		return
	}

	//2.使用调用控制器
	extraAddBindDataContext := data_transfer.DataAddContext(DTO, Consts.ValidatorPrefix, WebAuthnContext)
	if extraAddBindDataContext == nil {
		reposen.ErrorSystem(WebAuthnContext, errors.New("DataAddContext无法绑定到*gin.contex"))
	} else {
		(&AuthController.AuthController{}).VerifyWebAuthn(extraAddBindDataContext)
	}
}
