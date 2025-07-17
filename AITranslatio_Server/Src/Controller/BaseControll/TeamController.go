package BaseControll

import (
	"AITranslatio/Src/DTO"
	"github.com/gin-gonic/gin"
)

func (BaseController *BaseController) CreateTeam(CreateTeamCtx *gin.Context) {
	var CreateTeamDTO DTO.CreateTeamDTO

	//1.解析http请求,数据绑定到DTO
	err := CreateTeamCtx.ShouldBindBodyWithJSON(&CreateTeamDTO)
	if err != nil {
		BindingErr(CreateTeamCtx, err, 1001)
		return
	}

	err = BaseController.BaseService.CreateTeam(&CreateTeamDTO)

}

func (BaseController *BaseController) JoinTeam(Ctx *gin.Context) {

	var JoinTeamDTO DTO.JoinTeamDTO
	var JoinTeamCtx = Ctx

	//1.解析http请求,数据绑定到DTO
	err := JoinTeamCtx.ShouldBindBodyWithJSON(&JoinTeamDTO)
	if err != nil {
		BindingErr(JoinTeamCtx, err, 1001)
		return
	}

	err = BaseController.BaseService.JoinTeam(&JoinTeamDTO)

}
