package main

import (
	"github.com/peizhong/letsgo/playground/lagou_jobs/pipeline"
	"github.com/peizhong/letsgo/playground/lagou_jobs/spider"
	"log"
	"sync"
)

var (
	kds = []string{
		"golang",
	}
	citys = []string{
		"深圳",
	}

	initResults = []spider.InitResult{}
	loopResults = []spider.LoopResult{}
	jobPipeline = pipeline.NewJobPipeline()

	wg sync.WaitGroup
)

func main() {
	for _, kd := range kds {
		for _, city := range citys {
			wg.Add(1)
			go func(city string, kd string) {
				defer wg.Done()
				initResult, err := spider.InitJobs(city, 1, kd)
				if err != nil {
					log.Fatalln(err)
				}

				initResults = append(initResults, initResult...)
				loopResults = append(loopResults, spider.LoopJobs())
			}(city, kd)
		}
	}

	wg.Wait()

	jobPipeline.Push()

	log.Printf("Init Results: %v", initResults)
	log.Printf("Loop Results: %v", loopResults)
}
