package Reponse

type Auth struct {
	AccessToken  string `json:"AccessToken,omitempty"`
	RefreshToken string `json:"RefreshToken,omitempty"`
}

type LoginInfo struct {
	Auth Auth

	Nickname string
	UserID   int64
	Avatar   string
}
