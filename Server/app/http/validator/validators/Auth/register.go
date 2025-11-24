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

type RegisterDTO struct {
	UserName  string `json:"UserName"`
	Password  string `json:"Password" `
	Salt      string
	Email     string `json:"Email"`
	EamilCode string
	Timestamp int64
}

func (DTO RegisterDTO) CheckParams(RegisterContext *gin.Context) {

	//1.基础绑定验证
	if err := RegisterContext.ShouldBindJSON(&DTO); err != nil {
		reposen.ErrorParam(RegisterContext, fmt.Errorf("参数基础绑定验证", err))
	}

	//2.使用调用控制器
	extraAddBindDataContext := data_transfer.DataAddContext(DTO, Consts.ValidatorPrefix, RegisterContext)
	if extraAddBindDataContext == nil {
		reposen.ErrorSystem(RegisterContext, errors.New("DataAddContext无法绑定到*gin.contex"))
	} else {
		(&AuthController.AuthController{}).Register(extraAddBindDataContext)
	}

}
