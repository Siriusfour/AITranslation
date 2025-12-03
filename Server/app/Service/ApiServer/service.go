package ApiServer

import (
	"AITranslatio/Config/interf"
	"AITranslatio/app/DAO/ApiDAO"
	"go.uber.org/zap"
)

type ApiServer struct {
	cfg    interf.ConfigInterface
	logger *zap.Logger

	DAO *ApiDAO.ApiDAO
}

func NewService(cfg interf.ConfigInterface, logger *zap.Logger, DAO *ApiDAO.ApiDAO) *ApiServer {
	return &ApiServer{
		logger: logger,
		cfg:    cfg,
		DAO:    DAO,
	}

}
