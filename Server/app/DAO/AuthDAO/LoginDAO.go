package AuthDAO

import (
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/PasswordSecurity"
	"AITranslatio/app/Model/User"
	"AITranslatio/app/Model/webAuthn"
	"AITranslatio/app/types"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Inerf interface {
	LoginByPassword(Email string, password string) (*types.LoginInfo, error)
	FindUserByID(ID int64, IDtype string) (*User.User, error)
	CreateUser(UserInfo *types.RegisterDTO) error
	CreateUserByOAuth(user *User.User) error
	FindCredentialByUserID(UserID int64) ([]Credential, error)
	FindCredential(ctx *gin.Context) (*webAuthn.WebAuthnCredential, error)
}

type AuthDAO struct {
	DB_Client *gorm.DB
}

func NewDAOFactory(db *gorm.DB) *AuthDAO {
	return &AuthDAO{
		DB_Client: db,
	}
}

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

func (DAO *AuthDAO) CreateUser(UserInfo *types.RegisterDTO) error {
	return DAO.DB_Client.Create(UserInfo).Error
}

func (DAO *AuthDAO) CreateUserByOAuth(user *User.User) error {
	return DAO.DB_Client.Create(user).Error
}
