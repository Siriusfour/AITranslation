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

type LoginByOAuthDTO struct {
	OAuthProvider string
	Code          string `json:"code" mapstructure:"code"`
	State         string `json:"state" mapstructure:"state"`
}

func (DTO LoginByOAuthDTO) CheckParams(ctx *gin.Context) {

	if err := ctx.ShouldBindJSON(&DTO); err != nil {
		reposen.ErrorParam(ctx, fmt.Errorf("参数基础绑定验证错误：%w", err))
	}

	//2.使用调用控制器
	extraAddBindDataContext := data_transfer.DataAddContext(DTO, Consts.ValidatorPrefix, ctx)
	if extraAddBindDataContext == nil {
		reposen.ErrorSystem(ctx, errors.New("DataAddContext无法绑定到*gin.contex"))
	} else {
		provider := extraAddBindDataContext.GetString(Consts.ValidatorPrefix + "OAuthProvider")
		AuthController.CreateOAuthFactroy(provider).Login(extraAddBindDataContext)
	}
}
