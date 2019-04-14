package framework

import (
	"encoding/json"
	"os"
	"sync"
)

type logging struct {
	LogLevel map[string]string
}

type Appsettings struct {
	Logging           logging
	ConnectionStrings map[string]string
}

const defaultSettingFile = "config/appsettings.json"

var once sync.Once

var settings Appsettings

// GetConnectionString get connection string from appsettings.json
func GetConnectionString(key string) string {
	once.Do(func() {
		file, err := os.Open(defaultSettingFile)
		// 函数返回时，关闭文件
		defer file.Close()
		if err == nil {
			json.NewDecoder(file).Decode(&settings)
		}
	})
	// level := settings.Logging.LogLevel["Default"]
	// _ = level
	value := settings.ConnectionStrings[key]
	return value
}

func GetAppsettings(settingFile string) (s Appsettings) {
	if settingFile == "" {
		settingFile = defaultSettingFile
	}
	file, err := os.Open(settingFile)
	// 函数返回时，关闭文件
	defer file.Close()
	if err == nil {
		json.NewDecoder(file).Decode(&s)
	}
	return
}
