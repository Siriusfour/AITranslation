package AuthService

import (
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/PasswordSecurity"
	"AITranslatio/Utils/SnowFlak"
	"AITranslatio/Utils/token"
	"AITranslatio/app/DAO/AuthDAO"
	types2 "AITranslatio/app/types"
	"errors"
)

func (Service *AuthService) Register(DTO *types2.RegisterDTO) (*types2.Auth, error) {

	// TODO 1.验证邮箱验证码是否有效（是否存在于redis）

	// TODO 2.验证随机树是否有效 （是否存在于redis）
	//var Challenge = &WebAuthn.Challenge{
	//	Verify:    DTO.Verify,
	//	Timestamp: DTO.Timestamp,
	//	Domain:    DTO.Domain,
	//}

	//if ok, err := WebAuthn.VerifyChallenge(Challenge); err != nil || !ok {
	//	return nil, err
	//}

	p := PasswordSecurity.CreatePasswordGeneratorFactory(12)

	//为每个用户生成随机盐值
	salt, err := p.GenerateSalt()
	if err != nil {
		return nil, err
	}

	//password+salt进行bcrypt加密
	HashPasswordWithSalt, err := p.HashPasswordWithSalt(DTO.Password, salt)
	if err != nil {
		return nil, err
	}

	//调用雪花算法生成唯一UserID
	UserID := SnowFlak.CreateSnowflakeFactory().GetID()

	//补完DTO
	DTO.UserID = UserID
	DTO.Salt = salt
	DTO.Password = HashPasswordWithSalt

	//调用DAO存储在库中
	err = AuthDAO.CreateDAOFactory("mysql").Register(DTO)
	if err != nil {
		return nil, errors.New(MyErrors.ErrorRegisterIsFail + err.Error())
	}

	//生成token
	//AccessToken, ErrAK := token.CreateTokenFactory(Consts.AccessToken, UserID).GeneratedToken()
	//RefreshToken, ErrRK := token.CreateTokenFactory(Consts.RefreshToken, UserID).GeneratedToken()

	if ErrRK != nil || ErrAK != nil {
		return nil, errors.New(ErrAK.Error() + ErrRK.Error())
	}

	//----------------业务逻辑

	return &types2.Auth{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
	}, nil

}
