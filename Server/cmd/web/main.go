package main

import (
	"AITranslatio/Global"
	"AITranslatio/Router"
	_ "AITranslatio/bootstrap"

	"net/http"

	zipkinHTTP "github.com/openzipkin/zipkin-go/middleware/http"
)

func main() {

	router := Router.InitRouter()

	serverMiddleware := zipkinHTTP.NewServerMiddleware(
		Global.Tracing.Tracer,
	)

	addr := Global.Config.GetString("HttpServer.Port")

	if err := http.ListenAndServe(addr, serverMiddleware(router)); err != nil {
		panic("启动失败：" + err.Error())
	}
	//开启MQ消费者

}
