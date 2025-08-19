package reposen

import (
	"AITranslatio/Global"
	"AITranslatio/Src/Model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

type Branches struct {
	Branch  *Model.Branch
	Commits *[]Model.Commit
}
type Note struct {
	Note     *Model.Note
	Branches []Branches
}

type BaseController struct {
	Ctx    *gin.Context
	Logger *zap.SugaredLogger
}

// HttpResponse 设置响应的 JSON 数据和 HTTP 状态码，并向客户端返回默认的状态码
func ReturnResponse(ctx *gin.Context, HttpCode int, ServerCode int, Message string, Data interface{}) {
	ctx.JSON(HttpCode, Response{
		ServerCode,
		Message,
		Data,
	})
	ctx.Abort()
}

// Fail 通用错误,HttpCode=ServerCode=400
func Fail(context *gin.Context, data interface{}, Message string) {
	ReturnResponse(context, http.StatusBadRequest, http.StatusBadRequest, Message, data)
}

// OK 处理成功
func OK(context *gin.Context, data interface{}, Message string) {
	ReturnResponse(context, http.StatusOK, http.StatusOK, Message, data)
}

// ErrorSystem 服务器代码错误
func ErrorSystem(context *gin.Context, data interface{}, Message string) {
	ReturnResponse(context, http.StatusInternalServerError, Global.ServerOccurredErrorCode, Global.ServerOccurredErrorMsg+Message, data)
}

// TokenErrorParam token解析失败
func TokenErrorParam(context *gin.Context, data interface{}, Message string, wrongParam interface{}) {
	ReturnResponse(context, http.StatusUnauthorized, Global.ValidatorParamsCheckFailCode, Global.ValidatorParamsCheckFailMsg, wrongParam)
}
 
// ErrorTokenAuthFail token权限校验失败
func ErrorTokenAuthFail(c *gin.Context) {
	ReturnResponse(c, http.StatusUnauthorized, http.StatusUnauthorized, Global.ErrorsNoAuthorization, "")
}
