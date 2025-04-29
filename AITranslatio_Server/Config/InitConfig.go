package Config

import (
	"AITranslatio/Global"
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
)

func InitConfig() {

	_, filename, _, _ := runtime.Caller(0)
	absPath := filepath.Dir(filename)
	println(absPath)
	viper.SetConfigName("setting")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(absPath)
	err := viper.ReadInConfig()
	if err != nil {
		Global.Logger.Errorf(fmt.Sprintf("加载配置文件出错：%s", err))
		panic(fmt.Sprintf("加载配置文件出错：%s", err))
	}

	fmt.Printf("%s", viper.AllSettings())
}
