package UserDAO

import (
	"AITranslatio/app/Model/User"
	"AITranslatio/app/types"
)

func (DAO *UserDAO) Register(DTO *types.RegisterDTO) error {

	RegisterData := &User.User{
		UserID:   DTO.UserID,
		Nickname: DTO.UserName,
		Password: DTO.Password,
		Salt:     DTO.Salt,
		Email:    DTO.Email,
	}

	//在数据库增加用户
	result := DAO.DB_Client.Create(RegisterData)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
