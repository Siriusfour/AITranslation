package AuthService

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/WebAuthn"
	WebAuthn_Verify "AITranslatio/Utils/WebAuthn/Verify"
	"AITranslatio/app/DAO/AuthDAO"
	"AITranslatio/app/Model/User"
	"AITranslatio/app/types/DTO"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type CredentialOptions struct {
	Challenge  string               `json:"Challenge"`
	AllowCreds []AuthDAO.Credential `json:"AllowCreds"`
}

type RegisterWebAuthn struct {
	WebAuthn  *WebAuthn.WebAuthn
	UserName  string
	Challenge string
	UserID    int64
}

type LoginWebAuthn struct {
	Challenge string
	RPID      string
	RPName    string
}

// 获取服务端的webAuthn信息
func (s *AuthService) RegisterGetWebAuthnInfo(sessionID, userID int64) (*RegisterWebAuthn, error) {

	user, err := s.DAO.FindUserByID(userID, "UserID")
	if err != nil {
		return nil, err
	}

	challenge, err := s.WebAuthnGenerator.CreateChallenge(s.RedisClient, sessionID)
	if err != nil {
		return nil, errors.New(MyErrors.ErrorGetChallengeIsFail + err.Error())
	}

	rw := &RegisterWebAuthn{
		WebAuthn:  s.WebAuthnGenerator,
		UserName:  user.Nickname + "(" + user.Email + ")",
		UserID:    userID,
		Challenge: challenge,
	}

	return rw, nil
}

func (s *AuthService) RegisterWebAuthn(ctx *gin.Context, rw *DTO.RegisterWebAuthn) error {

	userID := ctx.GetInt64("UserID")
	sessionID := ctx.GetInt64("SessionID")

	err := WebAuthn_Verify.ClientDataJsonVerifyForRegister(s.cfg, s.RedisClient, rw.Response.ClientDataJSON, sessionID)
	if err != nil {
		return errors.New("ClientDataJson解析失败")
	}

	CredentialID, PublicKey, err := WebAuthn_Verify.AttestationObjectVerifyForRegister(s.cfg, rw.Response.AttestationObject)
	if err != nil {
		return errors.New("AttestationObject解析失败")
	}

	//公钥入库
	err = s.DAO.CreateCredential(userID, CredentialID, PublicKey)
	if err != nil {
		return errors.New("mysql创建凭证失败")
	}

	return nil
}

// 获取服务端的webAuthn信息
func (s *AuthService) LoginGetWebAuthnInfo(sessionID int64) (*LoginWebAuthn, error) {

	challenge, err := s.WebAuthnGenerator.CreateChallenge(s.RedisClient, sessionID)
	if err != nil {
		return nil, errors.New(MyErrors.ErrorGetChallengeIsFail + err.Error())
	}

	lw := &LoginWebAuthn{
		Challenge: challenge,
		RPID:      s.WebAuthnGenerator.RPID,
		RPName:    s.WebAuthnGenerator.RPName,
	}

	return lw, nil
}

// 登录时验证WebAuthn
func (s *AuthService) LoginByWebAuthn(clientDataJSON, AttestationObject, Signature, CredentialID string, sessionID, userID int64) (*User.User, string, string, error) {

	err := WebAuthn_Verify.ClientDataJsonVerifyForLogin(s.cfg, s.RedisClient, clientDataJSON, sessionID)
	if err != nil {
		return nil, "", "", fmt.Errorf("ClientData校验错误:%w", err)
	}
	err = WebAuthn_Verify.AttestationObjectVerifyForLogin(s.cfg, s.DAO, clientDataJSON, AttestationObject, Signature, CredentialID)
	if err != nil {
		return nil, "", "", fmt.Errorf("Attestation校验错误:%w", err)
	}

	//SignCount++

	//获取到userInfo
	userInfo, err := s.DAO.FindUserByID(userID, "UserID")
	if err != nil {
		return nil, "", "", fmt.Errorf("查库失败:%w", err)
	}

	//验证通过，生成ak，rk ，写入redis，返回请求
	AccessToken, errAk := s.TokenProvider.GeneratedToken(userID, Consts.AccessToken)
	RefreshToken, errRk := s.TokenProvider.GeneratedToken(userID, Consts.RefreshToken)

	if errAk != nil || errRk != nil {
		return nil, "", "", fmt.Errorf("：验证成功但生成token失败：%w,%w", errAk, errRk)
	}

	return userInfo, AccessToken, RefreshToken, nil
}
