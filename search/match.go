package search

import (
	"fmt"
	"log"
)

type Result struct {
	Field   string
	Content string
}

// 声明接口
type Matcher interface {
	Search(feed *Feed, searchTerm string) ([]*Result, error)
}

func Match(matcher Matcher, feed *Feed, searchTerm string, results chan<- *Result) {
	searchResult, err := matcher.Search(feed, searchTerm)
	if err != nil {
		log.Println(err)
		return
	}
	// 将结果写入通道
	for _, result := range searchResult {
		results <- result
	}
}

func Display(results chan *Result) {
	// 通道会一直阻塞，知道有结果写入
	// 一旦通道关闭，for循环会终止
	for result := range results {
		fmt.Println("%s: %s", result.Field, result.Content)
	}
}
