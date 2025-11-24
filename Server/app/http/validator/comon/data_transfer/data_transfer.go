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
	// 1. 把 validatorInterface 展平塞进 context
	var value map[string]interface{}
	if err := mapstructure.Decode(validatorInterface, &value); err == nil {
		flattenAndSetContext(ctx, extraAddDataPrefix, value)
	}

	// 2. 把 URL 的参数（query + path）也塞进 context
	urlData := make(map[string]interface{})

	// 2.1 query 参数，例如 ?page=1&name=xx
	query := ctx.Request.URL.Query()
	for k, vals := range query {
		if len(vals) == 1 {
			urlData[k] = vals[0]
		} else {
			// 有多个值的话，看你需求，可以直接存 slice
			urlData[k] = vals
		}
	}

	// 2.2 路由 path 参数，例如 /users/:id
	// gin 中是 ctx.Param("id") 或 ctx.Params
	for _, p := range ctx.Params {
		urlData[p.Key] = p.Value
	}

	// 展平 URL 参数
	flattenAndSetContext(ctx, extraAddDataPrefix, urlData)

	// 3. 公共字段：时间
	curDateTime := time.Now().Format(Global.DataFormt)
	ctx.Set(extraAddDataPrefix+"created_at", curDateTime)
	ctx.Set(extraAddDataPrefix+"updated_at", curDateTime)
	ctx.Set(extraAddDataPrefix+"deleted_at", curDateTime)

	// 4. 从 JWT 里取 UserID
	jwt := ctx.GetHeader("Authorization")
	if jwt != "" {
		if err, userID := token.GetDataFormToken[int64](jwt, "UserID"); err == nil {
			ctx.Set(extraAddDataPrefix+"UserID", userID)
		}
	}
	return ctx
}

func flattenAndSetContext(ctx *gin.Context, prefix string, data map[string]interface{}) {
	for k, v := range data {
		switch val := v.(type) {
		case map[string]interface{}:
			// 如果你想保留层级，可以改成 flattenAndSetContext(ctx, prefix+k+".", val)
			flattenAndSetContext(ctx, prefix, val)
		default:
			ctx.Set(prefix+k, val)
		}
	}
}
