package AuthController

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/Service/AuthService"
	"AITranslatio/app/http/reposen"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func (Controller *AuthController) Login(ctx *gin.Context) {

	//从ctx解析出需要的参数
	Email := ctx.GetString(Consts.ValidatorPrefix + "Email")
	passWord := ctx.GetString(Consts.ValidatorPrefix + "Password")

	err, Auth := AuthService.CreateAuthService().LoginByPassWord(Email, passWord)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		reposen.Fail(ctx, fmt.Errorf("登录失败，密码错误！%w", err))
		return
	}

	//=====加入业务逻辑===

	if err != nil {
		reposen.Fail(ctx, fmt.Errorf("登录失败：%w", err))
		return
	}

	reposen.OK(ctx, Auth)

	return
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
