package main

import (
	"AITranslatio/Router"
	_ "AITranslatio/bootstrap"
	"github.com/gin-gonic/gin"
)

func main() {

	router := Router.InitRouter()
	_ = router.Run(":3008")

	gin.ForceConsoleColor()
}
