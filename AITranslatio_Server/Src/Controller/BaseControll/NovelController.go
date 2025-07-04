package BaseControll

import (
	"AITranslatio/Src/DTO"
	"AITranslatio/Src/HTTP"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (BaseController *BaseController) CreateProgramming(Ctx *gin.Context) {

	var NovelDTO DTO.NovelDTO
	var NovelCtx = Ctx

	//1.解析http请求,把参数从HttpMessage.ctx绑定到HttpMessage.DTO
	err := NovelCtx.ShouldBindBodyWithJSON(&NovelDTO)
	if err != nil {
		HTTP.Fail(
			NovelCtx,
			HTTP.Response{
				Code:    10111, //数据绑定失败错误码
				Message: fmt.Errorf(" binding data is failed: %w", err).Error(),
			},
		)
		return
	}

	err = BaseController.BaseService.CreateNovelProgramming(NovelCtx, &NovelDTO)

}
