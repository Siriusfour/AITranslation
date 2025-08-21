package bootstrap

import (
	"AITranslatio/Global"
	"AITranslatio/config/interf"
	"AITranslatio/config/json"
	"AITranslatio/config/yaml"
	"fmt"
)

func InitConfig(configTYpe string, filename string) {

	var config interf.ConfigInterface

	switch configTYpe {
	case "ini":
		fmt.Println("====ini====")

	case "yaml":
		fmt.Println("====yaml====")

		y := &yaml.YamlType{}
		Global.Config = y.CreateConfig(filename)

	case "json":
		fmt.Println("====json====")
		j := &json.JsonType{}
		j.CreateConfig(filename)

	}

	//初始化DB文件
	Global.DB_Config = config.Clone("DB")

}
