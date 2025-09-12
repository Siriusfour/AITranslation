package DTO

import "AITranslatio/app/http/DTO/NotAuthDTO"

type NovelDTO struct {
	NotAuthDTO.Auth
	WriterID     int    `json:"WriterID"  binding:"required"`
	Introduction string `json:"Introduction"`
	NoteName     string `json:"NoteName" binding:"required"`
	Permissions  int    `json:"Permissions"  binding:"required"` // 1-公开所有人都能编辑，2-只有自己可以编辑
}

type Branch struct {
	NotAuthDTO.Auth
	WriterID     int    `json:"WriterID"  binding:"required"`
	Introduction string `json:"Introduction"`
	BranchName   string `json:"BranchName" binding:"required"`
	NoteID       int    `json:"NoteID"  binding:"required"`
	Permissions  int    `json:"Permissions"  binding:"required"`
	CommitID     int    `json:"CommitID"  binding:"required"`
}

type CommitDTO struct {
	NotAuthDTO.Auth
	CommitID     uint   `json:"CommitID"`
	WriterID     int    `json:"WriterID"  binding:"required"`
	Introduction string `json:"Introduction"`
	CommitName   string `json:"CommitName" binding:"required"`
	BranchID     int    `json:"BranchID" binding:"required"`
	LastNode     int    `json:"preNode"`
	NextNode     int    `json:"NextNode"`
}

type UserInfo struct {
	UserID int `json:"UserID"  binding:"required"`
}
