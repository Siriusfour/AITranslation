package Global

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/CustomErrors"
	"log"
	"os"
	"strings"
)

func init() {
	// 1.初始化程序根目录
	if curPath, err := os.Getwd(); err == nil {
		// 路径进行处理，兼容单元测试程序程序启动时的奇怪路径
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			Consts.BasePath = strings.Replace(strings.Replace(curPath, `\test`, "", 1), `/test`, "", 1)
		} else {
			Consts.BasePath = curPath
		}
	} else {
		log.Fatal(CustomErrors.ErrorsBasePath)
	}
}
