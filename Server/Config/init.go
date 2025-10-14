package Config

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/CustomErrors"
	"github.com/spf13/viper"
	"sync"
)

func CreateConfigFactory(FileName string, Type string) interf.ConfigInterface {

	yamlConfig := viper.New()
	yamlConfig.AddConfigPath(Consts.BasePath + "/Config")
	yamlConfig.SetConfigType(Type)
	yamlConfig.SetConfigName(FileName)

	if err := yamlConfig.ReadInConfig(); err != nil {
		panic(CustomErrors.ErrorsConfigYamlNotExists + err.Error())
	}

	return &interf.ConfigFile{
		Viper: yamlConfig,
		Mu:    new(sync.Mutex),
	}
}
