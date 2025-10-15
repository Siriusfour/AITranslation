package Programming

import (
	"gorm.io/gorm"
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

func (Commit) TableName() string { return "Commit" }
