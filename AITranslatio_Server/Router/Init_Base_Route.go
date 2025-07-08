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

}
