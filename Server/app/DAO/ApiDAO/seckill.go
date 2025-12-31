package ApiDAO

import (
	"AITranslatio/Global/MyErrors"
	"AITranslatio/app/Model/User"
	"AITranslatio/app/Model/goods"
	"errors"
	"gorm.io/gorm"
)

func (DAO *ApiDAO) FindUserByID(ID int64, IDtype string) (*User.User, error) {

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

func (DAO *ApiDAO) FindSeckillGoods(id int64) (*goods.SeckillGoods, error) {

	var Goods *goods.SeckillGoods
	err := DAO.DB_Client.Find(&Goods).Where("id = ?", id).Error
	return Goods, err

}

func (DAO *ApiDAO) FindAllSeckillGoods() ([]*goods.SeckillGoods, error) {
	var allGoods []*goods.SeckillGoods
	// GORM 会自动填充
	err := DAO.DB_Client.Find(&allGoods).Error
	return allGoods, err
}

func (DAO *ApiDAO) CreateSeckillOrder(seckillOrder goods.SeckillOrder) error {
	result := DAO.DB_Client.Create(&seckillOrder)
	return result.Error
}
