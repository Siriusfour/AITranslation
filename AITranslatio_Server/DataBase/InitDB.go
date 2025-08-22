package DataBase

import (
	"AITranslatio/Global"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"strings"
)

func initMySQL_Client() (client *gorm.DB, err error) {
	SQL_Type := "MySQL"
	IsOpenReadDB := Global.DB_Config.GetInt("DB." + SQL_Type + "IsOpenReadDB")
	return GetSqlDriver(SQL_Type, IsOpenReadDB)
}

func initPostgreSQL_Client() {

}

// 获取数据库驱动, 可以通过options 动态参数连接任意多个数据库
func GetSqlDriver(sqlType string, readDbIsOpen int, dbConf ...ConfigParams) (*gorm.DB, error) {

	var dbDialector gorm.Dialector
	if val, err := getDbDialector(sqlType, "Write", dbConf...); err != nil {
		Global.Logger.Error(Global.ErrorsDialectorDbInitFail+sqlType, zap.Error(err))
	} else {
		dbDialector = val
	}
	gormDb, err := gorm.Open(dbDialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 redefineLog(sqlType), //拦截、接管 gorm v2 自带日志
	})
	if err != nil {
		//gorm 数据库驱动初始化失败
		return nil, err
	}
	return gormDb, nil
}

// 根据不同的连接参数，获取具体的一类数据库的连接指针
func getDbDialector(sqlType, readWrite string, dbConf ...ConfigParams) (gorm.Dialector, error) {
	var dbDialector gorm.Dialector
	dsn := getDsn(sqlType, readWrite, dbConf...)
	switch strings.ToLower(sqlType) {
	case "mysql":
		dbDialector = mysql.Open(dsn)
	//case "postgres"

	default:
		return nil, errors.New(Global.ErrorsDbDriverNotExists + sqlType)
	}
	return dbDialector, nil
}

// 根据配置参数生成数据库驱动 dsn
func getDsn(sqlType, readWrite string, dbConf ...ConfigParams) string {
	Host := Global.DB_Config.GetString("Gormv2." + sqlType + "." + readWrite + ".Host")
	DataBase := Global.DB_Config.GetString("Gormv2." + sqlType + "." + readWrite + ".DataBase")
	Port := Global.DB_Config.GetInt("Gormv2." + sqlType + "." + readWrite + ".Port")
	User := Global.DB_Config.GetString("Gormv2." + sqlType + "." + readWrite + ".User")
	Pass := Global.DB_Config.GetString("Gormv2." + sqlType + "." + readWrite + ".Pass")
	Charset := Global.DB_Config.GetString("Gormv2." + sqlType + "." + readWrite + ".Charset")

	if len(dbConf) > 0 {
		if strings.ToLower(readWrite) == "write" {
			if len(dbConf[0].Write.Host) > 0 {
				Host = dbConf[0].Write.Host
			}
			if len(dbConf[0].Write.DataBase) > 0 {
				DataBase = dbConf[0].Write.DataBase
			}
			if dbConf[0].Write.Port > 0 {
				Port = dbConf[0].Write.Port
			}
			if len(dbConf[0].Write.User) > 0 {
				User = dbConf[0].Write.User
			}
			if len(dbConf[0].Write.Pass) > 0 {
				Pass = dbConf[0].Write.Pass
			}
			if len(dbConf[0].Write.Charset) > 0 {
				Charset = dbConf[0].Write.Charset
			}
		} else {
			if len(dbConf[0].Read.Host) > 0 {
				Host = dbConf[0].Read.Host
			}
			if len(dbConf[0].Read.DataBase) > 0 {
				DataBase = dbConf[0].Read.DataBase
			}
			if dbConf[0].Read.Port > 0 {
				Port = dbConf[0].Read.Port
			}
			if len(dbConf[0].Read.User) > 0 {
				User = dbConf[0].Read.User
			}
			if len(dbConf[0].Read.Pass) > 0 {
				Pass = dbConf[0].Read.Pass
			}
			if len(dbConf[0].Read.Charset) > 0 {
				Charset = dbConf[0].Read.Charset
			}
		}
	}

	switch strings.ToLower(sqlType) {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=false&loc=Local", User, Pass, Host, Port, DataBase, Charset)
	case "sqlserver", "mssql":
		return fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;encrypt=disable", Host, Port, DataBase, User, Pass)
	case "postgresql", "postgre", "postgres":
		return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable TimeZone=Asia/Shanghai", Host, Port, DataBase, User, Pass)
	}
	return ""
}

func redefineLog(sqlType string) gormLog.Interface {

	createCustomGormLog(sqlType)

}
