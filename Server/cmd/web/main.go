package main

import (
	"AITranslatio/Router"
	_ "AITranslatio/bootstrap"
)

func main() {

	router := Router.InitRouter()
	_ = router.Run(":3008")

}
