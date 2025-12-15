package data_transfer

import (
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
	curDateTime := time.Now().Format("2006-01-02 15:04:05")
	ctx.Set(extraAddDataPrefix+"created_at", curDateTime)
	ctx.Set(extraAddDataPrefix+"updated_at", curDateTime)

	// 4. 从 JWT 里取 UserID
	//jwt := ctx.GetHeader("Authorization")
	//SessionID, err := ctx.Cookie("SessionID")
	//if err != nil {
	//	// 对于大多数情况，可以直接当成“没有 SessionID”
	//	SessionID = ""
	//}
	//if jwt != "" {
	//	if err, userID := token.GetDataFormToken[int64](jwt, "UserID"); err == nil {
	//		ctx.Set(extraAddDataPrefix+"UserID", userID)
	//		ctx.Set(extraAddDataPrefix+"Token", jwt)
	//	}
	//
	//} else if SessionID != "" { //没有jwt，但是有sessionID ，说明来的是游客临时会话，
	//	ctx.Set(extraAddDataPrefix+"SessionID", SessionID)
	//
	//} else { //没有jwt，cookie里面也没有sessionID , 说明是该设备的第一个请求， 创建一个sessionID
	//
	//	SessionID = SnowFlak.CreateSnowflakeFactory().GetIDString()
	//
	//	ctx.Set(extraAddDataPrefix+"SessionID", SessionID)
	//
	//	ctx.SetCookie("SessionID", SessionID, 0, "/", "", false, true)
	//
	//}
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
