package UserModel

import "gorm.io/gorm"

type User struct {
	UserID    int64  `gorm:"type:int64;not null;column:UserID"`
	Nickname  string `gorm:"type:varchar(255);not null;column:Nickname"`
	Password  string `gorm:"type:varchar(255);not null;column:PasswordSecurity"`
	Salt      string `gorm:"type:varchar(255);not null;column:Salt"`
	Email     string `gorm:"type:varchar(255);not null;column:Email"`
	PublicKey string `gorm:"type:varchar(255);not null;column:PublicKey"`
	gorm.Model
}

func (User) TableName() string {
	return "userInfo"
}
