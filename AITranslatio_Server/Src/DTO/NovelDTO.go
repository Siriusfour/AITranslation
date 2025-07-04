package DTO

type NovelDTO struct {
	UserID          string `json:"UserID"  binding:"required" `
	AccessToken     string `json:"AccessToken" `
	RefreshToken    string `json:"RefreshToken" `
	ProgrammingName string `json:"ProgrammingNameID" `
	Introduction    string `json:"Introduction" `
}
