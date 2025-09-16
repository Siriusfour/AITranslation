package Router

import (
	"AITranslatio/Global"
	"fmt"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	gin.ForceConsoleColor()

	var router *gin.Engine
	router = gin.Default()
	fmt.Println(router)

	//路由分组
	//rgBase := r.Group("/Note/Api").Use(Middleware.Auth())   // 基础crud业务的路由组
	//rgFile := r.Group("/File").Use(Middleware.Auth())       //资源操 作相关路由组
	//rgCaptcha := r.Group("/Captcha").Use(Middleware.Auth()) // 验证码 相关操作的路由组
	rgNotAuth := router.Group("/NotAuth") // 不需要token验证的路由组

	//注册所有组别的路由
	router.POST("/hello", func(context *gin.Context) {
		s := Global.Config.GetString("HttpServer.Port")
		context.JSON(200, gin.H{
			"code": 200,
			"msg":  s,
		})
	})
	InitNotAuthRoute(rgNotAuth)

	return router

}
