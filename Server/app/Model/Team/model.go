package Team

import (
	"gorm.io/gorm"
	"time"
)

type Team struct {
	gorm.Model
	TeamName     string `gorm:"type:varchar(255);not null;column:TeamName"`
	NoteID       uint   `gorm:"type:varchar(255);not null;column:NoteID"`
	LeaderID     int64  `gorm:"type:int;not null;column:LeaderID"`
	Introduction string `gorm:"type:varchar(255);not null;column:Introduction"`
	TeamID       uint   `gorm:"type:int;not null;column:TeamID;primaryKey;autoIncrement"`
}

type Members struct {
	gorm.Model
	UserID       int       `gorm:"type:int;not null;column:UserID"`
	JoinTime     time.Time `gorm:"type:datetime;not null;column:joinTime"`
	CommitCount  uint      `gorm:"type:int;not null;column:CommitCount"`
	Introduction string    `gorm:"type:varchar(255);not null;column:Introduction"`
}

type JoinTeamApplication struct {
	gorm.Model
	ID           uint   `gorm:"type:int;not null;column:ID;primaryKey;autoIncrement"`
	FromUserID   int64  `gorm:"type:BIGINT;not null;column:FromUserID"`
	TeamID       int    `gorm:"type:varchar(255);not null;column:TeamID"`
	Introduction string `gorm:"type:varchar(255);not null;column:Introduction"`
	Status       int    `gorm:"type:int;not null;column:Status"`
}

func (JoinTeamApplication) TableName() string { return "Join_Team_Application" }
func (Team) TableName() string                { return "Team" }
func (Members) TableName() string             { return "Members" }
