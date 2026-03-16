package ApiController

import (
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/types/DTO"
	reponseTypes "AITranslatio/app/types/reponse"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (c *ApiController) AddSession(ctx *gin.Context) {

	userID := ctx.GetInt64("UserID")

	sessionID, err := c.Service.AddSession(userID)
	if err != nil {
		reposen.ErrorSystem(ctx, err)
		return
	}
	reposen.OK(ctx, sessionID)
}

func (c *ApiController) GetSessionContext(ctx *gin.Context) {

	UserID := ctx.GetInt64("UserID")
	var contextID DTO.RAGContextID

	err := ctx.ShouldBind(&contextID)
	if err != nil {

		reposen.ErrorParam(ctx, err)
		return
	}
	context, err := c.Service.GetSessionContext(strconv.FormatInt(UserID, 10), contextID.ContextID)
	if err != nil {
		reposen.ErrorSystem(ctx, err)
	}
	reposen.OK(ctx, context)
}

func (c *ApiController) Ask(ctx *gin.Context) {

	var Ask DTO.RAGAsk
	UserID := ctx.GetInt64("UserID")
	err := ctx.ShouldBind(&Ask)
	if err != nil {
		reposen.ErrorParam(ctx, err)
		return
	}

	answer, err := c.Service.Ask(Ask.Question, Ask.ContextID, strconv.FormatInt(UserID, 10))
	if err != nil {
		reposen.ErrorSystem(ctx, err)
		return
	}

	res := &reponseTypes.AskReponse{
		answer,
	}
	reposen.OK(ctx, res)

}
