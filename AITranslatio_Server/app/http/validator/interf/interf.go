package interf

import "github.com/gin-gonic/gin"

// 验证器接口
type ValidatorInterface interface {
	CheckParams(context *gin.Context)
}
