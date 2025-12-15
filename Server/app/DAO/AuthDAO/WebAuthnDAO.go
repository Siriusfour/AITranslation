package AuthDAO

import (
	"AITranslatio/app/Model/webAuthn"
	"encoding/base64"
	"fmt"
)

type Credential struct {
	CredentialID []byte `gorm:"column:credential_id"`
	Type         string `gorm:"column:type"`
}

// 查找某一条凭证
func (DAO *AuthDAO) FindCredential(CredentialID string) (*webAuthn.WebAuthnCredential, error) {

	var webAuthnCredential *webAuthn.WebAuthnCredential

	CredentialID_ByteArray, err := base64.RawURLEncoding.DecodeString(CredentialID)
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %w", err)
	}
	result := DAO.DB_Client.Table("webauthncredential").Where("credential_id = ?", CredentialID_ByteArray).First(&webAuthnCredential)
	if result.Error != nil {
		return nil, fmt.Errorf("FindCredential失败: %w", result.Error)
	}

	return webAuthnCredential, nil

}

// 创建凭证
func (DAO *AuthDAO) CreateCredential(userID int64, CredentialID, publicKey []byte) error {

	webAuthnCredential := &webAuthn.WebAuthnCredential{
		UserID:       userID,
		SignCount:    1,
		PublicKey:    publicKey,
		CredentialID: CredentialID,
	}

	result := DAO.DB_Client.Table("webauthncredential").Create(webAuthnCredential)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// 根据UserID查找其所有的CredentialID、type
func (DAO *AuthDAO) FindCredentialByUserID(UserID int64) ([]Credential, error) {
	var credentials []Credential

	result := DAO.DB_Client.
		Table("webauthncredential").
		Select("credential_id", "type").
		Where("user_id = ?", UserID).
		Find(&credentials)

	if result.Error != nil {
		return nil, fmt.Errorf("gorm Find 失败: %w", result.Error)
	}

	return credentials, nil
}
