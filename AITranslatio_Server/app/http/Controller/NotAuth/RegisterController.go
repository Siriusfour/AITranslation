package NotAuth

import (
	"AITranslatio/app/Service/NotAuthService"
	"AITranslatio/app/http/reposen"
	"github.com/gin-gonic/gin"
)

func (Controller *NotAuthController) Register(RegisterCtx *gin.Context) {

	UserName := RegisterCtx.GetString("UserName")
	Password := RegisterCtx.GetString("PasswordSecurity")
	Email := RegisterCtx.GetString("Email")

	Auth, err := NotAuthService.CreateNotAuthService().Register(UserName, Email, Password)
	if err != nil {
		reposen.ErrorSystem(RegisterCtx, "", "注册失败："+err.Error())
	}

	reposen.OK(RegisterCtx, Auth, "注册成功！")

}
