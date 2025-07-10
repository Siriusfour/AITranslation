package DTO

type JoinTeamDTO struct {
	Auth
	UserID       int    `json:"UserID"  binding:"required"`
	Introduction string `json:"Introduction"  binding:"required"`
	NickName     string `json:"NickName" binding:"required"`
}
