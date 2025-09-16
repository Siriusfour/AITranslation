package NotAuthService

import (
	"AITranslatio/Global"
	"AITranslatio/Global/CustomErrors"
	"AITranslatio/Utils/PasswordSecurity"
	"AITranslatio/Utils/SnowFlak"
	"AITranslatio/Utils/token"
	"AITranslatio/app/DAO/UserDAO"
	"AITranslatio/app/http/DTO/NotAuthDTO"
	"errors"
)

func (Service *NotAuthService) Register(DTO *NotAuthDTO.RegisterDTO) (*NotAuthDTO.Auth, error) {

	// TODO 1.验证邮箱验证码是否有效（是否存在于redis）

	// TODO 2.验证随机树是否有效 （是否存在于redis）
	//var Challenge = &WebAuthn.Challenge{
	//	Verify:    DTO.Verify,
	//	Timestamp: DTO.Timestamp,
	//	Domain:    DTO.Domain,
	//}
	//
	//if ok, err := WebAuthn.VerifyChallenge(Challenge); err != nil || !ok {
	//	return nil, err
	//}

	p := PasswordSecurity.CreatePasswordGeneratorFactory(12)

	//为每个用户生成随机盐值
	salt, err := p.GenerateSalt(32)
	if err != nil {
		return nil, err
	}

	//password+salt进行bcrypt加密
	HashPasswordWithSalt, err := p.HashPasswordWithSalt(DTO.Password, salt)
	if err != nil {
		return nil, err
	}

	//调用雪花算法生成唯一UserID
	UserID := SnowFlak.CreateSnowflakeFactory().GetId()

	//补完DTO
	DTO.UserID = UserID
	DTO.Salt = salt
	DTO.Password = HashPasswordWithSalt

	//调用DAO存储在库中
	err = UserDAO.CreateDAOfactory("mysql").Register(DTO)
	if err != nil {
		return nil, errors.New(CustomErrors.ErrorRegisterIsFail + err.Error())
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
