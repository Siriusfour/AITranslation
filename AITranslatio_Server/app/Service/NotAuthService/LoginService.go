package NotAuthService

import (
	"AITranslatio/Global"
	"AITranslatio/Utils/token"
	"AITranslatio/app/DAO/UserDAO"
	"AITranslatio/app/http/DTO/NotAuthDTO"
	"context"
	"errors"
	"strconv"
)

type UserInfo struct {
	Nickname string
	UserID   string
}

func CreateNotAuthService() *NotAuthService {

	return &NotAuthService{}
}

type NotAuthService struct{}

func (Service *NotAuthService) LoginByPassWord(UserID int, PassWord string) (error, *NotAuthDTO.Auth) {

	//验证PW,成功的话刷新内存里面的AK,RK，并向客户端返回新的AK.RK
	//查数据库校验PassWord
	DAO := UserDAO.CreateDAOfactory("mysql")
	err := DAO.LoginByPassword(UserID, PassWord)
	if err != nil {
		return errors.New("登录失败："), nil
	}

	//验证通过，生成ak，rk ，写入redis，返回请求
	AccessToken, errAk := token.CreateTokenFactory(Global.Config.GetInt("Token.AkOutTime")).GeneratedToken(UserID)
	RefreshToken, errRk := token.CreateTokenFactory(Global.Config.GetInt("Token.AkOutTime")).GeneratedToken(UserID)

	if errAk != nil && errRk != nil {
		return errors.New(Global.ErrorGeneratedTokenIsFail + ":" + errAk.Error() + "," + errRk.Error()), nil
	}

	TokenInfo := &token.Token{
		AccessToken:    AccessToken,
		RefreshToken:   RefreshToken,
		RegisteredTime: Global.DataFormt, //注册时间
	}

	err = Global.RedisClient.HMSet(context.Background(), "userID_"+strconv.Itoa(UserID), TokenInfo).Err()
	if err != nil {
		return err, nil
	}

	return nil, &NotAuthDTO.Auth{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
	}

}
