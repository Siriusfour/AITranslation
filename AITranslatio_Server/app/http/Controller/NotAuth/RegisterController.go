package NotAuth

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/CustomErrors"
	"AITranslatio/app/Service/NotAuthService"
	"AITranslatio/app/http/DTO/NotAuthDTO"
	"AITranslatio/app/http/reposen"
	"github.com/gin-gonic/gin"
)

func (Controller *NotAuthController) Register(RegisterCtx *gin.Context) {

	DTO := &NotAuthDTO.RegisterDTO{
		UserID:    0,
		UserName:  RegisterCtx.GetString(Consts.ValidatorPrefix + "UserName"),
		Password:  RegisterCtx.GetString(Consts.ValidatorPrefix + "Password"),
		Salt:      "123",
		Email:     RegisterCtx.GetString(Consts.ValidatorPrefix + "Email"),
		EamilCode: "0000",
		PublicKey: "",
	}
	Auth, err := NotAuthService.CreateNotAuthService().Register(DTO)
	if err != nil {
		reposen.ErrorSystem(RegisterCtx, "", CustomErrors.ErrorRegisterIsFail+err.Error())
		return
	}

	reposen.OK(RegisterCtx, Auth, "注册成功！")

}
