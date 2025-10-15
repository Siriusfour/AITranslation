package eventManage

import (
	"AITranslatio/Global"
	"AITranslatio/Global/MyErrors"
	"strings"
	"sync"
)

var sMap sync.Map

type EventManage interface {
	Set() bool
	Get() interface{}
	Call()
	Delete()
	FuzzyCall()
}

func createEventManageFactory() *eventManage {

	return &eventManage{}
}

type eventManage struct {
}

func (e *eventManage) Set(key string, value interface{}) bool {

	if _, exist := e.Get(key); exist == true {

		Global.Logger.Error(MyErrors.ErrorsFuncEventAlreadyExists, "keyName:"+key)

		return false

	} else {
		sMap.Store(key, value)
	}
	return false
}

func (e *eventManage) Get(key string) (interface{}, bool) {

	if value, ok := sMap.Load(key); ok {
		return value, true
	}

	return nil, false
}

func (e *eventManage) Delete(key string) {
	sMap.Delete(key)
}

func (e *eventManage) Call(key string, arg ...interface{}) {

	if fn, ok := e.Get(key); ok {
		if fn, ok := fn.(func(...interface{})); ok {
			fn(arg...)
		} else {
			Global.Logger.Error(MyErrors.ErrorsFuncEventNotCall, "key:"+key+"-相关函数无法调用")
		}

	} else {
		Global.Logger.Error(MyErrors.ErrorsFuncEventNotRegister, "key:"+key+"-相关函数未注册")

	}

}

func (e *eventManage) FuzzyCall(keyPre string) {

	sMap.Range(func(key, value interface{}) bool {
		if keyName, ok := key.(string); ok {
			if strings.HasPrefix(keyName, keyPre) {
				e.Call(keyName)
			}
		}
		return true
	})

}
