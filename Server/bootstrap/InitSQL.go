package bootstrap

import (
	"AITranslatio/Global"
	"fmt"
)
import "AITranslatio/DataBase"

func InitDB() {

	dbTYpe := Global.Config.GetInt("test")
	fmt.Println(dbTYpe)

	if Global.Config.GetInt("DB.MySQL.IsInitGlobalGormMysql") == 1 {
		MySQL_Client, err := DataBase.InitMySQL_Client()
		if err != nil {
			return
		}
		Global.MySQL_Client = MySQL_Client
	}

	if Global.Config.GetInt("DB.PostgreSql.IsInitGlobalGormPostgreSql") == 1 {
		PostgreSQL_Client, err := DataBase.InitPostgreSQL_Client()
		if err != nil {
			return
		}
		Global.PostgreSQL_Client = PostgreSQL_Client
	}

}
