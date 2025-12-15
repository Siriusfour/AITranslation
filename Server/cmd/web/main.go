package main

import (
	"AITranslatio/Global"
	"AITranslatio/Router"
	"flag"

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

	cfg := bootstrap.InitConfig(fileName)
	logger := bootstrap.InitLogger(cfg)
	db, err := bootstrap.InitDB(cfg, logger["DB"])
	redisClient := bootstrap.InitRedis()
	mqClinet := bootstrap.InitMQClient(cfg, logger["MQ"])
	tracing := bootstrap.InitZipkin(cfg)

	if err != nil {
		panic(err)
	}

	//传入基础设备，初始化全局变量
	infra := Global.InitInfrastructure(cfg, logger, db, redisClient, mqClinet, tracing)

	App := bootstrap.InitApp(
		infra.EncryptKey,
		infra.Config,
		infra.DbClient,
		infra.RedisClient,
		infra.Logger["Business"],
		infra.Tracing,
		infra.JwtManager,
	)

	router := Router.InitRouter(App)

	serverMiddleware := zipkinHTTP.NewServerMiddleware(
		tracing.Tracer,
	)

	addr := cfg.GetString("HttpServer.Port")

	if err := http.ListenAndServe(addr, serverMiddleware(router)); err != nil {
		panic("启动失败：" + err.Error())
	}

}
