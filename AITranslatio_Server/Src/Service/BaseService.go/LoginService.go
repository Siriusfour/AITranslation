package BaseService

import (
	"AITranslatio/Global"
	"AITranslatio/Src/DAO/BaseDAO"
	"AITranslatio/Src/DTO"
	"AITranslatio/Utils"
	"AITranslatio/Utils/UtilsStruct"
	"errors"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
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
	if LoginDTO.Password != "" {
		//查数据库校验PassWord
		err := BaseService.BaseDAO.LoginByPassword(LoginDTO.UserID, LoginDTO.Password)
		if err != nil {
			return errors.New("密码验证错误"), nil
		}

		//验证通过，生成ak，rk ，写入map，返回请求
		AccessToken, err := Utils.GeneratedToken(Global.PKEY, jwt.SigningMethodHS256, LoginDTO.UserID, time.Duration(10000))
		RefreshToken, err := Utils.GeneratedToken(Global.PKEY, jwt.SigningMethodHS256, LoginDTO.UserID, time.Duration(70000))
		if err != nil {
			return err, nil
		}

		Global.TokenMap.TokenMap[LoginDTO.UserID] = &UtilsStruct.TokenInfo{
			AccessToken:    AccessToken,
			RefreshToken:   RefreshToken,
			Revoked:        true,                                  //RK是否已经被吊销
			RegisteredTime: time.Now().Format("2006-01-02 15:04"), //注册时间
		}

		return nil, &DTO.Auth{
			AccessToken:  AccessToken,
			RefreshToken: RefreshToken,
		}

	} else {
		return errors.New("请输入密码"), nil
	}

}
