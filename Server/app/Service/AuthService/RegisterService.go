package AuthService

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/PasswordSecurity"
	"AITranslatio/app/types/DTO"

	"errors"
)

func (s *AuthService) Register(dto *DTO.RegisterDTO) (*DTO.Auth, error) {

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
	HashPasswordWithSalt, err := p.HashPasswordWithSalt(dto.Password, salt)
	if err != nil {
		return nil, err
	}

	//调用雪花算法生成唯一UserID
	UserID := s.SnowFlakeGenerator.GetID()

	//补完DTO
	dto.UserID = UserID
	dto.Salt = salt
	dto.Password = HashPasswordWithSalt

	//调用DAO存储在库中
	err = s.DAO.CreateUser(dto)
	if err != nil {
		return nil, errors.New(MyErrors.ErrorRegisterIsFail + err.Error())
	}

	//生成token
	AccessToken, ErrAK := s.TokenProvider.GeneratedToken(UserID, Consts.AccessToken)
	RefreshToken, ErrRK := s.TokenProvider.GeneratedToken(UserID, Consts.RefreshToken)

	if ErrRK != nil || ErrAK != nil {
		return nil, errors.New(ErrAK.Error() + ErrRK.Error())
	}

	//----------------业务逻辑

	return &DTO.Auth{AccessToken: AccessToken, RefreshToken: RefreshToken}, nil
}
