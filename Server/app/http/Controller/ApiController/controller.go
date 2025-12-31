package ApiController

import (
	"AITranslatio/Config/interf"
	"AITranslatio/app/Service/ApiServer"
	"go.uber.org/zap"
)

type ApiController struct {
	Cfg     interf.ConfigInterface
	Logger  *zap.Logger
	Service *ApiServer.ApiServer
}

func NewController(Cfg interf.ConfigInterface, Logger *zap.Logger, service *ApiServer.ApiServer) *ApiController {
	return &ApiController{
		Cfg:     Cfg,
		Logger:  Logger,
		Service: service,
	}

}
