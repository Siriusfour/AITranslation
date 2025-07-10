package HTTP

import (
	"AITranslatio/Src/DTO"
	"AITranslatio/Src/Model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"reflect"
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
	// 请求列表时候返回的总页数
	Sum int64 `json:"sum,omitempty"`
	//刷新的token信息
	Tokens DTO.Auth `json:"token"`
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
	Ctx *gin.Context
	//errors Utills.MyError
	Logger *zap.SugaredLogger
}

// HttpResponse 设置响应的 JSON 数据和 HTTP 状态码，并向客户端返回默认的状态码
func HttpResponse(ctx *gin.Context, response Response, Status int) {
	// 如果 response 为空（值类型为零值），终止请求并返回状态码
	if reflect.DeepEqual(response, Response{}) {
		ctx.AbortWithStatus(Status)
		return
	}
	//response不为空，将response序列化为json，返回给客户端，并终止本次通话
	ctx.AbortWithStatusJSON(Status, response)
}

// 向处理函数提供的接口，返回失败
func Fail(context *gin.Context, response Response) {
	HttpResponse(context, response, http.StatusBadRequest)
}

func OK(context *gin.Context, response Response) {
	HttpResponse(context, response, http.StatusOK)
}

func ServerFail(context *gin.Context, response Response) {
	HttpResponse(context, response, http.StatusInternalServerError)
}
