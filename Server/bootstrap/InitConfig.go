package bootstrap

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Global"

	"path"
)

// Configurator 是一个接口，定义了创建配置的方法

func InitConfig(FileName string) {

	//读取后缀
	FileType := path.Ext(FileName)

	Global.Config = interf.CreateConfigFactory(FileName, FileType[1:])

}
