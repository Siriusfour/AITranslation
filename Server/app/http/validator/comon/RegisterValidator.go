package comon

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/core/container"
	validators_Auth "AITranslatio/app/http/validator/validators/Auth"
	validators_Team "AITranslatio/app/http/validator/validators/Team"
)

func RegisterValidator() {

	containers := container.CreateContainersFactory()

	var key string

	key = "GetChallenge"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.GetChallengeDTO{})

	key = "Login"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.LoginDTO{})

	key = "Register"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.RegisterDTO{})

	key = "ApplicationWebAuthn"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.ApplicationWebAuthnDTO{})

	key = "RegisterWebAuthn"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.RegisterWebAuthnDTO{})

	key = "GetUserAllCredential"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.GetUserAllCredentialDTO{})

	key = "LoginByWebAuthn"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.LoginByWebAuthnDTO{})

	key = "LoginByGithub"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.LoginByOAuthDTO{"Github", "", ""})

	key = "LoginByWX"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.LoginByOAuthDTO{"WX", "", ""})

	key = "LoginByQQ"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.LoginByOAuthDTO{"QQ", "", ""})

	key = "TeamCreate"
	containers.Set(Consts.ValidatorPrefix+key, validators_Team.TeamDTO{})

	key = "Join"
	containers.Set(Consts.ValidatorPrefix+key, validators_Team.JoinDTO{})

}
