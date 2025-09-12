package NotAuthService

import (
	"AITranslatio/Global"
	"AITranslatio/Utils/PasswordSecurity"
	"AITranslatio/Utils/SnowFlak"
	"AITranslatio/Utils/token"
	"AITranslatio/app/DAO/UserDAO"
	"AITranslatio/app/http/DTO/NotAuthDTO"
	"errors"
)

func (Service *NotAuthService) Register(UserName, Email, EmailCode, Password string) (*NotAuthDTO.Auth, error) {

	// TODO 1.验证邮箱验证码是否有效（是否存在于redis）

	p := PasswordSecurity.CreatePasswordGeneratorFactory(12)

	//为每个用户生成随机盐值
	salt, err := p.GenerateSalt(32)
	if err != nil {
		return nil, err
	}

	//password+salt进行bcrypt加密
	withSalt, err := p.HashPasswordWithSalt(Password, salt)
	if err != nil {
		return nil, err
	}

	//调用雪花算法生成唯一UserID
	UserID := SnowFlak.CreateSnowflakeFactory().GetId()

	//调用DAO存储在库中
	err = UserDAO.CreateDAOfactory("mysql").Register(UserID, UserName, Email, EmailCode, withSalt, salt)
	if err != nil {
		return nil, errors.New(Global.ErrorRegisterIsFail + err.Error())
	}

	//生成token
	AccessToken, ErrAK := token.CreateTokenFactory(Global.Config.GetInt("Token.AkOutTime")).GeneratedToken(UserID)
	RefreshToken, ErrRK := token.CreateTokenFactory(Global.Config.GetInt("Token.RkOutTime")).GeneratedToken(UserID)

	if ErrRK != nil || ErrAK != nil {
		return nil, errors.New(ErrAK.Error() + ErrRK.Error())
	}

	return &NotAuthDTO.Auth{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
	}, nil

}
