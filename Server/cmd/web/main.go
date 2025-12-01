package main

import (
	"AITranslatio/Global"
	"AITranslatio/Router"

	"AITranslatio/bootstrap"
	_ "AITranslatio/bootstrap"

	"net/http"

	zipkinHTTP "github.com/openzipkin/zipkin-go/middleware/http"
)

func main() {

	App := bootstrap.InitApp()

	router := Router.InitRouter(App)

	serverMiddleware := zipkinHTTP.NewServerMiddleware(
		Global.Tracing.Tracer,
	)

	addr := Global.Config.GetString("HttpServer.Port")

	if err := http.ListenAndServe(addr, serverMiddleware(router)); err != nil {
		panic("启动失败：" + err.Error())
	}
	//开启MQ消费者

}
