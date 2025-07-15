package Model

import (
	"gorm.io/gorm"
	"time"
)

type Team struct {
	gorm.Model
	TeamName     string `gorm:"type:varchar(255);not null;column:TeamName"`
	NoteID       uint   `gorm:"type:varchar(255);not null;column:NoteID"`
	LeaderID     int    `gorm:"type:int;not null;column:LeaderID"`
	Introduction string `gorm:"type:varchar(255);not null;column:Introduction"`
}

type Members struct {
	gorm.Model
	MemberID     int       `gorm:"type:int;not null;column:MemberID"`
	JoinTime     time.Time `gorm:"type:datetime;not null;column:joinTime"`
	CommitCount  uint      `gorm:"type:int;not null;column:CommitCount"`
	Introduction string    `gorm:"type:varchar(255);not null;column:Introduction"`
}

type JoinTeamApplication struct {
	gorm.Model
	ApplicationID       int    `gorm:"type:varchar(255);not null;column:ApplicationID"`
	ApplicationNickName string `gorm:"type:varchar(255);not null;column:ApplicationNickName"`
	TeamID              int    `gorm:"type:varchar(255);not null;column:TeamID"`
	Introduction        string `gorm:"type:varchar(255);not null;column:Introduction"`
	Status              int    `gorm:"type:varchar(255);not null;column:Status"`
}

func (JoinTeamApplication) TableName() string { return "JoinTeamApplication" }
func (Team) TableName() string                { return "Team" }
func (Members) TableName() string             { return "Members" }
