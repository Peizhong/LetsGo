package framework

import (
	"encoding/json"
	"os"
)

type logging struct {
	LogLevel map[string]string
}

type appsettings struct {
	Logging           logging
	ConnectionStrings map[string]string
}

const dataFile = "config/appsettings.json"

func GetConnectionString(key string) (string, error) {
	var settings appsettings
	file, err := os.Open(dataFile)
	if err != nil {
		return "", err
	}
	// 函数返回时，关闭文件
	defer file.Close()

	// 文件解析成切片
	// ?func (dec *Decoder) Decode(v interface{}) error
	err = json.NewDecoder(file).Decode(&settings)
	if err != nil {
		return "", err
	}
	level := settings.Logging.LogLevel["Default"]
	_ = level
	value := settings.ConnectionStrings[key]
	return value, nil
}
