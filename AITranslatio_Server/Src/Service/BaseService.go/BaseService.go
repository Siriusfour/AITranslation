package BaseService

import (
	"AITranslatio/Global"
	"AITranslatio/Src/DAO/BaseDAO"
	"AITranslatio/Src/DTO"
	"AITranslatio/Src/HTTP"
	"AITranslatio/Utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"time"
)

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
	//验证ak
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

			//验证通过，生成ak，rk返回请求
			AccessToken, err := Utils.GeneratedToken(Global.PKEY, jwt.SigningMethodHS256, LoginDTO, time.Duration(1))
			RefreshToken, err := Utils.GeneratedToken(Global.PKEY, jwt.SigningMethodHS256, LoginDTO, time.Duration(7))
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

			return err
		}

	}

	return nil
}
