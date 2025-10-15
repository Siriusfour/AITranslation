package DTO

type JoinTeamDTO struct {
	FromUserID   int64  `json:"UserID"  binding:"required"`
	Introduction string `json:"Introduction"`
	NickName     string `json:"NickName" binding:"required"`
	TeamID       int    `json:"TeamID"   binding:"required"`
}

type CreateTeamDTO struct {
	UserID       int    `json:"UserID"  binding:"required"`
	TeamName     string `json:"NickName" binding:"required"`
	Introduction string `json:"Introduction"`
	NoteID       uint   `json:"NoteID"  binding:"required"`
}
