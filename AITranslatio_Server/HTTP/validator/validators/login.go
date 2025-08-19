package validators

import (
	"AITranslatio/Global"
	"AITranslatio/HTTP/validator/comon/data_transfer"
	"github.com/gin-gonic/gin"
)
import "AITranslatio/HTTP/reposen"

type Login struct{}

func (login *Login) CheckParams(context *gin.Context) {

	//1.基础绑定验证
	if err := context.ShouldBind(login); err != nil {
		reposen.ErrorTokenAuthFail(context)
	}

	//2.使用调用控制器
	extraAddBindDataContext := data_transfer.DataAddContext(login, Global.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		reposen.ErrorTokenAuthFail(context)
	} else {
		
	}

}
