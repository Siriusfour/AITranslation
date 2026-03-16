package User

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID      int64          `gorm:"type:BIGINT;not null;column:UserID;primaryKey"`
	Nickname    string         `gorm:"type:varchar(255);not null;column:NickName"`
	Avatar      string         `gorm:"type:blob;not null;column:Avatar"`
	Password    string         `gorm:"type:varchar(255);column:Password"`
	Salt        string         `gorm:"type:varchar(255);column:Salt"`
	Email       string         `gorm:"type:varchar(255);not null;column:Email"`
	Credential  bool           `gorm:"type:boolean;not null;column:credential"`
	Admin       bool           `gorm:"type:boolean;not null;column:Admin"`
	GithubID    int64          `gorm:"type:int;column:GithubId"`
	SessionList datatypes.JSON `gorm:"type:json;column:SessionList"`
}

type SessionInfo struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

func (u *User) TableName() string {
	return "user"
}
