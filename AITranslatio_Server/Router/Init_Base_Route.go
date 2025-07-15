package Router

import (
	"AITranslatio/Src/Controller/BaseControll"
	"github.com/gin-gonic/gin"
)

func Init_Base_Route(rgBase *gin.RouterGroup, BaseController *BaseControll.BaseController) {

	//rgBase.POST("/API/Translation", BaseController.Translation)
	rgBase.POST("API/Login", BaseController.Login)
	rgBase.POST("API/CreateNoteProgramming", BaseController.CreateProgramming)
	rgBase.POST("API/CreateBranch", BaseController.CreateBranch)
	rgBase.POST("API/CreateCommit", BaseController.CreateCommit)
	rgBase.GET("API/Programming", BaseController.Programming)
	rgBase.GET("API/ChangeCommit", BaseController.ChangeCommit)
	//rgBase.GET("API/UserInfo", BaseController.)
	//rgBase.POST("API/CreateTeam", BaseController.CreateTeam)
	rgBase.POST("API/JoinTeam", BaseController.JoinTeam) //申请加入一个团队
	rgBase.GET("SSE", BaseController.CreateSSE)
}
