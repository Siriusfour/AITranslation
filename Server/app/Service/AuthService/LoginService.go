package AuthService

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Utils/token"
	"AITranslatio/app/DAO/AuthDAO"
	"AITranslatio/app/types"
	"fmt"
	"go.uber.org/zap"
)

type AuthService struct {
	DAO           AuthDAO.Inerf
	TokenProvider token.TokenProvider
	Loggers       map[string]*zap.Logger
}

func NewAuthService(AuthDAO AuthDAO.Inerf, TokenProvider token.TokenProvider, loggers map[string]*zap.Logger) *AuthService {
	return &AuthService{
		AuthDAO,
		TokenProvider,
		loggers,
	}
}

func (Service *AuthService) LoginByPassWord(Email string, PassWord string) (*types.LoginInfo, error) {

	//验证PW,并向客户端返回新的AK.RK
	//查数据库校验PassWord

	LoginInfo, err := Service.DAO.LoginByPassword(Email, PassWord)
	if err != nil {
		return nil, fmt.Errorf("：数据库查询密码失败：%w", err)
	}

	//验证通过，生成ak，rk ，写入redis，返回请求
	AccessToken, errAk := Service.TokenProvider.GeneratedToken(LoginInfo.UserID, Consts.AccessToken)
	RefreshToken, errRk := Service.TokenProvider.GeneratedToken(LoginInfo.UserID, Consts.RefreshToken)

	if errAk != nil || errRk != nil {
		return nil, fmt.Errorf("：验证成功但生成token失败：%w,%w", errAk, errRk)
	}

	LoginInfo.Auth.AccessToken = AccessToken
	LoginInfo.Auth.RefreshToken = RefreshToken

	return LoginInfo, nil
}

func (Service *AuthService) FindUser(userID int64, IDType string) (*types.LoginInfo, error) {

	userInfo, err := Service.DAO.FindUserByID(userID, IDType)
	if err != nil {
		return nil, err
	}

	loginInfo := &types.LoginInfo{
		Nickname: userInfo.Nickname,
		UserID:   userInfo.UserID,
		Avatar:   userInfo.Avatar,
	}

	return loginInfo, nil
}
