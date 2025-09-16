package NotAuthDTO

type WebAuthnRegisterDTO struct {
	UserID int64 `json:"UserID"`

	Email string `json:"Email"`
}
