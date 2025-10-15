package Team

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/http/Controller/ApiController"
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/http/validator/comon/data_transfer"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type JoinDTO struct {
	TeamID       int    `json:"TeamID" binding:"required"`
	Introduction string `json:"Introduction" binding:"required"`
}

func (DTO JoinDTO) CheckParams(context *gin.Context) {

	//1.基础绑定验证
	if err := context.ShouldBindJSON(&DTO); err != nil {
		reposen.ErrorParam(context, fmt.Errorf("参数基础绑定验证", err))
		return
	}

	//2.在ctx里添加k=consts.ValidatorPrefix+ json标签对应的 键 => 值 形式绑定在上下文， 再传递给controller
	extraAddBindDataContext := data_transfer.DataAddContext(DTO, Consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		reposen.ErrorSystem(context, errors.New("DataAddContext无法绑定到*gin.contex"))
	} else {
		//调用Controller
		(&ApiController.ApiController{}).JoinTeam(extraAddBindDataContext)
	}

}
