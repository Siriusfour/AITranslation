package DTO

type Login struct {
	Password string `json:"password"`
	Email    string `json:"email" `
}

type OAuth struct {
	State string `json:"state" binding:"required"`
	Code  string `json:"code" binding:"required"`
}
type Challenge struct {
	Challenge string
}

type RegisterWebAuthn struct {
	ID         string   `json:"id" binding:"required"`
	RawID      string   `json:"rawId" binding:"required"`
	Type       string   `json:"type" binding:"required"`
	Transports []string `json:"transports"`

	Response struct {
		ClientDataJSON    string `json:"clientDataJSON" binding:"required"`
		AttestationObject string `json:"attestationObject" binding:"required"`
	} `json:"response" binding:"required"`
}

type LoginWebAuthn struct {
	ID    string `json:"id" binding:"required"`
	RawID string `json:"rawId" binding:"required"`
	Type  string `json:"type" binding:"required"`

	Response struct {
		ClientDataJSON    string `json:"clientDataJSON" binding:"required"`
		AuthenticatorData string `json:"AuthenticatorData" binding:"required"`
		Signature         string `json:"signature" binding:"required"`
		UserHandle        string `json:"userHandle" binding:"required"`
	} `json:"response" binding:"required"`
}
