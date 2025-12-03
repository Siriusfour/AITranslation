package bootstrap

import (
	"AITranslatio/Config/interf"
	"path"
)

// Configurator 是一个接口，定义了创建配置的方法

func InitConfig(FileName string) interf.ConfigInterface {

	//读取后缀
	FileType := path.Ext(FileName)

	return interf.CreateConfigFactory(FileName, FileType[1:])

}
