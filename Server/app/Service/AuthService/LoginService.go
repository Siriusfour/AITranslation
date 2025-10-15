package AuthService

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Utils/token"
	"AITranslatio/app/DAO/UserDAO"
	"AITranslatio/app/http/DTO/NotAuthDTO"
	"fmt"
)

func CreateAuthService() *AuthService {

	return &AuthService{}
}

type AuthService struct{}

func (Service *AuthService) LoginByPassWord(Email string, PassWord string) (error, *NotAuthDTO.Auth) {

	//验证PW,成功的话刷新内存里面的AK,RK，并向客户端返回新的AK.RK
	//查数据库校验PassWord
	DAO := UserDAO.CreateDAOFactory("mysql")
	UserID, err := DAO.LoginByPassword(Email, PassWord)
	if err != nil {
		return fmt.Errorf("登录失败：%w", err), nil
	}

	//验证通过，生成ak，rk ，写入redis，返回请求
	AccessToken, errAk := token.CreateTokenFactory(Consts.AccessToken, UserID).GeneratedToken()
	RefreshToken, errRk := token.CreateTokenFactory(Consts.RefreshToken, UserID).GeneratedToken()

	if errAk != nil && errRk != nil {
		return fmt.Errorf("登录失败：%w", errAk), nil
	}

	return nil, &NotAuthDTO.Auth{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
	}

}
