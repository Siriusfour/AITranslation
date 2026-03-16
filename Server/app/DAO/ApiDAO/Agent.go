package ApiDAO

import (
	"AITranslatio/app/Model/User"
	"gorm.io/datatypes"
)

func (DAO *ApiDAO) UpdateSessionList(UserID int64, data []byte) error {

	user := User.User{UserID: UserID}

	err := DAO.DB_Client.Model(&user).Update("SessionList", datatypes.JSON(data)).Error

	return err
}
