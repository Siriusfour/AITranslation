package BaseService

import (
	"AITranslatio/Global"
	"AITranslatio/Src/DAO/BaseDAO"
	"AITranslatio/Src/DTO"
	"AITranslatio/Src/HTTP"
	"AITranslatio/Utils"
	"AITranslatio/Utils/UtilsStruct"

	"github.com/gin-gonic/gin"
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

func (BaseService *BaseService) Login(ctx *gin.Context, LoginDTO *DTO.LoginDTO) error {

	switch {
	//验证ak，和原设备
	case LoginDTO.AccessToken != "":
		{
			_, err := Utils.ParseToken(Global.PKEY, LoginDTO.AccessToken)
			if err != nil {
				return err
			}
			return nil
		}
	//验证RK,成功的话返回新的AK
	case LoginDTO.RefreshToken != "":
		{

		}
	//验证PW,成功的话返回新的AK.RK
	case LoginDTO.Password != "":
		{
			//调用DAO查数据库校验PassWord
			err := BaseService.BaseDAO.LoginByPassword(LoginDTO.UserID, LoginDTO.Password)
			if err != nil {
				return err
			}

			//验证通过，生成ak，rk ，写入map，返回请求
			AccessToken, err := Utils.GeneratedToken(Global.PKEY, jwt.SigningMethodHS256, LoginDTO, time.Duration(1))
			RefreshToken, err := Utils.GeneratedToken(Global.PKEY, jwt.SigningMethodHS256, LoginDTO, time.Duration(7))
			if err != nil {
				return err
			}

			Global.TokenMap.TokenMap[LoginDTO.UserID] = &UtilsStruct.TokenInfo{
				AccessToken:    AccessToken,
				RefreshToken:   RefreshToken,
				Revoked:        true,
				RegisteredTime: time.Now().Format("2006-01-02 15:04"),
			}

			Tokens := HTTP.Tokens{
				AccessToken:  AccessToken,
				RefreshToken: RefreshToken,
			}
			reponse := HTTP.Response{
				Code:    10000,
				Message: "Login Success",
				Data:    nil,
				Token:   Tokens,
			}

			HTTP.OK(ctx, reponse)
		}
	}

	//查找用户信息并返回回给controller

	return nil
}
