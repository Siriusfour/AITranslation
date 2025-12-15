package Router

import (
	"AITranslatio/Global"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Middleware"
	"AITranslatio/Utils/GinRelease"
	"AITranslatio/Utils/metrics"
	"AITranslatio/bootstrap"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
)

func InitRouter(app *bootstrap.APP) *gin.Engine {

	var router *gin.Engine
	cfg := Global.GetInfra().Config
	Logger := Global.GetInfra().Logger

	//【生产模式】
	// 根据 gin 官方的说明：[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
	// 如果部署到生产环境，请使用以下模式：
	// 1.生产模式(release) 和开发模式的变化主要是禁用 gin 记录接口访问日志，
	// 2.go服务就必须使用nginx作为前置代理服务，这样也方便实现负载均衡
	// 3.如果程序发生 panic 等异常使用自定义的 panic 恢复中间件拦截、记录到日志
	product := cfg.GetBool("Mode.Product")
	if product {
		router = GinRelease.ReleaseRouter()

	} else { //开发模式，开启 pprof 包，便于开发阶段分析程序性能
		router = gin.Default()
		pprof.Register(router)
	}
	//打开调试面板的颜色
	gin.ForceConsoleColor()

	// 设置可信任的代理服务器列表,gin
	if cfg.GetInt("HttpServer.TrustProxies.IsOpen") == 1 {
		if err := router.SetTrustedProxies(cfg.GetStringSlice("HttpServer.TrustProxies.ProxyServerList")); err != nil {
			Logger["DB"].Error(MyErrors.ErrorGinSetTrustProxy, zap.Error(err))
		}
	} else {
		_ = router.SetTrustedProxies(nil)
	}

	//service := "Suis_Note"
	version := os.Getenv("SERVICE_VERSION")
	if version == "" {
		version = "dev"
	}

	metrics.MustRegister()

	//router.Use(Middleware.Prometheus(service, version)).GET("/metrics", gin.WrapH(promhttp.Handler()))
	//路由分组
	rgBase := router.Group("/Api").Use(Middleware.Auth(Global.GetInfra())).Use(Middleware.Cors()).Use(Middleware.HttpLog(Logger, "HTTP")).Use(Middleware.SessionID(Global.GetInfra()))  // 基础crud业务的路由组
	rgAuth := router.Group("/Auth").Use(Middleware.Auth(Global.GetInfra())).Use(Middleware.Cors()).Use(Middleware.HttpLog(Logger, "HTTP")).Use(Middleware.SessionID(Global.GetInfra())) // 鉴权相关的路由组
	rgNotAuth := router.Group("/NotAuth").Use(Middleware.Cors()).Use(Middleware.HttpLog(Logger, "HTTP")).Use(Middleware.SessionID(Global.GetInfra()))
	//注册所有组别的路由

	InitAuthRoute(rgAuth, app)
	InitBaseRoute(rgBase, app)
	InitNotAuthRoute(rgNotAuth, app)

	return router

}
