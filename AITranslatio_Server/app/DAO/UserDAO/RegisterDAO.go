package UserDAO

import (
	"AITranslatio/app/Model/UserModel"
	"AITranslatio/app/http/DTO/NotAuthDTO"
	"gorm.io/gorm"
)

func (DAO *UserDAO) Register(DTO *NotAuthDTO.RegisterDTO) error {

	RegisterData := &UserModel.User{
		UserID:    DTO.UserID,
		Nickname:  DTO.UserName,
		Password:  DTO.Password,
		Salt:      DTO.Salt,
		Email:     DTO.Email,
		Model:     gorm.Model{},
		PublicKey: DTO.PublicKey,
	}

	//在数据库增加用户
	result := DAO.DB_Client.Create(RegisterData)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
