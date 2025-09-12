package bootstrap

import "AITranslatio/Global"
import "AITranslatio/DataBase"

func InitDB() {

	if Global.DB_Config.GetInt("IsInitGlobalGormMysql") == 1 {
		MySQL_Client, err := DataBase.InitMySQL_Client()
		if err != nil {
			return
		}
		Global.MySQL_Client = MySQL_Client
	}

	if Global.DB_Config.GetInt("IsInitGlobalGormPostgreSql") == 1 {
		PostgreSQL_Client, err := DataBase.InitPostgreSQL_Client()
		if err != nil {
			return
		}
		Global.PostgreSQL_Client = PostgreSQL_Client
	}

}
