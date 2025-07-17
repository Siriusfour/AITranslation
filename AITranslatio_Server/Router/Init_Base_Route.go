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
	rgBase.GET("API/GetProgramming", BaseController.Programming)
	rgBase.GET("API/ChangeCommit", BaseController.ChangeCommit)
	//rgBase.GET("API/UserInfo", BaseController.)
	rgBase.POST("API/CreateTeam", BaseController.CreateTeam) //创建一个团队
	rgBase.POST("API/JoinTeam", BaseController.JoinTeam)     //申请加入一个团队 ok
	rgBase.GET("API/SSE", BaseController.CreateSSE)          //创建SSE链接，ok
}
