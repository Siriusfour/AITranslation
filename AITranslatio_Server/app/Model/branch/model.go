package branch

import "gorm.io/gorm"

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

func (Branch) TableName() string { return "Branch" }
