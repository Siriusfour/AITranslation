package Model

import (
	"gorm.io/gorm"
	"time"
)

type Note struct {
	gorm.Model
	WriterID     int    `gorm:"type:int;not null;column:WriterID"`          // 作者ID
	NoteName     string `gorm:"type:varchar(255);not null;column:NoteName"` //小说名
	BranchCount  int    `gorm:"type:int;not null;column:BranchCount"`       // 分支数
	Introduction string `gorm:"size:255;not null;column:Introduction"`      // 简介
	ReaderCount  int    `gorm:"type:int;not null;column:ReaderCount"`       // 订阅人数
	Permissions  int    `gorm:"type:uint;not null;column:Permissions"`
}

type Branch struct {
	gorm.Model
	BranchName   string `gorm:"type:varchar(255);not null;column:BranchName"`
	WriterID     int    `gorm:"type:int;not null;column:WriterID"`
	Introduction string `gorm:"type:varchar(255);not null;column:Introduction"`
	NoteID       uint   `gorm:"type:int;not null;column:NoteID"`
	Permissions  int    `gorm:"type:int;not null;column:Permissions"`
	CommitID     int    `gorm:"type:int;not null;column:CommitID"`
	ReaderCount  int    `gorm:"type:int;not null;column:ReaderCount"`
}

type Commit struct {
	gorm.Model
	CommitName   string `gorm:"type:varchar(255);not null;column:CommitName"`
	WriterID     int    `gorm:"type:int;not null;column:WriterID"`
	Introduction string `gorm:"type:varchar(255);not null;column:Introduction"`
	BranchID     int    `gorm:"type:int;not null;column:BranchID"`
	FilePath     string `gorm:"type:varchar(255);not null;column:FilePath"`
	LastNode     int    `gorm:"type:int;not null;column:LastNode"`
	NextNode     int    `gorm:"type:int;not null;column:NextNode"`
}

type Team struct {
	gorm.Model
	TeamName     string `gorm:"type:varchar(255);not null;column:TeamName"`
	NoteID       uint   `gorm:"type:varchar(255);not null;column:NoteID"`
	LeaderID     int    `gorm:"type:int;not null;column:LeaderID"`
	Introduction string `gorm:"type:varchar(255);not null;column:Introduction"`
}

type Members struct {
	TeamID       int       `gorm:"type:int;not null;column:TeamID"`
	MemberID     int       `gorm:"type:int;not null;column:MemberID"`
	JoinTime     time.Time `gorm:"type:datetime;not null;column:joinTime"`
	CommitCount  uint      `gorm:"type:int;not null;column:CommitCount"`
	Introduction string    `gorm:"type:varchar(255);not null;column:Introduction"`
}

func (Note) TableName() string    { return "Note" }
func (Branch) TableName() string  { return "Branch" }
func (Commit) TableName() string  { return "Commit" }
func (Team) TableName() string    { return "Team" }
func (Members) TableName() string { return "Members" }
