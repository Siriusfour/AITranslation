package AuthService

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Utils/token"
	"AITranslatio/app/DAO/UserDAO"
	"AITranslatio/app/types"
	"fmt"
)

func CreateAuthService() *AuthService {
	return &AuthService{}
}

type AuthService struct{}

func (Service *AuthService) LoginByPassWord(Email string, PassWord string) (*types.LoginInfo, error) {

	//验证PW,并向客户端返回新的AK.RK
	//查数据库校验PassWord

	DAO := UserDAO.CreateDAOFactory("mysql")
	LoginInfo, err := DAO.LoginByPassword(Email, PassWord)
	if err != nil {
		return nil, fmt.Errorf("登录失败：%w", err)
	}

	//验证通过，生成ak，rk ，写入redis，返回请求
	AccessToken, errAk := token.CreateTokenFactory(Consts.AccessToken, LoginInfo.UserID).GeneratedToken()
	RefreshToken, errRk := token.CreateTokenFactory(Consts.RefreshToken, LoginInfo.UserID).GeneratedToken()

	if errAk != nil && errRk != nil {
		return nil, fmt.Errorf("LoginByPassWord函数生成token失败：%w,%w", errAk, errRk)
	}

	LoginInfo.Auth.AccessToken = AccessToken
	LoginInfo.Auth.RefreshToken = RefreshToken

	return LoginInfo, nil

}
