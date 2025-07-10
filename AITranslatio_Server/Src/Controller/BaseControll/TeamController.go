package BaseControll

import (
	"AITranslatio/Src/DTO"
	"github.com/gin-gonic/gin"
)

func (BaseController *BaseController) JoinTeam(Ctx *gin.Context) {

	var JoinTeamDTO DTO.JoinTeamDTO
	var JoinTeamCtx = Ctx

	//1.解析http请求,把参数从HttpMessage.ctx绑定到HttpMessage.DTO
	err := JoinTeamCtx.ShouldBindBodyWithJSON(&JoinTeamDTO)
	if err != nil {
		BindingErr(JoinTeamCtx, err, 1001)
		return
	}

	err = BaseController.BaseService.JoinTeam(JoinTeamDTO)

}
