package Router

import (
	"AITranslatio/Src/Controller/BaseControll"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {

	//路由分组

	rgBase := r.Group("Attendance/Api/BaseControll")

	//注册所有组别的路由
	initBasePaltformRouter(rgBase)
}

func initBasePaltformRouter(rgBase *gin.RouterGroup) {

	BaseController := BaseControll.NewBaseController()
	Init_Base_Route(rgBase, BaseController)

}
