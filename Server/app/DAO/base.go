package DAO

import (
	"AITranslatio/Global"
	"AITranslatio/Global/MyErrors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

func ChooseDB_Conn(sqlType string) *gorm.DB {

	var DB_Client *gorm.DB
	sqlType = strings.Trim(sqlType, " ")
	switch strings.ToLower(sqlType) {

	case "mysql":

		if Global.MySQL_Client == nil {
			Global.Logger["db"].Error(MyErrors.ErrorsGormNotInitGlobalPointer, zap.String("sqlType", sqlType))
		}

		DB_Client = Global.MySQL_Client

	case "postgressql":
		if Global.PostgreSQL_Client == nil {
			Global.Logger["db"].Error(MyErrors.ErrorsGormNotInitGlobalPointer, zap.String("sqlType", sqlType))
		}

		DB_Client = Global.PostgreSQL_Client

	}
	return DB_Client
}
