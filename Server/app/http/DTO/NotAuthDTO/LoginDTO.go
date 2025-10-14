package NotAuthDTO

type TranslationDTO struct {
	Message string `json:"Message" binding:"required" Message:"提问内容为空！！" required_err:"用户名不能为空！！" `
}

type LoginDTO struct {
	Auth
	UserName int    `json:"UserName" binding:"required"`
	Password string `json:"Password" binding:"required"`
}

type Auth struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}
