package DTO

type JoinTeamDTO struct {
	Auth
	UserID       int    `json:"UserID"  binding:"required"`
	Introduction string `json:"Introduction"`
	NickName     string `json:"NickName" binding:"required"`
	JoinTeamID   int    `json:"JoinTeamID"  binding:"required"`
}
