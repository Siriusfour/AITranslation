package AuthService

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/WebAuthn"
	WebAuthn_Verify "AITranslatio/Utils/WebAuthn/Verify"
	"AITranslatio/app/DAO/AuthDAO"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type CredentialOptions struct {
	Challenge  string               `json:"Challenge"`
	AllowCreds []AuthDAO.Credential `json:"AllowCreds"`
}

// 申请一个WebAuthn密钥
func (Service *AuthService) ApplicationWebAuthn(UserID int64) (*WebAuthn.WebAuthn, error) {

	//获取userName和Email
	var UseName, Email string
	err := AuthDAO.CreateDAOFactory("mysql").
		DB_Client.
		Raw("SELECT NickName, Email FROM User WHERE UserID = ?", UserID).
		Row().
		Scan(&UseName, &Email)
	if err != nil {
		return nil, err
	}

	//生成Challenge，并把其置于redis五分钟
	w := WebAuthn.CreateWebAuthnConfigFactory(UserID, UseName, Email)
	err = w.CreateChallenge()
	if err != nil {
		return nil, errors.New(MyErrors.ErrorGetChallengeIsFail + err.Error())
	}

	return w, nil
}

// 注册时验证WebAuthn密钥
func (Service *AuthService) VerifyWebAuthnToRegister(WebAuthnCtx *gin.Context) error {

	err := WebAuthn_Verify.ClientDataJsonVerifyForRegister(WebAuthnCtx)
	if err != nil {
		return fmt.Errorf("WebAuthn注册,ClientData校验错误:%w", err)
	}
	err = WebAuthn_Verify.AttestationObjectVerifyForRegister(WebAuthnCtx)
	if err != nil {
		return fmt.Errorf("WebAuthn注册,Attestation校验错误:%w", err)
	}
	return nil
}

// 登录时验证WebAuthn
func (Service *AuthService) WebAuthnToLogin(WebAuthnCtx *gin.Context) error {

	err := WebAuthn_Verify.ClientDataJsonVerifyForLogin(WebAuthnCtx)
	if err != nil {
		return fmt.Errorf("WebAuthnd登录,ClientData校验错误:%w", err)
	}
	err = WebAuthn_Verify.AttestationObjectVerifyForLogin(WebAuthnCtx)
	if err != nil {
		return fmt.Errorf("WebAuthn登录,Attestation校验错误:%w", err)
	}

	//SignCount++

	return nil
}

func (Service *AuthService) GetUserAllCredentialDTO(WebAuthnCtx *gin.Context) (*CredentialOptions, error) {
	allowCreds, err := AuthDAO.CreateDAOFactory("mysql").FindCredentialByUserID(WebAuthnCtx)
	if err != nil {
		return nil, err
	}
	credentialOptions := &CredentialOptions{
		Challenge:  WebAuthnCtx.GetString(Consts.ValidatorPrefix + "challenge"),
		AllowCreds: allowCreds,
	}
	return credentialOptions, nil
}
