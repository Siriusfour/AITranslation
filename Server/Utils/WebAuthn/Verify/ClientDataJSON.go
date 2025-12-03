package Verify

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	"context"
	"github.com/redis/go-redis/v9"
	"time"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ClientData struct {
	Type        string `json:"type"`
	Challenge   string `json:"challenge"`
	Origin      string `json:"origin"`
	CrossOrigin bool   `json:"crossOrigin,omitempty"`
}

type User struct {
	TokenID          int64  `redis:"TokenID"`
	UserID           string `redis:"UserID"`
	RegisterTime     int64  `redis:"RegisterTime"`
	Challenge        string `redis:"challenge"`
	challengeOutTime int64  `redis:"challenge_OutTime"`
}

// ClientDataJsonVerifyForRegister 验证标准详见  https://w3c.github.io/webauthn/#dom-publickeycredentialcreationoptions-rp
// 处理注册新凭证时第 3 到第 6 步中对客户端数据的验证
//func ClientDataJsonVerifyForRegister(WebAuthnCtx *gin.Context) error {
//
//	ClientDataJSON := WebAuthnCtx.GetString(Consts.ValidatorPrefix + "ClientDataJSON")
//
//	//0.解析Base64
//	clientDataJSON, err := base64.RawURLEncoding.DecodeString(ClientDataJSON) //base64字符串->json字节数组
//	if err != nil {
//		return fmt.Errorf("base64URL解码错误: %v", err)
//	}
//
//	//计算ClinetDataJSON部分的hash，后面的attestationObject验证会用到（非none验证策略）
//
//	var clientData ClientData
//	if err := json.Unmarshal(clientDataJSON, &clientData); err != nil { //json字节数组 -> GO结构体
//		return fmt.Errorf("clientDataJSON解析失败: %w", err)
//	}
//	WebAuthnCtx.Set(Consts.ValidatorPrefix+"RPID", clientData.Origin)
//
//	//1.检查type是否为webauthn.create
//	//协议中的Registration Step 3. and Assertion Step 7.
//	if clientData.Type != "webauthn.create" {
//		return MyErrors.ErrorClientDataTypeIsFail
//	}
//
//	// 协议中的Assertion Step 8 and Registration Step 4.
//	//2.检查challenge是否匹配，
//	var userData User
//
//	UserID := WebAuthnCtx.GetInt64(Consts.ValidatorPrefix + "UserID")
//	err = Global.RedisClient.HGetAll(context.Background(), fmt.Sprintf("%d", UserID)).Scan(&userData)
//	if err != nil {
//		return fmt.Errorf("ClinetDataJSON_Verify获取Reids数据失败: %v", err)
//	}
//
//	if time.Now().Unix() <= userData.challengeOutTime && userData.Challenge == clientData.Challenge { //判断是否失效
//		return MyErrors.ErrorClientDataChallengeIsFail
//	}
//
//	//检查请求页面的域名是否在允许范围内，防止用户被人钓鱼 协议中的 Step 5 & Assertion Step 9
//	if clientData.Origin != cfg.GetString("WebAuthn.rp.ID") {
//		return MyErrors.ErrorClientDataRPID_IsFail
//	}
//
//	return nil
//}

func ClientDataJsonVerifyForLogin(cfg interf.ConfigInterface, redisClient *redis.Client, WebAuthnCtx *gin.Context) error {

	ClientDataJSON := WebAuthnCtx.GetString(Consts.ValidatorPrefix + "ClientDataJSON")

	//0.解析Base64
	clientDataJSON, err := base64.RawURLEncoding.DecodeString(ClientDataJSON) //base64字符串->json字节数组
	if err != nil {
		return fmt.Errorf("base64URL解码错误: %v", err)
	}

	var clientData ClientData
	if err := json.Unmarshal(clientDataJSON, &clientData); err != nil { //json字节数组 -> GO结构体
		return fmt.Errorf("clientDataJSON解析失败: %w", err)
	}

	//1.检查type是否为webauthn.create
	if clientData.Type != "webauthn.get" {
		return MyErrors.ErrorClientDataTypeIsFail
	}

	//2.检查challenge是否匹配，
	var userData User

	UserID := WebAuthnCtx.GetInt64(Consts.ValidatorPrefix + "UserID")
	err = redisClient.HGetAll(context.Background(), fmt.Sprintf("%d", UserID)).Scan(&userData)
	if err != nil {
		return fmt.Errorf("ClinetDataJSON_Verify获取Reids数据失败: %v", err)
	}

	if time.Now().Unix() <= userData.challengeOutTime && userData.Challenge == clientData.Challenge { //判断是否失效
		return MyErrors.ErrorClientDataChallengeIsFail
	}

	//检查请求页面的域名是否在允许范围内，防止用户被人钓鱼 协议中的 Step 5 & Assertion Step 9
	if clientData.Origin != cfg.GetString("WebAuthn.rp.ID") {
		return MyErrors.ErrorClientDataRPID_IsFail
	}

	return nil
}
