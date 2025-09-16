package UserDAO

import (
	"AITranslatio/app/DAO"
	"AITranslatio/app/Model/UserModel"
	"errors"
	"gorm.io/gorm"
)

type UserDAO struct {
	DB_Client *gorm.DB
}

func CreateDAOfactory(sqlType string) *UserDAO {
	return &UserDAO{
		DB_Client: DAO.ChooseDB_Conn(sqlType),
	}
}

// 通过密码登录
func (DAO *UserDAO) LoginByPassword(UserID int64, password string) error {

	var UserInfo UserModel.User

	result := DAO.DB_Client.Table("userinfo").Where("UserID = ?", UserID).First(&UserInfo)
	if result.Error != nil {
		return result.Error
	}

	if password != UserInfo.Password {
		return errors.New("密码错误！")
	}

	return nil

}

func (DAO *UserDAO) LoginByAccessToken(AccessToken string) error {
	return nil

}

func (DAO *UserDAO) LoginByRefreshToken(RefreshToken string) error {
	return nil
}