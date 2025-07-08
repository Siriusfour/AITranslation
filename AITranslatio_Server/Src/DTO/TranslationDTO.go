package DTO

type TranslationDTO struct {
	Message string `json:"Message" binding:"required" Message:"提问内容为空！！" required_err:"用户名不能为空！！" `
}

type LoginDTO struct {
	UserID       string `json:"UserID"`
	UUID         string `json:"UUID"` //设备唯一ID
	Password     string `json:"Password" `
	AccessToken  string `json:"AccessToken" `
	RefreshToken string `json:"RefreshToken" `
}

type Auth struct {
	AccessToken  string `json:"AccessToken" `
	RefreshToken string `json:"RefreshToken" `
}
