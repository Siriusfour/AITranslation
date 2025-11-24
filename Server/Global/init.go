package Global

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/Global/MyErrors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	// 1.初始化程序根目录
	if curPath, err := os.Getwd(); err == nil {

		var files []string

		root := curPath
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			fmt.Println("Walk：" + file)
		}

		// 路径进行处理，兼容单元测试程序程序启动时的奇怪路径
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			Consts.BasePath = strings.Replace(strings.Replace(curPath, `\test`, "", 1), `/test`, "", 1)
		} else {
			Consts.BasePath = curPath
		}
	} else {
		log.Fatal(MyErrors.ErrorsBasePath)
	}

}
