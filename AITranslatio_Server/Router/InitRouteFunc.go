package Router

import (
	"AITranslatio/app/http/Controller/NotAuth"
	"github.com/gin-gonic/gin"
)

func InitBaseRoute(rgBase *gin.RouterGroup, BaseController *NotAuth.BaseController) {
	rgBase.POST("API/Login", BaseController.Login)
	rgBase.POST("API/CreateNoteProgramming", BaseController.CreateProgramming)
	rgBase.POST("API/CreateBranch", BaseController.CreateBranch)
	rgBase.POST("API/CreateCommit", BaseController.CreateCommit)
	rgBase.GET("API/GetProgramming", BaseController.Programming)
	rgBase.GET("API/ChangeCommit", BaseController.ChangeCommit)
	rgBase.POST("API/CreateTeam", BaseController.CreateTeam) //创建一个团队
	rgBase.POST("API/JoinTeam", BaseController.JoinTeam)     //申请加入一个团队
	rgBase.GET("API/SSE", BaseController.CreateSSE)          //创建SSE链接，
}

func InitAuthRoute(rgBase *gin.RouterGroup) {}

func InitNotAuthRoute(rgBase *gin.RouterGroup) {}

func InitFilesRoute(rgBase *gin.RouterGroup) {}

func InitCaptchaRoute(rgBase *gin.RouterGroup) {}
