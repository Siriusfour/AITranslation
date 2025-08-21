package bootstrap

import (
	"AITranslatio/Global"
	"AITranslatio/HTTP/validator/comon"
	"flag"
)

func init() {

	//1.检查项目必须的非编译目录是否存在，避免编译后调用的时候缺失相关目录，接收命令行启动参数
	configType := *(flag.String("config-type", "yaml", "type of config file (e.g., yaml)"))
	configFile := *(flag.String("config", "default.yml", "path to config file"))
	flag.Parse()

	//2.注册表单校验容器
	comon.WebRegisterValidator()

	//3.初始化config文件,DBconfig文件
	InitConfig(configType, configFile)

	//初始化SQL数据库
	InitDB(Global.Config.GetInt("DB.Type"))

	//初始化日志组件
	InitLogger()

	//初始化redis客户端
	InitRedis()

}
