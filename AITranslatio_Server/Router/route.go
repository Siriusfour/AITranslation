package Router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {

	//路由分组
	rgBase := r.Group("/Note/Api")        // 基础crud业务的路由组
	rgNotAuth := rgBase.Group("/NotAuth") // 不需要token验证的路由组
	rgAuth := rgBase.Group("/Auth")       // token 相关操作的路由组
	rgFile := rgBase.Group("/File")       //资源操作相关路由组
	rgCaptcha := rgBase.Group("/Captcha") // 验证码 相关操作的路由组

	//注册所有组别的路由

}
