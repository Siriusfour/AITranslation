package bootstrap

import (
	"AITranslatio/app/http/validator/comon"
	"flag"
)

func init() {

	//todo 0.检查项目必须的非编译目录是否存在，避免编译后调用的时候缺失相关目录，

	//1. 接收命令行启动参数
	configType := *(flag.String("config-type", "yaml", "type of config file (e.g., yaml)"))
	configFile := *(flag.String("config", "setting.yaml", "path to config file"))
	flag.Parse()

	//2.注册表单校验容器
	comon.RegisterValidator()

	//3.初始化config文件,DBconfig文件
	InitConfig(configType, configFile)

	//初始化SQL数据库
	InitDB()

	//初始化日志组件
	InitLogger()

	//初始化redis客户端
	InitRedis()

}
