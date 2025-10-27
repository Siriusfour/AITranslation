package webAuthn

import "gorm.io/gorm"

type WebAuthnCredential struct {
	gorm.Model
	UserID       int64  `gorm:"type:BIGINT;not null;column:user_id;index"`
	CredentialID []byte `gorm:"type:BLOB;not null;column:credential_id;uniqueIndex:idx_credential_id,length:255;primaryKey"`
	PublicKey    []byte `gorm:"type:BLOB;not null;column:public_key"`
	SignCount    uint32 `gorm:"type:INT UNSIGNED;not null;default:0;column:sign_count"`
	Type         string `gorm:"type:VARCHAR(16);not null;default:'public-key';column:type"`
}

func (C *WebAuthnCredential) TableName() string {
	return "WebAuthnCredential"
}
