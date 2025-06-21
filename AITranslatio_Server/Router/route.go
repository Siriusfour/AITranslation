package Router

import (
	"AITranslatio/Src/Controller/Base"
	"fmt"
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
	fmt.Println("BaseController:", BaseController.Logger)
	fmt.Println("BaseController:", BaseController.BaseService)
	Init_Base_Route(rgBase, BaseController)

}
