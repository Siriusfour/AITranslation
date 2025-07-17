package cmd

import (
	"AITranslatio/Config"
	"AITranslatio/Global"
	"AITranslatio/Utils/SSE"
	"AITranslatio/Utils/UtilsStruct"
	"fmt"
)

//根据配置文件
//0.初始化配置文件
//1.初始化路由
//2.初始化日志组件
//3.初始化数据库
//4.初始化存储Token的map

func Start() {
	fmt.Println("=======")

	//=======初始化日志组件
	Global.Logger = Config.InitLogger()

	//=======初始化系统配置
	Config.InitConfig()

	//=======初始化redis客户端
	Config.InitRedis()

	//======初始化MySQL数据库
	var err error
	Global.DB, err = Config.InitDB()

	//======初始化存储token的map
	Global.TokenMap = UtilsStruct.InitTokenMap()

	Global.SSEClients = SSE.InitSSEClients()

	if err != nil {
		panic(err)
	}

	return
}
