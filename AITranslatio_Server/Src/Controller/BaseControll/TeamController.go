package BaseControll

import (
	"AITranslatio/Src/DTO"
	"errors"
	"github.com/gin-gonic/gin"
)

// CreateTeam 创建小组
func (BaseController *BaseController) CreateTeam(CreateTeamCtx *gin.Context) {
	var CreateTeamDTO DTO.CreateTeamDTO

	//1.解析http请求,数据绑定到DTO
	err := CreateTeamCtx.ShouldBindBodyWithJSON(&CreateTeamDTO)
	if err != nil {
		HTTPErr(CreateTeamCtx, err, 1001)
		return
	}

	err = BaseController.BaseService.CreateTeam(&CreateTeamDTO)

}

// JoinTeam 申请加入小组
func (BaseController *BaseController) JoinTeam(Ctx *gin.Context) {

	var JoinTeamDTO DTO.JoinTeamDTO
	var JoinTeamCtx = Ctx

	//1.解析http请求,数据绑定到DTO
	err := JoinTeamCtx.ShouldBindBodyWithJSON(&JoinTeamDTO)
	if err != nil {
		HTTPErr(JoinTeamCtx, err, 1001)
		return
	}

	err = BaseController.BaseService.JoinTeam(&JoinTeamDTO)
	if err != nil {
		HTTPErr(Ctx, errors.New("发起请求失败:"+err.Error()), 1014)
	} else {
		HTTPSuccess(Ctx, struct{}{}, "发起请求成功，请等待处理")
	}

}
