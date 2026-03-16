package DTO

type RAGContextID struct {
	ContextID string `json:"ContextID"`
}

type RAGAsk struct {
	Question  string `json:"Question"`
	ContextID string `json:"ContextID"`
}
