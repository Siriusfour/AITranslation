package main

import (
	"AITranslatio/Global"
	"AITranslatio/Router"
	"context"
	"flag"
	"golang.org/x/time/rate"
	"sync"

	"AITranslatio/bootstrap"

	"net/http"

	zipkinHTTP "github.com/openzipkin/zipkin-go/middleware/http"
)

// @title           Note
// @version         1.0
// @description     Suis
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  2605116008@qq.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3008
func main() {

	fileName := *(flag.String("File", "setting.yaml", "path to Config file"))
	flag.Parse()

	var one *sync.Once = &sync.Once{}

	cfg := bootstrap.InitConfig(fileName)
	logger := bootstrap.InitLogger(cfg)
	redisClient := bootstrap.InitRedis()
	scripts := bootstrap.InitScripts(one)
	mqClinet := bootstrap.InitMQClient(cfg, logger["MQ"])
	tracing := bootstrap.InitZipkin(cfg)
	db, err := bootstrap.InitDB(cfg, logger["DB"])
	l := rate.NewLimiter(rate.Limit(cfg.GetInt("limit.QPS")), cfg.GetInt("limit.QPS"))

	if err != nil {
		panic(err)
	}

	//传入基础设备，初始化全局变量
	infra := Global.InitInfrastructure(cfg, logger, db, redisClient, scripts, mqClinet, tracing, l, one)

	App := bootstrap.InitApp(
		infra.EncryptKey,
		infra.Config,
		infra.DbClient,
		infra.RedisClient,
		infra.RabbitmqClient,
		infra.Scripts,
		infra.Logger["Business"],
		infra.Tracing,
		infra.JwtManager,
	)

	//消费者启动
	// 创建一个全局的上下文，用于控制所有后台任务的退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 程序退出时取消
	go bootstrap.InitConsumer(ctx, App.Controller.ApiController.Service)

	//路由初始化
	router := Router.InitRouter(App)

	//zipkin启动
	serverMiddleware := zipkinHTTP.NewServerMiddleware(
		tracing.Tracer,
	)

	addr := cfg.GetString("HttpServer.Port")

	if err := http.ListenAndServe(addr, serverMiddleware(router)); err != nil {
		panic("启动失败：" + err.Error())
	}

}
