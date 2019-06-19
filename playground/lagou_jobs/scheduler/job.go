package scheduler

import (
	"github.com/emirpasic/gods/stacks/arraystack"
	"sync"
)

var (
	mutex     sync.Mutex
	JobParams = arraystack.New()
)

type JobParam struct {
	City string
	Pn   int
	Kd   string
}

func NewJobScheduler() *JobParam {
	return &JobParam{}
}

func (j *JobParam) Pop() *JobParam {
	mutex.Lock()
	job, _ := JobParams.Pop()
	mutex.Unlock()
	return job.(*JobParam)
}

func (j *JobParam) Append(city string, pn int, kd string) {
	mutex.Lock()
	JobParams.Push(&JobParam{
		City: city,
		Pn:   pn,
		Kd:   kd,
	})
	mutex.Unlock()
}

func (j *JobParam) Count() int {
	return JobParams.Size()
}
