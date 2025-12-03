package Middleware

import (
	"AITranslatio/Global"
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/SnowFlak"
	tokenUtil "AITranslatio/Utils/token"
	"AITranslatio/app/http/reposen"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

// token校验
func Auth(app *Global.Infrastructure) gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")

		if c.Request.URL.Path == "/Auth/Login" {
			c.Next()
		}

		Key := app.EncryptKey
		AkOutTime := app.Config.GetDuration("Token.AkOutTime")
		RkOutTime := app.Config.GetDuration("Token.RkOutTime")
		redis := app.RedisClient
		SnowFlakManager := SnowFlak.CreateSnowflakeFactory(app.Config, app.Logger["business"])

		ct := &tokenUtil.CreateToken{
			Key,
			AkOutTime,
			RkOutTime,
			SnowFlakManager,
			redis,
		}

		err := tokenUtil.CreateTokenFactory(ct).ParseToken(token)
		if err != nil {
			if errors.Is(err, MyErrors.ErrTokenExpired) {
				reposen.ErrorTokenAuthFail(c, fmt.Errorf("登录失败,登录信息已过期"), Consts.JwtTokenExpired)
				return
			}

			if errors.Is(err, MyErrors.ErrTokenExpired) {
				reposen.Fail(c, fmt.Errorf("登录失败,登录信息已过期"))
				return
			}

		}
		c.Next()
	}
}
