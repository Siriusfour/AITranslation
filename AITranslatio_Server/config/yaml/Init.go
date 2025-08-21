package yaml

import (
	"AITranslatio/Global"
	"AITranslatio/config/interf"
	"github.com/spf13/viper"
	"sync"
)

type YamlType struct{}

func (y *YamlType) CreateConfig(FileName ...string) interf.ConfigInterface {

	yamlConfig := viper.New()
	yamlConfig.AddConfigPath(Global.BasePath + "/config")
	yamlConfig.SetConfigType("yaml")
	if len(FileName) > 0 {
		yamlConfig.SetConfigFile(FileName[0])
	} else {
		yamlConfig.SetConfigName("config")
	}

	if err := yamlConfig.ReadInConfig(); err != nil {
		panic(Global.ErrorsConfigYamlNotExists + err.Error())
	}

	return &interf.ConfigFile{
		Viper: yamlConfig,
		Mu:    new(sync.Mutex),
	}
}
