package DTO

type TranslationDTO struct {
	Message string `json:"Message" binding:"required" Message:"提问内容为空！！" required_err:"用户名不能为空！！" `
}

type LoginDTO struct {
	UserID       string `json:"UserID" binding:"required"`
	Password     string `json:"Password" binding:"required"`
	AccessToken  string `json:"AccessToken" binding:"required"`
	RefreshToken string `json:"RefreshToken" binding:"required"`
}
