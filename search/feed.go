package search

import (
	"encoding/json"
	"os"
)

const dataFile = "config/data.json"

// 声明结构类型
type Feed struct {
	// ``里的内容是标记(tag), 描述了json
	Name string `json:"site"`
	URI  string `json:"link"`
	Type string `json:"type"`
}

func RetrieveFeeds() ([]*Feed, error) {
	var feeds []*Feed
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	// 函数返回时，关闭文件
	defer file.Close()

	// 文件解析成切片
	// ?func (dec *Decoder) Decode(v interface{}) error
	err = json.NewDecoder(file).Decode(&feeds)
	return feeds, err
}
