package AuthController

import (
	"AITranslatio/Global"
	"AITranslatio/Global/Consts"
	"AITranslatio/app/DAO/UserDAO"
	"AITranslatio/app/Service/AuthService"
	"AITranslatio/app/http/reposen"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type AuthController struct{}

func (Controller *AuthController) Login(ctx *gin.Context) {

	span := zipkin.SpanFromContext(ctx.Request.Context())

	if span != nil {
		span.Tag("UserID", strconv.FormatInt(ctx.GetInt64("UserID"), 10)) // 从 Token 解析出来的
		defer span.Finish()
	}

	//从ctx解析出需要的参数
	Email := ctx.GetString(Consts.ValidatorPrefix + "Email")
	passWord := ctx.GetString(Consts.ValidatorPrefix + "Password")

	//有账号密码
	if Email != "" || passWord != "" {
		LoginInfo, err := AuthService.CreateAuthService().LoginByPassWord(Email, passWord)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			reposen.Fail(ctx, fmt.Errorf("登录失败，密码错误！%w", err))
			return
		}
		reposen.OK(ctx, LoginInfo)
		return

		//尝试获取token登录
	} else {
		childSpan, newCtx := Global.Tracing.Tracer.StartSpanFromContext(ctx, "GetUserInfo")
		defer childSpan.Finish()
		ctx.Request = ctx.Request.WithContext(newCtx)

		childSpan.Tag("db.system", "MySQL")
		childSpan.Tag("sql:", "FindUser")
		loginInfo, err := UserDAO.CreateDAOFactory("mysql").FindUser(ctx.GetInt64(Consts.ValidatorPrefix + "UserID"))
		if err != nil {
			reposen.ErrorSystem(ctx, fmt.Errorf("获取用户信息失败：%w", err))
			return
		}

		reposen.OK(ctx, loginInfo)

	}

}

//func (Controller *AuthController) CreateSSE(CreateSSEctx *gin.Context) {
//
//	//获取到URL的参数，解析出UserID
//	ID, exists := CreateSSEctx.GetQuery("UserID")
//	if !exists {
//		HTTPErr(CreateSSEctx, errors.New("参数不存在"), 1012)
//	}
//
//	UserID, err := strconv.Atoi(ID)
//	if err != nil {
//		HTTPErr(CreateSSEctx, errors.New("参数解析失败！"), 1013)
//	}
//
//	err = Global.SSEClients.CreateSSE(CreateSSEctx, UserID)
//	if err != nil {
//		HTTPErr(CreateSSEctx, errors.New("创建SSE链接失败！"), 1014)
//	}
//}
