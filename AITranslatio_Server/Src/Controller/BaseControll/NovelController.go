package BaseControll

import (
	"AITranslatio/Src/DTO"
	"AITranslatio/Src/HTTP"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
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

	err = BaseController.BaseService.CreateNovelProgramming(&NovelDTO)
	if err != nil {
		HTTP.Fail(
			NovelCtx,
			HTTP.Response{
				Code:    10112, //数据绑定失败错误码
				Message: fmt.Errorf(" createNote is failed: %w", err).Error(),
			})
	} else {
		HTTP.Fail(
			NovelCtx,
			HTTP.Response{
				Code:    2000, //数据绑定失败错误码
				Message: "success",
			})
	}

}

func (BaseController *BaseController) CreateBranch(Ctx *gin.Context) {

	var BranchlDTO DTO.Branch
	var BranchCtx = Ctx

	//1.解析http请求,把参数从HttpMessage.ctx绑定到HttpMessage.DTO
	err := BranchCtx.ShouldBindBodyWithJSON(&BranchlDTO)
	if err != nil {
		HTTP.Fail(
			BranchCtx,
			HTTP.Response{
				Code:    10111, //数据绑定失败错误码
				Message: fmt.Errorf(" binding data is failed: %w", err).Error(),
			},
		)
		return
	}

	err = BaseController.BaseService.CreateBranch(&BranchlDTO)
	if err != nil {
		HTTP.Fail(
			BranchCtx,
			HTTP.Response{
				Code:    10112, //数据绑定失败错误码
				Message: fmt.Errorf(" createBranch is failed: %w", err).Error(),
			})
	} else {
		HTTP.Fail(
			BranchCtx,
			HTTP.Response{
				Code:    2000, //数据绑定失败错误码
				Message: "success",
			})
	}
}

func (BaseController *BaseController) CreateCommit(Ctx *gin.Context) {

	var CommitDTO DTO.CommitDTO
	var CommitCtx = Ctx
	err := CommitCtx.ShouldBind(&CommitDTO)
	if err != nil {
		HTTP.Fail(
			CommitCtx,
			HTTP.Response{
				Code:    10111, //数据绑定失败错误码
				Message: fmt.Errorf(" binding data is failed: %w", err).Error(),
			},
		)
		return
	}

	err = BaseController.BaseService.CreateCommit(CommitCtx, &CommitDTO)

}

func (BaseController *BaseController) Programming(Ctx *gin.Context) {

	NoteID, exists := Ctx.GetQuery("NoteID")
	if !exists {
		bindingErr(Ctx, errors.New("gender is required"))
		return
	}
	NoteIDNum, err := strconv.Atoi(NoteID)
	if err != nil {
		bindingErr(Ctx, errors.New("Failed to parse the data"))
		return
	}

	HTTPNote, err := BaseController.BaseService.Programming(NoteIDNum)

	HTTP.OK(Ctx, HTTP.Response{
		Code:    2000,
		Message: "Success",
		Data:    HTTPNote,
	})
}
