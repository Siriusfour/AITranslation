package NotAuth

import (
	"AITranslatio/Global"
	"AITranslatio/app/DAO/UserDAO"
	"AITranslatio/app/Service/NotAuthService"
	"AITranslatio/app/http/reposen"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

type NotAuthController struct{}

func (Controller *NotAuthController) Login(LoginCtx *gin.Context) {

	//从ctx解析出需要的参数
	userID := LoginCtx.GetInt(Global.ValidatorPrefix + "user_name")
	passWord := LoginCtx.GetString(Global.ValidatorPrefix + "pass")
	AccessToken := LoginCtx.GetString(Global.ValidatorPrefix + "accessToken")
	RefreshToken := LoginCtx.GetString(Global.ValidatorPrefix + "refreshToken")

	//请求有ak优先使用ak登录
	if AccessToken != "" {
		err := UserDAO.CreateDAOfactory().LoginByAccessToken(AccessToken)
		if err != nil {
		}

		return
	}

	//请求有Rk优先使用Rk登录
	if RefreshToken != "" {
		err := UserDAO.CreateDAOfactory().LoginByRefreshToken(RefreshToken)
		if err != nil {
		}
		return
	}

	//请求有PW优先使用PW登录
	if passWord != "" {
		 err, Auth := NotAuthService.CreateNotAuthService().LoginByPassWord(userID, passWord)
		 if err != nil {
			reposen.Fail(LoginCtx, "", err.Error())
		}

		reposen.OK(LoginCtx, Auth, "登录成功")
	}

	return

}

func (Controller *NotAuthController) CreateSSE(CreateSSEctx *gin.Context) {

	//获取到URL的参数，解析出UserID
	ID, exists := CreateSSEctx.GetQuery("UserID")
	if !exists {
		HTTPErr(CreateSSEctx, errors.New("参数不存在"), 1012)
	}

	UserID, err := strconv.Atoi(ID)
	if err != nil {
		HTTPErr(CreateSSEctx, errors.New("参数解析失败！"), 1013)
	}

	err = Global.SSEClients.CreateSSE(CreateSSEctx, UserID)
	if err != nil {
		HTTPErr(CreateSSEctx, errors.New("创建SSE链接失败！"), 1014)
	}


