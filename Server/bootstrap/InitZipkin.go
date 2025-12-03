package bootstrap

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Utils/zipkin"
)

func InitZipkin(cfg interf.ConfigInterface) *zipkin.Tracing {

	z := zipkin.CreateTracing(cfg.GetString("Zipkin.ServerName"), cfg.GetString("Zipkin.URL"), cfg.GetString("Zipkin.Port"))

	return z

}
