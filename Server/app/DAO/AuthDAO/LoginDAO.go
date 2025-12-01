package AuthDAO

import (
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/PasswordSecurity"
	base "AITranslatio/app/DAO"
	"AITranslatio/app/Model/User"
	"AITranslatio/app/types"
	"errors"
	"gorm.io/gorm"
)

type Inerf interface {
	LoginByPassword(Email string, password string) (*types.LoginInfo, error)
	FindUserByID(ID int64, IDtype string) (*User.User, error)

	CreateUser(UserInfo *User.User) error
}

type AuthDAO struct {
	DB_Client *gorm.DB
}

func CreateDAOFactory(sqlType string) *AuthDAO {
	return &AuthDAO{
		DB_Client: base.ChooseDB_Conn(sqlType),
	}
}

func (DAO *AuthDAO) CreateUser(user *User.User) error {
	return DAO.DB_Client.Create(user).Error
}

// CheckOAuthID 根据ID查找用户是否存在（OAuthID/UserID）
func (DAO *AuthDAO) FindUserByID(ID int64, IDtype string) (*User.User, error) {

	var UserInfo *User.User

	result := DAO.DB_Client.First(&UserInfo, IDtype+" = ?", ID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {

		return nil, MyErrors.ErrorOAuthIDrNotFound

	} else if result.Error != nil {

		return nil, result.Error

	} else {
		return UserInfo, nil
	}
}

//func (DAO *AuthDAO) FindUser(UserID int64) (*types.LoginInfo, error) {
//	var UserInfo User.User
//
//	result := DAO.DB_Client.Table("User").Where("UserID = ?", UserID).First(&UserInfo)
//	if result.Error != nil {
//		return nil, result.Error
//	}
//
//	loginInfo := &types.LoginInfo{
//		UserID:   UserInfo.UserID,
//		Nickname: UserInfo.Nickname,
//		Avatar:   UserInfo.Avatar,
//	}
//
//	return loginInfo, nil
//
//}

// 通过密码登录
func (DAO *AuthDAO) LoginByPassword(Email string, password string) (*types.LoginInfo, error) {

	var UserInfo User.User

	result := DAO.DB_Client.Table("User").Where("Email = ?", Email).First(&UserInfo)
	if result.Error != nil {
		return nil, result.Error
	}

	if err := PasswordSecurity.CreatePasswordGeneratorFactory(12).ValidatePasswordWithSalt(UserInfo.Password, password, UserInfo.Salt); err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &types.LoginInfo{
		Auth:     types.Auth{},
		Nickname: UserInfo.Nickname,
		UserID:   UserInfo.UserID,
		Avatar:   UserInfo.Avatar,
	}, nil

}
