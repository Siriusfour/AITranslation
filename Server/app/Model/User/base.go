package User

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID     int64  `gorm:"type:BIGINT;not null;column:UserID;primaryKey"`
	Nickname   string `gorm:"type:varchar(255);not null;column:NickName"`
	Password   string `gorm:"type:varchar(255);not null;column:Password"`
	Salt       string `gorm:"type:varchar(255);not null;column:Salt"`
	Email      string `gorm:"type:varchar(255);not null;column:Email"`
	Credential bool   `gorm:"type:boolean;not null;column:credential"`
}

func (u *User) TableName() string {
	return "user"
}
