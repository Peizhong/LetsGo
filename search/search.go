package search

import (
	"log"
	"sync"
)

var matchers = make(map[string]Matcher)

// Run 执行搜索逻辑
func Run(searchTerm string) {
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}

	// 无缓冲的通道，接受匹配后的结果
	results := make(chan *Result)

	// 处理所有的数据源
	var waitGroup sync.WaitGroup

	waitGroup.Add(len(feeds))

	for _, feed := range feeds {
		matcher, exists := matchers[feed.Type]
		if !exists {
			matcher = matchers["default"]
		}
		// 启动goroutine执行搜索
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)
			waitGroup.Done()
		}(matcher, feed)
	}

	// 启动goroutine监控是否所有工作完成
	go func() {
		// 等待所有任务完成
		waitGroup.Wait()
		// 关闭通道
		close(results)
	}()

	Display(results)
}

// Register dasd
func Register(feedType string, matcher Matcher) {
	if _, exist := matchers[feedType]; exist {
		log.Fatalln(feedType, "Matcher already registerd")
	}
	log.Println("Register", feedType, "matcher")
	matchers[feedType] = matcher
}
