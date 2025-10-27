package Router

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/http/validator/comon/factory"
	"github.com/gin-gonic/gin"
)

func InitAuthRoute(rg gin.IRoutes) {
	rg.POST("/Login", factory.Create(Consts.ValidatorPrefix+"Login")) //在全局容器里面找到Login验证器
	rg.POST("/Register", factory.Create(Consts.ValidatorPrefix+"Register"))
	rg.GET("/ApplicationWebAuthn", factory.Create(Consts.ValidatorPrefix+"ApplicationWebAuthn"))
	rg.POST("/RegisterWebAuthn", factory.Create(Consts.ValidatorPrefix+"RegisterWebAuthn"))

	rg.GET("/GetUserAllCredential", factory.Create(Consts.ValidatorPrefix+"GetUserAllCredential"))
	rg.POST("/LoginByWebAuthn", factory.Create(Consts.ValidatorPrefix+"LoginByWebAuthn"))
}

func InitBaseRoute(rg gin.IRoutes) {
	rg.POST("/Team/Create", factory.Create(Consts.ValidatorPrefix+"TeamCreate"))
	rg.POST("/Team/Join", factory.Create(Consts.ValidatorPrefix+"Join"))
}

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
