package interf

import (
	"AITranslatio/Global/Consts"
	"AITranslatio/app/core/container"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

type ConfigInterface interface {
	ConfigFileChangeListen()
	Clone(fileName string) ConfigInterface
	Get(keyName string) interface{}
	GetString(keyName string) string
	GetBool(keyName string) bool
	GetInt(keyName string) int
	GetInt32(keyName string) int32
	GetInt64(keyName string) int64
	GetFloat64(keyName string) float64
	GetDuration(keyName string) time.Duration
	GetStringSlice(keyName string) []string
}

type ConfigFile struct {
	Viper *viper.Viper
	Mu    *sync.Mutex
}

var configContainer = container.CreateContainersFactory()
var lastChangeTime = time.Now()

func (Confige *ConfigFile) ConfigFileChangeListen() {

	Confige.Viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if changeEvent.Op.String() == "WRITE" {
				Confige.clearCache()
				lastChangeTime = time.Now()
			}
		}
	})
	Confige.Viper.WatchConfig()

}

func (Config *ConfigFile) Clone(fileName string) ConfigInterface {

	CloneConfigFile := *Config
	CloneConfigFile.Viper.SetConfigName(fileName)
	CloneConfigFile.Mu = &sync.Mutex{}
	if err := CloneConfigFile.Viper.ReadInConfig(); err != nil {
		log.Fatal("clone config is fail" + err.Error())
	}

	return &CloneConfigFile
}

func (Confige *ConfigFile) Get(keyName string) interface{} {
	if Confige.KeyIsExistsCache(keyName) {
		value := configContainer.Get(Consts.ConfigKeyPrefix + keyName)
		return value
	}

	value := Confige.Viper.Get(keyName)
	Confige.Cache(keyName, value)
	return value
}

// GetString 字符串格式返回值
func (y *ConfigFile) GetString(keyName string) string {
	if y.KeyIsExistsCache(keyName) {
		return configContainer.Get(keyName).(string)

	} else {
		value := y.Viper.GetString(keyName)
		y.Cache(keyName, value)
		return value
	}

}

// GetBool 布尔格式返回值
func (y *ConfigFile) GetBool(keyName string) bool {
	if y.KeyIsExistsCache(keyName) {
		return configContainer.Get(keyName).(bool)
	} else {
		value := y.Viper.GetBool(keyName)
		y.Cache(keyName, value)
		return value
	}
}

// GetInt 整数格式返回值
func (y *ConfigFile) GetInt(keyName string) int {
	if y.KeyIsExistsCache(keyName) {
		return configContainer.Get(keyName).(int)
	} else {
		value := y.Viper.GetInt(keyName)
		y.Cache(keyName, value)
		return value
	}
}

// GetInt32 整数格式返回值
func (y *ConfigFile) GetInt32(keyName string) int32 {
	if y.KeyIsExistsCache(keyName) {
		return configContainer.Get(keyName).(int32)
	} else {
		value := y.Viper.GetInt32(keyName)
		y.Cache(keyName, value)
		return value
	}
}

// GetInt64 整数格式返回值
func (y *ConfigFile) GetInt64(keyName string) int64 {
	if y.KeyIsExistsCache(keyName) {
		return configContainer.Get(keyName).(int64)
	} else {
		value := y.Viper.GetInt64(keyName)
		y.Cache(keyName, value)
		return value
	}
}

// GetFloat64 小数格式返回值
func (y *ConfigFile) GetFloat64(keyName string) float64 {
	if y.KeyIsExistsCache(keyName) {
		return configContainer.Get(keyName).(float64)
	} else {
		value := y.Viper.GetFloat64(keyName)
		y.Cache(keyName, value)
		return value
	}
}

// GetDuration 时间单位格式返回值
func (y *ConfigFile) GetDuration(keyName string) time.Duration {
	if y.KeyIsExistsCache(keyName) {
		return configContainer.Get(keyName).(time.Duration)
	} else {
		value := y.Viper.GetDuration(keyName)
		y.Cache(keyName, value)
		return value
	}
}

// GetStringSlice 字符串切片数格式返回值
func (y *ConfigFile) GetStringSlice(keyName string) []string {
	if y.KeyIsExistsCache(keyName) {
		return configContainer.Get(keyName).([]string)
	} else {
		value := y.Viper.GetStringSlice(keyName)
		y.Cache(keyName, value)
		return value
	}
}

//=================================================

func (Confige *ConfigFile) KeyIsExistsCache(keyName string) bool {

	_, ok := configContainer.KeyIsExists(Consts.ConfigKeyPrefix + keyName)

	if ok {
		return true
	}
	return false
}

func (Confige *ConfigFile) Cache(keyName string, value interface{}) bool {

	Confige.Mu.Lock()
	defer Confige.Mu.Unlock()

	//如果不存在于SMap，则调用set写入该key-value
	ok := Confige.KeyIsExistsCache(keyName)
	if !ok {
		configContainer.Set(Consts.ConfigKeyPrefix+keyName, value)
	}
	return true
}

func (Confige *ConfigFile) clearCache() {
	configContainer.FuzzyDelete(Consts.ConfigKeyPrefix)
}
