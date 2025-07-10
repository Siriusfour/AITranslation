package BaseDAO

import (
	"AITranslatio/Global"
	"AITranslatio/Src/Model"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseDAO struct {
	orm    *gorm.DB
	Logger *zap.SugaredLogger
}

func New_Base_DAO() *BaseDAO {
	fmt.Println("========= Global.Logger:", Global.Logger)
	return &BaseDAO{
		orm:    Global.DB,
		Logger: Global.Logger,
	}
}

func (BaseDAO *BaseDAO) LoginByPassword(UserID int, password string) error {

	var UserInfo Model.User

	fmt.Println(BaseDAO)

	result := BaseDAO.orm.Table("userinfo").Where("UserID = ?", UserID).First(&UserInfo)
	if result.Error != nil {

		return result.Error
	}

	if password != UserInfo.Password {
		return errors.New("密码错误！")
	}

	return nil

}

func (BaseDAO *BaseDAO) FindUserInfo(UserID string) error {
	return nil
}
