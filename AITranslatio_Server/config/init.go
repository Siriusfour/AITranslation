package config

import (
	"AITranslatio/config/interf"
	"AITranslatio/config/json"
	"AITranslatio/config/yaml"
	"fmt"
)

func InitConfig(configTYpe string, filename string) interf.ConfigInterface {

	var config interf.ConfigInterface

	switch configTYpe {
	case "ini":
		fmt.Println("====ini====")

	case "yaml":
		fmt.Println("====yaml====")

		y := &yaml.YamlType{}
		config = y.CreateConfig(filename)

	case "json":
		fmt.Println("====json====")
		j := &json.JsonType{}
		j.CreateConfig(filename)

	}

	return config
}
