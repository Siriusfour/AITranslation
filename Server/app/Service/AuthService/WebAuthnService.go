package AuthService

import (
	"AITranslatio/Global/CustomErrors"
	"AITranslatio/Utils/WebAuthn"
	"AITranslatio/app/DAO/UserDAO"
	"errors"
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
	w := WebAuthn.CreateWebAuthnConfigFactory(UseName, Email)
	err = w.CreateChallenge(UserID)
	if err != nil {
		return nil, errors.New(CustomErrors.ErrorGetChallengeIsFail + err.Error())
	}

	return w, nil
}
