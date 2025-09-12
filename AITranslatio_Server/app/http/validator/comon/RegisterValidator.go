package comon

import (
	"AITranslatio/Global"
	"AITranslatio/app/core/container"
	validators2 "AITranslatio/app/http/validator/validators"
)

func WebRegisterValidator() {

	containers := container.CreateContainersFactory()

	var key string

	key = Global.ValidatorPrefix + "login"
	containers.Set(key, validators2.Login{})

	key = Global.ValidatorPrefix + "register"
	containers.Set(key, validators2.Register{})

}
