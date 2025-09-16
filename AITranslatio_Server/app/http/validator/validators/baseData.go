package validators

type BaseField struct {
	UserID   int
	PassWord string
	Auth
}

type Auth struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}
