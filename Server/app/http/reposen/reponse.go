package reposen

import (
	"AITranslatio/Global/Consts"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Request 从https里解析出来的结构
type Request struct {
	Ctx *gin.Context
	DTO any
}

// Response 返回给客户端的结构体
type Response struct {
	//HTTP状态码的扩展，自定义的扩展码，
	Code int `json:"code,omitempty"`
	//本次请求结果的详细描述
	Message string `json:"message,omitempty"`
	//返回的数据
	Data interface{} `json:"data,omitempty"`
}

// HttpResponse 设置响应的 JSON 数据和 http 状态码，并向客户端返回默认的状态码
func ReturnResponse(ctx *gin.Context, HttpCode int, ServerCode int, Message error, Data interface{}) {
	ctx.JSON(HttpCode, Response{
		ServerCode,
		MessageHandle(Message),
		Data,
	})
	ctx.Abort()
}

// Fail 通用错误,HttpCode=ServerCode=400
func Fail(Context *gin.Context, err error, data ...interface{}) {
	ReturnResponse(Context, http.StatusBadRequest, http.StatusBadRequest, err, data)
}

// success handller
func OK(Context *gin.Context, data interface{}) {
	ReturnResponse(Context, http.StatusOK, http.StatusOK, errors.New("success"), data)
}

// ErrorSystem 服务器代码错误
func ErrorSystem(Context *gin.Context, err error) {
	ReturnResponse(Context, http.StatusInternalServerError, Consts.ServerOccurredErrorCode, err, nil)
}

// ErrorTokenAuthFail token解析失败、该用户权限不足
func ErrorTokenAuthFail(Context *gin.Context, err error, Code ...int) {

	if len(Code) == 0 {
		ReturnResponse(Context, http.StatusUnauthorized, http.StatusUnauthorized, err, nil)
	} else {
		ReturnResponse(Context, http.StatusUnauthorized, Code[1], err, nil)
	}
}

// ErrorParam 参数校验错误
func ErrorParam(Context *gin.Context, err error) {
	ReturnResponse(Context, http.StatusBadRequest, Consts.ValidatorParamsCheckFailCode, err, nil)

}
