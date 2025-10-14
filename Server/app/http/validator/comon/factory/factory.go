package factory

import (
	"AITranslatio/Global"
	"AITranslatio/Global/CustomErrors"
	"AITranslatio/app/core/container"
	"AITranslatio/app/http/validator/interf"
	"github.com/gin-gonic/gin"
)

// 表单参数验证器工厂 （请勿修改）
func Create(key string) func(context *gin.Context) {

	if value := container.CreateContainersFactory().Get(key); value != nil {
		if val, isOk := value.(interf.ValidatorInterface); isOk {
			return val.CheckParams
		}
	}

	Global.Logger.Error(CustomErrors.ErrorsValidatorNotExists + ", 验证器模块：" + key)
	return nil
}
