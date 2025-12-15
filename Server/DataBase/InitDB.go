package DataBase

import (
	"AITranslatio/Config/interf"
	"go.uber.org/zap"

	"AITranslatio/Global/MyErrors"
	"AITranslatio/app/Model/Team"
	"AITranslatio/app/Model/User"
	"AITranslatio/app/Model/webAuthn"
	"log"

	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"strings"
)

func InitMySQL_Client(cfg interf.ConfigInterface, logger *zap.Logger) (client *gorm.DB, err error) {
	SQL_Type := "MySQL"
	//IsOpenReadDB := cfg.GetInt("DB." + SQL_Type + ".IsOpenReadDB")
	return GetSqlDriver(cfg, logger, SQL_Type)
}
func InitPostgreSQL_Client(cfg interf.ConfigInterface, logger *zap.Logger) (client *gorm.DB, err error) {
	SQL_Type := "PostgreSQL"
	//IsOpenReadDB := cfg.GetInt("MySQL_DB." + SQL_Type + "IsOpenReadDB")
	return GetSqlDriver(cfg, logger, SQL_Type)
}

// 获取数据库驱动, 可以通过options 动态参数连接任意多个数据库
func GetSqlDriver(cfg interf.ConfigInterface, logger *zap.Logger, sqlType string, dbConf ...ConfigParams) (*gorm.DB, error) {

	var dbDialector gorm.Dialector
	if val, err := getDbDialector(sqlType, "Write", cfg, dbConf...); err != nil {
		log.Fatal("GetSqlDriver失败")
	} else {
		dbDialector = val
	}
	MyLogger := redefineLog(cfg, logger, sqlType)

	gormDb, err := gorm.Open(dbDialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 MyLogger,
	})

	if err != nil {
		//gorm 数据库驱动初始化失败
		return nil, err
	}

	err = gormDb.AutoMigrate(&Team.Team{}, &User.User{}, &Team.JoinTeamApplication{}, &webAuthn.WebAuthnCredential{})
	if err != nil {
		return nil, fmt.Errorf("gorm自动建表失败:%w", err)
	}

	return gormDb, nil
}

// 根据不同的连接参数，获取具体的一类数据库的连接指针
func getDbDialector(sqlType, readWrite string, cfg interf.ConfigInterface, dbConf ...ConfigParams) (gorm.Dialector, error) {
	var dbDialector gorm.Dialector
	dsn := getDsn(sqlType, readWrite, cfg, dbConf...)
	switch strings.ToLower(sqlType) {

	case "mysql":
		dbDialector = mysql.Open(dsn)

	default:
		return nil, errors.New(MyErrors.ErrorsDbDriverNotExists + sqlType)
	}
	return dbDialector, nil
}

// 根据配置参数生成数据库驱动 dsn
func getDsn(sqlType, readWrite string, cfg interf.ConfigInterface, dbConf ...ConfigParams) string {
	Host := cfg.GetString("DB." + sqlType + "." + readWrite + ".Host")
	DataBase := cfg.GetString("DB." + sqlType + "." + readWrite + ".DataBase")
	Port := cfg.GetInt("DB." + sqlType + "." + readWrite + ".Port")
	User := cfg.GetString("DB." + sqlType + "." + readWrite + ".User")
	Pass := cfg.GetString("DB." + sqlType + "." + readWrite + ".Password")
	Charset := cfg.GetString("DB." + sqlType + "." + readWrite + ".Charset")

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
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local", User, Pass, Host, Port, DataBase, Charset)
	case "sqlserver", "mssql":
		return fmt.Sprintf("server=%s;port=%d;database=%s;UserModel id=%s;password=%s;encrypt=disable", Host, Port, DataBase, User, Pass)
	case "postgresql", "postgre", "postgres":
		return fmt.Sprintf("host=%s port=%d dbname=%s UserModel=%s password=%s sslmode=disable TimeZone=Asia/Shanghai", Host, Port, DataBase, User, Pass)
	}
	return ""
}
func redefineLog(cfg interf.ConfigInterface, logger *zap.Logger, sqlType string) gormLog.Interface {
	return NewDBLogger(
		cfg,

		logger,

		sqlType,

		WithLogLevel(gormLog.Info),

		// 慢 SQL 阈值
		WithSlowThreshold(cfg.GetDuration("SlowThreshold")),

		// 是否忽略 RecordNotFound（推荐忽略，避免日志被灌爆）
		WithIgnoreRecordNotFound(true),
	)
}
