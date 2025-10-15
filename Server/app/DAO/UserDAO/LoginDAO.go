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

	return UserInfo.UserID, nil

}

func (DAO *UserDAO) LoginByAccessToken(AccessToken string) error {
	return nil

}

func (DAO *UserDAO) LoginByRefreshToken(RefreshToken string) error {
	return nil
}
