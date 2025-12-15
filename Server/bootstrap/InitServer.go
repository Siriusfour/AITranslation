package bootstrap

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Utils/SnowFlak"
	"AITranslatio/Utils/WebAuthn"
	"AITranslatio/Utils/token"
	"AITranslatio/Utils/zipkin"
	"AITranslatio/app/DAO/ApiDAO"
	"AITranslatio/app/DAO/AuthDAO"
	"AITranslatio/app/Service/ApiServer"
	"AITranslatio/app/Service/AuthService"
	"AITranslatio/app/Service/AuthService/OAuthService"
	"AITranslatio/app/http/Controller/ApiController"
	"AITranslatio/app/http/Controller/AuthController"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Controller struct {
	ApiController  *ApiController.ApiController
	AuthController *AuthController.AuthController
}

type APP struct {
	Controller *Controller
}

func InitApp(EncryptKey []byte, cfg interf.ConfigInterface, db *gorm.DB, redisClient *redis.Client, logger *zap.Logger, tracing *zipkin.Tracing, jwtGenerator *token.JWTGenerator) *APP {

	t := token.CreateTokenFactory(&token.CreateToken{
		EncryptKey,
		cfg.GetDuration("Token.AkOutTime"),
		cfg.GetDuration("Token.RkOutTime"),
		SnowFlak.CreateSnowflakeFactory(cfg, logger),
		redisClient,
	})

	s := SnowFlak.CreateSnowflakeFactory(cfg, logger)
	oauthDAO := AuthDAO.NewDAOFactory(db)
	GithubService := OAuthService.CreateOAuthServiceFactroy(cfg, logger, t, s, redisClient, oauthDAO, "Github")
	WxService := OAuthService.CreateOAuthServiceFactroy(cfg, logger, t, s, redisClient, oauthDAO, "WX")
	QQService := OAuthService.CreateOAuthServiceFactroy(cfg, logger, t, s, redisClient, oauthDAO, "QQ")
	oauthMap := map[string]AuthController.OAuthController{
		"Github": AuthController.NewOAuthControllerFactroy(GithubService, logger),
		"WX":     AuthController.NewOAuthControllerFactroy(WxService, logger),
		"QQ":     AuthController.NewOAuthControllerFactroy(QQService, logger),
	}

	w := WebAuthn.CreateWebAuthnConfigFactory(cfg)
	//AuthController的创建
	authDAO := AuthDAO.NewDAOFactory(db)
	authService := AuthService.NewService(cfg, logger, t, s, w, redisClient, authDAO)
	authController := AuthController.NewController(cfg, logger, authService, tracing.Tracer, oauthMap, jwtGenerator, w)

	//ApiController的创建
	apiDAO := ApiDAO.NewDAOFactory(db)
	apiService := ApiServer.NewService(cfg, logger, apiDAO)
	apiController := ApiController.NewController(cfg, logger, apiService)

	App := &APP{
		Controller: &Controller{
			ApiController:  apiController,
			AuthController: authController,
		},
	}

	return App
}
