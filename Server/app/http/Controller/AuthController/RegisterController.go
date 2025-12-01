package AuthController

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/Service/AuthService"
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/types"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (Controller *AuthController) Register(ctx *gin.Context) {

	DTO := &types.RegisterDTO{
		UserID:    0,
		UserName:  ctx.GetString(Consts.ValidatorPrefix + "UserName"),
		Password:  ctx.GetString(Consts.ValidatorPrefix + "Password"),
		Salt:      "",
		Email:     ctx.GetString(Consts.ValidatorPrefix + "Email"),
		EamilCode: "0000",
	}
	Auth, err := AuthService.NewAuthService().Register(DTO)
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("创建AuthService失败: %w", err))
		return
	}

	reposen.OK(ctx, Auth)

}
