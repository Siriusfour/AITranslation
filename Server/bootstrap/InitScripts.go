package bootstrap

import (
	"embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"path/filepath"
	"strings"
	"sync"
)

// 使用 embed.FS 嵌入整个 scripts 目录
//
//go:embed scripts/*.lua
var scriptFS embed.FS

func InitScripts(once *sync.Once) map[string]*redis.Script {
	Manager := make(map[string]*redis.Script)

	if err := loadScripts(Manager); err != nil {
		panic(fmt.Sprintf("Failed to load lua scripts: %v", err))
	}

	return Manager
}

// loadScripts 遍历嵌入的文件系统并解析
func loadScripts(scripts map[string]*redis.Script) error {
	// 读取 scripts 目录下的所有条目
	entries, err := scriptFS.ReadDir("scripts")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		if filepath.Ext(fileName) != ".lua" {
			continue
		}

		// 读取文件内容
		content, err := scriptFS.ReadFile("scripts/" + fileName)
		if err != nil {
			return err
		}

		// 生成 Key：比如 "seckill.lua" -> "seckill"
		key := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		// 创建 Redis Script 对象并存入 Map
		scripts[key] = redis.NewScript(string(content))

	}
	return nil
}
