package bootstrap

import "AITranslatio/HTTP/validator/comon"

func init() {
	//1.检查项目必须的非编译目录是否存在，避免编译后调用的时候缺失相关目录

	//2.注册表单校验容器
	comon.WebRegisterValidator()

}
