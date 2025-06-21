package Model

import "gorm.io/gorm"

type User struct {
	UserID   int    `gorm:"type:int;not null;column:UserID"`
	Password string `gorm:"type:varchar(255);not null;column:Password"`
	gorm.Model
}

func (User) TableName() string {
	return "userInfo"
}
