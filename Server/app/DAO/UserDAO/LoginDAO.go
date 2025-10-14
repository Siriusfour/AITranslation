package UserDAO

import (
	"AITranslatio/Utils/PasswordSecurity"
	"AITranslatio/app/DAO"
	"AITranslatio/app/Model/UserModel"
	"errors"
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

	var UserInfo UserModel.User

	result := DAO.DB_Client.Table("User").Raw("select * from User where Email = ? limit 1", Email).Scan(&UserInfo)
	if result.Error != nil {
		return 0, result.Error
	}

	if PasswordSecurity.CreatePasswordGeneratorFactory().ValidatePasswordWithSalt(UserInfo.Password, password, UserInfo.Salt) {
		return 0, errors.New("密码错误！")
	}

	return UserInfo.UserID, nil

}

func (DAO *UserDAO) LoginByAccessToken(AccessToken string) error {
	return nil

}

func (DAO *UserDAO) LoginByRefreshToken(RefreshToken string) error {
	return nil
}
