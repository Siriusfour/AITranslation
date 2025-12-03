package AuthService

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Global/Consts"
	"AITranslatio/Utils/SnowFlak"
	"AITranslatio/Utils/token"
	"AITranslatio/app/DAO/AuthDAO"
	"AITranslatio/app/types"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type AuthService struct {
	cfg                interf.ConfigInterface
	Loggers            *zap.Logger
	TokenProvider      token.TokenProvider
	SnowFlakeGenerator *SnowFlak.SnowFlakeGenerator
	RedisClient        *redis.Client
	DAO                AuthDAO.Inerf
}

func NewService(cfg interf.ConfigInterface, logger *zap.Logger, TokenProvider token.TokenProvider, SnowFlakeGenerator *SnowFlak.SnowFlakeGenerator, redisClient *redis.Client, AuthDAO AuthDAO.Inerf) *AuthService {
	return &AuthService{
		cfg,
		logger,
		TokenProvider,
		SnowFlakeGenerator,
		redisClient,
		AuthDAO,
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
