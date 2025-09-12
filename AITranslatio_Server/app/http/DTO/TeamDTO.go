package DTO

type JoinTeamDTO struct {
	UserID       int    `json:"UserID"  binding:"required"`
	Introduction string `json:"Introduction"`
	NickName     string `json:"NickName" binding:"required"`
	JoinTeamID   int    `json:"JoinTeamID"  binding:"required"`
}

type CreateTeamDTO struct {
	UserID       int    `json:"UserID"  binding:"required"`
	TeamName     string `json:"NickName" binding:"required"`
	Introduction string `json:"Introduction"`
	NoteID       uint   `json:"NoteID"  binding:"required"`
}
