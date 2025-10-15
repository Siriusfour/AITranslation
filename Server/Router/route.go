package Router

import (
	"AITranslatio/Global"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Middleware"
	"AITranslatio/Utils/GinRelease"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter() *gin.Engine {

	var router *gin.Engine
	//【生产模式】
	// 根据 gin 官方的说明：[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
	// 如果部署到生产环境，请使用以下模式：
	// 1.生产模式(release) 和开发模式的变化主要是禁用 gin 记录接口访问日志，
	// 2.go服务就必须使用nginx作为前置代理服务，这样也方便实现负载均衡
	// 3.如果程序发生 panic 等异常使用自定义的 panic 恢复中间件拦截、记录到日志
	product := Global.Config.GetBool("Mode.Product")
	if product {
		router = GinRelease.ReleaseRouter()

	} else { //开发模式，开启 pprof 包，便于开发阶段分析程序性能
		router = gin.Default()
		pprof.Register(router)
	}
	//打开调试面板的颜色
	gin.ForceConsoleColor()

	// 设置可信任的代理服务器列表,gin
	if Global.Config.GetInt("HttpServer.TrustProxies.IsOpen") == 1 {
		if err := router.SetTrustedProxies(Global.Config.GetStringSlice("HttpServer.TrustProxies.ProxyServerList")); err != nil {
			Global.Logger.Error(MyErrors.ErrorGinSetTrustProxy, zap.Error(err))
		}
	} else {
		_ = router.SetTrustedProxies(nil)
	}
	//路由分组
	rgBase := router.Group("/Api").Use(Middleware.Auth()).Use(Middleware.Cors())  // 基础crud业务的路由组
	rgAuth := router.Group("/Auth").Use(Middleware.Auth()).Use(Middleware.Cors()) // 鉴权相关的路由组

	//注册所有组别的路由

	InitAuthRoute(rgAuth)
	InitBaseRoute(rgBase)

	return router

}
