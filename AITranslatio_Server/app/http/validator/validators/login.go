package validators

import (
	"AITranslatio/Global"
	"AITranslatio/app/http/Controller/NotAuth"
	"AITranslatio/app/http/DTO/NotAuthDTO"
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/http/validator/comon/data_transfer"
	"github.com/gin-gonic/gin"
)

type Login struct {
	BaseField
}

func (login *Login) CheckParams(context *gin.Context) {

	var LoginDTO NotAuthDTO.LoginDTO

	//1.基础绑定验证
	if err := context.ShouldBind(LoginDTO); err != nil {
		reposen.ErrorParam(context, "基础绑定验证不通过")
	}

	//2.在ctx里添加k=consts.ValidatorPrefix+ json标签对应的 键 => 值 形式绑定在上下文， 再传递给controller
	extraAddBindDataContext := data_transfer.DataAddContext(login, Global.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		reposen.ErrorSystem(context, "", Global.ServerOccurredErrorMsg)
	} else {

		//调用Controller
		(&NotAuth.NotAuthController{}).Login(extraAddBindDataContext)
	}

}
