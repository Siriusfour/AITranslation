package bootstrap

import (
	"AITranslatio/Global"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
)

func InitConfig() {

	//获取当前完整路径
	_, filename, _, _ := runtime.Caller(0)
	absPath := filepath.Dir(filename)

	viper.SetConfigName("setting")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(absPath)
	err := viper.ReadInConfig()
	if err != nil {
		Global.Logger.Errorf(fmt.Sprintf("加载配置文件出错：%s", err))
		panic(fmt.Sprintf("加载配置文件出错：%s", err))
	}

	fmt.Printf("AllSettings:%s", viper.AllSettings())

	//热重载
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Printf("Error reloading config file: %s", err)
		} else {
			// 重新加载配置后，输出新的配置值
			fmt.Println("Updated config value:", viper.AllSettings())
		}
	})
}
