package validators

import (
	"AITranslatio/Global"
	"AITranslatio/HTTP/reposen"
	"AITranslatio/HTTP/validator/comon/data_transfer"
	"github.com/gin-gonic/gin"
)

type Register struct{}

func (register *Register) CheckParams(context *gin.Context) {

	//1.基础绑定验证
	if err := context.ShouldBind(register); err != nil {
		reposen.ErrorTokenAuthFail(context)
	}

	//2.使用调用控制器
	extraAddBindDataContext := data_transfer.DataAddContext(register, Global.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		reposen.ErrorTokenAuthFail(context)
	} else {

	}

}
