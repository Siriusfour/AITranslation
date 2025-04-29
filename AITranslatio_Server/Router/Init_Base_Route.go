package Router

import (
	"AITranslatio/Src/Controller/Base"
	"github.com/gin-gonic/gin"
)

func Init_Base_Route(rgBase *gin.RouterGroup, BaseController *Base.BaseController) {

	rgBase.POST("/API/Translation", BaseController.Translation)

}
