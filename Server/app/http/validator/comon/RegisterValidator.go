package comon

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/core/container"
	validators_Auth "AITranslatio/app/http/validator/validators"
	validators_Team "AITranslatio/app/http/validator/validators/Team"
)

func RegisterValidator() {

	containers := container.CreateContainersFactory()

	var key string

	key = "Login"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.LoginDTO{})

	key = "Register"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.RegisterDTO{})

	key = "WebAuthn"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.WebAuthnDTO{})

	key = "VerifyWebAuthn"
	containers.Set(Consts.ValidatorPrefix+key, validators_Auth.VerifyWebAuthnDTO{})

	key = "TeamCreate"
	containers.Set(Consts.ValidatorPrefix+key, validators_Team.TeamDTO{})

	key = "Join"
	containers.Set(Consts.ValidatorPrefix+key, validators_Team.JoinDTO{})
}
