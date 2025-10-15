package AuthService

import (
	"AITranslatio/Global/MyErrors"
	"AITranslatio/Utils/WebAuthn"
	WebAuthn_Verify "AITranslatio/Utils/WebAuthn/Verify"
	"AITranslatio/app/DAO/UserDAO"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

func (Service *AuthService) ApplicationWebAuthn(UserID int64) (*WebAuthn.WebAuthn, error) {

	//获取userName和Email
	var UseName, Email string
	err := UserDAO.CreateDAOFactory("mysql").
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

func (Service *AuthService) VerifyWebAuthn(WebAuthnCtx *gin.Context) error {

	clientDataHash, err := WebAuthn_Verify.ClinetDataJSON_Verify(WebAuthnCtx)
	fmt.Println(clientDataHash)
	if err != nil {
		return fmt.Errorf("ClinetDataJSON校验错误%w", err)
	}
	return nil
}
