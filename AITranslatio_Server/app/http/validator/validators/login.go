package validators

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/CustomErrors"
	"AITranslatio/app/http/Controller/NotAuth"
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/http/validator/comon/data_transfer"
	"github.com/gin-gonic/gin"
)

type LoginDTO struct {
	UserName string `json:"userName"`
	Password string `json:"Password"`
	Auth
}

func (login LoginDTO) CheckParams(context *gin.Context) {

	var LoginDTO LoginDTO

	//1.基础绑定验证
	if err := context.ShouldBind(LoginDTO); err != nil {
		reposen.ErrorParam(context, "基础绑定验证不通过")
	}

	//2.在ctx里添加k=consts.ValidatorPrefix+ json标签对应的 键 => 值 形式绑定在上下文， 再传递给controller
	extraAddBindDataContext := data_transfer.DataAddContext(login, Consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		reposen.ErrorSystem(context, "", CustomErrors.ServerOccurredErrorMsg)
	} else {
		//调用Controller
		(&NotAuth.NotAuthController{}).Login(extraAddBindDataContext)
	}

}
