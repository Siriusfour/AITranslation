package BaseDAO

import (
	"AITranslatio/Global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseDAO struct {
	orm    *gorm.DB
	Logger *zap.SugaredLogger
}

func New_Base_DAO() *BaseDAO {
	return &BaseDAO{
		orm:    Global.DB,
		Logger: Global.Logger,
	}
}
