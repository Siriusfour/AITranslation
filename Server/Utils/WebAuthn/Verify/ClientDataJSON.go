package Verify

import (
	"AITranslatio/Global"
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	"context"
	"crypto/sha256"
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

// 验证标准详见  https://w3c.github.io/webauthn/#dom-publickeycredentialcreationoptions-rp
// 处理注册新凭证时第 3 到第 6 步中对客户端数据的验证
func ClinetDataJSON_Verify(WebAuthnCtx *gin.Context) ([32]byte, error) {

	ClientDataJSON := WebAuthnCtx.GetString(Consts.ValidatorPrefix + "clientDataJSON")

	//0.解析Base64
	clientDataJSON, err := base64.RawURLEncoding.DecodeString(ClientDataJSON)
	if err != nil {
		return [32]byte{}, fmt.Errorf("base64URL解码错误: %v", err)
	}
	clientDataHash := sha256.Sum256(clientDataJSON)

	var clientData ClientData
	if err := json.Unmarshal(clientDataJSON, &clientData); err != nil {
		return [32]byte{}, fmt.Errorf("clientDataJSON解析失败: %w", err)
	}

	//1.检查type是否为webauthn.create
	//协议中的Registration Step 3. and Assertion Step 7.
	if clientData.Type != "webauthn.create" {
		return [32]byte{}, MyErrors.ErrorClientDataTypeIsFail
	}

	// 协议中的Assertion Step 8 and Registration Step 4.
	//2.检查challenge是否匹配，
	var userData User

	UserID := WebAuthnCtx.GetInt64(Consts.ValidatorPrefix + "UserID")
	err = Global.RedisClient.HGetAll(context.Background(), fmt.Sprintf("%d", UserID)).Scan(&userData)
	if err != nil {
		return [32]byte{}, fmt.Errorf("ClinetDataJSON_Verify获取Reids数据失败: %v", err)
	}

	if time.Now().Unix() <= userData.challengeOutTime && userData.Challenge == clientData.Challenge { //判断是否失效
		return [32]byte{}, MyErrors.ErrorClientDataChallengeIsFail
	}

	//检查请求页面的域名是否在允许范围内，防止用户被人钓鱼 协议中的 Step 5 & Assertion Step 9
	if clientData.Origin != Global.Config.GetString("WebAuthn.rp.ID") {
		return [32]byte{}, MyErrors.ErrorClientDataRPID_IsFail
	}

	return clientDataHash, nil
}
