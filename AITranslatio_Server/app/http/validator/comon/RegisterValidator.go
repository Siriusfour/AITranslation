package comon

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/core/container"
	validators "AITranslatio/app/http/validator/validators"
)

func RegisterValidator() {

	containers := container.CreateContainersFactory()

	var key string

	key = Consts.ValidatorPrefix + "Login"
	containers.Set(key, validators.LoginDTO{})

	key = Consts.ValidatorPrefix + "Register"
	containers.Set(key, validators.RegisterDTO{})

	key = Consts.ValidatorPrefix + "WebAuthn"
	//containers.Set(key, validators.WebAuthn{})
}
