package validators

type BaseField struct {
	UserID       int
	PassWord     string
	AccessToken  string
	RefreshToken string
}

type Auth struct {
	AccessToken  string
	RefreshToken string
}
