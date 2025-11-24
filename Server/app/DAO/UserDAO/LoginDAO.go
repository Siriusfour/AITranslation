package UserDAO

import (
	"AITranslatio/Utils/PasswordSecurity"
	"AITranslatio/app/DAO"
	"AITranslatio/app/Model/User"
	"gorm.io/gorm"
)

type UserDAO struct {
	DB_Client *gorm.DB
}

func CreateDAOFactory(sqlType string) *UserDAO {
	return &UserDAO{
		DB_Client: DAO.ChooseDB_Conn(sqlType),
	}
}

func (DAO *UserDAO) CreateUser(user *User.User) error {
	return DAO.DB_Client.Create(user).Error
}

func (DAO *UserDAO) CheckOAuthID(ID int) (bool *User.User) {

}

func (DAO *UserDAO) FindUser(UserID int64) (*User.User, error) {
	var UserInfo User.User

	result := DAO.DB_Client.Table("User").Where("UserID = ?", UserID).First(&UserInfo)
	if result.Error != nil {
		return nil, result.Error
	}

	return &UserInfo, nil

}

// 通过密码登录
func (DAO *UserDAO) LoginByPassword(Email string, password string) (int64, error) {

	var UserInfo User.User

	result := DAO.DB_Client.Table("User").Where("Email = ?", Email).First(&UserInfo)
	if result.Error != nil {
		return 0, result.Error
	}

	if err := PasswordSecurity.CreatePasswordGeneratorFactory(12).ValidatePasswordWithSalt(UserInfo.Password, password, UserInfo.Salt); err != nil {
		return 0, err
	}

	if result.Error != nil {
		return 0, result.Error
	}

}

func (DAO *UserDAO) LoginByAccessToken(AccessToken string) error {
	return nil

}

func (DAO *UserDAO) LoginByRefreshToken(RefreshToken string) error {
	return nil
}
