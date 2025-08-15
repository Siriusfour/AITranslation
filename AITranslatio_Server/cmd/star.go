package cmd

import (
	"AITranslatio/Global"
	"AITranslatio/Utils/SSE"
	"AITranslatio/bootstrap"
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
	Global.Logger = bootstrap.InitLogger()

	//=======初始化系统配置
	bootstrap.InitConfig()

	//=======初始化redis客户端
	bootstrap.InitRedis()

	//======初始化MySQL数据库
	var err error
	Global.DB, err = bootstrap.InitDB()

	Global.SSEClients = SSE.InitSSEClients()

	if err != nil {
		panic(err)
	}

	return
}
