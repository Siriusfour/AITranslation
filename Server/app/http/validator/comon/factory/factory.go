package factory

import (
	"AITranslatio/Global"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/app/core/container"
	"AITranslatio/app/http/validator/interf"
	"github.com/gin-gonic/gin"
)

// 表单参数验证器工厂
func Create(key string) func(context *gin.Context) {
	if value := container.CreateContainersFactory().Get(key); value != nil {
		if val, isOk := value.(interf.ValidatorInterface); isOk {
			return val.CheckParams
		}
	}

	Global.Logger.Error(MyErrors.ErrorsValidatorNotExists + ", 验证器模块：" + key)
	return nil

}
