package BaseService

import (
	"AITranslatio/Global"
	"AITranslatio/Src/DAO/BaseDAO"
	"go.uber.org/zap"
)

type BaseService struct {
	Logger  *zap.SugaredLogger
	BaseDAO *BaseDAO.BaseDAO
}

func NewBaseService() *BaseService {
	return &BaseService{
		Logger:  Global.Logger,
		BaseDAO: BaseDAO.New_Base_DAO(),
	}
}

func (service *BaseService) Login() {

}
