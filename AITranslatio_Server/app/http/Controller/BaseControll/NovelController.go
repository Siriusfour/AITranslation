package BaseControll

import (
	"AITranslatio/app/http/Controller/NotAuth"
	"AITranslatio/app/http/DTO"
	"AITranslatio/app/http/reposen"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (BaseController *NotAuth.BaseController) CreateProgramming(Ctx *gin.Context) {

	var NovelDTO DTO.NovelDTO
	var NovelCtx = Ctx

	//1.解析http请求,把参数从HttpMessage.ctx绑定到HttpMessage.DTO
	err := NovelCtx.ShouldBindBodyWithJSON(&NovelDTO)
	if err != nil {
		NotAuth.HTTPErr(NovelCtx, err, 1001)
		return
	}

	err = BaseController.BaseService.CreateNovelProgramming(&NovelDTO)
	if err != nil {
		reposen.Fail(
			NovelCtx,
			reposen.Response{
				Code:    1101,
				Message: fmt.Errorf(" createNote is failed: %w", err).Error(),
			})
	} else {
		reposen.Fail(
			NovelCtx,
			reposen.Response{
				Code:    2000,
				Message: "success",
			})
	}

}

func (BaseController *NotAuth.BaseController) CreateBranch(Ctx *gin.Context) {

	var BranchlDTO DTO.Branch
	var BranchCtx = Ctx

	//1.解析http请求,把参数从HttpMessage.ctx绑定到HttpMessage.DTO
	err := BranchCtx.ShouldBindBodyWithJSON(&BranchlDTO)
	if err != nil {
		reposen.Fail(
			BranchCtx,
			reposen.Response{
				Code:    10111, //数据绑定失败错误码
				Message: fmt.Errorf(" binding data is failed: %w", err).Error(),
			},
		)
		return
	}

	err = BaseController.BaseService.CreateBranch(&BranchlDTO)
	if err != nil {
		reposen.Fail(
			BranchCtx,
			reposen.Response{
				Code:    10112, //数据绑定失败错误码
				Message: fmt.Errorf(" createBranch is failed: %w", err).Error(),
			})
	} else {
		reposen.Fail(
			BranchCtx,
			reposen.Response{
				Code:    2000, //数据绑定失败错误码
				Message: "success",
			})
	}
}

func (BaseController *NotAuth.BaseController) CreateCommit(Ctx *gin.Context) {

	var CommitDTO DTO.CommitDTO
	var CommitCtx = Ctx
	err := CommitCtx.ShouldBind(&CommitDTO)
	if err != nil {
		NotAuth.HTTPErr(Ctx, err, 1001)
		return
	}
	err = BaseController.BaseService.CreateCommit(CommitCtx, &CommitDTO)
}

func (BaseController *NotAuth.BaseController) Programming(Ctx *gin.Context) {

	NoteID, exists := Ctx.GetQuery("NoteID")
	if !exists {
		NotAuth.HTTPErr(Ctx, errors.New("NoteID is required"), 1001)
		return
	}
	NoteIDNum, err := strconv.Atoi(NoteID)
	if err != nil {
		NotAuth.HTTPErr(Ctx, errors.New("Failed to parse the data"), 1001)
		return
	}

	HTTPNote, err := BaseController.BaseService.Programming(NoteIDNum)

	reposen.OK(Ctx, reposen.Response{
		Code:    2000,
		Message: "Success",
		Data:    HTTPNote,
	})
}

func (BaseController *NotAuth.BaseController) ChangeCommit(Ctx *gin.Context) {

	var CommitDTO DTO.CommitDTO
	var CommitCtx = Ctx
	err := CommitCtx.ShouldBind(&CommitDTO)
	if err != nil {
		NotAuth.HTTPErr(Ctx, err, 1001)
		return
	}

	err = BaseController.BaseService.ChangeCommit(Ctx, &CommitDTO)
	if err != nil {
		return
	}

}

func (BaseController *NotAuth.BaseController) GetUserInfo(GetUserInfoCtx *gin.Context) {

	//1.解析http请求,把参数从HttpMessage.ctx绑定到HttpMessage.DTO
	UserID, exists := GetUserInfoCtx.GetQuery("UserID")
	if !exists {
		NotAuth.HTTPErr(GetUserInfoCtx, errors.New("NoteID is required"), 1001)
		return
	}
	_, err := strconv.Atoi(UserID)
	if err != nil {
		NotAuth.HTTPErr(GetUserInfoCtx, errors.New("Failed to parse the data"), 1001)
		return
	}
	//BaseController.BaseService.GetUserInfo(NoteIDNum)

}
