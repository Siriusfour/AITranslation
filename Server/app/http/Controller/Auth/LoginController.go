package Auth

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/DAO/UserDAO"
	//"AITranslatio/app/DAO/UserDAO"
	"AITranslatio/app/Service/AuthService"
	"AITranslatio/app/http/reposen"
	"github.com/gin-gonic/gin"
)

type NotAuthController struct{}

func (Controller *NotAuthController) Login(LoginCtx *gin.Context) {

	//从ctx解析出需要的参数
	Email := LoginCtx.GetString(Consts.ValidatorPrefix + "Email")
	passWord := LoginCtx.GetString(Consts.ValidatorPrefix + "Password")
	AccessToken := LoginCtx.GetString(Consts.ValidatorPrefix + "AccessToken")
	RefreshToken := LoginCtx.GetString(Consts.ValidatorPrefix + "RefreshToken")

	//请求有ak优先使用ak登录
	if AccessToken != "" {
		err := UserDAO.CreateDAOFactory("mysql").LoginByAccessToken(AccessToken)
		if err != nil {
		}

		return
	}

	//请求有Rk优先使用Rk刷新登录
	if RefreshToken != "" {
		err := UserDAO.CreateDAOFactory("mysql").LoginByRefreshToken(RefreshToken)
		if err != nil {
		}
		return
	}

	//请求有PW优先使用PW登录
	if passWord != "" {
		err, Auth := AuthService.CreateNotAuthService().LoginByPassWord(Email, passWord)
		if err != nil {
			reposen.Fail(LoginCtx, "", err.Error())
		}

		reposen.OK(LoginCtx, Auth)

	}

	return

}

func (Controller *NotAuthController) LoginByWebAuthn(LoginCtx *gin.Context) {

}

//func (Controller *NotAuthController) CreateSSE(CreateSSEctx *gin.Context) {
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
