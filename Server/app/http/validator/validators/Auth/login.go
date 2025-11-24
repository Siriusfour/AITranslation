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

type LoginDTO struct {
	Email    string `json:"Email" binding:"required"`
	Password string `json:"Password" binding:"required"`
}

func (login LoginDTO) CheckParams(context *gin.Context) {

	//1.基础绑定验证
	if err := context.ShouldBindJSON(&login); err != nil {
		reposen.ErrorParam(context, fmt.Errorf("参数基础绑定验证", err))
	}

	//2.在ctx里添加k=consts.ValidatorPrefix+ json标签对应的 键 => 值 形式绑定在上下文， 再传递给controller
	extraAddBindDataContext := data_transfer.DataAddContext(login, Consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		reposen.ErrorSystem(context, errors.New("DataAddContext无法绑定到*gin.contex"))
	} else {
		//调用Controller
		(&AuthController.AuthController{}).Login(extraAddBindDataContext)
	}

}
