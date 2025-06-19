package BaseService

import (
	"AITranslatio/Global"
	"AITranslatio/Src/DAO/BaseDAO"
	"AITranslatio/Src/DTO"

	"fmt"
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

func (service *BaseService) Login(LoginDTO DTO.LoginDTO) error {
	fmt.Println(LoginDTO)
	fmt.Println(LoginDTO.UserID)
	fmt.Print(LoginDTO.Password)
	return nil
}
