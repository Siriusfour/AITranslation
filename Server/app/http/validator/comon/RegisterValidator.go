package comon

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/core/container"
	validators "AITranslatio/app/http/validator/validators"
)

func RegisterValidator() {

	containers := container.CreateContainersFactory()

	var key string

	key = "Login"
	containers.Set(Consts.ValidatorPrefix+key, validators.LoginDTO{})

	key = "Register"
	containers.Set(Consts.ValidatorPrefix+key, validators.RegisterDTO{})

	key = "WebAuthn"
	containers.Set(Consts.ValidatorPrefix+key, validators.WebAuthnDTO{})
}
