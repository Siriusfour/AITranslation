package Model

type User struct {
	UserID   int    `gorm:"type:int;not null"`
	NickName string `gorm:"type:varchar(20);not null"`
	Password string `gorm:"type:varchar(255);not null"`
}
