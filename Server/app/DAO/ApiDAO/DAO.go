package ApiDAO

import "gorm.io/gorm"

type ApiDAO struct {
	DB_Client *gorm.DB
}

func NewDAOFactory(db *gorm.DB) *ApiDAO {
	return &ApiDAO{
		DB_Client: db,
	}
}
