package Router

import (
	"AITranslatio/Src/Controller/Base"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {

	//路由分组

	rgBase := r.Group("Attendance/Api/Base")

	//注册所有组别的路由
	initBasePaltformRouter(rgBase)
}

func initBasePaltformRouter(rgBase *gin.RouterGroup) {

	BaseController := Base.NewBaseController()
	Init_Base_Route(rgBase, BaseController)

}
