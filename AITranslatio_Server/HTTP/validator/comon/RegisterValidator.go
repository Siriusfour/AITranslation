package comon

import (
	"AITranslatio/Global"
	"AITranslatio/HTTP/validator/validators"
	"AITranslatio/app/core/container"
)

func WebRegisterValidator() {

	containers := container.CreateContainersFactory()

	var key string

	key = Global.ValidatorPrefix + "login"
	containers.Set(key, validators.Login{})

	key = Global.ValidatorPrefix + "register"
	containers.Set(key, validators.Register{})

}
