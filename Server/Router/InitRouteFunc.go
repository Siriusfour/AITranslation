package Router

import (
	"AITranslatio/bootstrap"
	"github.com/gin-gonic/gin"
)

func InitAuthRoute(rg gin.IRoutes, app *bootstrap.APP) {
	rg.POST("/Login", app.Controller.AuthController.Login)
	rg.POST("/Register", app.Controller.AuthController.Register)
	rg.GET("/ApplicationWebAuthn", app.Controller.AuthController.ApplicationWebAuthn)

	rg.GET("/GetUserAllCredential", app.Controller.AuthController.GetUserAllCredential)
	rg.POST("/LoginByWebAuthn", app.Controller.AuthController.WebAuthnByLogin)

	//OAuth相关
	rg.GET("/GetChallenge", app.Controller.AuthController.OAuthMap["Github"].GetChallenge)
	rg.POST("LoginByGithub", app.Controller.AuthController.OAuthMap["Github"].Login)
}

func InitBaseRoute(rg gin.IRoutes, app *bootstrap.APP) {
	rg.POST("/Team/Create", app.Controller.ApiController.CreateTeam)
	rg.POST("/Team/Join", app.Controller.ApiController.JoinTeam)
}

func InitFilesRoute(rgBase *gin.RouterGroup) {}

func InitCaptchaRoute(rgBase *gin.RouterGroup) {}
