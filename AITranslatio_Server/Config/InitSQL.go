package Config

import (
	"AITranslatio/Src/Model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func InitDB() (*gorm.DB, error) {
	LogMode := logger.Info
	if !viper.GetBool("Mode.develop") {
		LogMode = logger.Error
	}
	db, err := gorm.Open(mysql.Open(viper.GetString("DB.DNS")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   "sys_",
		},
		Logger: logger.Default.LogMode(LogMode),
	})

	if err != nil {
		return nil, err
	}

	SqlDB, _ := db.DB()
	SqlDB.SetMaxIdleConns(viper.GetInt("DB.MaxIdleConns"))
	SqlDB.SetMaxOpenConns(viper.GetInt("DB.MaxOpenConns"))
	SqlDB.SetConnMaxLifetime(time.Hour)

	//将数据库的表与user对象同步，有则修改，没有创建
	err = db.AutoMigrate(&Model.User{}, &Model.Note{}, &Model.Branch{}, &Model.Commit{})
	if err != nil {
		panic(err)
	}
	return db, nil
}
