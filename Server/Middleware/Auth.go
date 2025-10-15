package Middleware

import (
	tokenUtil "AITranslatio/Utils/token"
	"github.com/gin-gonic/gin"
)

// token校验
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("Authorization")

		err := tokenUtil.ParseToken(token)
		if err != nil {
			return
		}
		c.Next()
	}
}
