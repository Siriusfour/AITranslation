package UserModel

type User struct {
	UserID    int64  `gorm:"type:int64;not null;column:UserID"`
	Nickname  string `gorm:"type:varchar(255);not null;column:Nickname"`
	Password  string `gorm:"type:varchar(255);not null;column:Password"`
	Salt      string `gorm:"type:varchar(255);not null;column:Salt"`
	Email     string `gorm:"type:varchar(255);not null;column:Email"`
	PublicKey string `gorm:"type:varchar(255);not null;column:PublicKey"`
}

func (u *User) TableName() string {
	return "User"
}
