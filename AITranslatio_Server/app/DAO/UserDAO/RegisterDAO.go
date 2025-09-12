package UserDAO

import (
	"AITranslatio/Global"
	"AITranslatio/Utils/token"
	"AITranslatio/app/Model/UserModel"
	"AITranslatio/app/http/DTO/NotAuthDTO"
	"errors"
	"gorm.io/gorm"
)

func (DAO *UserDAO) Register(UserID int64, UserName, Email, EmailCode, HashPassword, Salt string) error {

	RegisterData := &UserModel.User{
		UserID:   UserID,
		Nickname: UserName,
		Password: HashPassword,
		Salt:     Salt,
		Email:    Email,
		Model:    gorm.Model{},
	}

	//在数据库增加用户
	result := DAO.DB_Client.Create(RegisterData)
	if result.Error != nil {
		return result.Error
	}

}
