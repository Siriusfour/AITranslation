package comon

import (
	"AITranslatio/Global"
	"AITranslatio/app/core/container"
)

func webRegisterValidator() {

	containers := container.CreateContainersFactory()

	var key string

	key = Global.ValidatorPrefix + "login"
	containers.Set(key, 1)

}
