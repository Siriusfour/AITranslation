package Auth

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/http/Controller/AuthController"
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/http/validator/comon/data_transfer"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ApplicationWebAuthnDTO struct{}

type GetUserAllCredentialDTO struct{}

type LoginByWebAuthnDTO struct {
	RawID    string `json:"rawId" required:"true"`
	ID       string `json:"id" required:"true"`
	Type     string `json:"type" required:"true"`
	Response struct {
		ClientDataJSON    string `json:"clientDataJSON" required:"true"`
		AttestationObject string `json:"authenticatorData" required:"true"`
		Signature         string `json:"signature" required:"true"`
	} `json:"response" required:"true"`
}

type RegisterWebAuthnDTO struct {
	RawID    string `json:"rawId"`
	ID       string `json:"id"`
	Type     string `json:"type"`
	Response struct {
		ClientDataJSON    string `json:"clientDataJSON"`
		AttestationObject string `json:"attestationObject"`
	} `json:"response"`
}

func (DTO ApplicationWebAuthnDTO) CheckParams(WebAuthnContext *gin.Context) {
	//使用调用控制器
	extraAddBindDataContext := data_transfer.DataAddContext(DTO, Consts.ValidatorPrefix, WebAuthnContext)
	if extraAddBindDataContext == nil {
		reposen.ErrorSystem(WebAuthnContext, errors.New("DataAddContext无法绑定到*gin.contex"))
	} else {
		(&AuthController.AuthController{}).ApplicationWebAuthn(extraAddBindDataContext)
	}
}

func (DTO RegisterWebAuthnDTO) CheckParams(WebAuthnContext *gin.Context) {
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
		(&AuthController.AuthController{}).WebAuthnToRegister(extraAddBindDataContext)
	}
}

func (DTO LoginByWebAuthnDTO) CheckParams(WebAuthnContext *gin.Context) {
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
		(&AuthController.AuthController{}).WebAuthnByLogin(extraAddBindDataContext)
	}
}

func (DTO GetUserAllCredentialDTO) CheckParams(WebAuthnContext *gin.Context) {

	//2.使用调用控制器
	extraAddBindDataContext := data_transfer.DataAddContext(DTO, Consts.ValidatorPrefix, WebAuthnContext)
	if extraAddBindDataContext == nil {
		reposen.ErrorSystem(WebAuthnContext, errors.New("DataAddContext无法绑定到*gin.contex"))
	} else {
		(&AuthController.AuthController{}).GetUserAllCredential(extraAddBindDataContext)
	}
}
