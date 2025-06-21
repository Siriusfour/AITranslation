package Base

import (
	"AITranslatio/Global"
	"AITranslatio/Src/DTO"
	"AITranslatio/Src/HTTP"
	"AITranslatio/Src/Service/BaseService.go"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yankeguo/zhipu"
	"go.uber.org/zap"
)

type BaseController struct {
	Ctx         *gin.Context
	Logger      *zap.SugaredLogger
	BaseService *BaseService.BaseService
}

func NewBaseController() *BaseController {
	return &BaseController{
		Logger:      Global.Logger,
		BaseService: BaseService.NewBaseService(),
	}
}

func (BaseController *BaseController) Login(Ctx *gin.Context) {

	var LoginDTO DTO.LoginDTO
	var LoginCtx = Ctx

	//1.解析http请求,把参数从HttpMessage.ctx绑定到HttpMessage.DTO
	err := LoginCtx.ShouldBindBodyWithJSON(&LoginDTO)
	if err != nil {
		HTTP.Fail(
			LoginCtx,
			HTTP.Response{
				Code:    10111, //数据绑定失败错误码
				Message: fmt.Errorf(" binding data is failed: %w", err).Error(),
			},
		)
		return
	}

	err = BaseController.BaseService.Login(LoginCtx, &LoginDTO)

}

func (BaseController *BaseController) Translation(ctx *gin.Context) {
	//0.实例化接受ctx参数的DTO结构体
	var AITranslation DTO.TranslationDTO

	HttpMessage := &HTTP.Request{
		Ctx: ctx,
		DTO: &AITranslation,
	}

	//1.解析http请求,把参数ctx
	err := HttpMessage.Ctx.ShouldBindBodyWithJSON(&HttpMessage.DTO)
	if err != nil {
		HTTP.Fail(
			HttpMessage.Ctx,
			HTTP.Response{
				Code:    10111, //数据绑定失败错误码
				Message: fmt.Errorf("open config failed: %w", err).Error(),
			},
		)
		return
	}

	client, err := zhipu.NewClient(zhipu.WithAPIKey("14cc2eb752714fba9b55a681793edfde.m0yOWQwvz8psxqZK"))
	service := client.ChatCompletion("glm-4-flash").
		AddMessage(zhipu.ChatCompletionMessage{
			Role:    "user",
			Content: HttpMessage.DTO.(*DTO.TranslationDTO).Message,
		})

	res, err := service.Do(context.Background())

	if err != nil {
		zhipu.GetAPIErrorCode(err) // get the API error code
	} else {
		println(res.Choices[0].Message.Content)
	}
}
