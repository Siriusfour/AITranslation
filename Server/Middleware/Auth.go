package Middleware

import (
	"AITranslatio/Global"
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/app/http/reposen"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

// Auth token校验
func Auth(app *Global.Infrastructure) gin.HandlerFunc {
	return func(c *gin.Context) {

		jwt := c.GetHeader("Authorization")

		if jwt != "" {
			jwtInfo, err := app.JwtManager.ParseToken(jwt)
			if err != nil {
				if errors.Is(err, MyErrors.ErrTokenExpired) {
					reposen.ErrorTokenAuthFail(c, fmt.Errorf("登录失败:登录信息已过期%w", err), Consts.JwtTokenExpired)
					return
				} else {
					reposen.ErrorTokenAuthFail(c, fmt.Errorf("解析token失败！%w", err), Consts.JwtTokenInvalid)
					return
				}
			}
			c.Set("UserID", jwtInfo.UserID)
		} else {
			reposen.ErrorTokenAuthFail(c, fmt.Errorf("缺失登录相关信息!"), Consts.JwtTokenInvalid)
			return
		}

		c.Next()

	}
}
