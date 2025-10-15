package main

import (
	"AITranslatio/Global"
	"AITranslatio/Router"
	_ "AITranslatio/bootstrap"
)

func main() {
	router := Router.InitRouter()
	_ = router.Run(Global.Config.GetString("HttpServer.Port"))

	//开启MQ消费者

}
