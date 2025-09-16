package yaml

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/CustomErrors"
	"AITranslatio/config/interf"
	"github.com/spf13/viper"
	"sync"
)

type YamlType struct{}

func (y *YamlType) CreateConfig(FileName ...string) interf.ConfigInterface {

	yamlConfig := viper.New()
	yamlConfig.AddConfigPath(Consts.BasePath + "/config")
	viper.AddConfigPath(".")
	yamlConfig.SetConfigType("yaml")
	if len(FileName) > 0 {
		yamlConfig.SetConfigName(FileName[0])
	} else {
		yamlConfig.SetConfigName("config")
	}

	if err := yamlConfig.ReadInConfig(); err != nil {
		panic(CustomErrors.ErrorsConfigYamlNotExists + err.Error())
	}

	return &interf.ConfigFile{
		Viper: yamlConfig,
		Mu:    new(sync.Mutex),
	}
}
