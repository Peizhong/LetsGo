package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

// todo: 文件监视: watch file, delivering events to a channel

type logging struct {
	LogLevel map[string]string
}

type Appsettings struct {
	Logging           logging
	ConnectionStrings map[string]string
}

const defaultSettingFile = "config/appsettings.json"

var once sync.Once

var settings atomic.Value

// GetConnectionString get connection string from appsettings.json
func GetConnectionString(key string) string {
	s := GetAppsettings(defaultSettingFile)
	value := s.ConnectionStrings[key]
	return value
}

func GetAppsettings(settingFile string) (s Appsettings) {
	once.Do(func() {
		if settingFile == "" {
			settingFile = defaultSettingFile
		}
		settingFile = filepath.ToSlash(settingFile)
		file, err := os.Open(settingFile)
		// 函数返回时，关闭文件
		defer file.Close()
		if err == nil {
			s := Appsettings{}
			err := json.NewDecoder(file).Decode(&s)
			if err == nil {
				settings.Store(s)
			}
		}
	})
	return settings.Load().(Appsettings)
}
