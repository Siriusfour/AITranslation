package bootstrap

import (
	"AITranslatio/Config/interf"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)
import "AITranslatio/DataBase"

func InitDB(cfg interf.ConfigInterface, logger *zap.Logger) (*gorm.DB, error) {

	if cfg.GetInt("DB.MySQL.IsInitGlobalGormMysql") == 1 {
		MySQL_Client, err := DataBase.InitMySQL_Client(cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("MySQL初始化失败 %w", err)
		}
		return MySQL_Client, nil
	}

	if cfg.GetInt("DB.PostgreSql.IsInitGlobalGormPostgreSql") == 1 {
		PostgreSQL_Client, err := DataBase.InitPostgreSQL_Client(cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("PostgreSQL初始化失败 %w", err)
		}
		return PostgreSQL_Client, nil
	}

	return nil, fmt.Errorf("不允许的数据库类型")
}
