package Router

import (
	_ "AITranslatio/app/docs"
	"AITranslatio/bootstrap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitNotAuthRoute(rg gin.IRoutes, app *bootstrap.APP) {
	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	rg.POST("/Login", app.Controller.AuthController.Login)

	rg.GET("/LoginGetWebAuthnInfo", app.Controller.AuthController.LoginGetWebAuthnInfo)
	rg.POST("/LoginByWebAuthn", app.Controller.AuthController.LoginByWebAuthn)

	//OAuth相关
	rg.GET("/GetChallenge", app.Controller.AuthController.OAuthMap["Github"].GetChallenge)
	rg.POST("LoginByGithub", app.Controller.AuthController.OAuthMap["Github"].Login)
}

func InitAuthRoute(rg gin.IRoutes, app *bootstrap.APP) {

	rg.POST("/Register", app.Controller.AuthController.Register)
	rg.GET("/RegisterGetWebAuthnInfo", app.Controller.AuthController.RegisterGetWebAuthnInfo)
	rg.POST("/RegisterWebAuthn", app.Controller.AuthController.RegisterWebAuthn)

	rg.GET("/StartSeckill", app.Controller.ApiController.StartSeckill)
	rg.POST("/SeckillBuy", app.Controller.ApiController.SeckillBuy)

}
