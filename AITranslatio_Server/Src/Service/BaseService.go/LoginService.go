package BaseService

import (
	"AITranslatio/Global"
	"AITranslatio/Src/DAO/BaseDAO"
	"AITranslatio/Src/DTO"
	"AITranslatio/Utils"
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type UserInfo struct {
	Nickname string
	UserID   string
}

type BaseService struct {
	Logger  *zap.SugaredLogger
	BaseDAO *BaseDAO.BaseDAO
}

func NewBaseService() *BaseService {
	return &BaseService{
		Logger:  Global.Logger,
		BaseDAO: BaseDAO.New_Base_DAO(),
	}
}

func (BaseService *BaseService) Login(LoginDTO *DTO.LoginDTO) (error, *DTO.Auth) {

	//验证PW,成功的话刷新内存里面的AK,RK，并向客户端返回新的AK.RK
	if LoginDTO.Password != "" || LoginDTO.UserID != 0 {
		//查数据库校验PassWord
		err := BaseService.BaseDAO.LoginByPassword(LoginDTO.UserID, LoginDTO.Password)
		if err != nil {
			return errors.New("密码验证错误"), nil
		}

		//验证通过，生成ak，rk ，写入redis，返回请求
		AccessToken, err := Utils.GeneratedToken(Global.PKEY, jwt.SigningMethodHS256, LoginDTO.UserID, time.Duration(10000))
		RefreshToken, err := Utils.GeneratedToken(Global.PKEY, jwt.SigningMethodHS256, LoginDTO.UserID, time.Duration(70000))
		if err != nil {
			return err, nil
		}

		TokenInfo := &Utils.TokenInfo{
			AccessToken:    AccessToken,
			RefreshToken:   RefreshToken,
			RegisteredTime: time.Now().Format("2006-01-02 15:04"), //注册时间
		}

		err = Global.RedisClient.HMSet(context.Background(), "userID:"+strconv.Itoa(LoginDTO.UserID), TokenInfo).Err()
		if err != nil {
			return err, nil
		}

		return nil, &DTO.Auth{
			AccessToken:  AccessToken,
			RefreshToken: RefreshToken,
		}

	} else {
		return errors.New("请输入密码"), nil
	}

}
