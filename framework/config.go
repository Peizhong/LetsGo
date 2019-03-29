package framework

import (
	"encoding/json"
	"os"
	"sync"
)

type logging struct {
	LogLevel map[string]string
}

type appsettings struct {
	Logging           logging
	ConnectionStrings map[string]string
}

const dataFile = "config/appsettings.json"

var once sync.Once

var settings appsettings

// GetConnectionString get connection string from appsettings.json
func GetConnectionString(key string) string {
	once.Do(func() {
		file, err := os.Open(dataFile)
		// 函数返回时，关闭文件
		defer file.Close()
		if err != nil {
		} else {
			json.NewDecoder(file).Decode(&settings)
		}
	})
	// level := settings.Logging.LogLevel["Default"]
	// _ = level
	value := settings.ConnectionStrings[key]
	return value
}
