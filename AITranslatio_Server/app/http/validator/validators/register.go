package validators

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/http/Controller/NotAuth"
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/http/validator/comon/data_transfer"
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
		reposen.ErrorTokenAuthFail(RegisterContext)
	}

	//2.使用调用控制器
	extraAddBindDataContext := data_transfer.DataAddContext(DTO, Consts.ValidatorPrefix, RegisterContext)
	if extraAddBindDataContext == nil {
		reposen.ErrorTokenAuthFail(RegisterContext)
	} else {
		(&NotAuth.NotAuthController{}).Register(extraAddBindDataContext)
	}

}
