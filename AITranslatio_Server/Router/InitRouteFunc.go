package Router

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/http/validator/comon/factory"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitNotAuthRoute(rg *gin.RouterGroup) {
	rg.POST("/Login", factory.Create(Consts.ValidatorPrefix+"Login"))
	rg.POST("/Register", factory.Create(Consts.ValidatorPrefix+"Register"))
	//rg.POST("/WebAuthn", factory.Create(Consts.ValidatorPrefix+"login"))
	rg.POST("/", func(context *gin.Context) {
		context.String(http.StatusOK, "HelloWorld,这是后端模块")
	})
}

func InitAuthRoute(rgBase *gin.RouterGroup) {}

func InitFilesRoute(rgBase *gin.RouterGroup) {}

func InitCaptchaRoute(rgBase *gin.RouterGroup) {}

//rg.POST("/CreateNoteProgramming", Controller.CreateProgramming)
//rg.POST("/CreateBranch", Controller.CreateBranch)
//rg.POST("/CreateCommit", Controller.CreateCommit)
//rg.GET("/GetProgramming", Controller.Programming)
//rg.GET("/ChangeCommit", Controller.ChangeCommit)
//rg.POST("/CreateTeam", Controller.CreateTeam) //创建一个团队
//rg.POST("/JoinTeam", Controller.JoinTeam)     //申请加入一个团队
//rg.GET("/SSE", Controller.CreateSSE)
