package bootstrap

import (
	"AITranslatio/Global"
	"AITranslatio/Utils/token"
	"AITranslatio/app/DAO/AuthDAO"
	"AITranslatio/app/Service/AuthService"
	"AITranslatio/app/http/Controller/ApiController"
	"AITranslatio/app/http/Controller/AuthController"
)

type Controller struct {
	ApiController  *ApiController.ApiController
	AuthController *AuthController.AuthController
}

type APP struct {
	Controller *Controller
}

func InitApp() *APP {

	t := token.CreateTokenFactory(&token.CreateToken{
		Global.EncryptKey,
		Global.Config.GetDuration("Token.AkOutTime"),
		Global.Config.GetDuration("Token.RkOutTime"),

		Global.SnowflakeManage,
		Global.RedisClient,
	})

	authDAO := AuthDAO.CreateDAOFactory("mysql")
	authService := AuthService.NewAuthService(authDAO, t, Global.Logger)
	authController := AuthController.NewController(authService, Global.Tracing.Tracer, Global.Logger)

	App := &APP{
		Controller: &Controller{
			ApiController:  nil,
			AuthController: authController,
		},
	}

	return App
}
