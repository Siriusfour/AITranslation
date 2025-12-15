package DTO

type Auth struct {
	AccessToken  string `json:"AccessToken,omitempty"`
	RefreshToken string `json:"RefreshToken,omitempty"`
}
