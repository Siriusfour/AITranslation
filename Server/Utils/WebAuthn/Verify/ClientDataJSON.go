package Verify

import (
	"AITranslatio/Config/interf"
	"AITranslatio/Global/MyErrors"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"

	"encoding/base64"
	"encoding/json"
	"fmt"
)

type ClientData struct {
	Type        string `json:"type"`
	Challenge   string `json:"challenge"`
	Origin      string `json:"origin"`
	CrossOrigin bool   `json:"crossOrigin,omitempty"`
}

type User struct {
	Challenge        string `redis:"webAuthn_challenge"`
	ChallengeOutTime string `redis:"webAuthn_challenge_OutTime"`
}

// ClientDataJsonVerifyForRegister 验证标准详见  https://w3c.github.io/webauthn/#dom-publickeycredentialcreationoptions-rp
// 处理注册新凭证时第 3 到第 6 步中对客户端数据的验证
func ClientDataJsonVerifyForRegister(cfg interf.ConfigInterface, redisClient *redis.Client, clientDataJSON string, SessionID int64) error {

	//0.解析Base64
	c, err := base64.RawURLEncoding.DecodeString(clientDataJSON) //json字符串->二进制字节数组
	if err != nil {
		return fmt.Errorf("base64URL解码错误: %v", err)
	}

	//计算ClinetDataJSON部分的hash，后面的attestationObject验证会用到（非none验证策略）

	var clientData ClientData
	if err := json.Unmarshal(c, &clientData); err != nil { //json字节数组 -> GO结构体
		return fmt.Errorf("JSON反序列化失败: %w", err)
	}

	//1.检查type是否为webauthn.create
	//协议中的Registration Step 3. and Assertion Step 7.
	if clientData.Type != "webauthn.create" {
		return errors.New("type不为webauthn.create！")
	}

	// 协议中的Assertion Step 8 and Registration Step 4.
	//2.检查challenge是否匹配，
	var userData User

	err = redisClient.HGetAll(context.Background(), fmt.Sprintf("SessionID:%d", SessionID)).Scan(&userData)
	if err != nil {
		return fmt.Errorf("Reids HGetAll失败: %v", err)
	}

	outTime, _ := strconv.ParseInt(userData.ChallengeOutTime, 10, 64)

	if userData.Challenge != clientData.Challenge { //判断是否失效
		return errors.New("challenge校验失败！")
	}

	fmt.Println(time.Now().Unix())
	if time.Now().Unix() > outTime {
		return errors.New("challenge过期！")
	}

	//检查请求页面的域名是否在允许范围内，防止用户被人钓鱼 协议中的 Step 5 & Assertion Step 9
	if clientData.Origin != cfg.GetString("WebAuthn.rp.origin") {
		return errors.New("域名错误，当前页面域名与允许的RPID不符！")
	}

	return nil
}

func ClientDataJsonVerifyForLogin(cfg interf.ConfigInterface, redisClient *redis.Client, clientDataJSON string, sessionID int64) error {

	//0.解析Base64
	cj, err := base64.RawURLEncoding.DecodeString(clientDataJSON) //base64字符串->json字节数组
	if err != nil {
		return fmt.Errorf("base64URL解码错误: %v", err)
	}

	var clientData ClientData
	if err := json.Unmarshal(cj, &clientData); err != nil { //json字节数组 -> GO结构体
		return fmt.Errorf("clientDataJSON解析失败: %w", err)
	}

	//1.检查type是否为webauthn.get
	if clientData.Type != "webauthn.get" {
		return MyErrors.ErrorClientDataTypeIsFail
	}

	//2.检查challenge是否匹配，
	var userData User

	err = redisClient.HGetAll(context.Background(), fmt.Sprintf("SessionID:%d", sessionID)).Scan(&userData)
	if err != nil {
		return fmt.Errorf("ClinetDataJSON_Verify获取Reids数据失败: %v", err)
	}

	outTime, _ := strconv.ParseInt(userData.ChallengeOutTime, 10, 64)

	if userData.Challenge != clientData.Challenge { //判断是否失效
		return errors.New("challenge校验失败！")
	}

	fmt.Println(time.Now().Unix())
	if time.Now().Unix() > outTime {
		return errors.New("challenge过期！")
	}

	if clientData.Origin != cfg.GetString("WebAuthn.rp.origin") {
		return errors.New("域名错误，当前页面域名与允许的RPID不符！")
	}

	return nil
}
