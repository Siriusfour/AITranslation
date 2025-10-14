package validators

type BaseField struct {
	UserID   int
	PassWord string
	AuthDTO
}

type AuthDTO struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}
