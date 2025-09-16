package NotAuthDTO

type TranslationDTO struct {
	Message string `json:"Message" binding:"required" Message:"提问内容为空！！" required_err:"用户名不能为空！！" `
}

type LoginDTO struct {
	Auth
	UserID   int    `json:"UserID" binding:"required"`
	Password string `json:"PasswordSecurity"`
}

type Auth struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}
