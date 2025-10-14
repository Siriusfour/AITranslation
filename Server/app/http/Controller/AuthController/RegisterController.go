package Auth

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/Service/AuthService"
	"AITranslatio/app/http/DTO/NotAuthDTO"
	"AITranslatio/app/http/reposen"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (Controller *AuthController) Register(ctx *gin.Context) {

	DTO := &NotAuthDTO.RegisterDTO{
		UserID:    0,
		UserName:  ctx.GetString(Consts.ValidatorPrefix + "UserName"),
		Password:  ctx.GetString(Consts.ValidatorPrefix + "Password"),
		Salt:      "",
		Email:     ctx.GetString(Consts.ValidatorPrefix + "Email"),
		EamilCode: "0000",
	}
	Auth, err := AuthService.CreateAuthService().Register(DTO)
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("创建AuthService失败: %w", err))
		return
	}

	reposen.OK(ctx, Auth)

}
