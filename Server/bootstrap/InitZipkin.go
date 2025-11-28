package bootstrap

import (
	"AITranslatio/Global"
	"AITranslatio/Utils/zipkin"
)

func InitZipkin() {

	z := zipkin.CreateTracing(Global.Config.GetString("Zipkin.ServerName"), Global.Config.GetString("Zipkin.URL"), Global.Config.GetString("Zipkin.Port"))

	Global.Tracing = z

}
