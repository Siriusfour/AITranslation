package data_transfer

import (
	"AITranslatio/Global"
	"AITranslatio/Utils/token"
	"AITranslatio/app/http/validator/interf"
	"github.com/gin-gonic/gin"
	"github.com/go-viper/mapstructure/v2"

	"time"
)

// 将验证器成员(字段)绑定到数据传输上下文，方便控制器获取
/**
本函数参数说明：
validatorInterface 实现了验证器接口的结构体
extra_add_data_prefix  验证器绑定参数传递给控制器的数据前缀
context  gin上下文
*/

func DataAddContext(validatorInterface interf.ValidatorInterface, extraAddDataPrefix string, ctx *gin.Context) *gin.Context {
	var value map[string]interface{}
	// 🚀 mapstructure 保留类型，不会把 int 转 float64
	if err := mapstructure.Decode(validatorInterface, &value); err == nil {
		for k, v := range value {
			ctx.Set(extraAddDataPrefix+k, v)
		}
	}

	curDateTime := time.Now().Format(Global.DataFormt)
	ctx.Set(extraAddDataPrefix+"created_at", curDateTime)
	ctx.Set(extraAddDataPrefix+"updated_at", curDateTime)
	ctx.Set(extraAddDataPrefix+"deleted_at", curDateTime)

	jwt := ctx.GetHeader("Authorization")
	if jwt != "" {
		if err, userID := token.GetDataFormToken[int64](jwt, "UserID"); err == nil {
			ctx.Set(extraAddDataPrefix+"UserID", userID)
		}
	}
	return ctx
}
