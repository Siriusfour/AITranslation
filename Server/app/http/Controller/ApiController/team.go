package ApiController

import (
	"AITranslatio/app/Service/ApiServer"
	"AITranslatio/app/http/reposen"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (Controller *ApiController) CreateTeam(ctx *gin.Context) {

	err := ApiServer.CreateApiServer().CreateTeam(ctx)
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("创建新群聊失败: %w", err))
		return
	}

	reposen.OK(ctx, gin.H{})

}

func (Controller *ApiController) JoinTeam(ctx *gin.Context) {

	//1.消息入库
	err := ApiServer.CreateApiServer().JoinTeam(ctx)
	if err != nil {
		reposen.ErrorSystem(ctx, fmt.Errorf("申请失败!: %w", err))
	}

	reposen.OK(ctx, "Success")

}

func (Controller *ApiController) DeleteTeam(ctx *gin.Context) {}
