package main

import (
	"AITranslatio/Global"
	"AITranslatio/Router"
	"flag"

	"AITranslatio/bootstrap"

	"net/http"

	zipkinHTTP "github.com/openzipkin/zipkin-go/middleware/http"
)

func main() {

	fileName := *(flag.String("File", "setting.yaml", "path to Config file"))
	flag.Parse()

	cfg := bootstrap.InitConfig(fileName)
	logger := bootstrap.InitLogger(cfg)
	db, err := bootstrap.InitDB(cfg, logger["db"])
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
		infra.Logger["business"],
		infra.Tracing,
	)

	router := Router.InitRouter(App)

	serverMiddleware := zipkinHTTP.NewServerMiddleware(
		tracing.Tracer,
	)

	addr := cfg.GetString("HttpServer.Port")

	if err := http.ListenAndServe(addr, serverMiddleware(router)); err != nil {
		panic("启动失败：" + err.Error())
	}
	//开启MQ消费者

}
