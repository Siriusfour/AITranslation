package Middleware

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	tokenUtil "AITranslatio/Utils/token"
	"AITranslatio/app/http/reposen"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

// token校验
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")

		if c.Request.URL.Path == "/Auth/Login" {
			c.Next()
		}

		err := tokenUtil.ParseToken(token)
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
