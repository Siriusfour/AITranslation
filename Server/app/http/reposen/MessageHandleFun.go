package reposen

import (
	"AITranslatio/Global"
	"errors"
	"fmt"
	"strings"
)

func MessageHandle(err error) string {

	//1.判断当前模式
	if Global.GetInfra().Config.GetBool("Mode.Product") { //生产模式，只返回顶层信息 ，如："登录失败！“

		return topMessage(err)
	} else { //开发模式，返回全链错误 ，如："登录失败！: service error:  XXXX :authDAO error : XXXX“

		return fullChain(err)
	}

}

func topMessage(err error) string {
	if err == nil {
		return ""
	}
	// 如果有内层错误，把它切掉，只保留外层信息
	if inner := errors.Unwrap(err); inner != nil {
		// %v 会返回 "外层信息: 内层错误"
		s := fmt.Sprintf("%v", err)
		innerMsg := fmt.Sprintf("%v", inner)
		// 去掉拼接上去的部分，只保留外层信息
		if idx := strings.Index(s, innerMsg); idx > 0 {
			return strings.TrimSpace(strings.TrimSuffix(s[:idx], ":"))
		}
		return s
	}
	// 没有内层错误，直接返回本身
	return fmt.Sprintf("%v", err)

}

func fullChain(err error) string {
	var parts []string
	for e := err; e != nil; e = errors.Unwrap(e) {
		parts = append(parts, firstLine(e))
	}
	// 从外到内，或反过来都可；这里外->内
	return strings.Join(parts, ": ")
}

func firstLine(err error) string {
	// 避免多行错误污染输出（有些库会返回多行）
	s := fmt.Sprintf("%v", err)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
