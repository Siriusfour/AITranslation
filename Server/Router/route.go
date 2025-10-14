package Router

import (
	"AITranslatio/Global"
	"AITranslatio/Utils/GinRelease"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	var router *gin.Engine
	//【生产模式】
	// 根据 gin 官方的说明：[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
	// 如果部署到生产环境，请使用以下模式：
	// 1.生产模式(release) 和开发模式的变化主要是禁用 gin 记录接口访问日志，
	// 2.go服务就必须使用nginx作为前置代理服务，这样也方便实现负载均衡
	// 3.如果程序发生 panic 等异常使用自定义的 panic 恢复中间件拦截、记录到日志
	product := Global.Config.GetBool("Mode.Develop")
	if product {
		router = GinRelease.ReleaseRouter()

	} else { //开发模式
		router = gin.Default()
	}

	gin.ForceConsoleColor()

	//路由分组
	rgBase := router.Group("/Api") // 基础crud业务的路由组
	//rgFile := r.Group("/File").Use(Middleware.AuthDTO())       //资源操 作相关路由组
	//rgCaptcha := r.Group("/Captcha").Use(Middleware.AuthDTO()) // 验证码 相关操作的路由组
	rgNotAuth := router.Group("/AuthDTO") // 不需要token验证的路由组

	//注册所有组别的路由

	InitNotAuthRoute(rgNotAuth)
	InitBaseRoute(rgBase)

	return router

}
