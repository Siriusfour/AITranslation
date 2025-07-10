package Middleware

import (
	"AITranslatio/Src/HTTP"
	"AITranslatio/Utils"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		err := Utils.Verify(token)
		if err != nil {

			HTTP.Fail(c, HTTP.Response{
				Code:    1001,
				Message: "认证失败：" + err.Error(),
			})

			return
		}
		c.Next()
	}
}
