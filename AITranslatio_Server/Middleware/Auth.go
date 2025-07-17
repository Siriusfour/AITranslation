package Middleware

import (
	"AITranslatio/Src/HTTP"
	"AITranslatio/Utils"
	"github.com/gin-gonic/gin"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 定义不需要认证的路由
		publicRoutes := []string{"/Attendance/Api/BaseControll/API/Login"}
		currentPath := c.Request.URL.Path

		// 检查当前路由是否在公共路由列表中
		for _, route := range publicRoutes {
			if strings.HasPrefix(currentPath, route) {
				c.Next() // 跳过认证，直接继续
				return
			}
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			token2, _ := c.Cookie("Authorization")

			token = token2
		}

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
