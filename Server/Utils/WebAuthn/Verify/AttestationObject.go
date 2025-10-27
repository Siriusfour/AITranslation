package Verify

import (
	"AITranslatio/Global"
	"AITranslatio/Global/Consts"
	"AITranslatio/app/DAO/UserDAO"
	"AITranslatio/app/http/reposen"
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/fxamacker/cbor/v2"
	"github.com/gin-gonic/gin"
	"math/big"
)

type AttestationObject struct {
	Fmt      string                 `cbor:"fmt"`
	AuthData []byte                 `cbor:"authData"`
	AttStmt  map[string]interface{} `cbor:"attStmt"`
}

type AuthenticatorData struct {
	RPIDHash               []byte // 32 bytes
	Flags                  byte
	SignCount              uint32
	AttestedCredentialData *AttestedCredentialData
	Extensions             []byte
	Raw                    []byte // 原始字节
}

type AttestedCredentialData struct {
	AAGUID              []byte
	CredentialID        []byte
	CredentialPublicKey []byte
}

func AttestationObjectVerifyForRegister(WebAuthnCtx *gin.Context) error {

	AttestationObjectBase64 := WebAuthnCtx.GetString(Consts.ValidatorPrefix + "AttestationObject")

	//0.解析Base64
	AttestationObjectJSON, err := base64.RawURLEncoding.DecodeString(AttestationObjectBase64) //Base64字符串->json字节数组
	if err != nil {
		return fmt.Errorf("base64URL解码错误: %v", err)
	}

	var attestationObject AttestationObject
	if err := cbor.Unmarshal(AttestationObjectJSON, &attestationObject); err != nil {
		return fmt.Errorf("AttestationObject CBOR解析失败: %w", err)
	}

	// ===== 0. 提取 authData 各部分 ,RPIDHash,flags,signCount=====
	authData := attestationObject.AuthData

	if len(authData) < 37 {
		return fmt.Errorf("AuthData 太短")
	}

	// 提取各部分
	// =====  验证 RP ID Hash 是否与配置文件的rpID的hash值相同
	RPIDHash := authData[0:32]
	flags := authData[32]
	signCount := binary.BigEndian.Uint32(authData[33:37])
	UserID := WebAuthnCtx.GetInt64(Consts.ValidatorPrefix + "UserID")
	attestedCredentialData, err := extractCredentialData(authData)
	if err != nil {
		return fmt.Errorf("public解析失败: %w", err)
	}

	RPID := Global.Config.GetString("WebAuthn.RPID")
	if err := verifyRPID(RPID, RPIDHash); err != nil {
		return err
	}

	if err = verifyFlags(flags); err != nil {
		return err
	}

	if attestationObject.Fmt == "none" {
		if len(attestationObject.AttStmt) != 0 {
			return errors.New("格式有误！Fmt为none时AttStmt应该为零")
		}
	}

	//数据库插入数据
	err = UserDAO.CreateDAOFactory("mysql").CreateCredential(UserID, signCount, attestedCredentialData.CredentialID, attestedCredentialData.CredentialPublicKey)
	if err != nil {
		return fmt.Errorf("DAO层CreateCredential调用失败: %w", err)
	}

	reposen.OK(WebAuthnCtx, nil)
	return nil
}

func AttestationObjectVerifyForLogin(WebAuthnCtx *gin.Context) error {

	AuthenticatorDataBase64 := WebAuthnCtx.GetString(Consts.ValidatorPrefix + "AttestationObject")

	//0.解析Base64
	authData, err := base64.RawURLEncoding.DecodeString(AuthenticatorDataBase64) //Base64字符串->字节数组
	if err != nil {
		return fmt.Errorf("base64URL解码错误: %v", err)
	}

	// ===== 0. 提取 authData 各部分 ,RPIDHash,flags,signCount=====

	// 提取各部分
	RPIDHash := authData[0:32]
	flags := authData[32]
	signCount := binary.BigEndian.Uint32(authData[33:37]) //把4字节按照大端序转化成一个uint32整数

	//由UserID和Credential ID数据库提取公钥
	Credential, err := UserDAO.CreateDAOFactory("mysql").FindCredential(WebAuthnCtx)
	if err != nil {
		return fmt.Errorf("webAuthn中根据CredentialID查找Credential失败： %w", err)
	}

	// =====  验证 RP ID Hash 是否与配置文件的rpID的hash值相同
	RPID := Global.Config.GetString("WebAuthn.RPID")
	if err := verifyRPID(RPID, RPIDHash); err != nil {
		return err
	}

	//验证flag位
	if err = verifyFlags(flags); err != nil {
		return err
	}

	//验证signCount
	if err = verifySignCount(signCount, Credential.SignCount); err != nil {
		return err
	}

	//验证签名
	ClientData, err := base64.RawURLEncoding.DecodeString(WebAuthnCtx.GetString(Consts.ValidatorPrefix + "ClientDataJSON"))
	if err != nil {
		return err
	}

	clientDataHash := sha256.Sum256(ClientData)
	messageHash := append(authData, clientDataHash[:]...)
	signature, err := base64.RawURLEncoding.DecodeString(WebAuthnCtx.GetString(Consts.ValidatorPrefix + "Signature"))
	if err != nil {
		return err
	}

	publicKey, alg, err := parseCOSEPublicKey(Credential.PublicKey)
	if err != nil {
		return err
	}

	// 3) 按算法验证
	switch k := publicKey.(type) {
	case *ecdsa.PublicKey: // ES256 (-7)
		var d = messageHash
		// 首选 DER
		if signature[0] == 0x30 {
			if !ecdsa.VerifyASN1(k, d[:], signature) {
				return errors.New("ecdsa verify failed")
			}
		} else if len(signature) == 64 { // 兜底 raw r||s
			r := new(big.Int).SetBytes(signature[:32])
			s := new(big.Int).SetBytes(signature[32:])
			if !ecdsa.Verify(k, d[:], r, s) {
				return errors.New("ecdsa raw verify failed")
			}
		} else {
			return errors.New("unexpected ECDSA signature format")
		}

	case ed25519.PublicKey: // Ed25519 (-8)
		if !ed25519.Verify(k, signature, signature) {
			return errors.New("ed25519 verify failed")
		}

	case *rsa.PublicKey: // RS256 (-257)
		d := sha256.Sum256(signature)
		if err := rsa.VerifyPKCS1v15(k, crypto.SHA256, d[:], signature); err != nil {
			return fmt.Errorf("rsa verify failed: %w", err)
		}

	default:
		return fmt.Errorf("unsupported key type/alg=%d", alg)
	}

	return nil
}

func verifyRPID(rpID string, rpIdHash []byte) error {
	expectedRPIDHash := sha256.Sum256([]byte(rpID))
	if !bytes.Equal(expectedRPIDHash[:], rpIdHash) {
		return errors.New("RPIDHash不一致")
	}
	return nil
}

// Bit 0 (0x01): UP - User Present（用户在场）
// Bit 2 (0x04): UV - User Verified（用户验证/生物识别）
// Bit 6 (0x40): AT - Attested Credential Data（有凭证数据，仅注册时）
// Bit 7 (0x80): ED - Extension Data（有扩展数据
func verifyFlags(flags byte) error {

	RequireUserVerification := Global.Config.GetBool("WebAuthn.RequireUserVerification")

	// ===== 1. 验证 User Present (UP) - 必须检查 =====
	if (flags & 0x01) == 0 {
		return errors.New("User Present (UP) 标志未设置 - 用户不在场")
	}

	// ===== 2. 验证 User Verified (UV) - 根据策略 =====
	if RequireUserVerification {
		if (flags & 0x04) == 0 {
			return errors.New("User Verified (UV) 标志未设置 - 用户未验证")
		}
	}

	// ===== 3. 验证 Attested Credential Data (AT) - 注册时必须有 =====
	/*	if isRegistration {
		if (flags & 0x40) == 0 {
			return errors.New("Attested Credential Data (AT) 标志未设置 - 缺少凭证数据")
		}
		fmt.Println("✅ Attested Credential Data (AT) 存在")
	}*/

	return nil
}
func verifySignCount(currentCount, storedCount uint32) error {
	if currentCount != 0 && storedCount != 0 && currentCount <= storedCount {
		return fmt.Errorf("SignCount 错误，currentCount 值为 %d，storedCount 值为 %d", currentCount, storedCount)
	}
	return nil
}
