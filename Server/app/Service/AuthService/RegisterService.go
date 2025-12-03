package AuthService

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/PasswordSecurity"

	"AITranslatio/app/types"
	"errors"
)

func (Service *AuthService) Register(DTO *types.RegisterDTO) (*types.Auth, error) {

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
	UserID := Service.SnowFlakeGenerator.GetID()

	//补完DTO
	DTO.UserID = UserID
	DTO.Salt = salt
	DTO.Password = HashPasswordWithSalt

	//调用DAO存储在库中
	err = Service.DAO.CreateUser(DTO)
	if err != nil {
		return nil, errors.New(MyErrors.ErrorRegisterIsFail + err.Error())
	}

	//生成token
	AccessToken, ErrAK := Service.TokenProvider.GeneratedToken(UserID, Consts.AccessToken)
	RefreshToken, ErrRK := Service.TokenProvider.GeneratedToken(UserID, Consts.RefreshToken)

	if ErrRK != nil || ErrAK != nil {
		return nil, errors.New(ErrAK.Error() + ErrRK.Error())
	}

	//----------------业务逻辑

	return &types.Auth{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
	}, nil

}
